package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
	vabra "github.com/lunfardo314/goq/abra/validate"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/utils"
)

type Function struct {
	module            *QuplaModule
	Analyzed          bool // finished analysis
	Joins             map[string]int
	Affects           map[string]int
	Name              string
	retSize           int
	RetExpr           ExpressionInterface
	Sites             []*QuplaSite
	NumParams         int  // idx < NumParams represents parameter, idx >= represents local var (assign)
	BufLen            int  // total length of the local var buffer
	HasStateVariables bool // if has state vars itself
	hasState          bool // if directly or indirectly references those with state vars
	InSize            int
	ParamSizes        []int
	traceLevel        int
	nextCallIndex     byte
	StateHashMap      *StateHashMap
	expandedInline    utils.StringSet // needed for optimisation, prevention of recursions while expanding inline
}

func NewFunction(name string, size int, module *QuplaModule) *Function {
	return &Function{
		module:         module,
		Name:           name,
		retSize:        size,
		Sites:          make([]*QuplaSite, 0, 10),
		Joins:          make(map[string]int),
		Affects:        make(map[string]int),
		ParamSizes:     make([]int, 0, 5),
		expandedInline: make(utils.StringSet),
	}
}

func (def *Function) AppendInline(s string) {
	def.expandedInline.Append(s)
}

func (def *Function) WasInline(s string) bool {
	return def.expandedInline.Contains(s)
}

func (def *Function) NextCallIndex() byte {
	ret := def.nextCallIndex
	if ret == 0xFF {
		panic("can't be more than 2^16 function calls within function body")
	}
	def.nextCallIndex++
	return ret
}

func (def *Function) SetTraceLevel(traceLevel int) {
	def.traceLevel = traceLevel
}

func (def *Function) HasState() bool {
	return def.HasStateVariables || def.hasState
}

func (def *Function) References(funName string) bool {
	for _, vi := range def.Sites {
		if vi.Assign != nil && vi.Assign.References(funName) {
			return true
		}
	}
	return def.RetExpr.References(funName)
}

func (def *Function) Size() int {
	return def.retSize
}

func (def *Function) ArgSize() int {
	return def.InSize
}

func (def *Function) HasEnvStmt() bool {
	return len(def.Joins) > 0 || len(def.Affects) > 0
}

func (def *Function) GetJoinEnv() map[string]int {
	return def.Joins
}

func (def *Function) GetAffectEnv() map[string]int {
	return def.Affects
}

func (def *Function) GetVarIdx(name string) int {
	for i, lv := range def.Sites {
		if lv.Name == name {
			return i
		}
	}
	return -1
}

func (def *Function) VarByIdx(idx int) (*QuplaSite, error) {
	if idx < 0 || idx >= len(def.Sites) {
		return nil, fmt.Errorf("worng var idx %v", idx)
	}
	return def.Sites[idx], nil
}

func (def *Function) VarByName(name string) (*QuplaSite, error) {
	idx := def.GetVarIdx(name)
	if idx < 0 {
		return nil, fmt.Errorf("can't finc variabe with name '%v'", name)
	}
	return def.VarByIdx(idx)
}

func (def *Function) CheckArgSizes(args []ExpressionInterface) error {
	for i := range args {
		if i >= def.NumParams || args[i].Size() != def.Sites[i].Size {
			return fmt.Errorf("param and arg # %v mismach in %v", i, def.Name)
		}
	}
	return nil
}

// mock expression with all null arguments
func (def *Function) NewFuncExpressionWithNulls(callIndex byte) *FunctionExpr {
	ret := NewFunctionExpr("", def, callIndex)

	offset := 0
	for _, sz := range def.ParamSizes {
		ret.AppendSubExpr(NewNullExpr(sz))
		offset += sz
	}
	return ret
}

func (def *Function) Eval(frame *EvalFrame, result Trits) bool {
	//if def.Name == "golGen" {
	//	fmt.Printf("++++ kuku\n")
	//}
	null := def.RetExpr.Eval(frame, result)
	if def.traceLevel > 0 {
		if !null {
			//bi, _ := utils.TritsToBigInt(result)
			res := utils.TritsToString(result)
			reslen := len(res)
			if reslen > 100 {
				res = res[:100] + "..."
			}
			Logf(def.traceLevel, "trace '%v': returned '%v' len = %v",
				def.Name, res, reslen)
		} else {
			Logf(2+def.traceLevel, "trace '%v': returned null", def.Name)
		}
	}
	return null
}

func (def *Function) Stats() map[string]int {
	ret := make(map[string]int)
	for _, site := range def.Sites {
		if !site.IsParam {
			countTypesInExpression(site.Assign, ret)
		}
	}
	countTypesInExpression(def.RetExpr, ret)
	return ret
}

func countTypesInExpression(expr ExpressionInterface, stats map[string]int) {
	t := fmt.Sprintf("%T", expr)
	if _, ok := stats[t]; !ok {
		stats[t] = 0
	}
	stats[t]++
	for _, se := range expr.GetSubexpressions() {
		countTypesInExpression(se, stats)
	}
}

func (def *Function) GetLookupName() string {
	return "qupla_function_" + def.Name
}

func (def *Function) GetAbraBranchBlock(codeUnit *abra.CodeUnit) *abra.Block {
	//if def.Name == "arcRadixLeaf_243_8019" {
	//	fmt.Printf("kuku\n")
	//}
	lookupName := def.GetLookupName()
	ret := cabra.FindBranchBlock(codeUnit, lookupName)
	if ret != nil {
		return ret
	}
	ret = cabra.MustAddNewBranchBlock(codeUnit, lookupName, def.Size())

	for _, vi := range def.Sites {
		if vi.IsParam {
			cabra.AddInputSite(ret.Branch, vi.Size)
		}
	}
	if concatExpr, ok := def.RetExpr.(*ConcatExpr); ok {
		for _, se := range concatExpr.subExpr {
			site := se.GetAbraSite(ret.Branch, codeUnit, "")
			cabra.ChangeSiteType(site, abra.SITE_OUTPUT)
		}
	} else {
		singleOutput := def.RetExpr.GetAbraSite(ret.Branch, codeUnit, "")
		if singleOutput.SiteType == abra.SITE_BODY {
			cabra.ChangeSiteType(singleOutput, abra.SITE_OUTPUT)
		} else {
			singleOutput = cabra.NewMergeSite(def.RetExpr.Size(), "", singleOutput)
			cabra.ChangeSiteType(singleOutput, abra.SITE_OUTPUT)
			cabra.AddOrUpdateSite(ret.Branch, singleOutput)
		}
	}
	// finalize with state sites
	// for each state var site generate abra site and update temporary site with the new
	for _, vi := range def.Sites {
		if vi.IsState {
			vi.Assign.GetAbraSite(ret.Branch, codeUnit, vi.GetAbraLookupName())
		}
	}
	vabra.AssertValid(ret.Branch)
	return ret
}

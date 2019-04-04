package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
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
	LocalVars         []*QuplaSite
	NumParams         int  // idx < NumParams represents parameter, idx >= represents local var (assign)
	BufLen            int  // total length of the local var buffer
	HasStateVariables bool // if has state vars itself
	hasState          bool // if directly or indirectly references those with state vars
	isRecursive       bool // is directly or indirectly recursive
	InSize            int
	ParamSizes        []int
	traceLevel        int
	nextCallIndex     uint8
	StateHashMap      *StateHashMap
}

func NewFunction(name string, size int, module *QuplaModule) *Function {
	return &Function{
		module:     module,
		Name:       name,
		retSize:    size,
		LocalVars:  make([]*QuplaSite, 0, 10),
		Joins:      make(map[string]int),
		Affects:    make(map[string]int),
		ParamSizes: make([]int, 0, 5),
	}
}

func (def *Function) NextCallIndex() uint8 {
	if def == nil {
		return 0
	}
	ret := def.nextCallIndex
	if ret == 0xFF {
		panic("can't be more than 256 function calls within function body")
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
	for _, vi := range def.LocalVars {
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
	for i, lv := range def.LocalVars {
		if lv.Name == name {
			return i
		}
	}
	return -1
}

func (def *Function) VarByIdx(idx int) (*QuplaSite, error) {
	if idx < 0 || idx >= len(def.LocalVars) {
		return nil, fmt.Errorf("worng var idx %v", idx)
	}
	return def.LocalVars[idx], nil
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
		if i >= def.NumParams || args[i].Size() != def.LocalVars[i].Size {
			return fmt.Errorf("param and arg # %v mismach in %v", i, def.Name)
		}
	}
	return nil
}

// mock expression with all null arguments
func (def *Function) NewFuncExpressionWithNulls(callIndex uint8) *FunctionExpr {
	ret := NewFunctionExpr("", def, callIndex)

	offset := 0
	for _, sz := range def.ParamSizes {
		ret.AppendSubExpr(NewNullExpr(sz))
		offset += sz
	}
	return ret
}

func (def *Function) Eval(frame *EvalFrame, result Trits) bool {
	null := def.RetExpr.Eval(frame, result)
	if def.traceLevel > 0 {
		if !null {
			bi, _ := utils.TritsToBigInt(result)
			Logf(def.traceLevel, "trace '%v': returned %v, '%v'",
				def.Name, bi, utils.TritsToString(result))
		} else {
			Logf(2+def.traceLevel, "trace '%v': returned null")
		}
	}
	return null
}

func (def *Function) Optimize() {
	if Config.OptimizeOneTimeSites {
		before := def.ZeroInternalSites()
		def.RetExpr = def.optimizeOneTimeSites(def.RetExpr)

		_, _, _, numVars, numUnusedVars := def.NumSites()
		Logf(5, "Optimized %v sites out of %v in '%v'", numUnusedVars, numVars, def.Name)
		after := def.ZeroInternalSites()
		if !before && after {
			Logf(5, "'%v' became inlineable", def.Name)
		}
	}
	if Config.OptimizeInlineSlices {
		def.RetExpr = def.optimizeInlineSlices(def.RetExpr)
	}
}

func (def *Function) optimizeOneTimeSites(expr ExpressionInterface) ExpressionInterface {
	sliceExpr, ok := expr.(*SliceExpr)
	if !ok {
		subExpr := make([]ExpressionInterface, 0)
		for _, se := range expr.GetSubexpressions() {
			opt := def.optimizeOneTimeSites(se)
			subExpr = append(subExpr, opt)
		}
		expr.SetSubexpressions(subExpr)
		return expr
	}
	if sliceExpr.vi.IsState || sliceExpr.vi.IsParam || sliceExpr.vi.numUses > 1 {
		return expr
	}
	// slice expressions optimized to SliceInline
	opt := def.optimizeOneTimeSites(def.LocalVars[sliceExpr.vi.Idx].Assign)
	sliceExpr.vi.notUsed = true
	return NewSliceInline(sliceExpr, opt)
}

func (def *Function) optimizeInlineSlices(expr ExpressionInterface) ExpressionInterface {
	if _, ok := expr.(*SliceInline); ok {
		def.module.IncStat("numOptimizedInlineSlices")
	}
	return optimizeInlineSlicesExpr(expr)
}

// returns numSites, numParam, numState, numVars, numUnusedVars
func (def *Function) NumSites() (int, int, int, int, int) {
	var numSites, numParam, numState, numVars, numUnusedVars int
	for _, vi := range def.LocalVars {
		numSites++
		if vi.IsParam {
			numParam++
		}
		if vi.IsState {
			numState++
		}
		if !vi.IsParam && !vi.IsState {
			numVars++
			if vi.notUsed {
				numUnusedVars++
			}
		}
	}
	return numSites, numParam, numState, numVars, numUnusedVars
}

func (def *Function) ZeroInternalSites() bool {
	_, _, _, numVars, numUnusedVars := def.NumSites()
	return numVars == numUnusedVars
}

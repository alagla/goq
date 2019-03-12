package qupla

import (
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/utils"
	. "github.com/lunfardo314/quplayaml/go"
)

type QuplaFuncDef struct {
	yamlSource        *QuplaFuncDefYAML // needed for analysis phase only
	module            ModuleInterface
	joins             StringSet
	affects           StringSet
	name              string
	retSize           int64
	retExpr           ExpressionInterface
	localVars         []*VarInfo
	numParams         int64 // idx < numParams represents parameter, idx >= represents local var (assign)
	bufLen            int64 // total length of the local var buffer
	hasStateVariables bool
	hasState          bool
	argSize           int64
	argSizes          []int64
}

func (def *QuplaFuncDef) SetName(name string) {
	def.name = name
}

func (def *QuplaFuncDef) GetName() string {
	if def == nil {
		return "(undef)"
	}
	return def.name
}

func (def *QuplaFuncDef) HasState() bool {
	return def.hasStateVariables || def.hasState
}

func (def *QuplaFuncDef) References(funName string) bool {
	for _, vi := range def.localVars {
		if vi.Assign != nil && vi.Assign.References(funName) {
			return true
		}
	}
	return def.retExpr.References(funName)
}

func AnalyzeFuncDef(name string, defYAML *QuplaFuncDefYAML, module *QuplaModule) error {
	var err error
	defer func(perr *error) {
		if *perr != nil {
			errorf("Error while analyzing func def '%v': %v", name, *perr)
		}
	}(&err)

	module.IncStat("numFuncDef")

	def := &QuplaFuncDef{
		yamlSource: defYAML,
		module:     module,
		name:       name,
		argSizes:   make([]int64, 0, len(defYAML.Params)),
	}
	def.AnalyzeEnvironmentStatements()
	if def.HasEnvStmt() {
		def.module.IncStat("numEnvFundef")
	}
	// return size. Must be const expression
	// this must be first because in recursive calls return size must be known
	// scope must be nil because const value do not have scope
	ce, err := module.AnalyzeExpression(defYAML.ReturnType, nil)
	if err != nil {
		return err
	}
	var sz int64
	if sz, err = GetConstValue(ce); err != nil {
		return err
	}
	// function def in the module is always with valid retSize (other parts not analyzed yet)
	def.retSize = sz
	module.AddFuncDef(name, def)

	// build var scope
	if err = def.createVarScope(); err != nil {
		return err
	}
	if err = def.analyzeAssigns(); err != nil {
		return err
	}
	if err = def.finalizeLocalVars(); err != nil {
		return err
	}
	// return expression
	if def.retExpr, err = module.AnalyzeExpression(defYAML.ReturnExpr, def); err != nil {
		return err
	}
	if def.retExpr == nil {
		return fmt.Errorf("in funcdef '%v': return expression can't be nil", def.name)
	}
	return nil
}

func (def *QuplaFuncDef) Size() int64 {
	return def.retSize
}

func (def *QuplaFuncDef) ArgSize() int64 {
	return def.argSize
}

func (def *QuplaFuncDef) AnalyzeEnvironmentStatements() {
	for _, envYAML := range def.yamlSource.Env {
		if envYAML.Join {
			if def.joins == nil {
				def.joins = make(StringSet)
				def.joins.Append(envYAML.Name)
			}
			def.module.IncStat("numEnvJoin")
		} else {
			if def.affects == nil {
				def.affects = make(StringSet)
				def.affects.Append(envYAML.Name)
			}
			def.module.IncStat("numEnvAffect")
		}
	}
}

func (def *QuplaFuncDef) HasEnvStmt() bool {
	return len(def.joins) > 0 || len(def.affects) > 0
}

func (def *QuplaFuncDef) GetJoinEnv() StringSet {
	return def.joins
}

func (def *QuplaFuncDef) GetAffectEnv() StringSet {
	return def.affects
}

func (def *QuplaFuncDef) GetVarIdx(name string) int64 {
	for i, lv := range def.localVars {
		if lv.Name == name {
			return int64(i)
		}
	}
	return -1
}

func (def *QuplaFuncDef) VarByIdx(idx int64) *VarInfo {
	if idx < 0 || idx >= int64(len(def.localVars)) {
		return nil
	}
	return def.localVars[idx]
}

func (def *QuplaFuncDef) VarByName(name string) *VarInfo {
	return def.VarByIdx(def.GetVarIdx(name))
}

func (def *QuplaFuncDef) GetVarInfo(name string) (*VarInfo, error) {
	ret := def.VarByName(name)
	if ret == nil {
		return nil, nil
	}
	if ret.Analyzed {
		return ret, nil
	}
	var err error
	ret.Analyzed = true

	if ret.IsParam {
		// param
		ret.Assign = nil
	} else {
		// local var (can be state)
		realVarName := name
		e, ok := def.yamlSource.Assigns[realVarName]
		if !ok {
			return nil, fmt.Errorf("inconsistency with vars")
		}
		if ret.Assign, err = def.module.AnalyzeExpression(e, def); err != nil {
			return ret, err
		}
		if ret.IsState {
			if ret.Size != ret.Assign.Size() {
				return nil, fmt.Errorf("expression and state variable has different sizes in the assign")
			}
		} else {
			ret.Size = ret.Assign.Size()
		}
	}
	return ret, nil
}

func (def *QuplaFuncDef) createVarScope() error {
	src := def.yamlSource
	def.localVars = make([]*VarInfo, 0, len(src.Params)+len(src.Assigns))
	// function parameters (first numParams)
	def.numParams = int64(len(src.Params))
	for idx, arg := range src.Params {
		if def.GetVarIdx(arg.ArgName) >= 0 {
			return fmt.Errorf("duplicate arg name '%v'", arg.ArgName)
		}
		def.localVars = append(def.localVars, &VarInfo{
			Idx:      int64(idx),
			Name:     arg.ArgName,
			Size:     arg.Size,
			Analyzed: true,
			IsParam:  true,
			IsState:  false,
		})
	}
	// the rest of indices belong to local vars (incl state)
	// state variables
	var idx int64
	def.hasStateVariables = len(src.State) > 0
	for name, s := range src.State {
		idx = def.GetVarIdx(name)
		if idx >= 0 {
			return fmt.Errorf("wrong declared state variable: '%v' in '%v'", name, def.name)
		} else {
			// for old value
			def.localVars = append(def.localVars, &VarInfo{
				Idx:     int64(len(def.localVars)),
				Name:    name,
				Size:    s.Size,
				IsState: true,
			})
		}
		def.module.IncStat("numStateVars")
	}
	// variables defined by assigns
	var vi *VarInfo
	for name := range src.Assigns {
		vi = def.VarByName(name)
		if vi != nil {
			if vi.IsParam {
				return fmt.Errorf("cannot assign to function parameter: '%v' in '%v'", name, def.name)
			}
			if !vi.IsState {
				return fmt.Errorf("several assignment to the same var '%v' in '%v' is not allowed", name, def.name)
			}
		} else {
			def.localVars = append(def.localVars, &VarInfo{
				Idx:     int64(len(def.localVars)),
				Name:    name,
				Size:    0, // unknown yet
				IsState: false,
				IsParam: false,
			})
		}
	}
	return nil
}

func (def *QuplaFuncDef) analyzeAssigns() error {
	var err error
	for name := range def.yamlSource.Assigns {
		// GetVarInfo analyzes expression if necessary
		if _, err = def.GetVarInfo(name); err != nil {
			return err
		}
	}
	return nil
}

func (def *QuplaFuncDef) finalizeLocalVars() error {
	var curOffset int64
	def.argSize = 0
	for _, v := range def.localVars {
		if v.Size == 0 {
			v.Size = v.Assign.Size()
		}
		if v.Size == 0 {
			return fmt.Errorf("can't determine var size '%v': '%v'", v.Name, def.GetName())
		}
		v.Offset = curOffset
		curOffset += v.Size
		if !v.IsParam {
			if v.Assign == nil {
				return fmt.Errorf("variable '%v' in '%v' is not assigned", v.Name, def.GetName())
			}
		} else {
			def.argSize += v.Size
			def.argSizes = append(def.argSizes, v.Size)
		}
		if v.Assign != nil && v.Assign.Size() != v.Size {
			return fmt.Errorf("sizes doesn't match for var '%v' in '%v'", v.Name, def.GetName())
		}
	}
	def.bufLen = int64(curOffset)

	if def.hasStateVariables {
		def.module.IncStat("numStatefulFunDef")
	}
	return nil
}

func (def *QuplaFuncDef) checkArgSizes(args []ExpressionInterface) error {
	for i := range args {
		if int64(i) >= def.numParams || args[i].Size() != def.localVars[i].Size {
			return fmt.Errorf("param and arg # %v mismach in %v", i, def.GetName())
		}
	}
	return nil
}

func (def *QuplaFuncDef) NewExpressionWithArgs(args trinary.Trits) (ExpressionInterface, error) {
	if def.argSize != int64(len(args)) {
		return nil, fmt.Errorf("size mismatch: fundef '%v' has arg size %v, trit vector's size = %v",
			def.GetName(), def.ArgSize(), len(args))
	}
	ret := NewQuplaFuncExpr("", def)

	offset := int64(0)
	for _, sz := range def.argSizes {
		e := NewQuplaValueExpr(args[offset : offset+sz])
		ret.AppendSubExpr(e)
		offset += sz
	}
	return ret, nil
}

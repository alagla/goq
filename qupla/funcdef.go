package qupla

import (
	"fmt"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaFuncDef struct {
	yamlSource *QuplaFuncDefYAML // needed for analysis phase only
	name       string
	retSize    int64
	retExpr    ExpressionInterface
	localVars  []*LocalVariable
	numParams  int   // idx < numParams represents parameter, idx >= represents local var (assign)
	bufLen     int64 // total length of the local var buffer
}

// represents local variable in func def
type LocalVariable struct {
	name     string
	isState  bool
	offset   int64               // offset in call context buffer
	size     int64               // size of the variable
	expr     ExpressionInterface // for assigns only
	analyzed bool
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
		name:       name,
	}
	// return size. Must be const expression
	// this must be first because in recursive calls return size must be known
	// scope must be nil because const value do not scope
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
	if err = def.analyzeAssigns(defYAML, module); err != nil {
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

func (def *QuplaFuncDef) GetVarIdx(name string) int {
	for i, lv := range def.localVars {
		if lv.name == name {
			return i
		}
	}
	return -1
}

func (def *QuplaFuncDef) InScope(name string) bool {
	return def.GetVarIdx(name) >= 0
}

func (def *QuplaFuncDef) VarByIdx(idx int) *LocalVariable {
	if idx < 0 {
		return nil
	}
	return def.localVars[idx]
}

// it tries to find var idx. Analyzes var if not analyzed yet
func (def *QuplaFuncDef) FindVarIdx(name string, module *QuplaModule) (int, error) {
	idx := def.GetVarIdx(name)
	if idx < 0 {
		return -1, nil
	}
	ret := def.localVars[idx]
	if ret.analyzed {
		return idx, nil
	}
	var err error

	ret.analyzed = true

	if idx >= def.numParams {
		// local var
		v := def.localVars[idx]
		e, ok := def.yamlSource.Assigns[name]
		if !ok {
			return -1, fmt.Errorf("inconsistency with vars")
		}
		if ret.expr, err = module.AnalyzeExpression(e, def); err != nil {
			return idx, err
		}
		if v.isState {
			if ret.size != ret.expr.Size() {
				return -1, fmt.Errorf("expression and state variable has different sizes in the assign")
			}
		} else {
			ret.size = ret.expr.Size()
		}
	} else {
		// param
		def.localVars[idx].expr = nil
	}
	return idx, nil
}

func (def *QuplaFuncDef) createVarScope() error {
	src := def.yamlSource
	def.localVars = make([]*LocalVariable, 0, len(src.Params)+len(src.Assigns))
	// first numParams indices belong to parameters
	def.numParams = len(src.Params)
	for _, arg := range src.Params {
		if def.InScope(arg.ArgName) {
			return fmt.Errorf("duplicate arg name '%v'", arg.ArgName)
		}
		def.localVars = append(def.localVars, &LocalVariable{
			name:     arg.ArgName,
			size:     arg.Size,
			analyzed: true,
		})
	}
	// the rest indices belong to local vars (incl state)
	var idx int
	for name, s := range src.State {
		idx = def.GetVarIdx(name)
		if idx >= 0 {
			return fmt.Errorf("wrong declared state variable: '%v' in '%v'", name, def.name)
		} else {
			def.localVars = append(def.localVars, &LocalVariable{
				name:    name,
				size:    s.Size,
				isState: true,
			})
		}
	}
	for name := range src.Assigns {
		idx = def.GetVarIdx(name)
		if idx >= 0 {
			if idx < def.numParams {
				return fmt.Errorf("cannot assign to function parameter: '%v' in '%v'", name, def.name)
			} else {
				v := def.localVars[idx]
				if !v.isState {
					return fmt.Errorf("several assignment to the same var '%v' in '%v' is not allowed", name, def.name)
				}
			}
		} else {
			def.localVars = append(def.localVars, &LocalVariable{
				name: name,
			})
		}
	}
	return nil
}

func (def *QuplaFuncDef) analyzeAssigns(defYAML *QuplaFuncDefYAML, module *QuplaModule) error {
	var err error
	var idx int
	for name := range defYAML.Assigns {
		if idx, err = def.FindVarIdx(name, module); err != nil {
			return err
		}
		s := def.localVars[idx].expr.Size()
		if def.localVars[idx].isState && s != def.localVars[idx].size {
			return fmt.Errorf("sizes doesn't match for var '%v' in '%v'", name, def.GetName())
		}
		def.localVars[idx].size = s
	}
	return nil
}

func (def *QuplaFuncDef) finalizeLocalVars() error {
	var curOffset int64
	for _, v := range def.localVars {
		if v.size == 0 {
			return fmt.Errorf("can't determine var size '%v': '%v'", v.name, def.GetName())
		}
		v.offset = curOffset
		curOffset += v.size
	}
	def.bufLen = int64(curOffset)
	return nil
}

func (def *QuplaFuncDef) checkArgSizes(args []ExpressionInterface) error {
	for i := range args {
		if i >= def.numParams || args[i].Size() != def.localVars[i].size {
			return fmt.Errorf("param and arg # %v mismach in %v", i, def.GetName())
		}
	}
	return nil
}

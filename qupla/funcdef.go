package qupla

import (
	"fmt"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaFuncDef struct {
	yamlSource *QuplaFuncDefYAML // needed for analysis phase only
	name       string
	retSize    int64
	retExpr    ExpressionInterface
	localVars  []*VarInfo
	numParams  int64 // idx < numParams represents parameter, idx >= represents local var (assign)
	bufLen     int64 // total length of the local var buffer
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

// it tries to find var idx, offset, size. Analyzes var if not analyzed yet
func (def *QuplaFuncDef) GetVarInfo(name string, module ModuleInterface) (*VarInfo, error) {
	idx := def.GetVarIdx(name)
	if idx < 0 {
		return nil, nil
	}
	ret := def.localVars[idx]
	if ret.Analyzed {
		return ret, nil
	}
	var err error
	ret.Analyzed = true

	if ret.IsParam {
		// param
		ret.Expr = nil
	} else {
		// local var
		e, ok := def.yamlSource.Assigns[name]
		if !ok {
			return nil, fmt.Errorf("inconsistency with vars")
		}
		if ret.Expr, err = module.AnalyzeExpression(e, def); err != nil {
			return ret, err
		}
		if ret.IsState {
			if ret.Size != ret.Expr.Size() {
				return nil, fmt.Errorf("expression and state variable has different sizes in the assign")
			}
		} else {
			ret.Size = ret.Expr.Size()
		}
	}
	return ret, nil
}

func (def *QuplaFuncDef) createVarScope() error {
	src := def.yamlSource
	def.localVars = make([]*VarInfo, 0, len(src.Params)+len(src.Assigns))
	// first numParams indices belong to parameters
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
		})
	}
	// the rest of indices belong to local vars (incl state)
	var idx int64
	for name, s := range src.State {
		idx = def.GetVarIdx(name)
		if idx >= 0 {
			return fmt.Errorf("wrong declared state variable: '%v' in '%v'", name, def.name)
		} else {
			def.localVars = append(def.localVars, &VarInfo{
				Idx:     int64(len(def.localVars)),
				Name:    name,
				Size:    s.Size,
				IsState: true,
				IsParam: false,
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
				if !v.IsState {
					return fmt.Errorf("several assignment to the same var '%v' in '%v' is not allowed", name, def.name)
				}
			}
		} else {
			def.localVars = append(def.localVars, &VarInfo{
				Idx:  int64(len(def.localVars)),
				Name: name,
			})
		}
	}
	return nil
}

func (def *QuplaFuncDef) analyzeAssigns(defYAML *QuplaFuncDefYAML, module *QuplaModule) error {
	var err error
	var vi *VarInfo
	for name := range defYAML.Assigns {
		if vi, err = def.GetVarInfo(name, module); err != nil {
			return err
		}
		s := vi.Expr.Size()
		if vi.IsState && s != vi.Size {
			return fmt.Errorf("sizes doesn't match for var '%v' in '%v'", name, def.GetName())
		}
		vi.Size = s
	}
	return nil
}

func (def *QuplaFuncDef) finalizeLocalVars() error {
	var curOffset int64
	for _, v := range def.localVars {
		if v.Size == 0 {
			return fmt.Errorf("can't determine var size '%v': '%v'", v.Name, def.GetName())
		}
		v.Offset = curOffset
		curOffset += v.Size
	}
	def.bufLen = int64(curOffset)
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

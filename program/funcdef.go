package program

import (
	"fmt"
)

type QuplaEnvStmt struct {
	Name string `yaml:"name"`
	Join bool   `yaml:"join"`
}

type QuplaFuncArg struct {
	ArgName string                  `yaml:"argName"`
	Size    int64                   `yaml:"size"`
	Type    *QuplaExpressionWrapper `yaml:"type"` // not used
}

type QuplaStateVar struct {
	Size int64  `yaml:"size"`
	Type string `yaml:"type"`
}

type QuplaFuncDef struct {
	ReturnType     *QuplaExpressionWrapper            `yaml:"returnType"` // only size is necessary
	Params         []*QuplaFuncArg                    `yaml:"params"`
	State          map[string]*QuplaStateVar          `yaml:"state"`
	Env            []*QuplaEnvStmt                    `yaml:"env,omitempty"`
	Assigns        map[string]*QuplaExpressionWrapper `yaml:"assigns,omitempty"`
	ReturnExprWrap *QuplaExpressionWrapper            `yaml:"return"`
	//-------
	analyzed  bool
	name      string
	retSize   int64
	retExpr   ExpressionInterface
	localVars []*LocalVariable
	numParams int // idx < numParams represents parameter, idx >= represents local var (assign)
}

// represents local variablein func def
type LocalVariable struct {
	name     string
	isState  bool
	size     int64
	expr     ExpressionInterface
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

func (def *QuplaFuncDef) Analyze(module *QuplaModule) (*QuplaFuncDef, error) {
	if def.analyzed {
		return def, nil
	}
	def.analyzed = true
	var err error
	module.IncStat("numFuncDef")

	//debugf("Analyzing func def '%v'", def.Name)
	defer func(perr *error) {
		if *perr != nil {
			errorf("Error while analyzing func def '%v': %v", def.name, *perr)
		}
	}(&err)

	// return size. Must be const expression
	ce, err := def.ReturnType.Analyze(module, def)
	if err != nil {
		return nil, err
	}
	var sz int64
	if sz, err = GetConstValue(ce); err != nil {
		return nil, err
	}
	def.retSize = sz

	// build var scope
	if err = def.createVarScope(); err != nil {
		return nil, err
	}
	if err = def.analyzeAssigns(module); err != nil {
		return nil, err
	}
	// return expression
	if def.retExpr, err = def.ReturnExprWrap.Analyze(module, def); err != nil {
		return nil, err
	}
	if def.retExpr == nil {
		return nil, fmt.Errorf("in funcdef '%v': return expression can't be nil", def.name)
	}
	return def, nil
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

func (def *QuplaFuncDef) FindVarIdx(name string, module *QuplaModule) (int, error) {
	idx := def.GetVarIdx(name)
	if idx < 0 {
		return -1, nil
	}
	ret := def.localVars[idx]
	if ret.analyzed {
		return idx, nil
	}
	if ret.expr == nil {
		ret.analyzed = true
		ret.size = 0
		return idx, nil // ???
	}
	var err error

	if idx >= def.numParams {
		v := def.localVars[idx]
		ret.analyzed = true
		if ret.expr, err = ret.expr.Analyze(module, def); err != nil {
			return idx, err
		}
		if v.isState {
			if ret.size != ret.expr.Size() {
				return -1, fmt.Errorf("expression and state variable has different sizes in the assign")
			}
		} else {
			ret.size = ret.expr.Size()
		}
	}
	return idx, nil // ???
}

func (def *QuplaFuncDef) createVarScope() error {
	def.localVars = make([]*LocalVariable, 0, len(def.Params)+len(def.Assigns))
	// first numParams indices belong to parameters
	def.numParams = len(def.Params)
	for _, arg := range def.Params {
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
	for name, s := range def.State {
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
	for name, a := range def.Assigns {
		idx = def.GetVarIdx(name)
		if idx >= 0 {
			if idx < def.numParams {
				return fmt.Errorf("cannot assign to function parameter: '%v' in '%v'", name, def.name)
			} else {
				v := def.localVars[idx]
				if !v.isState {
					return fmt.Errorf("several assignment to the same var '%v' in '%v' is not allowed", name, def.name)
				} else {
					v.expr = a
				}
			}
		} else {
			def.localVars = append(def.localVars, &LocalVariable{
				name: name,
				expr: a,
			})
		}
	}
	return nil
}

func (def *QuplaFuncDef) analyzeAssigns(module *QuplaModule) error {
	var err error
	for name := range def.Assigns {
		if _, err = def.FindVarIdx(name, module); err != nil {
			return err
		}
	}
	return nil
}

func (def *QuplaFuncDef) checkArgSizes(args []ExpressionInterface) error {
	if len(args) != len(def.Params) {
		return fmt.Errorf("sizes of param and arg lists mismach in %v", def.GetName())
	}
	for i := range args {
		if args[i].Size() != def.Params[i].Size {
			return fmt.Errorf("size of param and arg # %v mismach in %v", i, def.GetName())
		}
	}
	return nil
}

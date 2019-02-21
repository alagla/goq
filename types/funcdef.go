package types

import "fmt"

type QuplaEnvStmt struct {
	Name string `yaml:"Name"`
	Join bool   `yaml:"join"`
}

type QuplaFuncDef struct {
	ReturnType     *QuplaExpressionWrapper            `yaml:"returnType"` // only size is necessary
	Params         map[string]*QuplaExpressionWrapper `yaml:"params"`
	Env            []*QuplaEnvStmt                    `yaml:"env,omitempty"`
	Assigns        map[string]*QuplaExpressionWrapper `yaml:"assigns,omitempty"`
	ReturnExprWrap *QuplaExpressionWrapper            `yaml:"return"`
	//-------
	Name     string
	retSize  int
	retExpr  ExpressionInterface
	varScope map[string]*LocalVariable
}

// represents local variable or parameter in func def
type LocalVariable struct {
	name  string
	isArg bool
	size  int64
	expr  ExpressionInterface
}

func (def *QuplaFuncDef) SetName(name string) {
	def.Name = name
}

func (def *QuplaFuncDef) Analyze(module *QuplaModule, _ *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	// build var scope
	if err = def.buildVarScope(module); err != nil {
		return nil, err
	}
	if err = def.analyzeAssigns(module); err != nil {
		return nil, err
	}

	// return size. Must be const expression
	ce, err := def.ReturnType.Analyze(module, def)
	if err != nil {
		return nil, err
	}
	var sz int64
	if sz, err = GetConstValue(ce); err != nil {
		return nil, err
	}
	def.retSize = int(sz)

	// return expression
	if def.retExpr, err = def.ReturnExprWrap.Analyze(module, def); err != nil {
		return nil, err
	}
	if def.retExpr == nil {
		return nil, fmt.Errorf("in funcdef '%v': return expression can't be nil", def.Name)
	}
	return def, nil
}

func (def *QuplaFuncDef) FindVar(name string) *LocalVariable {
	if def == nil || def.varScope == nil {
		return nil
	}
	ret, ok := def.varScope[name]
	if !ok {
		return nil
	}
	return ret
}

func (def *QuplaFuncDef) buildVarScope(module *QuplaModule) error {
	var numParams, numVars int
	def.varScope = make(map[string]*LocalVariable)
	// params
	for name, t := range def.Params {
		if def.FindVar(name) != nil {
			return fmt.Errorf("duplicate Name '%v'", name)
		}
		t, err := t.Analyze(module, nil)
		if err != nil {
			return err
		}
		size, err := GetConstValue(t)
		if err != nil {
			return err
		}
		def.varScope[name] = &LocalVariable{
			name:  name,
			isArg: true,
			size:  size,
		}
		numParams++
	}
	// local variables

	for name, a := range def.Assigns {
		if def.FindVar(name) != nil {
			return fmt.Errorf("duplicate Name '%v'", name)
		}
		def.varScope[name] = &LocalVariable{
			name:  name,
			isArg: false,
			size:  0, // todo
			expr:  a, // not analyzed yet
		}
		numVars++
	}
	debugf("FuncDef '%v': %v params, %v local variables", def.Name, numParams, numVars)
	return nil
}

func (def *QuplaFuncDef) analyzeAssigns(module *QuplaModule) error {
	var err error
	for name := range def.Assigns {
		v := def.FindVar(name)
		if v == nil {
			return fmt.Errorf("inconsistency: variable '%v' is not in the scope", name)
		}
		if !v.isArg {
			if v.expr, err = v.expr.Analyze(module, def); err != nil {
				return err
			}
		}
	}
	return nil
}

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
	analyzed bool
	Name     string
	retSize  int64
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

func (def *QuplaFuncDef) Analyze(module *QuplaModule) (*QuplaFuncDef, error) {
	if def.analyzed {
		return def, nil
	}
	def.analyzed = true
	var err error

	debugf("Analyzing func def '%v'", def.Name)
	defer func(perr *error) {
		if *perr == nil {
			debugf("Finished analyzing func def '%v'", def.Name)
		} else {
			errorf("Error while analyzing func def '%v': %v", def.Name, *perr)
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
	if err = def.analyzeVarScope(module); err != nil {
		return nil, err
	}

	// return expression
	if def.retExpr, err = def.ReturnExprWrap.Analyze(module, def); err != nil {
		return nil, err
	}
	if def.retExpr == nil {
		return nil, fmt.Errorf("in funcdef '%v': return expression can't be nil", def.Name)
	}
	return def, nil
}

func (def *QuplaFuncDef) Size() int64 {
	return def.retSize
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

func (def *QuplaFuncDef) analyzeVarScope(module *QuplaModule) error {
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
			expr:  a,
		}
		numVars++
	}
	// analyze rhs
	for _, v := range def.varScope {
		if !v.isArg {
			a, err := v.expr.Analyze(module, def)
			if err != nil {
				return err
			}
			v.expr = a
			v.size = a.Size()
		}
	}
	return nil
}

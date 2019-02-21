package types

import "fmt"

type QuplaFuncParam struct {
	Name string                  `yaml:"name"`
	Type *QuplaExpressionWrapper `yaml:"type"`
}

type QuplaEnvStmt struct {
	Name string `yaml:"name"`
	Join bool   `yaml:"join"`
}

type QuplaFuncDef struct {
	ReturnType     *QuplaExpressionWrapper            `yaml:"returnType"` // only size is necessary
	Params         []*QuplaFuncParam                  `yaml:"params"`
	Env            []*QuplaEnvStmt                    `yaml:"env,omitempty"`
	Assigns        map[string]*QuplaExpressionWrapper `yaml:"assigns,omitempty"`
	ReturnExprWrap *QuplaExpressionWrapper            `yaml:"return"`
	//-------
	name    string
	retSize int
	retExpr ExpressionInterface
	varMap  map[string]*localVariable
}

// represents local variable or parameter in func def
type localVariable struct {
	name  string
	isArg bool
	size  int64
	expr  ExpressionInterface
}

func (def *QuplaFuncDef) SetName(name string) {
	def.name = name
}

func (def *QuplaFuncDef) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	// return size. Must be const expression
	ce, err := def.ReturnType.Analyze(module)
	if err != nil {
		return nil, err
	}
	var sz int64
	if sz, err = GetConstValue(ce); err != nil {
		return nil, err
	}
	def.retSize = int(sz)

	// return expression
	if def.retExpr, err = def.ReturnExprWrap.Analyze(module); err != nil {
		return nil, err
	}
	if def.retExpr == nil {
		return nil, fmt.Errorf("in funcdef '%v': return expression can't be nil", def.name)
	}

	if err = def.buildVarMap(module); err != nil {
		return nil, err
	}
	return def, nil
}

func (def *QuplaFuncDef) checkVar(name string) bool {
	_, ok := def.varMap[name]
	return ok
}

func (def *QuplaFuncDef) buildVarMap(module *QuplaModule) error {
	var numParams, numVars int
	def.varMap = make(map[string]*localVariable)
	// params
	for _, a := range def.Params {
		if def.checkVar(a.Name) {
			return fmt.Errorf("duplicate name '%v'", a.Name)
		}
		t, err := a.Type.Analyze(module)
		if err != nil {
			return err
		}
		size, err := GetConstValue(t)
		if err != nil {
			return err
		}
		def.varMap[a.Name] = &localVariable{
			name:  a.Name,
			isArg: true,
			size:  size,
		}
		numParams++
	}
	// local variables

	for name, a := range def.Assigns {
		if def.checkVar(name) {
			return fmt.Errorf("duplicate name '%v'", name)
		}
		ae, err := a.Analyze(module)
		if err != nil {
			return err
		}
		def.varMap[name] = &localVariable{
			name:  name,
			isArg: false,
			size:  0, // todo
			expr:  ae,
		}
		numVars++
	}

	infof("FuncDef '%v': %v params, %v local variables", def.name, numParams, numVars)
	return nil
}

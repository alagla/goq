package types

import (
	"fmt"
)

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
	name     string
	isArg    bool
	size     int64
	expr     ExpressionInterface
	analyzed bool
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

	//debugf("Analyzing func def '%v'", def.Name)
	defer func(perr *error) {
		if *perr != nil {
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
	def.createVarScope()

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

func (def *QuplaFuncDef) InScope(name string) bool {
	if def.varScope == nil {
		panic(fmt.Errorf("var scope not ready for %v", def.Name))
	}
	_, ok := def.varScope[name]
	return ok
}

func (def *QuplaFuncDef) FindVar(name string, module *QuplaModule) (*LocalVariable, error) {
	if !def.InScope(name) {
		return nil, nil
	}
	ret := def.varScope[name]
	if ret.analyzed {
		return ret, nil
	}
	var err error
	if ret.isArg {
		// param
		if ret.expr == nil {
			panic("ret.expr == nil")
		}
		t, err := ret.expr.Analyze(module, nil)
		if err != nil {
			return nil, err
		}
		if ret.size, err = GetConstValue(t); err != nil {
			return nil, err
		}
		ret.expr = nil
		ret.analyzed = true
		return ret, nil
	}
	// local var
	if ret.expr == nil {
		ret.analyzed = true
		return ret, nil // ???
	}
	if ret.expr, err = ret.expr.Analyze(module, def); err != nil {
		return nil, err
	}
	ret.size = ret.expr.Size()
	ret.analyzed = true
	return ret, nil
}

func (def *QuplaFuncDef) createVarScope() {
	def.varScope = make(map[string]*LocalVariable)
	for name, t := range def.Params {
		def.varScope[name] = &LocalVariable{
			name:  name,
			isArg: true,
			expr:  t,
		}
	}
	for name, a := range def.Assigns {
		def.varScope[name] = &LocalVariable{
			name:  name,
			isArg: false,
			expr:  a,
		}
	}
}

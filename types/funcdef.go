package types

import (
	"fmt"
)

type QuplaEnvStmt struct {
	Name string `yaml:"name"`
	Join bool   `yaml:"join"`
}

type QuplaNameExpr struct {
	Size int64                   `yaml:"size"`
	Type *QuplaExpressionWrapper `yaml:"type"`
}

type QuplaStateVar struct {
	Size int64  `yaml:"size"`
	Type string `yaml:"type"`
}

type QuplaFuncDef struct {
	ReturnType     *QuplaExpressionWrapper            `yaml:"returnType"` // only size is necessary
	Params         map[string]*QuplaNameExpr          `yaml:"params"`
	State          map[string]*QuplaStateVar          `yaml:"state"`
	Env            []*QuplaEnvStmt                    `yaml:"env,omitempty"`
	Assigns        map[string]*QuplaExpressionWrapper `yaml:"assigns,omitempty"`
	ReturnExprWrap *QuplaExpressionWrapper            `yaml:"return"`
	//-------
	analyzed bool
	name     string
	retSize  int64
	retExpr  ExpressionInterface
	varScope map[string]*LocalVariable
}

const (
	VARTYPE_ARG   = 0
	VARTYPE_STATE = 1
	VARTYPE_LOCAL = 2
)

// represents local variable or parameter in func def
type LocalVariable struct {
	name     string
	vartype  int
	size     int64
	expr     ExpressionInterface
	analyzed bool
}

func (def *QuplaFuncDef) SetName(name string) {
	def.name = name
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

func (def *QuplaFuncDef) InScope(name string) bool {
	if def.varScope == nil {
		panic(fmt.Errorf("var scope not ready for %v", def.name))
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
	if ret.expr == nil {
		ret.analyzed = true
		ret.size = 0
		return ret, nil // ???
	}
	var err error

	switch ret.vartype {
	case VARTYPE_LOCAL:
		ret.analyzed = true
		if ret.expr, err = ret.expr.Analyze(module, def); err != nil {
			return nil, err
		}
		ret.size = ret.expr.Size()
	case VARTYPE_STATE:
		ret.analyzed = true
		if ret.expr, err = ret.expr.Analyze(module, def); err != nil {
			return nil, err
		}
		if ret.size != ret.expr.Size() {
			return nil, fmt.Errorf("expression and state variable has different sizes in the assign")
		}
	}
	return ret, nil // ???
}

func (def *QuplaFuncDef) createVarScope() error {
	def.varScope = make(map[string]*LocalVariable)
	for name, t := range def.Params {
		def.varScope[name] = &LocalVariable{
			name:     name,
			vartype:  VARTYPE_ARG,
			size:     t.Size,
			analyzed: true,
		}
	}
	for name, a := range def.Assigns {
		def.varScope[name] = &LocalVariable{
			name:    name,
			vartype: VARTYPE_LOCAL,
			expr:    a,
		}
	}
	for name, s := range def.State {
		v, ok := def.varScope[name]
		if ok && v.vartype == VARTYPE_ARG {
			return fmt.Errorf("function parameter can't be declared state variable: '%v' in '%v'", name, def.name)
		}
		v.vartype = VARTYPE_STATE
		v.size = s.Size
	}
	return nil
}

func (def *QuplaFuncDef) analyzeAssigns(module *QuplaModule) error {
	var err error
	for name := range def.Assigns {
		if _, err = def.FindVar(name, module); err != nil {
			return err
		}
	}
	return nil
}

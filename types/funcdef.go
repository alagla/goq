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
}

func (def *QuplaFuncDef) SetName(name string) {
	def.name = name
}

func (def *QuplaFuncDef) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	// return size. Must be cont expression
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
	return def, nil
}

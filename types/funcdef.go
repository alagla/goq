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

func (def *QuplaFuncDef) Analyze(module *QuplaModule) error {
	var err error
	// return size
	ei, err := def.ReturnType.Unwarp()
	if err != nil {
		return err
	}
	if err = ei.Analyze(module); err != nil {
		return err
	}
	if _, ok := ei.(*QuplaConstTypeName); !ok {
		return fmt.Errorf("in funcdef '%v': return type must be ConstTypeName", def.name)
	}
	ce, _ := ToConstExpression(ei)
	def.retSize = ce.GetConstValue()

	// return expression
	def.retExpr, err = def.ReturnExprWrap.Unwarp()
	if def.retExpr == nil {
		return fmt.Errorf("in funcdef '%v': return expression can't be nil", def.name)
	}
	if err = def.retExpr.Analyze(module); err != nil {
		return err
	}

	return nil
}

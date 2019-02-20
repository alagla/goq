package types

type QuplaFieldExpr struct {
	FieldName       string                  `yaml:"fieldName"`
	CondExprWrapper *QuplaExpressionWrapper `yaml:"condExpr"`
	//---
	condExpr ExpressionInterface
}

func (e *QuplaFieldExpr) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	var err error
	e.condExpr, err = e.CondExprWrapper.Analyze(module)
	if err != nil {
		return nil, err
	}
	return e, nil
}

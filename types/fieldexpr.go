package types

type QuplaFieldExpr struct {
	FieldName       string                  `yaml:"fieldName"`
	CondExprWrapper *QuplaExpressionWrapper `yaml:"condExpr"`
	//---
	condExpr ExpressionInterface
}

func (e *QuplaFieldExpr) Analyze(module *QuplaModule) error {
	var err error
	e.condExpr, err = e.CondExprWrapper.Unwarp()
	if err != nil {
		return err
	}
	return e.condExpr.Analyze(module)
}

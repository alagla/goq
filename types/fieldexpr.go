package types

type QuplaFieldExpr struct {
	FieldName string           `yaml:"fieldName"`
	CondExpr  *QuplaExpression `yaml:"condExpr"`
}

func (fieldExpr *QuplaFieldExpr) Analyze(mod *QuplaModule) error {
	return fieldExpr.CondExpr.Analyze(mod)
}

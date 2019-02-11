package expr

type QuplaFieldExpr struct {
	FieldName string           `yaml:"fieldName"`
	CondExpr  *QuplaExpression `yaml:"condExpr"`
}

func (fieldExpr *QuplaFieldExpr) Analyze() error {
	return fieldExpr.CondExpr.Analyze()
}

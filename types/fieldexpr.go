package types

type QuplaFieldExpr struct {
	FieldName string                  `yaml:"fieldName"`
	CondExpr  *QuplaExpressionWrapper `yaml:"condExpr"`
}

func (fieldExpr *QuplaFieldExpr) Analyze(mod *QuplaModule) error {
	return fieldExpr.CondExpr.Analyze(mod)
}

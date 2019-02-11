package expr

type QuplaCondExpr struct {
	If   *QuplaExpression `yaml:"if"`
	Then *QuplaExpression `yaml:"then"`
	Else *QuplaExpression `yaml:"else"`
}

func (conExpr *QuplaCondExpr) Analyze() error {
	return nil
}

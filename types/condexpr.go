package types

type QuplaCondExpr struct {
	If   *QuplaExpressionWrapper `yaml:"if"`
	Then *QuplaExpressionWrapper `yaml:"then"`
	Else *QuplaExpressionWrapper `yaml:"else"`
}

func (conExpr *QuplaCondExpr) Analyze(module *QuplaModule) error {
	return nil
}

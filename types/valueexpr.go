package types

type QuplaValueExpr struct {
	Trits  string `yaml:"Trits"`
	Trytes string `yaml:"trytes"`
}

func (valueExpr *QuplaValueExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	return valueExpr, nil
}

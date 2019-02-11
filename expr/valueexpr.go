package expr

type QuplaValueExpr struct {
	Trits  string `yaml:"trits"`
	Trytes string `yaml:"trytes"`
}

func (valueExpr *QuplaValueExpr) Analyze() error {
	return nil
}

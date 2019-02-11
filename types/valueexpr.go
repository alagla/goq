package types

type QuplaValueExpr struct {
	Trits  string `yaml:"trits"`
	Trytes string `yaml:"trytes"`
}

func (valueExpr *QuplaValueExpr) Analyze(module *QuplaModule) error {
	return nil
}

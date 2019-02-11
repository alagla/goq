package types

type QuplaLutExpr struct {
	Name string             `yaml:"name"`
	Args []*QuplaExpression `yaml:"args"`
}

func (lutExpr *QuplaLutExpr) Analyze(module *QuplaModule) error {
	for _, arg := range lutExpr.Args {
		if err := arg.Analyze(module); err != nil {
			return err
		}
	}
	return nil
}

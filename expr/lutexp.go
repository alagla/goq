package expr

type QuplaLutExpr struct {
	Name string             `yaml:"name"`
	Args []*QuplaExpression `yaml:"args"`
}

func (lutExpr *QuplaLutExpr) Analyze() error {
	for _, arg := range lutExpr.Args {
		if err := arg.Analyze(); err != nil {
			return err
		}
	}
	return nil
}

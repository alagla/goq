package types

type QuplaFuncExpr struct {
	Name string             `yaml:"name"`
	Args []*QuplaExpression `yaml:"args"`
}

func (funcExpr *QuplaFuncExpr) Analyze(mod *QuplaModule) error {
	for _, arg := range funcExpr.Args {
		if err := arg.Analyze(mod); err != nil {
			return err
		}
	}
	return nil
}

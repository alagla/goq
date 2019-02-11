package expr

type QuplaFuncExpr struct {
	Name string             `yaml:"name"`
	Args []*QuplaExpression `yaml:"args"`
}

func (funcExpr *QuplaFuncExpr) Analyze() error {
	for _, arg := range funcExpr.Args {
		if err := arg.Analyze(); err != nil {
			return err
		}
	}
	return nil
}

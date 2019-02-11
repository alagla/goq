package types

type QuplaSliceExpr struct {
	Name  string           `yaml:"name"`
	Start *QuplaExpression `yaml:"start,omitempty"`
	End   *QuplaExpression `yaml:"end,omitempty"`
}

func (sliceExpr *QuplaSliceExpr) Analyze(module *QuplaModule) error {
	if err := sliceExpr.Start.Analyze(module); err != nil {
		return err
	}
	return sliceExpr.End.Analyze(module)
}

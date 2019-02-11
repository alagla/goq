package expr

type QuplaSliceExpr struct {
	Name  string           `yaml:"name"`
	Start *QuplaExpression `yaml:"start,omitempty"`
	End   *QuplaExpression `yaml:"end,omitempty"`
}

func (sliceExpr *QuplaSliceExpr) Analyze() error {
	if err := sliceExpr.Start.Analyze(); err != nil {
		return err
	}
	return sliceExpr.End.Analyze()
}

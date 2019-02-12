package types

type QuplaSliceExpr struct {
	Name  string                  `yaml:"name"`
	Start *QuplaExpressionWrapper `yaml:"start,omitempty"`
	End   *QuplaExpressionWrapper `yaml:"end,omitempty"`
}

func (sliceExpr *QuplaSliceExpr) Analyze(module *QuplaModule) error {
	if err := sliceExpr.Start.Analyze(module); err != nil {
		return err
	}
	return sliceExpr.End.Analyze(module)
}

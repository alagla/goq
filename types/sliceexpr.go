package types

type QuplaSliceExpr struct {
	Name          string                  `yaml:"name"`
	StartExprWrap *QuplaExpressionWrapper `yaml:"start,omitempty"`
	EndExprWrap   *QuplaExpressionWrapper `yaml:"end,omitempty"`
	//----
	startExpr ExpressionInterface
	endExpr   ExpressionInterface
}

func (e *QuplaSliceExpr) Analyze(module *QuplaModule) error {
	var err error
	if e.startExpr, err = e.StartExprWrap.Unwarp(); err != nil {
		return err
	}
	if e.endExpr, err = e.EndExprWrap.Unwarp(); err != nil {
		return err
	}
	if err := e.startExpr.Analyze(module); err != nil {
		return err
	}
	return e.endExpr.Analyze(module)
}

package types

type QuplaSliceExpr struct {
	Name          string                  `yaml:"name"`
	StartExprWrap *QuplaExpressionWrapper `yaml:"start,omitempty"`
	EndExprWrap   *QuplaExpressionWrapper `yaml:"end,omitempty"`
	//----
	startExpr ExpressionInterface
	endExpr   ExpressionInterface
}

func (e *QuplaSliceExpr) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	var err error
	if e.startExpr != nil {
		if e.startExpr, err = e.StartExprWrap.Analyze(module); err != nil {
			return nil, err
		}
	} else {
		return e, nil
	}
	if e.endExpr != nil {
		if e.endExpr, err = e.EndExprWrap.Analyze(module); err != nil {
			return nil, err
		}
	}
	return e, nil
}

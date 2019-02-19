package types

type QuplaMergeExpr struct {
	LhsWrap *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap *QuplaExpressionWrapper `yaml:"rhs"`
	//----
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func (e *QuplaMergeExpr) Analyze(module *QuplaModule) error {
	var err error
	e.lhsExpr, err = e.LhsWrap.Unwarp()
	if err != nil {
		return err
	}
	e.rhsExpr, err = e.RhsWrap.Unwarp()
	if err != nil {
		return err
	}
	if err := e.lhsExpr.Analyze(module); err != nil {
		return err
	}
	return e.rhsExpr.Analyze(module)
}

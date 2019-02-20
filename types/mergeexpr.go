package types

type QuplaMergeExpr struct {
	LhsWrap *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap *QuplaExpressionWrapper `yaml:"rhs"`
	//----
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func (e *QuplaMergeExpr) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	var err error
	e.lhsExpr, err = e.LhsWrap.Analyze(module)
	if err != nil {
		return nil, err
	}
	e.rhsExpr, err = e.RhsWrap.Analyze(module)
	if err != nil {
		return nil, err
	}
	return e, nil
}

package types

type QuplaConcatExpr struct {
	LhsWrap *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap *QuplaExpressionWrapper `yaml:"rhs"`
	//---
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func (e *QuplaConcatExpr) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	var err error
	if e.lhsExpr, err = e.LhsWrap.Analyze(module); err != nil {
		return nil, err
	}
	if e.rhsExpr, err = e.RhsWrap.Analyze(module); err != nil {
		return nil, err
	}
	return e, nil
}

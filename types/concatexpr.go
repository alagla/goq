package types

type QuplaConcatExpr struct {
	LhsWrap *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap *QuplaExpressionWrapper `yaml:"rhs"`
	//---
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func (e *QuplaConcatExpr) Analyze(module *QuplaModule) error {
	var err error
	if e.lhsExpr, err = e.LhsWrap.Unwarp(); err != nil {
		return err
	}
	if e.rhsExpr, err = e.RhsWrap.Unwarp(); err != nil {
		return err
	}
	if err := e.rhsExpr.Analyze(module); err != nil {
		return err
	}
	return e.rhsExpr.Analyze(module)
}

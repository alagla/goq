package types

type QuplaMergeExpr struct {
	LhsWrap *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap *QuplaExpressionWrapper `yaml:"rhs"`
	//----
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func (e *QuplaMergeExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	e.lhsExpr, err = e.LhsWrap.Analyze(module, scope)
	if err != nil {
		return nil, err
	}
	e.rhsExpr, err = e.RhsWrap.Analyze(module, scope)
	if err != nil {
		return nil, err
	}
	return e, nil
}

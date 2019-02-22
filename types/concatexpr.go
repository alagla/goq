package types

type QuplaConcatExpr struct {
	LhsWrap *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap *QuplaExpressionWrapper `yaml:"rhs"`
	//---
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func (e *QuplaConcatExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	if e.lhsExpr, err = e.LhsWrap.Analyze(module, scope); err != nil {
		return nil, err
	}
	if e.rhsExpr, err = e.RhsWrap.Analyze(module, scope); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *QuplaConcatExpr) Size() int64 {
	if e == nil {
		return 0
	}
	ls := e.lhsExpr.Size()
	rs := e.rhsExpr.Size()
	if ls < rs {
		return rs
	}
	return ls
}

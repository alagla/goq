package program

import "fmt"

type QuplaConcatExpr struct {
	LhsWrap *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap *QuplaExpressionWrapper `yaml:"rhs"`
	//---
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func (e *QuplaConcatExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	module.IncStat("numConcat")

	if e.lhsExpr, err = e.LhsWrap.Analyze(module, scope); err != nil {
		return nil, err
	}
	if e.rhsExpr, err = e.RhsWrap.Analyze(module, scope); err != nil {
		return nil, err
	}
	if e.rhsExpr.Size() == 0 || e.lhsExpr.Size() == 0 {
		return nil, fmt.Errorf("size of concat opeation can't be 0: scope '%v'", scope.GetName())
	}
	return e, nil
}

func (e *QuplaConcatExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.lhsExpr.Size() + e.rhsExpr.Size()
}

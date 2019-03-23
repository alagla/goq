package qupla

import . "github.com/lunfardo314/goq/abstract"

type QuplaExprBase struct {
	source  string
	subexpr []ExpressionInterface
}

func NewQuplaExprBase(source string) QuplaExprBase {
	return QuplaExprBase{
		source:  source,
		subexpr: make([]ExpressionInterface, 0, 5),
	}
}

func (e *QuplaExprBase) GetSource() string {
	return e.source
}

func (e *QuplaExprBase) GetSubExpr(idx int) ExpressionInterface {
	return e.subexpr[idx]
}

func (e *QuplaExprBase) NumSubExpr() int {
	return len(e.subexpr)
}

func (e *QuplaExprBase) AppendSubExpr(se ExpressionInterface) {
	e.subexpr = append(e.subexpr, se)
}

func (e *QuplaExprBase) HasSubExpr() bool {
	return len(e.subexpr) > 0
}

func (e *QuplaExprBase) ReferencesSubExprs(funName string) bool {
	for _, se := range e.subexpr {
		if se.References(funName) {
			return true
		}
	}
	return false
}

func (e *QuplaExprBase) References(funName string) bool {
	return e.ReferencesSubExprs(funName)
}

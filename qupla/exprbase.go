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

func (e *QuplaExprBase) AppendSubExpr(se ExpressionInterface) {
	e.subexpr = append(e.subexpr, se)
}

func (e *QuplaExprBase) HasSubExpr() bool {
	return len(e.subexpr) > 0
}

func (e *QuplaExprBase) HasState() bool {
	return false
}

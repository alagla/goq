package qupla

type ExpressionBase struct {
	source  string
	subexpr []ExpressionInterface
}

func NewExpressionBase(source string) ExpressionBase {
	return ExpressionBase{
		source:  source,
		subexpr: make([]ExpressionInterface, 0, 5),
	}
}

func (e *ExpressionBase) GetSource() string {
	return e.source
}

func (e *ExpressionBase) GetSubExpr(idx int) ExpressionInterface {
	return e.subexpr[idx]
}

func (e *ExpressionBase) NumSubExpr() int {
	return len(e.subexpr)
}

func (e *ExpressionBase) AppendSubExpr(se ExpressionInterface) {
	e.subexpr = append(e.subexpr, se)
}

func (e *ExpressionBase) HasSubExpr() bool {
	return len(e.subexpr) > 0
}

func (e *ExpressionBase) ReferencesSubExprs(funName string) bool {
	for _, se := range e.subexpr {
		if se.References(funName) {
			return true
		}
	}
	return false
}

func (e *ExpressionBase) References(funName string) bool {
	return e.ReferencesSubExprs(funName)
}

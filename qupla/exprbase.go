package qupla

type ExpressionBase struct {
	source  string
	subExpr []ExpressionInterface
}

func NewExpressionBase(source string) ExpressionBase {
	return ExpressionBase{
		source:  source,
		subExpr: make([]ExpressionInterface, 0, 5),
	}
}

// shallow copy
func (base *ExpressionBase) copyBase() ExpressionBase {
	return ExpressionBase{
		source:  base.source,
		subExpr: base.subExpr,
	}
}

func (e *ExpressionBase) GetSubexpressions() []ExpressionInterface {
	return e.subExpr
}

func (e *ExpressionBase) SetSubexpressions(subExpr []ExpressionInterface) {
	e.subExpr = subExpr
}

func (e *ExpressionBase) GetSource() string {
	if e == nil {
		return ""
	}
	return e.source
}

func (e *ExpressionBase) GetSubExpr(idx int) ExpressionInterface {
	return e.subExpr[idx]
}

func (e *ExpressionBase) NumSubExpr() int {
	return len(e.subExpr)
}

func (e *ExpressionBase) AppendSubExpr(se ExpressionInterface) {
	e.subExpr = append(e.subExpr, se)
}

func (e *ExpressionBase) HasSubExpr() bool {
	return len(e.subExpr) > 0
}

func (e *ExpressionBase) ReferencesSubExprs(funName string) bool {
	for _, se := range e.subExpr {
		if se.References(funName) {
			return true
		}
	}
	return false
}

func (e *ExpressionBase) References(funName string) bool {
	return e.ReferencesSubExprs(funName)
}

func (e *ExpressionBase) hasStateSubexpr() bool {
	for _, se := range e.subExpr {
		if se.HasState() {
			return true
		}
	}
	return false
}

func (e *ExpressionBase) HasState() bool {
	return e.hasStateSubexpr()
}

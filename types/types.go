package types

type ExpressionInterface interface {
	Analyze(*QuplaModule, *QuplaFuncDef) (ExpressionInterface, error)
	Size() int64
	RequireSize(int64) error
}

type QuplaExecStmt struct {
	ExprWrap     *QuplaExpressionWrapper `yaml:"expr"`
	ExpectedWrap *QuplaExpressionWrapper `yaml:"expected,omitempty"`
	//---
	isTest       bool
	expr         ExpressionInterface
	exprExpected ExpressionInterface
}

type QuplaTypeDef struct {
	Size   string                            `yaml:"size"`
	Fields map[string]*struct{ Size string } `yaml:"fields,omitempty"`
}

type QuplaNullExpr struct{}

func (e *QuplaNullExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	return e, nil
}

func (e *QuplaNullExpr) Size() int64 {
	return 0
}

func (e *QuplaNullExpr) RequireSize(_ int64) error {
	return nil
}

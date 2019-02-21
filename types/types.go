package types

type ExpressionInterface interface {
	Analyze(*QuplaModule, *QuplaFuncDef) (ExpressionInterface, error)
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

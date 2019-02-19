package types

type ExpressionInterface interface {
	Analyze(*QuplaModule) error
}

type QuplaExecStmt struct {
	ExprWrap     *QuplaExpressionWrapper `yaml:"expr"`
	ExpectedWrap *QuplaExpressionWrapper `yaml:"expected,omitempty"`
	//---
	isTest       bool
	expr         ExpressionInterface
	exprExpected ExpressionInterface
}

type QuplaFuncParam struct {
	Name string                  `yaml:"name"`
	Type *QuplaExpressionWrapper `yaml:"type"`
}

type QuplaEnvStmt struct {
	Name string `yaml:"name"`
	Join bool   `yaml:"join"`
}

type QuplaFuncDef struct {
	ReturnType *QuplaExpressionWrapper            `yaml:"returnType"`
	Params     []*QuplaFuncParam                  `yaml:"params"`
	Env        []*QuplaEnvStmt                    `yaml:"env,omitempty"`
	Assigns    map[string]*QuplaExpressionWrapper `yaml:"assigns,omitempty"`
	Return     *QuplaExpressionWrapper            `yaml:"return"`
}

type QuplaTypeDef struct {
	Size   string                            `yaml:"size"`
	Fields map[string]*struct{ Size string } `yaml:"fields,omitempty"`
}

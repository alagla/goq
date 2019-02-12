package types

type QuplaExecStmt struct {
	Expr     *QuplaExpressionWrapper `yaml:"expr"`
	Expected *QuplaExpressionWrapper `yaml:"expected,omitempty"`
	//---
	isTest bool
}

type QuplaAssignStmt struct {
	Lhs string                  `yaml:"lhs"`
	Rhs *QuplaExpressionWrapper `yaml:"rhs"`
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
	ReturnType *QuplaExpressionWrapper `yaml:"returnType"`
	Params     []*QuplaFuncParam       `yaml:"params"`
	Env        []*QuplaEnvStmt         `yaml:"env,omitempty"`
	Assigns    []*QuplaAssignStmt      `yaml:"assigns,omitempty"`
	Return     *QuplaExpressionWrapper `yaml:"return"`
}

type QuplaTypeDef struct {
	Size   string                            `yaml:"size"`
	Fields map[string]*struct{ Size string } `yaml:"fields,omitempty"`
}

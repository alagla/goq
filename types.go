package main

import . "github.com/lunfardo314/goq/expr"

type QuplaExecStmt struct {
	Expr     *QuplaExpression `yaml:"expr"`
	Expected *QuplaExpression `yaml:"expected,omitempty"`
	//---
	isTest bool
}

type QuplaAssignStmt struct {
	Lhs string           `yaml:"lhs"`
	Rhs *QuplaExpression `yaml:"rhs"`
}

type QuplaFuncParam struct {
	Name string           `yaml:"name"`
	Type *QuplaExpression `yaml:"type"`
}

type QuplaEnvStmt struct {
	Name string `yaml:"name"`
	Join bool   `yaml:"join"`
}

type QuplaFuncDef struct {
	ReturnType *QuplaExpression   `yaml:"returnType"`
	Params     []*QuplaFuncParam  `yaml:"params"`
	Env        []*QuplaEnvStmt    `yaml:"env,omitempty"`
	Assigns    []*QuplaAssignStmt `yaml:"assigns,omitempty"`
	Return     *QuplaExpression   `yaml:"return"`
}

type QuplaTypeDef struct {
	Size   string                            `yaml:"size"`
	Fields map[string]*struct{ Size string } `yaml:"fields"`
}

package main

type QuplaExecStmt struct {
	Expr     *QuplaExpression `yaml:"expr"`
	Expected *QuplaExpression `yaml:"expected,omitempty"`
}

type QuplaTypeExpr struct {
	Type   *QuplaExpression   `yaml:"type"`
	Fields []*QuplaExpression `yaml:"fields"`
}

type QuplaMergeExpr struct {
	Lhs *QuplaExpression `yaml:"lhs"`
	Rhs *QuplaExpression `yaml:"rhs"`
}

type QuplaConcatExpr struct {
	Lhs *QuplaExpression `yaml:"lhs"`
	Rhs *QuplaExpression `yaml:"rhs"`
}

type QuplaConstExpr struct {
	Operator string           `yaml:"operator"`
	Lhs      *QuplaExpression `yaml:"lhs"`
	Rhs      *QuplaExpression `yaml:"rhs"`
}

type QuplaConstTerm struct {
	Operator string           `yaml:"operator"`
	Lhs      *QuplaExpression `yaml:"lhs"`
	Rhs      *QuplaExpression `yaml:"rhs"`
}

type QuplaValueExpr struct {
	Trits  string `yaml:"trits"`
	Trytes string `yaml:"trytes"`
}

type QuplaFieldExpr struct {
	FieldName string           `yaml:"fieldName"`
	CondExpr  *QuplaExpression `yaml:"condExpr"`
}

type QuplaSliceExpr struct {
	Name  string           `yaml:"name"`
	Start *QuplaExpression `yaml:"start,omitempty"`
	End   *QuplaExpression `yaml:"end,omitempty"`
}

type QuplaLutExpr struct {
	Name string             `yaml:"name"`
	Args []*QuplaExpression `yaml:"args"`
}

type QuplaCondExpr struct {
	If   *QuplaExpression `yaml:"if"`
	Then *QuplaExpression `yaml:"then"`
	Else *QuplaExpression `yaml:"else"`
}

type QuplaFuncExpr struct {
	Name string             `yaml:"name"`
	Args []*QuplaExpression `yaml:"args"`
}

type QuplaExpression struct {
	QuplaCondExpr *QuplaCondExpr   `yaml:"CondExpr,omitempty"`
	LutExpr       *QuplaLutExpr    `yaml:"LutExpr,omitempty"`
	SliceExpr     *QuplaSliceExpr  `yaml:"SliceExpr,omitempty"`
	ValueExpr     *QuplaValueExpr  `yaml:"ValueExpr,omitempty"`
	FuncExpr      *QuplaFuncExpr   `yaml:"FuncExpr,omitempty"`
	FieldExpr     *QuplaFieldExpr  `yaml:"FieldExpr,omitempty"`
	ConstNumber   string           `yaml:"ConstNumber,omitempty"`
	ConstTypeName string           `yaml:"ConstTypeName,omitempty"`
	ConstTerm     *QuplaConstTerm  `yaml:"ConstTerm,omitempty"`
	ConstExpr     *QuplaConstExpr  `yaml:"ConstExpr,omitempty"`
	ConcatExpr    *QuplaConcatExpr `yaml:"ConcatExpr,omitempty"`
	MergeExpr     *QuplaMergeExpr  `yaml:"MergeExpr,omitempty"`
	TypeExpr      *QuplaTypeExpr   `yaml:"TypeExpr,omitempty"`
}

type QuplaAssignStmt struct {
	Lhs string           `yaml:"lhs"`
	Rhs *QuplaExpression `yaml:"rhs"`
}

type QuplaFuncParam struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type QuplaEnvStmt struct {
	Name string `yaml:"name"`
	Join bool   `yaml:"join"`
}

type QuplaFuncDef struct {
	Type    string             `yaml:"type"`
	Params  []*QuplaFuncParam  `yaml:"params"`
	Env     []*QuplaEnvStmt    `yaml:"env,omitempty"`
	Assigns []*QuplaAssignStmt `yaml:"assigns,omitempty"`
	Return  *QuplaExpression   `yaml:"return"`
}

type QuplaTypeDef struct {
	Size   string                            `yaml:"size"`
	Fields map[string]*struct{ Size string } `yaml:"fields"`
}

type QuplaModuleYAML struct {
	Types     map[string]*QuplaTypeDef `yaml:"types"`
	Luts      map[string][]string      `yaml:"luts"`
	Functions map[string]*QuplaFuncDef `yaml:"functions"`
	Execs     []*QuplaExecStmt         `yaml:"execs"`
}

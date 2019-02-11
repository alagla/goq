package types

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

type QuplaConstTypeName struct {
	TypeName string `yaml:"typeName"`
	Size     string `yaml:"size"`
}

type QuplaConstNumber struct {
	Value string `yaml:"value"`
}

func (e *QuplaConstExpr) Analyze(module *QuplaModule) error {
	if err := e.Lhs.Analyze(module); err != nil {
		return err
	}
	return e.Rhs.Analyze(module)
}

func (e *QuplaConstTerm) Analyze(module *QuplaModule) error {
	if err := e.Lhs.Analyze(module); err != nil {
		return err
	}
	return e.Rhs.Analyze(module)
}

func (e *QuplaConstTypeName) Analyze(module *QuplaModule) error {
	return nil
}

func (e *QuplaConstNumber) Analyze(module *QuplaModule) error {
	return nil
}

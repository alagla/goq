package expr

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

func (e *QuplaConstExpr) Analyze() error {
	if err := e.Lhs.Analyze(); err != nil {
		return err
	}
	return e.Rhs.Analyze()
}

func (e *QuplaConstTerm) Analyze() error {
	if err := e.Lhs.Analyze(); err != nil {
		return err
	}
	return e.Rhs.Analyze()
}

func (e *QuplaConstTypeName) Analyze() error {
	return nil
}

func (e *QuplaConstNumber) Analyze() error {
	return nil
}

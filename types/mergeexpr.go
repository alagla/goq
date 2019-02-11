package types

type QuplaMergeExpr struct {
	Lhs *QuplaExpression `yaml:"lhs"`
	Rhs *QuplaExpression `yaml:"rhs"`
}

func (e *QuplaMergeExpr) Analyze(module *QuplaModule) error {
	if err := e.Lhs.Analyze(module); err != nil {
		return err
	}
	return e.Rhs.Analyze(module)
}

package expr

type QuplaMergeExpr struct {
	Lhs *QuplaExpression `yaml:"lhs"`
	Rhs *QuplaExpression `yaml:"rhs"`
}

func (e *QuplaMergeExpr) Analyze() error {
	if err := e.Lhs.Analyze(); err != nil {
		return err
	}
	return e.Rhs.Analyze()
}

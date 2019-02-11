package expr

type QuplaConcatExpr struct {
	Lhs *QuplaExpression `yaml:"lhs"`
	Rhs *QuplaExpression `yaml:"rhs"`
}

func (e *QuplaConcatExpr) Analyze() error {
	if err := e.Lhs.Analyze(); err != nil {
		return err
	}
	return e.Rhs.Analyze()
}

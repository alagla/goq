package expr

type QuplaTypeExpr struct {
	Type   *QuplaExpression   `yaml:"type"`
	Fields []*QuplaExpression `yaml:"fields"`
}

func (e *QuplaTypeExpr) Analyze() error {
	if err := e.Type.Analyze(); err != nil {
		return err
	}
	for _, fld := range e.Fields {
		if err := fld.Analyze(); err != nil {
			return err
		}
	}
	return nil
}

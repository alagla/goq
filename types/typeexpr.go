package types

type QuplaTypeExpr struct {
	Type   *QuplaExpression   `yaml:"type"`
	Fields []*QuplaExpression `yaml:"fields"`
}

func (e *QuplaTypeExpr) Analyze(mod *QuplaModule) error {
	if err := e.Type.Analyze(mod); err != nil {
		return err
	}
	for _, fld := range e.Fields {
		if err := fld.Analyze(mod); err != nil {
			return err
		}
	}
	return nil
}

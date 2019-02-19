package types

// ----- ?????? do we need it?
type QuplaTypeExpr struct {
	TypeExprWrap *QuplaExpressionWrapper   `yaml:"type"`
	Fields       []*QuplaExpressionWrapper `yaml:"fields"`
	//---
	typeExpr ExpressionInterface
	fields   []ExpressionInterface
}

func (e *QuplaTypeExpr) Analyze(module *QuplaModule) error {
	e.fields = make([]ExpressionInterface, 0, len(e.Fields))
	var err error
	if e.typeExpr, err = e.TypeExprWrap.Unwarp(); err != nil {
		return err
	}
	if err := e.typeExpr.Analyze(module); err != nil {
		return err
	}
	var fe ExpressionInterface
	for _, fld := range e.Fields {
		if fe, err = fld.Unwarp(); err != nil {
			return err
		}
		if err := fe.Analyze(module); err != nil {
			return err
		}
		e.fields = append(e.fields, fe)
	}
	return nil
}

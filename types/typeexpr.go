package types

// ----- ?????? do we need it?
type QuplaTypeExpr struct {
	TypeExprWrap *QuplaExpressionWrapper   `yaml:"type"`
	Fields       []*QuplaExpressionWrapper `yaml:"fields"`
	//---
	typeExpr ExpressionInterface
	fields   []ExpressionInterface
}

func (e *QuplaTypeExpr) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	e.fields = make([]ExpressionInterface, 0, len(e.Fields))
	var err error
	if e.typeExpr, err = e.TypeExprWrap.Analyze(module); err != nil {
		return nil, err
	}
	var fe ExpressionInterface
	for _, fld := range e.Fields {
		if fe, err = fld.Analyze(module); err != nil {
			return nil, err
		}
		e.fields = append(e.fields, fe)
	}
	return e, nil
}

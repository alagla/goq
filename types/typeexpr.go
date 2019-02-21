package types

// ----- ?????? do we need it?
type QuplaTypeExpr struct {
	TypeExprWrap *QuplaExpressionWrapper            `yaml:"type"`
	Fields       map[string]*QuplaExpressionWrapper `yaml:"fields"`
	//---
	typeExpr ExpressionInterface
	fields   map[string]ExpressionInterface
}

func (e *QuplaTypeExpr) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	e.fields = make(map[string]ExpressionInterface)
	var err error
	if e.typeExpr, err = e.TypeExprWrap.Analyze(module); err != nil {
		return nil, err
	}
	var fe ExpressionInterface
	for name, expr := range e.Fields {
		if fe, err = expr.Analyze(module); err != nil {
			return nil, err
		}
		e.fields[name] = fe
	}
	return e, nil
}

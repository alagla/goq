package types

import "fmt"

// ----- ?????? do we need it?
type QuplaTypeExpr struct {
	TypeExprWrap *QuplaExpressionWrapper            `yaml:"type"`
	Fields       map[string]*QuplaExpressionWrapper `yaml:"fields"`
	//---
	typeExpr ExpressionInterface
	size     int64
	fields   map[string]ExpressionInterface
}

func (e *QuplaTypeExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	e.fields = make(map[string]ExpressionInterface)
	var err error
	if e.typeExpr, err = e.TypeExprWrap.Analyze(module, scope); err != nil {
		return nil, err
	}

	if e.size, err = GetConstValue(e.typeExpr); err != nil {
		return nil, err
	}

	var fe ExpressionInterface
	for name, expr := range e.Fields {
		if fe, err = expr.Analyze(module, scope); err != nil {
			return nil, err
		}
		e.fields[name] = fe
	}
	if err = e.CheckSize(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *QuplaTypeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *QuplaTypeExpr) CheckSize() error {
	var sumFld int64

	for _, f := range e.fields {
		sumFld += f.Size()
	}
	if sumFld != e.size {
		return fmt.Errorf("sum of field sizes != type end")
	}
	return nil
}

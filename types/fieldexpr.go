package types

import "fmt"

type QuplaFieldExpr struct {
	FieldName       string                  `yaml:"fieldName"`
	CondExprWrapper *QuplaExpressionWrapper `yaml:"condExpr"`
	//---
	condExpr ExpressionInterface
}

func (e *QuplaFieldExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	e.condExpr, err = e.CondExprWrapper.Analyze(module, scope)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *QuplaFieldExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.condExpr.Size()
}

func (e *QuplaFieldExpr) RequireSize(size int64) error {
	if size != e.Size() {
		return fmt.Errorf("size mismatch in FieldExpr")
	}
	return nil
}

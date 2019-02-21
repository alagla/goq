package types

import "fmt"

type QuplaSliceExpr struct {
	Name          string                  `yaml:"Name"`
	StartExprWrap *QuplaExpressionWrapper `yaml:"start,omitempty"`
	EndExprWrap   *QuplaExpressionWrapper `yaml:"end,omitempty"`
	//----
	localVar  *LocalVariable
	startExpr ExpressionInterface
	endExpr   ExpressionInterface
}

func (e *QuplaSliceExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error

	e.localVar = scope.FindVar(e.Name)
	if e.localVar == nil {
		var sc string
		if scope != nil {
			sc = scope.Name
		} else {
			sc = "nil"
		}
		return nil, fmt.Errorf("can't find local varable '%v' in scope '%v'", e.Name, sc)
	}
	if e.startExpr != nil {
		if e.startExpr, err = e.StartExprWrap.Analyze(module, scope); err != nil {
			return nil, err
		}
	} else {
		return e, nil
	}
	if e.endExpr != nil {
		if e.endExpr, err = e.EndExprWrap.Analyze(module, scope); err != nil {
			return nil, err
		}
	}
	return e, nil
}

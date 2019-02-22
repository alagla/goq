package types

import "fmt"

type QuplaSliceExpr struct {
	Name          string                  `yaml:"name"`
	StartExprWrap *QuplaExpressionWrapper `yaml:"start,omitempty"`
	EndExprWrap   *QuplaExpressionWrapper `yaml:"end,omitempty"`
	//----
	localVar  *LocalVariable
	startExpr ExpressionInterface
	sizeExpr  ExpressionInterface
	offset    int64
	size      int64
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
		return nil, fmt.Errorf("can't find local variable '%v' in scope '%v'", e.Name, sc)
	}
	if e.StartExprWrap != nil {
		if e.startExpr, err = e.StartExprWrap.Analyze(module, scope); err != nil {
			return nil, err
		}
		if e.offset, err = GetConstValue(e.startExpr); err != nil || e.offset >= e.localVar.size {
			return nil, fmt.Errorf("offset must be constant expression less than var size in SliceExpr")
		}
		e.size = 1
	} else {
		e.startExpr = nil
		e.offset = 0
		e.size = e.localVar.size
		return e, nil
	}
	if e.EndExprWrap != nil {
		if e.sizeExpr, err = e.EndExprWrap.Analyze(module, scope); err != nil {
			return nil, err
		}
		if e.size, err = GetConstValue(e.sizeExpr); err != nil || e.size <= 0 {
			return nil, fmt.Errorf("size must be positive constant expression in SliceExpr")
		}
	} else {
		e.sizeExpr = nil
	}
	return e, nil
}

func (e *QuplaSliceExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

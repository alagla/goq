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
	start     int64
	end       int64
}

func (e *QuplaSliceExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error

	if e.localVar, err = scope.FindVar(e.Name, module); err != nil {
		return nil, err
	}

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
		if e.start, err = GetConstValue(e.startExpr); err != nil || e.start < 0 || e.start >= e.localVar.size {
			return nil, fmt.Errorf("wrong start offset in SliceExpr")
		}
		e.end = e.start + 1
	} else {
		e.startExpr = nil
		e.start = 0
		e.end = e.localVar.size
		return e, nil
	}
	if e.EndExprWrap != nil {
		if e.sizeExpr, err = e.EndExprWrap.Analyze(module, scope); err != nil {
			return nil, err
		}
		var s int64
		if s, err = GetConstValue(e.sizeExpr); err != nil || s > e.localVar.size {
			return nil, fmt.Errorf("wrong slice size in SliceExpr")
		}
		e.end = e.start + s
	} else {
		e.sizeExpr = nil
		e.end = e.start + 1
	}
	return e, nil
}

func (e *QuplaSliceExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.end - e.start
}

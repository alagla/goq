package types

import "fmt"

type QuplaSliceExpr struct {
	Name          string                  `yaml:"name"`
	StartExprWrap *QuplaExpressionWrapper `yaml:"start,omitempty"`
	EndExprWrap   *QuplaExpressionWrapper `yaml:"end,omitempty"`
	//----
	localVar  *LocalVariable
	startExpr ExpressionInterface
	endExpr   ExpressionInterface
	fromIdx   int64
	toIdx     int64
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
	if e.startExpr != nil {
		if e.startExpr, err = e.StartExprWrap.Analyze(module, scope); err != nil {
			return nil, err
		}
		if !IsConstExpression(e.startExpr) {
			return nil, fmt.Errorf("must be constant expression in SliceExpr")
		}
	} else {
		return e, nil
	}
	if e.endExpr != nil {
		if e.endExpr, err = e.EndExprWrap.Analyze(module, scope); err != nil {
			return nil, err
		}
		if !IsConstExpression(e.endExpr) {
			return nil, fmt.Errorf("must be constant expression in SliceExpr")
		}
		e.toIdx, _ = GetConstValue(e.endExpr)
		e.fromIdx, _ = GetConstValue(e.startExpr)
		if e.toIdx <= e.fromIdx {
			sc := "nil"
			if scope != nil {
				sc = scope.Name
			}
			return nil, fmt.Errorf("scope '%v': wrong slice range in SliceExpr", sc)
		}
	}
	return e, nil
}

func (e *QuplaSliceExpr) Size() int64 {
	if e == nil {
		return 0
	}
	switch {
	case e.startExpr == nil && e.endExpr == nil:
		return e.localVar.size
	case e.startExpr != nil && e.endExpr == nil:
		return 1
	case e.startExpr != nil && e.endExpr != nil:
		return e.toIdx - e.fromIdx
	}
	panic("inconsistency")
}

func (e *QuplaSliceExpr) RequireSize(size int64) error {
	sz := e.Size()
	if size != sz || e.fromIdx >= sz || e.toIdx >= sz {
		return fmt.Errorf("wrong or size mismatch in SliceExpr")
	}
	return nil
}

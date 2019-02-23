package types

import "fmt"

type QuplaSliceExpr struct {
	Var           string                  `yaml:"var"`
	Offset        int64                   `yaml:"offset"`
	SliceSize     int64                   `yaml:"size"`
	StartExprWrap *QuplaExpressionWrapper `yaml:"start,omitempty"` // not used
	EndExprWrap   *QuplaExpressionWrapper `yaml:"end,omitempty"`   // not used
	//----
	localVar *LocalVariable
}

func (e *QuplaSliceExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error

	if e.localVar, err = scope.FindVar(e.Var, module); err != nil {
		return nil, err
	}
	module.IncStat("numSliceExpr")
	if e.localVar == nil {
		var sc string
		if scope != nil {
			sc = scope.name
		} else {
			sc = "nil"
		}
		return nil, fmt.Errorf("can't find local variable '%v' in scope '%v'", e.Var, sc)
	}
	if e.Offset+e.SliceSize > e.localVar.size {
		return nil, fmt.Errorf("wrong offset/size for the slice of '%v'", e.Var)
	}
	return e, nil
}

func (e *QuplaSliceExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.SliceSize
}

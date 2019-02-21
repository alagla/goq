package types

import "fmt"

type QuplaLutExpr struct {
	Name string                    `yaml:"name"`
	Args []*QuplaExpressionWrapper `yaml:"args"`
	//----
	argExpr []ExpressionInterface
	lutDef  *QuplaLutDef
}

func (e *QuplaLutExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	var ae ExpressionInterface
	e.argExpr = make([]ExpressionInterface, 0, len(e.Args))
	for _, a := range e.Args {
		ae, err = a.Analyze(module, scope)
		if err != nil {
			return nil, err
		}
		e.argExpr = append(e.argExpr, ae)
	}
	e.lutDef, err = module.FindLUTDef(e.Name)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *QuplaLutExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.lutDef.Size()
}

func (e *QuplaLutExpr) RequireSize(size int64) error {
	if size != e.Size() {
		return fmt.Errorf("size mismatch in LutExpr")
	}
	return nil
}

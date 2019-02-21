package types

import "fmt"

type QuplaLutExpr struct {
	Name string                    `yaml:"Name"`
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
	e.lutDef = module.FindLUTDef(e.Name)
	if e.lutDef == nil {
		return nil, fmt.Errorf("can't find LUT definition for '%v'", e.Name)
	}
	return e, nil
}

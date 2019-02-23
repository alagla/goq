package program

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
	e.lutDef, err = module.FindLUTDef(e.Name)
	if err != nil {
		return nil, err
	}
	module.IncStat("numLUTExpr")

	e.argExpr = make([]ExpressionInterface, 0, len(e.Args))
	for _, a := range e.Args {
		ae, err = a.Analyze(module, scope)
		if err != nil {
			return nil, err
		}
		if err = RequireSize(ae, 1); err != nil {
			return nil, fmt.Errorf("LUT expression with '%v': %v", e.lutDef.name, err)
		}
		e.argExpr = append(e.argExpr, ae)
	}
	if e.lutDef.inputSize != len(e.argExpr) {
		return nil, fmt.Errorf("num arg doesnt't match input dimension of the LUT %v", e.lutDef.name)
	}
	return e, nil
}

func (e *QuplaLutExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.lutDef.Size()
}

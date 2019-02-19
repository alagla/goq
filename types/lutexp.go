package types

import "fmt"

type QuplaLutExpr struct {
	Name string                    `yaml:"name"`
	Args []*QuplaExpressionWrapper `yaml:"args"`
	//----
	argExpr []ExpressionInterface
	lutDef  *QuplaLutDef
}

func (e *QuplaLutExpr) Analyze(module *QuplaModule) error {
	var err error
	var ae ExpressionInterface
	e.argExpr = make([]ExpressionInterface, 0, len(e.Args))
	for _, a := range e.Args {
		ae, err = a.Unwarp()
		if err != nil {
			return err
		}
		err = ae.Analyze(module)
		if err != nil {
			return err
		}
		e.argExpr = append(e.argExpr, ae)
	}
	e.lutDef = module.FindLUTDef(e.Name)
	if e.lutDef == nil {
		return fmt.Errorf("can't find LUT definition for '%v'", e.Name)
	}
	return nil
}

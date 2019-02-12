package types

import "fmt"

type QuplaLutExpr struct {
	Name string                    `yaml:"name"`
	Args []*QuplaExpressionWrapper `yaml:"args"`
	//----
	lutDef *QuplaLutDef
}

func (lutExpr *QuplaLutExpr) Analyze(module *QuplaModule) error {
	lutExpr.lutDef = module.FindLUTDef(lutExpr.Name)
	if lutExpr.lutDef == nil {
		return fmt.Errorf("can't find LUT definition for '%v'", lutExpr.Name)
	}
	for _, arg := range lutExpr.Args {
		if err := arg.Analyze(module); err != nil {
			return err
		}
	}
	return nil
}

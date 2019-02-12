package types

import "fmt"

type QuplaFuncExpr struct {
	Name string                    `yaml:"name"`
	Args []*QuplaExpressionWrapper `yaml:"args"`
	//---
	funcDef *QuplaFuncDef
}

func (funcExpr *QuplaFuncExpr) Analyze(module *QuplaModule) error {
	for _, arg := range funcExpr.Args {
		if err := arg.Analyze(module); err != nil {
			return err
		}
	}
	funcExpr.funcDef = module.FindFuncDef(funcExpr.Name)
	if funcExpr.funcDef == nil {
		return fmt.Errorf("can't find function definition '%v'", funcExpr.Name)
	}
	return nil
}

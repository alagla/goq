package types

import "fmt"

type QuplaFuncExpr struct {
	Name     string                    `yaml:"name"`
	ArgsWrap []*QuplaExpressionWrapper `yaml:"args"`
	//---
	funcDef *QuplaFuncDef
	args    []ExpressionInterface
}

func (e *QuplaFuncExpr) Analyze(module *QuplaModule) error {
	e.funcDef = module.FindFuncDef(e.Name)
	if e.funcDef == nil {
		return fmt.Errorf("can't find function definition '%v'", e.Name)
	}
	var err error
	var fe ExpressionInterface

	e.args = make([]ExpressionInterface, 0, len(e.ArgsWrap))
	for _, arg := range e.ArgsWrap {
		if fe, err = arg.Unwarp(); err != nil {
			return err
		}
		if err := fe.Analyze(module); err != nil {
			return err
		}
		e.args = append(e.args, fe)
	}
	return nil
}

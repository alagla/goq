package types

type QuplaFuncExpr struct {
	Name     string                    `yaml:"name"`
	ArgsWrap []*QuplaExpressionWrapper `yaml:"args"`
	//---
	funcDef *QuplaFuncDef
	args    []ExpressionInterface
}

func (e *QuplaFuncExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	e.funcDef, err = module.FindFuncDef(e.Name)
	if err != nil {
		return nil, err
	}
	var fe ExpressionInterface
	module.IncStat("numFuncExpr")

	e.args = make([]ExpressionInterface, 0, len(e.ArgsWrap))
	for _, arg := range e.ArgsWrap {
		if fe, err = arg.Analyze(module, scope); err != nil {
			return nil, err
		}
		e.args = append(e.args, fe)
	}
	err = e.funcDef.checkArgSizes(e.args)
	return e, err
}

func (e *QuplaFuncExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.funcDef.Size()
}

package types

type QuplaCondExpr struct {
	If   *QuplaExpressionWrapper `yaml:"if"`
	Then *QuplaExpressionWrapper `yaml:"then"`
	Else *QuplaExpressionWrapper `yaml:"else"`
	//--
	ifExpr   ExpressionInterface
	thenExpr ExpressionInterface
	elseExpr ExpressionInterface
}

func (e *QuplaCondExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error

	if e.ifExpr, err = e.If.Analyze(module, scope); err != nil {
		return nil, err
	}
	if e.thenExpr, err = e.Then.Analyze(module, scope); err != nil {
		return nil, err
	}
	if e.elseExpr, err = e.Else.Analyze(module, scope); err != nil {
		return nil, err
	}
	return e, nil
}

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

func (e *QuplaCondExpr) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	var err error

	if e.ifExpr, err = e.If.Analyze(module); err != nil {
		return nil, err
	}
	if e.thenExpr, err = e.Then.Analyze(module); err != nil {
		return nil, err
	}
	if e.elseExpr, err = e.Else.Analyze(module); err != nil {
		return nil, err
	}
	return e, nil
}

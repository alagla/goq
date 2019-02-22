package types

import "fmt"

type ExpressionInterface interface {
	Analyze(*QuplaModule, *QuplaFuncDef) (ExpressionInterface, error)
	Size() int64
}

type QuplaExecStmt struct {
	ExprWrap     *QuplaExpressionWrapper `yaml:"expr"`
	ExpectedWrap *QuplaExpressionWrapper `yaml:"expected,omitempty"`
	//---
	isTest       bool
	expr         ExpressionInterface
	exprExpected ExpressionInterface
}

type QuplaTypeDef struct {
	Size   string                            `yaml:"end"`
	Fields map[string]*struct{ Size string } `yaml:"fields,omitempty"`
}

type QuplaNullExpr struct{}

func (e *QuplaNullExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	return e, nil
}

func (e *QuplaNullExpr) Size() int64 {
	return 0
}

func MatchSizes(e1, e2 ExpressionInterface) error {
	s1 := e1.Size()
	s2 := e2.Size()

	if s1 != s2 {
		return fmt.Errorf("sizes doesn't match: %v != %v", s1, s2)
	}
	return nil
}

func RequireSize(e ExpressionInterface, size int64) error {
	s := e.Size()

	if s != size {
		return fmt.Errorf("sizes doesn't match: required %v != %v", size, s)
	}
	return nil
}

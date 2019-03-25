package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type ConstExpression interface {
	ExpressionInterface
	GetConstValue() int
	GetConstName() string
}

type ConstValue struct {
	ExpressionBase
	value int
	name  string
	size  int
}

func (e *ConstValue) Size() int {
	return 0 //todo ??
}

func (_ *ConstValue) Eval(_ *EvalFrame, _ Trits) bool {
	return true // todo ??
}

func (e *ConstValue) GetConstValue() int {
	return e.value
}

func (e *ConstValue) GetConstName() string {
	return e.name
}

func NewConstValue(name string, value int) *ConstValue {
	return &ConstValue{
		name:  name,
		value: value,
	}
}

func (e *ConstValue) GetTrits() Trits {
	t := IntToTrits(int64(e.value))
	if e.size == 0 {
		return t
	}
	if e.size == len(t) {
		return t
	}
	ret := make(Trits, 0, e.size)
	copy(ret, t)
	return ret
}

func GetConstValue(expr ExpressionInterface) (int, error) {
	ce, ok := expr.(ConstExpression)
	if !ok {
		return 0, fmt.Errorf("not a constant value")
	}
	return ce.GetConstValue(), nil
}

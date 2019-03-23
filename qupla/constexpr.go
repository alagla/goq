package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type ConstExpression interface {
	ExpressionInterface
	GetConstValue() int64
	GetConstName() string
}

type ConstValue struct {
	QuplaExprBase
	value int64
	name  string
	size  int64
}

func (e *ConstValue) Size() int64 {
	return 0 //todo
}

func (_ *ConstValue) Eval(_ ProcessorInterface, _ Trits) bool {
	return true // todo
}

func (e *ConstValue) GetConstValue() int64 {
	return e.value
}

func (e *ConstValue) GetConstName() string {
	return e.name
}

func NewConstValue(name string, value int64) ConstExpression {
	return &ConstValue{
		name:  name,
		value: value,
	}
}

func (e *ConstValue) GetTrits() Trits {
	t := IntToTrits(e.value)
	if e.size == 0 {
		return t
	}
	if e.size == int64(len(t)) {
		return t
	}
	ret := make(Trits, 0, e.size)
	copy(ret, t)
	return ret
}

func GetConstValue(expr ExpressionInterface) (int64, error) {
	ce, ok := expr.(ConstExpression)
	if !ok {
		return 0, fmt.Errorf("not a constant value")
	}
	return ce.GetConstValue(), nil
}

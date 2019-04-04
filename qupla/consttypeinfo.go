package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type ConstTypeFieldInfo struct {
	Offset int
	Size   int
}

type ConstTypeInfo struct {
	name   string
	size   int
	Fields map[string]*ConstTypeFieldInfo
}

func NewConstTypeInfo(name string, size int) *ConstTypeInfo {
	return &ConstTypeInfo{
		name:   name,
		size:   size,
		Fields: make(map[string]*ConstTypeFieldInfo, 5),
	}
}

func (e *ConstTypeInfo) GetConstValue() int {
	return e.size
}

func (e *ConstTypeInfo) GetConstName() string {
	return e.name
}

func (e *ConstTypeInfo) GetSource() string {
	return ""
}

func (e *ConstTypeInfo) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	return &ConstTypeInfo{
		name:   e.name,
		size:   e.size,
		Fields: e.Fields,
	}
}

func (e *ConstTypeInfo) Size() int {
	return 0
}

func (e *ConstTypeInfo) Eval(_ *EvalFrame, _ Trits) bool {
	return true
}

func (e *ConstTypeInfo) References(_ string) bool {
	return false
}

func (e *ConstTypeInfo) GetFieldInfo(fldname string) (*ConstTypeFieldInfo, error) {
	fi, ok := e.Fields[fldname]
	if !ok {
		return nil, fmt.Errorf("can;t find field '%v' in type '%v'", fldname, e.name)
	}
	return fi, nil
}

func (e *ConstTypeInfo) HasState() bool {
	return false
}

func (e *ConstTypeInfo) GetSubexpressions() []ExpressionInterface {
	return nil
}

func (e *ConstTypeInfo) SetSubexpressions(_ []ExpressionInterface) {
}

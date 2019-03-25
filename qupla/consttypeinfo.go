package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type ConstTypeFieldInfo struct {
	Offset int64
	Size   int64
}

type ConstTypeInfo struct {
	name   string
	size   int64
	Fields map[string]*ConstTypeFieldInfo
}

func NewConstTypeInfo(name string, size int64) *ConstTypeInfo {
	return &ConstTypeInfo{
		name:   name,
		size:   size,
		Fields: make(map[string]*ConstTypeFieldInfo, 5),
	}
}

func (e *ConstTypeInfo) GetConstValue() int64 {
	return e.size
}

func (e *ConstTypeInfo) GetConstName() string {
	return e.name
}

func (e *ConstTypeInfo) GetSource() string {
	return ""
}

func (e *ConstTypeInfo) Size() int64 {
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

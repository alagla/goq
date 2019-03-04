package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
	"strconv"
)

type ConstTypeFieldInfo struct {
	offset int64
	size   int64
}

type ConstTypeInfo struct {
	name   string
	size   int64
	fields map[string]*ConstTypeFieldInfo
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

func (e *ConstTypeInfo) Eval(_ ProcessorInterface, _ Trits) bool {
	return true
}

func (e *ConstTypeInfo) References(_ string) bool {
	return false
}

func (e *ConstTypeInfo) GetFieldInfo(fldname string) (*ConstTypeFieldInfo, error) {
	fi, ok := e.fields[fldname]
	if !ok {
		return nil, fmt.Errorf("can;t find field '%v' in type '%v'", fldname, e.name)
	}
	return fi, nil
}

func AnalyzeConstTypeName(exprYAML *QuplaConstTypeNameYAML, _ ModuleInterface, funcDef FuncDefInterface) (ConstExpression, error) {
	ret := &ConstTypeInfo{
		name:   exprYAML.TypeName,
		fields: make(map[string]*ConstTypeFieldInfo),
	}
	if sz, err := strconv.Atoi(exprYAML.SizeString); err != nil {
		return nil, err
	} else {
		ret.size = int64(sz)
	}
	var err error
	var size, offset int

	for fldname, fld := range exprYAML.Fields {
		if size, err = strconv.Atoi(fld.SizeString); err == nil {
			offset, err = strconv.Atoi(fld.OffsetString)
		}
		if err != nil {
			return nil, fmt.Errorf("wrong size or offset in field '%v' in type info '%v in func def '%v': %v",
				fldname, exprYAML.TypeName, funcDef.GetName(), err)
		}
		ret.fields[fldname] = &ConstTypeFieldInfo{
			offset: int64(offset),
			size:   int64(size),
		}
	}
	return ret, nil
}

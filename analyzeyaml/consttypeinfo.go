package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"strconv"
)

func AnalyzeConstTypeName(exprYAML *QuplaConstTypeNameYAML, _ *QuplaModule, funcDef *Function) (ConstExpression, error) {
	sz, err := strconv.Atoi(exprYAML.SizeString)
	if err != nil {
		return nil, err
	}
	ret := NewConstTypeInfo(exprYAML.TypeName, sz)

	var offset int

	for fldname, fld := range exprYAML.Fields {
		if sz, err = strconv.Atoi(fld.SizeString); err == nil {
			offset, err = strconv.Atoi(fld.OffsetString)
		}
		if err != nil {
			return nil, fmt.Errorf("wrong size or offset in field '%v' in type info '%v in func def '%v': %v",
				fldname, exprYAML.TypeName, funcDef.Name, err)
		}
		ret.Fields[fldname] = &ConstTypeFieldInfo{
			Offset: offset,
			Size:   sz,
		}
	}
	return ret, nil
}

package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type fieldExpr struct {
	offset int64
	size   int64
}
type QuplaTypeExpr struct {
	QuplaExprBase
	size   int64
	fields []fieldExpr
}

func AnalyzeTypeExpr(exprYAML *QuplaTypeExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaTypeExpr, error) {
	ret := &QuplaTypeExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
		fields:        make([]fieldExpr, 0, len(exprYAML.Fields)),
	}
	var err error
	module.IncStat("numTypeExpr")

	var constexpr ExpressionInterface
	if constexpr, err = module.AnalyzeExpression(exprYAML.TypeInfo, scope); err != nil {
		return nil, err
	}

	typeInfo, ok := constexpr.(*ConstTypeInfo)
	if !ok {
		return nil, fmt.Errorf("type info expected in '%v': '%v'", scope.GetName(), exprYAML.Source)
	}
	ret.size = typeInfo.GetConstValue()

	var fe ExpressionInterface
	var fi *ConstTypeFieldInfo
	var sumFld int64
	for fldName, expr := range exprYAML.Fields {
		fi, err = typeInfo.GetFieldInfo(fldName)
		if err != nil {
			return nil, err
		}
		if fe, err = module.AnalyzeExpression(expr, scope); err != nil {
			return nil, err
		}
		if fe.Size() != fi.size {
			return nil, fmt.Errorf("field '%v' size mismatch in type expression '%v'", fldName, exprYAML.Source)
		}
		ret.fields = append(ret.fields, fieldExpr{
			offset: fi.offset,
			size:   fi.size,
		})
		sumFld += fi.size
	}
	if sumFld != ret.size {
		return nil, fmt.Errorf("sum of field sizes != type size in field expression '%v'", ret.source)
	}
	return ret, nil
}

func (e *QuplaTypeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *QuplaTypeExpr) Eval(proc ProcessorInterface, result Trits) bool {
	for idx, subExpr := range e.subexpr {
		if proc.Eval(subExpr, result[e.fields[idx].offset:e.fields[idx].offset+e.fields[idx].size]) {
			return true
		}
	}
	return false
}

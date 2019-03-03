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
	expr   ExpressionInterface
}
type QuplaTypeExpr struct {
	QuplaExprBase
	size   int64
	fields map[string]*fieldExpr
}

func AnalyzeTypeExpr(exprYAML *QuplaTypeExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaTypeExpr, error) {
	ret := &QuplaTypeExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
		fields:        make(map[string]*fieldExpr),
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

		ret.fields[fldName] = &fieldExpr{
			offset: fi.offset,
			size:   fi.size,
			expr:   fe, // TODO must be condExpr by syntax. Not exactly ConditionExpression
		}
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
	for _, fi := range e.fields {
		if proc.Eval(fi.expr, result[fi.offset:fi.offset+fi.size]) {
			return true
		}
	}
	return false
}

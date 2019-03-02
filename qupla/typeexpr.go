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
	if constexpr, err = module.AnalyzeExpression(exprYAML.TypeNameConst, scope); err != nil {
		return nil, err
	}

	if ret.size, err = GetConstValue(constexpr); err != nil {
		return nil, err
	}
	var typeName string
	if typeName, err = GetConstName(constexpr); err != nil {
		return nil, err
	}

	var fe ExpressionInterface
	var offset, size int64

	for fldName, expr := range exprYAML.Fields {
		offset, size, err = module.GetTypeFieldInfo(typeName, fldName)

		if fe, err = module.AnalyzeExpression(expr, scope); err != nil {
			return nil, err
		}
		if fe.Size() != size {
			return nil, fmt.Errorf("field '%v' size mismatch in type expression '%v'", fldName, exprYAML.Source)
		}

		ret.fields[fldName] = &fieldExpr{
			offset: offset,
			size:   size,
			expr:   fe, // TODO must be condExpr by syntaxis. Not exactly ConditionExpression
		}
	}
	if err = ret.CheckSize(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *QuplaTypeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *QuplaTypeExpr) CheckSize() error {
	var sumFld int64

	for _, f := range e.fields {
		sumFld += f.size
	}
	if sumFld != e.size {
		return fmt.Errorf("sum of field sizes != type end in field expression")
	}
	return nil
}

func (e *QuplaTypeExpr) Eval(proc ProcessorInterface, result Trits) bool {
	for _, fi := range e.fields {
		if proc.Eval(fi.expr, result[fi.offset:fi.size]) {

		}
	}
	return true
}

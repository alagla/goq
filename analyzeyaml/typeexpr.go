package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeTypeExpr(exprYAML *QuplaTypeExprYAML, module *QuplaModule, scope *Function) (*TypeExpr, error) {
	var err error
	module.IncStat("numTypeExpr")

	var constexpr ExpressionInterface
	if constexpr, err = AnalyzeExpression(exprYAML.TypeInfo, module, scope); err != nil {
		return nil, err
	}

	typeInfo, ok := constexpr.(*ConstTypeInfo)
	if !ok {
		return nil, fmt.Errorf("type info expected in '%v': '%v'", scope.Name, exprYAML.Source)
	}
	ret := NewQuplaTypeExpr(exprYAML.Source, typeInfo.GetConstValue())

	var fe ExpressionInterface
	var fi *ConstTypeFieldInfo
	var sumFld int64
	for fldName, expr := range exprYAML.Fields {
		fi, err = typeInfo.GetFieldInfo(fldName)
		if err != nil {
			return nil, err
		}
		if fe, err = AnalyzeExpression(expr, module, scope); err != nil {
			return nil, err
		}
		if fe.Size() != fi.Size {
			return nil, fmt.Errorf("field '%v' Size mismatch in type expression '%v'", fldName, exprYAML.Source)
		}
		ret.Fields = append(ret.Fields, FieldExpr{
			Offset: fi.Offset,
			Size:   fi.Size,
		})
		ret.AppendSubExpr(fe)
		sumFld += fi.Size
	}
	if sumFld != ret.Size() {
		return nil, fmt.Errorf("sum of field sizes != type Size in field expression '%v'", exprYAML.Source)
	}
	return ret, nil
}

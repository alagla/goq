package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeSliceExpr(exprYAML *QuplaSliceExprYAML, module *QuplaModule, scope *Function) (*SliceExpr, error) {
	ret := NewQuplaSliceExpr(exprYAML.Source, exprYAML.Offset, exprYAML.SliceSize)
	module.IncStat("numSliceExpr")
	ret.VarScope = scope
	vi, err := scope.VarByName(exprYAML.Var)
	if err != nil {
		return nil, fmt.Errorf("'%v': %v", scope.Name, err)
	}
	ret.LocalVarIdx = vi.Idx
	if ret.LocalVarIdx < 0 {
		return nil, fmt.Errorf("can't find local variable '%v' in scope '%v'", exprYAML.Var, scope.Name)
	}
	// can't do it because in recursive situations var can be not analysed yet
	//if ret.offset+ret.size > vi.Size {
	//	return nil, fmt.Errorf("wrong offset/size for the slice of '%v'", exprYAML.Var)
	//}
	return ret, nil
}

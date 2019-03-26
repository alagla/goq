package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeSliceExpr(exprYAML *QuplaSliceExprYAML, module *QuplaModule, scope *Function) (*SliceExpr, error) {
	module.IncStat("numSliceExpr")
	vi, err := scope.VarByName(exprYAML.Var)
	if err != nil {
		return nil, fmt.Errorf("'%v': %v", scope.Name, err)
	}
	ret := NewQuplaSliceExpr(vi, exprYAML.Source, exprYAML.Offset, exprYAML.SliceSize)

	// can't do it because in recursive situations var can be not analysed yet
	//if ret.offset+ret.size > vi.Size {
	//	return nil, fmt.Errorf("wrong offset/size for the slice of '%v'", exprYAML.Var)
	//}
	return ret, nil
}

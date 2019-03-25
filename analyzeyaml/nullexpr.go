package analyzeyaml

import (
	. "github.com/lunfardo314/goq/qupla"
)

func AnalyzeNullExpr(module *QuplaModule) (*NullExpr, error) {
	module.IncStat("nullExpr")
	return &NullExpr{}, nil
}

package analyzeyaml

import (
	. "github.com/lunfardo314/goq/qupla"
)

func AnalyzeNullExpr(module *QuplaModule) (*QuplaNullExpr, error) {
	module.IncStat("nullExpr")
	return &QuplaNullExpr{}, nil
}

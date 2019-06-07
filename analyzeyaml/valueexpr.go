package analyzeyaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeValueExpr(exprYAML *QuplaValueExprYAML, module *QuplaModule) (*ValueExpr, error) {
	module.IncStat("numValueExpr")

	if exprYAML.Trits == "" {
		return nil, fmt.Errorf("invalid empty trit string in ValueExpr")
	}
	t := make([]int8, len(exprYAML.Trits))
	for i := range exprYAML.Trits {
		switch exprYAML.Trits[i] {
		case '-':
			t[i] = -1
		case '0':
			t[i] = 0
		case '1':
			t[i] = 1
		default:
			return nil, fmt.Errorf("invalid trit string '%v'", exprYAML.Trits)
		}
	}
	var err error
	if t, err = NewTrits(t); err != nil {
		return nil, err
	}
	return NewValueExpr(t), nil

}

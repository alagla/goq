package analyzeyaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	"math/big"
)

func AnalyzeValueExpr(exprYAML *QuplaValueExprYAML, module *QuplaModule) (*ValueExpr, error) {
	module.IncStat("numValueExpr")

	if exprYAML.Trits == "" {
		return nil, fmt.Errorf("invalid trit string n ValueExpr '%v'", exprYAML.Trits)
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
	orig := big.NewFloat(0)
	var err error
	var ok bool
	if exprYAML.Value == "-" {
		orig.SetInt64(-1)
	} else {
		orig, ok = orig.SetString(exprYAML.Value)
		if !ok {
			return nil, fmt.Errorf("wrong 'value' field '%v' in ValueExpr", exprYAML.Value)
		}
	}
	if t, err = NewTrits(t); err != nil {
		return nil, err
	}
	return NewValueExpr(t), nil

	// Todo checking float values

	//bi, err := utils.TritsToBigInt(ret.TritValue)
	//bif := big.NewFloat(0)
	//bif.SetInt(bi)
	//
	//if err != nil{
	//	return nil, fmt.Errorf("can't convert trits to BigInt")
	//}
	//if orig.Cmp(bif) != 0{
	//	return nil, fmt.Errorf("not equal values between trits and decimal '%v' != '%v'", orig, bif)
	//}
}

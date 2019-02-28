package utils

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"math/big"
)

func TritsToString(trits Trits) string {
	b := make([]byte, len(trits), len(trits))
	for i := range trits {
		switch trits[i] {
		case -1:
			b[i] = '-'
		case 0:
			b[i] = '0'
		case 1:
			b[i] = '1'
		default:
			b[i] = '?'
		}
	}
	return string(b)
}

var (
	a0 = big.NewInt(0)
	a1 = big.NewInt(1)
	a2 = big.NewInt(2)
)

func TritToBigInt(t int8) (*big.Int, error) {
	switch t {
	case -1:
		return a0, nil
	case 0:
		return a1, nil
	case 1:
		return a2, nil
	}
	return nil, fmt.Errorf("wrong trit value")
}

var big3 = big.NewInt(3)

func TritsToBigInt(t Trits) (*big.Int, error) {
	if err := ValidTrits(t); err != nil {
		return nil, err
	}
	ret := big.NewInt(0)
	var err error
	var trit *big.Int
	for i := len(t) - 1; i >= 0; i-- {
		if trit, err = TritToBigInt(t[i]); err != nil {
			return nil, err
		}
		ret.Mul(big3, ret)
		ret.Add(ret, trit)
	}
	return ret, nil
}

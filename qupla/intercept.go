package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"strings"
)

type interception struct {
	namePrefix string
	onReturn   func(Trits)
}

var interceptions = []*interception{
	{"print", print},
	{"threeXplusOne", print},
}

func print(value Trits) {
	bi, _ := utils.TritsToBigInt(value)
	logf(2, "  value %v, '%v'", bi, utils.TritsToString(value))
}

func getOnReturnInterceptions(funcName string) []func(Trits) {
	ret := make([]func(Trits), 0)
	for _, ic := range interceptions {
		if strings.HasPrefix(funcName, ic.namePrefix) {
			ret = append(ret, ic.onReturn)
		}
	}
	return ret
}

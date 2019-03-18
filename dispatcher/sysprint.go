package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
)

type sysPrintEntityCore struct {
}

func (fc *sysPrintEntityCore) Call(args Trits, _ Trits) bool {
	fmt.Println(utils.TritsToString(args))
	return true
}

func NewSysPrintEntity(disp *Dispatcher) *Entity {
	return disp.NewEntity(EntityOpts{
		Name:     "sysprint",
		InSize:   0,
		OutSize:  0,
		Core:     &sysPrintEntityCore{},
		Terminal: true,
	})
}

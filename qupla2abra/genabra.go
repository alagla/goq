package qupla2abra

import (
	"github.com/lunfardo314/goq/qupla"
)

type AbraIR struct {
	module *qupla.QuplaModule
	luts   []*qupla.LutDef // luts for Tritcode (not include size > 1)
}

func PrepareQupla4Abra(module *qupla.QuplaModule) *AbraIR {
	ret := AbraIR{
		module: module,
		luts:   make([]*qupla.LutDef, 0, len(module.Luts)+10),
	}
	ret.PrepareLuts()
	return &ret
}

func (air *AbraIR) PrepareLuts() {
	for _, lutDef := range air.module.Luts {
		if lutDef.Size() == 1 {
			air.luts = append(air.luts, lutDef)
		} else {

		}
	}
}

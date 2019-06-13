package qupla2abra

import (
	"github.com/lunfardo314/goq/qupla"
)

type AbraIR struct {
	module *qupla.QuplaModule
}

func PrepareQupla4Abra(module *qupla.QuplaModule) *AbraIR {
	ret := AbraIR{
		module: module,
	}

	return &ret
}

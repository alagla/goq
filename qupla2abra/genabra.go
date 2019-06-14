package qupla2abra

import (
	"github.com/lunfardo314/goq/qupla"
)

type AbraIR struct {
	module    *qupla.QuplaModule
	luts      []string // string representation
	functions []string
}

func (air *AbraIR) FindLutBlockIndex(lutDef *qupla.LutDef) int {
	strRepr := lutDef.GetStringRepr()
	for i, sr := range air.luts {
		if strRepr == sr {
			return i
		}
	}
	return -1
}

func (air *AbraIR) EnsureLutIndex(lutDef *qupla.LutDef) int {
	if ret := air.FindLutBlockIndex(lutDef); ret >= 0 {
		return ret
	}
	air.luts = append(air.luts, lutDef.GetStringRepr())
	return len(air.luts)
}

func (air *AbraIR) FindFunctionBranchBlockIndex(fun *qupla.Function) int {
	for i := range air.functions {
		if air.functions[i] == fun.Name {
			return i
		}
	}
	return -1
}

func (air *AbraIR) EnsureFunctionIndex(fun *qupla.Function) int {
	if ret := air.FindFunctionBranchBlockIndex(fun); ret >= 0 {
		return ret
	}
	air.luts = append(air.luts, fun.Name)
	return len(air.luts)
}

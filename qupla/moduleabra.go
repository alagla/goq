package qupla

import (
	"github.com/lunfardo314/goq/abra"
	. "github.com/lunfardo314/goq/cfg"
	"sort"
)

func (module *QuplaModule) GetAbra(codeUnit *abra.CodeUnit) {
	// sort LUTs by name
	names := make([]string, 0)
	for n := range module.Luts {
		names = append(names, n)
	}
	sort.Strings(names)

	// TODO environments etc
	for _, n := range names {
		strRepr := module.Luts[n].GetStringRepr()
		Logf(2, "generating LUT block: '%s' -> '%s'", n, strRepr)
		codeUnit.NewLUTBlock(strRepr, abra.BinaryEncodedLUTFromString(strRepr))
	}

	// sort functions by name
	names = make([]string, 0)
	for n := range module.Functions {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		Logf(2, "generating branch for '%s'", n)
		module.Functions[n].GetAbraBranchBlock(codeUnit)
	}
}

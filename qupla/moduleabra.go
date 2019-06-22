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
	Logf(2, "---- generating LUT blocks")
	count := 0
	for _, n := range names {
		strRepr := module.Luts[n].GetStringRepr()
		if codeUnit.FindLUTBlock(strRepr) != nil {
			Logf(2, "%20s -> '%s'  repeated", n, strRepr)
			continue
		}
		Logf(2, "%20s -> '%s'", n, strRepr)
		codeUnit.GetLUTBlock(strRepr)
		count++
	}
	Logf(2, "total generated %d LUT blocks out of %d LUT definitions", count, len(names))

	// sort functions by name
	names = make([]string, 0)
	for n := range module.Functions {
		names = append(names, n)
	}
	sort.Strings(names)
	Logf(2, "---- generating branch blocks")
	for _, n := range names {
		b := module.Functions[n].GetAbraBranchBlock(codeUnit)
		st := b.Branch.GetStats()
		Logf(2, "%s -> inputs: %d bodySites: %d, stateSites: %d, outputs: %d, knots: %d, merges: %d",
			n, st.NumInputs, st.NumBodySites, st.NumStateSites, st.NumOutputs, st.NumKnots, st.NumMerges)
	}
}

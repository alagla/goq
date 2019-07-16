package qupla

import (
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
	vabra "github.com/lunfardo314/goq/abra/validate"
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
		if cabra.FindLUTBlock(codeUnit, strRepr) != nil {
			continue
		}
		Logf(2, "%20s -> '%s'", n, strRepr)
		cabra.GetLUTBlock(codeUnit, strRepr)
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
		st := vabra.GetStats(b.Branch)
		Logf(2, "%30s -> inputs: %2d, outputs: %2d, bodySites: %2d, stateSites: %2d, knots: %2d, merges: %2d inSizes: %v=%d",
			n, st.NumInputs, st.NumOutputs, st.NumBodySites, st.NumStateSites, st.NumKnots, st.NumMerges, st.InputSizes, st.InputSize)
	}
}

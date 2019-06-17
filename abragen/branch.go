package abragen

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/qupla"
	"strings"
)

func equalSites(inds1, inds2 []*SiteIR) bool {
	if len(inds1) != len(inds2) {
		return false
	}
	for i, idx := range inds1 {
		if idx != inds2[i] {
			return false
		}
	}
	return true
}

func (branch *BranchIR) GetConcatKnotForInputs(air *AbraIR, inputs []*SiteIR) *SiteIR {
	for _, site := range branch.bodySites {
		if site.siteType == SITE_KNOT_CONCAT && equalSites(inputs, site.inputs) {
			return site
		}
	}
	size := 0
	for _, s := range inputs {
		size += s.Size()
	}
	ret := &SiteIR{
		siteType: SITE_KNOT_CONCAT,
		inputs:   inputs,
		branch:   air.GetConcatBranchForSize(size),
	}
	branch.bodySites = append(branch.bodySites, ret)
	return ret
}

func (branch *BranchIR) Get1TritConstLutSite(air *AbraIR, val int8) *SiteIR {
	// find the only lut for 1 trit constant
	lutRepr := strings.Repeat(qupla.Get1TritConstLutRepr(val), 27)
	lutValConst := air.GetLut(lutRepr)

	for _, site := range branch.bodySites {
		if site.siteType == SITE_KNOT_LUT_1_INP && site.lut == lutValConst {
			return site
		}
	}
	ret := &SiteIR{
		idx:      len(branch.bodySites) + 1,
		siteType: SITE_KNOT_LUT_1_INP,
		inputs:   []*SiteIR{branch.bodySites[0]}, // #0 input site of the branch, always exist -> "any" input -> LUT output is constant anyway
		lut:      lutValConst,
	}
	branch.bodySites = append(branch.bodySites, ret)
	return ret
}

func (branch *BranchIR) GetTritConstSite(air *AbraIR, val Trits) *SiteIR {
	if len(val) == 1 {
		return branch.Get1TritConstLutSite(air, val[0])
	}
	inputs := make([]*SiteIR, len(val))
	for i, valtrit := range val {
		inputs[i] = branch.Get1TritConstLutSite(air, valtrit)
	}
	return branch.GetConcatKnotForInputs(air, inputs)
}

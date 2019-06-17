package abragen

// Interim representation for generating Abra tritcode form qupla
// Mainly needed to assign indices to Abra object

type AbraIR struct {
	luts     []*LutIR // string representation
	branches []*BranchIR
}

type LutIR struct {
	idx     int // block index within module
	strRepr string
}

type BranchIRType int

const (
	BRANCH_1ST_INPUT_TRIT = BranchIRType(0)
	BRANCH_QUPLA_CONCAT   = BranchIRType(1)
	BRANCH_QUPLA_SLICE    = BranchIRType(2)
)

type BranchIR struct {
	idx         int // block index within module
	inputs      []*SiteIR
	bodySites   []*SiteIR
	stateSites  []*SiteIR
	outputSites []*SiteIR
	// metadata for search
	branchType BranchIRType
	size       int // BRANCH_QUPLA_CONCAT, BRANCH_QUPLA_SLICE
	offset     int // BRANCH_QUPLA_SLICE
}

type SiteIRType int

const (
	SITE_INPUT          = SiteIRType(0)
	SITE_MERGE          = SiteIRType(1) // the rest are knot sites
	SITE_KNOT_LUT_1_INP = SiteIRType(2) //
	SITE_KNOT_CONCAT    = SiteIRType(3)
)

type SiteIR struct {
	idx      int        // index within branch
	siteType SiteIRType // metatag to search
	size     int        // used by SITE_INPUT only
	inputs   []*SiteIR
	lut      *LutIR    // if site is lut-knot
	branch   *BranchIR // if site branch-knot
}

// ret -> index
func (air *AbraIR) GetLut(strRepr string) *LutIR {
	for _, lut := range air.luts {
		if lut.strRepr == strRepr {
			return lut
		}
	}
	ret := &LutIR{
		strRepr: strRepr,
	}
	air.luts = append(air.luts, ret)
	return ret
}

func (air *AbraIR) GetConcatBranchForSize(size int) *BranchIR {
	for _, b := range air.branches {
		if b.branchType == BRANCH_QUPLA_CONCAT && b.size == size {
			return b
		}
	}
	inputs := []*SiteIR{{siteType: SITE_INPUT, size: size}}
	ret := &BranchIR{
		branchType: BRANCH_QUPLA_CONCAT,
		size:       size,
		inputs:     inputs,
		bodySites: []*SiteIR{{
			siteType: SITE_MERGE, // merges the only input site
			inputs:   inputs,
		}},
	}
	air.branches = append(air.branches, ret)
	return ret
}

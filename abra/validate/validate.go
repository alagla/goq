package validate

import (
	"fmt"
	. "github.com/lunfardo314/goq/abra"
)

func Validate(codeUnit *CodeUnit, assumeSizes bool) []error {
	ret := make([]error, 0, 10)
	for _, block := range codeUnit.Code.Blocks {
		if (assumeSizes && block.AssumedSize != block.Size) || block.Size == 0 {
			ret = append(ret, fmt.Errorf("AssumedSize (%d) != Size (%d) in block '%s'",
				block.AssumedSize, block.Size, block.LookupName))
		}
		var err error
		switch block.BlockType {
		case BLOCK_LUT:
			if block.Size != 1 {
				ret = append(ret, fmt.Errorf("LUT size != 1 in '%s'", block.LookupName))
			}
		case BLOCK_BRANCH:
			if err = ValidateBranch(block.Branch, block.LookupName); err != nil {
				ret = append(ret, fmt.Errorf("ValidateBranch for '%s': '%s'", block.LookupName, err))
			} else if err = ValidateBranchSizes(block.Branch, block.LookupName, assumeSizes); err != nil {
				ret = append(ret, fmt.Errorf("ValidateBranchSizes for '%s': '%s'", block.LookupName, err))
			}
		case BLOCK_EXTERNAL:
			panic("implement me")
		}
	}
	return ret
}

type BranchStats struct {
	NumSites      int
	NumInputs     int
	NumBodySites  int
	NumStateSites int
	NumOutputs    int
	NumKnots      int
	NumMerges     int
	InputSizes    []int
	InputSize     int
}

func GetStats(branch *Branch) *BranchStats {
	ret := &BranchStats{
		InputSizes: make([]int, 0, 5),
	}
	for _, s := range branch.AllSites {
		switch s.SiteType {
		case SITE_INPUT:
			ret.NumInputs++
			ret.InputSizes = append(ret.InputSizes, s.Size)
		case SITE_BODY:
			ret.NumBodySites++
		case SITE_STATE:
			ret.NumStateSites++
		case SITE_OUTPUT:
			ret.NumOutputs++
		}
		if s.IsKnot {
			ret.NumKnots++
		} else {
			ret.NumMerges++
		}
	}
	ret.NumSites = len(branch.AllSites)
	for _, s := range ret.InputSizes {
		ret.InputSize += s
	}
	return ret
}

func ValidateBranch(branch *Branch, lookupName string) error {
	if branch.NumInputs+branch.NumBodySites+branch.NumStateSites+branch.NumOutputs != branch.NumSites {
		return fmt.Errorf("something wrong with enumerating sites in branch '%s'", lookupName)
	}
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_OUTPUT && s.Merge == nil && s.Knot == nil {
			return fmt.Errorf("wrong output in branch '%s'", lookupName)
		}
	}
	for _, s := range branch.AllSites {
		err := ValidateSite(branch, s)
		if err != nil {
			return fmt.Errorf("ValidateSite for %s: '%v'", lookupName, err)
		}
	}
	return nil
}

func ValidateBranchSizes(branch *Branch, lookupName string, assumeSizes bool) error {
	if branch.Size == 0 || GetBranchInputSize(branch) == 0 {
		return fmt.Errorf("wrong branch size %d in '%s'", branch.Size, lookupName)
	}
	for _, s := range branch.AllSites {

		if s.Size == 0 || (assumeSizes && s.Size != s.AssumedSize) {
			return fmt.Errorf("site.AssumedSize (%d) != site.Size (%d) in site '%s' of block '%s'",
				s.AssumedSize, s.Size, s.LookupName, lookupName)
		}
		if s.SiteType != SITE_INPUT && s.Knot == nil && s.Merge == nil {
			return fmt.Errorf("inconsistent site '%s' in branch '%s'", s.LookupName, lookupName)
		}
		if s.SiteType != SITE_INPUT && s.IsKnot {
			switch {
			case s.Knot.Block.BlockType == BLOCK_LUT:
				if 3 != GetKnotInputSize(s.Knot) {
					return fmt.Errorf("sum of sizes of the inputs (%d) != branch size (%d) in knot '%s' of branch '%s'",
						GetKnotInputSize(s.Knot), GetBranchInputSize(s.Knot.Block.Branch), s.LookupName, lookupName)
				}
			case s.Knot.Block.BlockType == BLOCK_BRANCH:
				if s.Knot.Block.Branch.Size != GetKnotInputSize(s.Knot) {
					return fmt.Errorf("sum of sizes of the inputs (%d) != branch out size (%d) in knot '%s' of branch '%s'",
						GetKnotInputSize(s.Knot), s.Knot.Block.Size, s.LookupName, lookupName)
				}
			case s.Knot.Block.BlockType == BLOCK_EXTERNAL:
				panic("implement me")
			}
		}
	}
	return nil
}

func ValidateSite(branch *Branch, site *Site) error {
	if site.SiteType != SITE_INPUT {
		if (site.Merge == nil) == (site.Knot == nil) {
			return fmt.Errorf("invalid site '(site.Merge == nil) == (site.Knot == nil)'")
		}
		if site.IsKnot {
			if len(site.Knot.Sites) == 0 {
				return fmt.Errorf("invalid body knot site 'len(site.Knot.Sites) == 0'")
			}
		} else {
			if len(site.Merge.Sites) == 0 {
				return fmt.Errorf("invalid body merge site 'len(site.Merge.Sites) == 0'")
			}
		}
	}
	if site.CalculatedSize && site.Size <= 0 {
		return fmt.Errorf("invalid input site 'site.Size <= 0'")
	}
	return checkSiteIndices(branch, site)
}

func checkSiteIndices(branch *Branch, site *Site) error {
	if site.SiteType == SITE_INPUT {
		if site.Index < 0 || site.Index >= branch.NumInputs {
			return fmt.Errorf("invalid input site: 'site.Index < 0 || site.Index >= branch.NumInputs'")
		}
		return nil
	}
	if site.Index < 0 {
		return fmt.Errorf("invalid site: 'site.Index < 0'")
	}
	if site.Index < 0 || site.Index >= branch.NumSites {
		return fmt.Errorf("invalid site: 'site.Index < 0 || site.Index >= branch.NumSites'")
	}
	if site.IsKnot {
		for _, s := range site.Knot.Sites {
			if site.SiteType != SITE_STATE && s.SiteType != SITE_STATE && s.Index >= site.Index {
				return fmt.Errorf("invalid knot site '%s': 's.Index >= site.Index'", site.LookupName)
			}
		}
	} else {
		for _, s := range site.Merge.Sites {
			if site.SiteType != SITE_STATE && s.SiteType != SITE_STATE && s.Index >= site.Index {
				return fmt.Errorf("invalid merge site '%s': 's.Index >= site.Index'", site.LookupName)
			}
		}
	}
	return nil
}

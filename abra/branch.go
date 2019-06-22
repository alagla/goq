package abra

import "fmt"

func (codeUnit *CodeUnit) AddNewBranchBlock(lookupName string, size int) *Block {
	retbranch := &Branch{
		InputSites:  make([]*Site, 0, 10),
		BodySites:   make([]*Site, 0, 10),
		OutputSites: make([]*Site, 0, 10),
		StateSites:  make([]*Site, 0, 10),
		AllSites:    make([]*Site, 0, 10),
		Size:        size,
	}
	ret := retbranch.NewBlock(lookupName)
	if codeUnit.AddNewBlock(ret) {
		return ret
	}
	panic(fmt.Errorf("branch block '%s' already exists", lookupName))
}

func (branch *Branch) NewBlock(lookupName string) *Block {
	return &Block{
		BlockType:  BLOCK_BRANCH,
		Branch:     branch,
		LookupName: lookupName,
	}
}

func (branch *Branch) FindSite(lookupName string) *Site {
	if lookupName == "" {
		return nil
	}
	for _, site := range branch.AllSites {
		if site.LookupName != "" && site.LookupName == lookupName {
			return site
		}
	}
	return nil
}

func (branch *Branch) AddInputSite(size int) *Site {
	ret := &Site{
		SiteType: SITE_INPUT,
		Size:     size,
	}
	branch.InputSites = append(branch.InputSites, ret)
	return ret
}

// if find site with same lookup name, updates its isKnot, Knot and merge field with new
// returns found site.
// this is needed for generation of state sites in two steps
// therefore all site lookup names must be unique (if not "")

func (branch *Branch) GenOrUpdateSite(site *Site) *Site {
	ret := branch.FindSite(site.LookupName)
	if ret != nil {
		ret.IsKnot = site.IsKnot
		ret.Knot = site.Knot
		ret.Merge = site.Merge
		return ret
	}
	branch.AllSites = append(branch.BodySites, site)
	return site
}

func (branch *Branch) AddUnfinishedStateSite(lookupName string) *Site {
	ret := &Site{
		LookupName: lookupName,
		SiteType:   SITE_STATE,
	}
	return branch.GenOrUpdateSite(ret)
}

type BranchStats struct {
	NumSites      int
	NumInputs     int
	NumBodySites  int
	NumStateSites int
	NumOutputs    int
	NumKnots      int
	NumMerges     int
}

func (branch *Branch) GetStats() *BranchStats {
	ret := &BranchStats{}
	for _, s := range branch.AllSites {
		switch s.SiteType {
		case SITE_INPUT:
			ret.NumInputs++
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
	return ret
}

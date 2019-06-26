package abra

import (
	"fmt"
	"strings"
)

func (codeUnit *CodeUnit) AddNewBranchBlock(lookupName string, assumedSize int) *Block {
	retbranch := &Branch{
		inputSites:  make([]*Site, 0, 10),
		bodySites:   make([]*Site, 0, 10),
		outputSites: make([]*Site, 0, 10),
		stateSites:  make([]*Site, 0, 10),
		AllSites:    make([]*Site, 0, 10),
		AssumedSize: assumedSize,
	}
	ret := retbranch.NewBlock(lookupName, assumedSize)
	if codeUnit.AddNewBlock(ret) {
		return ret
	}
	panic(fmt.Errorf("branch block '%s' already exists", lookupName))
}

func (branch *Branch) NewBlock(lookupName string, assumedSize int) *Block {
	return &Block{
		BlockType:   BLOCK_BRANCH,
		Branch:      branch,
		LookupName:  lookupName,
		AssumedSize: assumedSize,
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
	branch.AllSites = append(branch.AllSites, ret)
	return ret
}

func (branch *Branch) GetInputSite(idx int) *Site {
	counter := 0
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_INPUT {
			if counter == idx {
				return s
			}
			counter++
		}
	}
	panic("input site index out of bound")
}

// if find site with same lookup name, updates its isKnot, Knot and merge field with new
// returns found site.
// this is needed for generation of state sites in two steps
// therefore all site lookup names must be unique (if not "")

func (branch *Branch) AddOrUpdateSite(site *Site) *Site {
	ret := branch.FindSite(site.LookupName)
	if ret != nil {
		ret.IsKnot = site.IsKnot
		ret.Knot = site.Knot
		ret.Merge = site.Merge
		return ret
	}
	branch.AllSites = append(branch.AllSites, site)
	return site
}

func (branch *Branch) AddUnfinishedStateSite(lookupName string) *Site {
	ret := &Site{
		LookupName: lookupName,
		SiteType:   SITE_STATE,
	}
	return branch.AddOrUpdateSite(ret)
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

func (branch *Branch) GetStats() *BranchStats {
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

func (branch *Branch) GetSize() (int, error) {
	if branch.Size < 0 {
		return 0, RecursionDetectedError
	}
	branch.Size = -1
	ret := 0
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_OUTPUT {
			sz, err := s.GetSize()
			if err != nil {
				return 0, err
			}
			ret += sz
		}
	}
	branch.Size = ret
	if branch.Size != branch.AssumedSize {
		return 0, fmt.Errorf("branch: Size %d != AssumedSize %d", branch.Size, branch.AssumedSize)
	}
	return ret, nil
}

func (branch *Branch) GetInputSize() int {
	ret := 0
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_INPUT {
			ret += s.Size
		}
	}
	return ret
}

func (branch *Branch) AssertValid() {
	// assert
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_OUTPUT && s.Merge == nil && s.Knot == nil {
			panic("assert: wrong output")
		}
	}
	for _, s := range branch.AllSites {
		s.AssertValidSite()
	}
}

func (block *Block) GetSize() (int, error) {
	if strings.Contains(block.LookupName, "arcRadixLeaf_243_8019") {
		fmt.Printf("kuku")
	}
	switch block.BlockType {
	case BLOCK_LUT:
		return 1, nil
	case BLOCK_BRANCH:
		return block.Branch.GetSize()
	case BLOCK_EXTERNAL:
		panic("implement me")
	}
	panic("wrong block type")
}

func (block *Block) GetInputSize() int {
	switch block.BlockType {
	case BLOCK_LUT:
		return 3
	case BLOCK_BRANCH:
		return block.Branch.GetInputSize()
	}
	panic("implement me")
}

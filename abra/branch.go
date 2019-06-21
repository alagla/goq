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
	if codeUnit.addBlock(ret) {
		return ret
	}
	return nil
}

func (branch *Branch) NewBlock(lookupName string) *Block {
	return &Block{
		BlockType:  BLOCK_BRANCH,
		Branch:     branch,
		LookupName: lookupName,
	}
}

func (branch *Branch) FindSite(lookupName string) *Site {
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

func (branch *Branch) AddNewSite(site *Site, lookupName string) {
	ret := branch.FindSite(site.LookupName)
	if ret != nil {
		panic(fmt.Errorf("duplicate site lookup name '%s'", lookupName))
	}
	branch.AllSites = append(branch.BodySites, site)
}

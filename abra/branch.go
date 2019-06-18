package abra

func (codeUnit *CodeUnit) AddNewBranchBlock(lookupName string, size int) *Block {
	retbranch := &Branch{
		InputSites:  make([]*Site, 0, 10),
		BodySites:   make([]*Site, 0, 10),
		OutputSites: make([]*Site, 0, 10),
		StateSites:  make([]*Site, 0, 10),
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

func (branch *Branch) FindBodySite(lookupName string) *Site {
	for _, site := range branch.BodySites {
		if site.LookupName == lookupName {
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

func (branch *Branch) AddBodySite(site *Site) (*Site, bool) {
	for _, bs := range branch.BodySites {
		if site.LookupName == "" && site.LookupName == bs.LookupName {
			return bs, false
		}
	}
	branch.BodySites = append(branch.BodySites, site)
	return site, true
}

func (branch *Branch) AddStateSite(site *Site) (*Site, bool) {
	if site.LookupName == "" {
		panic("state sites have names!")
	}
	for _, bs := range branch.StateSites {
		if site.LookupName == bs.LookupName {
			return bs, false
		}
	}
	branch.StateSites = append(branch.StateSites, site)
	return site, true
}

func (branch *Branch) AddOutputSite(site *Site) bool {
	branch.OutputSites = append(branch.OutputSites, site)
	return true
}

package abra

func (site *Site) SetLookupName(ln string) *Site {
	site.LookupName = ln
	return site
}

func NewMerge(sites ...*Site) *Merge {
	return &Merge{
		Sites: sites,
	}
}

func NewKnot(block *Block, sites ...*Site) *Knot {
	return &Knot{
		Sites: sites,
		Block: block,
	}
}

func (merge *Merge) NewSite(lookupName string) *Site {
	return &Site{
		SiteType:   SITE_MERGE,
		Merge:      merge,
		LookupName: lookupName,
	}
}

func (knot *Knot) NewSite(lookupName string) *Site {
	return &Site{
		SiteType:   SITE_KNOT,
		Knot:       knot,
		LookupName: lookupName,
	}
}

func (branch *Branch) AddKnotSiteForInputs(knotBlock *Block, lookupName string, inputs ...*Site) *Site {
	ret := NewKnot(knotBlock, inputs...).NewSite(lookupName)
	branch.AddBodySite(ret)
	return ret
}

package abra

func (site *Site) SetLookupName(ln string) *Site {
	site.LookupName = ln
	return site
}

func (site *Site) SetType(t SiteType) *Site {
	site.SiteType = t
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

func (merge *Merge) NewSite() *Site {
	return &Site{
		IsKnot:   false,
		Merge:    merge,
		SiteType: SITE_BODY,
	}
}

func (knot *Knot) NewSite() *Site {
	return &Site{
		IsKnot:   true,
		Knot:     knot,
		SiteType: SITE_BODY,
	}
}

package abra

func (site *Site) SetLookupName(ln string) *Site {
	site.LookupName = ln
	return site
}

func (site *Site) ChangeType(t SiteType) *Site {
	if site.SiteType != SITE_BODY {
		panic("only type of the body site can be changed")
	}
	site.SiteType = t
	return site
}

func (site *Site) AssertValidSite() {
	switch site.SiteType {
	case SITE_INPUT:
		if site.Knot != nil {
			panic("invalid site 1")
		}
		if site.Merge != nil {
			panic("invalid site 2")
		}
	case SITE_BODY:
		if (site.Merge == nil) == (site.Knot == nil) {
			panic("invalid site 3")
		}
		if site.IsKnot {
			if len(site.Knot.Sites) == 0 {
				panic("invalid site 4")
			}
		} else {
			if len(site.Merge.Sites) == 0 {
				panic("invalid site 5")
			}
		}
	case SITE_OUTPUT:
		if (site.Merge == nil) == (site.Knot == nil) {
			panic("invalid site 6")
		}
		if site.IsKnot {
			if len(site.Knot.Sites) == 0 {
				panic("invalid site 7")
			}
		} else {
			if len(site.Merge.Sites) == 0 {
				panic("invalid site 8")
			}
		}
	case SITE_STATE:
		if (site.Merge == nil) == (site.Knot == nil) {
			panic("invalid site 9")
		}
		if site.IsKnot {
			if len(site.Knot.Sites) == 0 {
				panic("invalid site 10")
			}
		} else {
			if len(site.Merge.Sites) == 0 {
				panic("invalid site 11")
			}
		}
	}
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

func (merge *Merge) NewSite(assumedSize int) *Site {
	return &Site{
		IsKnot:      false,
		Merge:       merge,
		SiteType:    SITE_BODY,
		AssumedSize: assumedSize,
	}
}

func (knot *Knot) NewSite(assumedSize int) *Site {
	return &Site{
		IsKnot:      true,
		Knot:        knot,
		SiteType:    SITE_BODY,
		AssumedSize: assumedSize,
	}
}

func (site *Site) CalcSize() {
	if site.SiteType == SITE_INPUT {
		return
	}
	if site.IsKnot {
		site.Size = site.Knot.CalcSize()
	} else {
		site.Size = site.Merge.CalcSize()
	}
}

func (knot *Knot) CalcSize() int {
	return knot.Block.Size
}

func (merge *Merge) CalcSize() int {
	ret := 0
	for _, s := range merge.Sites {
		if s.Size != 0 {
			ret = s.Size
		}
	}
	return ret
}

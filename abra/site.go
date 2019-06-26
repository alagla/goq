package abra

import "fmt"

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

func (knot *Knot) Size() (int, error) {
	bsize, err := knot.Block.GetSize()
	if err != nil {
		return 0, err
	}
	argsz := 0
	var sz int
	for _, s := range knot.Sites {
		sz, err = s.GetSize()
		if err != nil {
			return 0, err
		}
		argsz += sz
	}
	isize := knot.Block.GetInputSize()
	if isize != argsz {
		return 0, fmt.Errorf("mismatch of block input size %d and args size %d in the knot", isize, argsz)
	}
	return bsize, nil
}

func (merge *Merge) Size() (int, error) {
	var lastsz, sz int
	var err error
	for _, s := range merge.Sites {
		sz, err = s.GetSize()
		if err == RecursionRetected {
			// recursion will be resolved with another merge patch
			continue
		}
		if err != nil {
			return 0, err
		}
		if lastsz != 0 && lastsz != sz {
			// not 100% correct
			return 0, fmt.Errorf("inputs of merge must be same size")
		}
		lastsz = sz
	}
	return lastsz, nil
}

// special way to determine size of state site
func (site *Site) GetSize() (int, error) {
	switch site.SiteType {
	case SITE_INPUT:
		return site.Size, nil
	case SITE_STATE:
		if site.Size < 0 {
			return 0, RecursionRetected
		}
	}
	site.Size = -1
	var err error
	if site.IsKnot {
		site.Size, err = site.Knot.Size()
	} else {
		site.Size, err = site.Merge.Size()
	}
	if site.Size != site.AssumedSize {
		return 0, fmt.Errorf("site Size %d != site AssumedSize %d", site.Size, site.AssumedSize)
	}
	return site.Size, err
}

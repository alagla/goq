package abra

import "fmt"

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

// special type of error to resolve cycles in state sites
type ContainsState struct{}

func (e *ContainsState) Error() string {
	return "Contains state"
}

var ErrorContainsState = &ContainsState{}

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
	if bsize != argsz {
		return 0, fmt.Errorf("mismatch of block size %d and args size %d in the knot", bsize, argsz)
	}
	return bsize, nil
}

func (merge *Merge) Size() (int, error) {
	var lastsz, sz int
	var err error
	for _, s := range merge.Sites {
		sz, err = s.GetSize()
		if err == ErrorContainsState {
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
	return 0, fmt.Errorf("undetermined size of merge")
}

// special way to determine size of state site
func (site *Site) GetSize() (int, error) {
	if site.SiteType == SITE_INPUT {
		return site.Size, nil
	}
	if site.SiteType == SITE_OUTPUT {
		if site.Size < 0 {
			return 0, ErrorContainsState
		} else if site.Size > 0 {
			return site.Size, nil
		}
		// site.Size == 0 -> first time evaluated
		site.Size = -1
	}
	var err error
	if site.IsKnot {
		site.Size, err = site.Knot.Size()
	} else {
		site.Size, err = site.Merge.Size()
	}
	return site.Size, err
}

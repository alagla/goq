package construct

import (
	"fmt"
	. "github.com/lunfardo314/goq/abra"
)

const TRITCODE_VERSION = 0

func NewCodeUnit() *CodeUnit {
	return &CodeUnit{
		EntityAttachments: make([]*EntityAttachment, 0, 10),
		Code: &Code{
			TritcodeVersion: TRITCODE_VERSION,
			Blocks:          make([]*Block, 0, 100),
		},
	}
}

func addBlock(codeUnit *CodeUnit, block *Block) {
	codeUnit.Code.Blocks = append(codeUnit.Code.Blocks, block)
}

func MustAddNewLUTBlock(codeUnit *CodeUnit, strRepr string, name string) *Block {
	if FindLUTBlock(codeUnit, strRepr) != nil {
		panic(fmt.Errorf("repeating LUT block '%s'", strRepr))
	}
	block := &Block{
		BlockType: BLOCK_LUT,
		LUT: &LUT{
			Binary: BinaryEncodedLUTFromString(strRepr),
			Name:   name,
		},
		LookupName:  strRepr,
		AssumedSize: 1,
	}
	addBlock(codeUnit, block)
	return block
}

func FindLUTBlock(codeUnit *CodeUnit, lookupName string) *Block {
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == BLOCK_LUT && b.LookupName == lookupName {
			return b
		}
	}
	return nil
}

func FindBranchBlock(codeUnit *CodeUnit, lookupName string) *Block {
	if lookupName == "" {
		return nil
	}
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == BLOCK_BRANCH && b.LookupName != "" && b.LookupName == lookupName {
			return b
		}
	}
	return nil
}

func MustAddNewBranchBlock(codeUnit *CodeUnit, lookupName string, assumedSize int) *Block {
	if FindBranchBlock(codeUnit, lookupName) != nil {
		panic(fmt.Errorf("repeating branch block with lookupName = '%s'", lookupName))
	}
	retbranch := &Branch{
		InputSites:  make([]*Site, 0, 10),
		BodySites:   make([]*Site, 0, 10),
		OutputSites: make([]*Site, 0, 10),
		StateSites:  make([]*Site, 0, 10),
		AllSites:    make([]*Site, 0, 10),
		AssumedSize: assumedSize,
	}
	ret := &Block{
		BlockType:   BLOCK_BRANCH,
		Branch:      retbranch,
		LookupName:  lookupName,
		AssumedSize: assumedSize,
	}

	addBlock(codeUnit, ret)
	return ret
}

func FindSite(branch *Branch, lookupName string) *Site {
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

func AddInputSite(branch *Branch, size int) *Site {
	ret := &Site{
		SiteType:    SITE_INPUT,
		Size:        size,
		AssumedSize: size,
	}
	branch.AllSites = append(branch.AllSites, ret)
	return ret
}

func GetInputSite(branch *Branch, idx int) *Site {
	counter := 0
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_INPUT {
			if counter == idx {
				return s
			}
			counter++
		}
	}
	panic("input site index out of bounds")
}

func MustAddNewSite(branch *Branch, site *Site) {
	if FindSite(branch, site.LookupName) != nil {
		panic(fmt.Errorf("repeated site '%s'", site.LookupName))
	}
	branch.AllSites = append(branch.AllSites, site)
}

// if find site with same lookup name, updates its isKnot, Knot and Merge fields with new
// returns found site.
// this is needed for generation of state sites in two steps
// therefore all site lookup names must be unique (if not "")

func AddOrUpdateSite(branch *Branch, site *Site) *Site {
	ret := FindSite(branch, site.LookupName)
	if ret == nil {
		branch.AllSites = append(branch.AllSites, site)
		return site
	}
	ret.IsKnot = site.IsKnot
	ret.Knot = site.Knot
	ret.Merge = site.Merge
	return ret
}

func MustAddUnfinishedStateSite(branch *Branch, lookupName string, assumedSize int) *Site {
	if FindSite(branch, lookupName) != nil {
		panic(fmt.Errorf("repeated state site '%s'", lookupName))
	}
	ret := &Site{
		LookupName:  lookupName,
		SiteType:    SITE_STATE,
		AssumedSize: assumedSize,
	}
	branch.AllSites = append(branch.AllSites, ret)
	return ret
}

func ChangeSiteType(site *Site, t SiteType) *Site {
	if site.SiteType != SITE_BODY {
		panic("only type of the body site can be changed")
	}
	site.SiteType = t
	return site
}

func NewMergeSite(assumedSize int, lookupName string, sites ...*Site) *Site {
	return &Site{
		IsKnot:      false,
		Merge:       &Merge{Sites: sites},
		SiteType:    SITE_BODY,
		LookupName:  lookupName,
		AssumedSize: assumedSize,
	}
}

func NewKnotSite(assumedSize int, lookupName string, block *Block, sites ...*Site) *Site {
	return &Site{
		IsKnot:      true,
		Knot:        &Knot{Sites: sites, Block: block},
		SiteType:    SITE_BODY,
		LookupName:  lookupName,
		AssumedSize: assumedSize,
	}
}

func NewExternalBlock(external *ExternalBlock) *Block {
	return &Block{
		BlockType:     BLOCK_EXTERNAL,
		ExternalBlock: external,
	}
}

package abra

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/utils"
)

// finds or creates branch for concatenation of inputs with size

func (branch *Branch) FindSite(lookupName string) *Site {
	for _, site := range branch.BodySites {
		if site.LookupName == lookupName {
			return site
		}
	}
	return nil
}

func (branch *Branch) GenKnotSiteForInputs(knotBlock *Block, codeUnit *CodeUnit, lookupName string, inputs ...*Site) *Site {
	return NewKnot(knotBlock, inputs...).NewSite(lookupName)
}

func (branch *Branch) Get1TritConstLutSite(codeUnit *CodeUnit, val int8) *Site {
	// first try to find if there's constant lut site for val
	// if not, create one
	lookupName := fmt.Sprintf("1trit_const_site_%s", qupla.TritName(val))
	ret := branch.FindSite(lookupName)
	if ret != nil {
		return ret
	}
	// didn't find. Need to create one
	// first find or create the only lut for 1 trit constant
	lutRepr := qupla.Get1TritConstLutRepr(val)
	lutValConstBlock := codeUnit.FindLUTBlock(lutRepr)
	if lutValConstBlock == nil {
		lut := qupla.BinaryEncodedLUTFromString(lutRepr)
		lutValConstBlock = codeUnit.NewLUTBlock(lutRepr, lut)
	}
	// now create site in the branch
	// it will always generate constant trit
	// the input for lut is 3 repeated 1-trit sites from lsb of the branches input
	any := branch.GetAnyTritSite(codeUnit)
	ret = branch.GenKnotSiteForInputs(lutValConstBlock, codeUnit, lookupName, any, any, any)
	return ret
}

func (branch *Branch) GetAnyTritSite(codeUnit *CodeUnit) *Site {
	lookupName := "any_input_site" // each branch will have site with this name
	ret := branch.FindSite(lookupName)
	if ret != nil {
		return ret
	}
	lsbBlock := codeUnit.GetLsbSliceBlock()
	ret = NewKnot(lsbBlock, branch.InputSites[0]).NewSite(lookupName + "_knot")
	branch.AddBodySite(ret)
	return ret
}

func (branch *Branch) GetTritConstSite(codeUnit *CodeUnit, val Trits) *Site {
	lookupName := utils.TritsToString(val) + "_const_site"
	ret := branch.FindSite(lookupName)
	if ret != nil {
		return ret
	}
	inputs := make([]*Site, len(val))
	for i, trit := range val {
		inputs[i] = branch.Get1TritConstLutSite(codeUnit, trit)
	}

	concatBlock := codeUnit.GetConcatBlockForSize(len(val))
	ret = branch.GenKnotSiteForInputs(concatBlock, codeUnit, lookupName, inputs...)
	return ret
}

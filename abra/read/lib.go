package read

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/abra/construct"
)

type tritReader struct {
	trits  Trits
	curPos int
}

// return trit, false or -100, true (eof)
func (tr *tritReader) readTrit() (int8, bool) {
	if tr.curPos >= len(tr.trits) {
		return -100, true
	}
	ret := tr.trits[tr.curPos]
	tr.curPos++
	return ret, false
}

func readNTrits(tReader *tritReader, n int) (Trits, error) {
	ret := make(Trits, n)
	var eof bool
	for i := 0; i < n; i++ {
		ret[i], eof = tReader.readTrit()
		if eof {
			return nil, fmt.Errorf("unexpected EOF at pos %d", tReader.curPos)
		}
	}
	return ret, nil
}

func ParseTritcode(trits Trits) (*abra.CodeUnit, error) {
	ret := construct.NewCodeUnit()
	tReader := &tritReader{trits: trits}

	err := ParseCode(tReader, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func ParseCode(tReader *tritReader, codeUnit *abra.CodeUnit) error {
	var err error
	var ver int
	ver, err = ParsePosInt(tReader)
	if err != nil {
		return err
	}
	if ver != codeUnit.Code.TritcodeVersion {
		return fmt.Errorf("expected tritcode version %d, got %d", codeUnit.Code.TritcodeVersion, ver)
	}
	// read LUTs
	codeUnit.Code.NumLUTs, err = ParsePosInt(tReader)
	if err != nil {
		return err
	}
	for i := 0; i < codeUnit.Code.NumLUTs; i++ {
		err = ParseLUTBlock(tReader, codeUnit)
		if err != nil {
			return err
		}
	}
	// read Branche blocks
	codeUnit.Code.NumBranches, err = ParsePosInt(tReader)
	if err != nil {
		return err
	}
	branchDefs := make([]*branchSiteDefinitions, 0, codeUnit.Code.NumBranches)
	var bd *branchSiteDefinitions
	for i := 0; i < codeUnit.Code.NumBranches; i++ {
		bd, err = ParseBranchBlock(tReader, codeUnit)
		if err != nil {
			return err
		}
		branchDefs = append(branchDefs, bd)
	}
	// read External blocks
	codeUnit.Code.NumExternalBlocks, err = ParsePosInt(tReader)
	if err != nil {
		return err
	}
	for i := 0; i < codeUnit.Code.NumExternalBlocks; i++ {
		err = ParseExternalBlock(tReader, codeUnit)
		if err != nil {
			return err
		}
	}
	err = finalizeBranches(codeUnit, branchDefs)
	return err
}

func ParseLUTBlock(tReader *tritReader, codeUnit *abra.CodeUnit) error {
	n, err := ParsePosInt(tReader)
	if err != nil {
		return err
	}
	if n != 35 {
		return fmt.Errorf("expected PosInt == 35 at position %d", tReader.curPos)
	}
	var trits Trits
	trits, err = readNTrits(tReader, 35)
	if err != nil {
		return err
	}
	strRepr := abra.StringFromBinaryEncodedLUT(uint64(TritsToInt(trits)))
	_, err = construct.AddNewLUTBlock(codeUnit, strRepr, "")
	if err != nil {
		return err
	}
	return nil
}

func ParseBranchBlock(tReader *tritReader, codeUnit *abra.CodeUnit) (*branchSiteDefinitions, error) {
	blen, err := ParsePosInt(tReader)
	if err != nil {
		return nil, err
	}
	var trits Trits
	trits, err = readNTrits(tReader, blen)
	if err != nil {
		return nil, err
	}
	bReader := &tritReader{trits: trits}
	var numInputs, numBodySites, numOutputSites, numStateSites int

	numInputs, err = ParsePosInt(bReader)
	if err != nil {
		return nil, err
	}
	inputLengths := make([]int, numInputs)
	for i := 0; i < numInputs; i++ {
		inputLengths[i], err = ParsePosInt(bReader)
		if err != nil {
			return nil, err
		}
	}
	numBodySites, err = ParsePosInt(bReader)
	if err != nil {
		return nil, err
	}
	numOutputSites, err = ParsePosInt(bReader)
	if err != nil {
		return nil, err
	}
	numStateSites, err = ParsePosInt(bReader)
	if err != nil {
		return nil, err
	}
	block := construct.AddNewBranchBlock(codeUnit, numInputs, numBodySites, numOutputSites, numStateSites)
	ret := &branchSiteDefinitions{
		inputLengths:   inputLengths,
		bodySiteDefs:   make([]*siteDefinition, numBodySites),
		outputSiteDefs: make([]*siteDefinition, numOutputSites),
		stateSiteDefs:  make([]*siteDefinition, numStateSites),
	}

	//for i := 0; i < numInputs; i++ {
	//	block.Branch.AllSites[i] = construct.NewInputSite(inputLengths[i], i)
	//}

	// create unfinished sites
	for i := 0; i < numBodySites; i++ {
		ret.bodySiteDefs[i], err = ReadSiteDefinition(bReader)
		if err != nil {
			return nil, err
		}
		err = createUnfinishedSite(block.Branch, ret.bodySiteDefs[i], numInputs+i, abra.SITE_BODY)
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < numOutputSites; i++ {
		ret.outputSiteDefs[i], err = ReadSiteDefinition(bReader)
		if err != nil {
			return nil, err
		}
		err = createUnfinishedSite(block.Branch, ret.outputSiteDefs[i], numInputs+numBodySites+i, abra.SITE_OUTPUT)
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < numStateSites; i++ {
		ret.stateSiteDefs[i], err = ReadSiteDefinition(bReader)
		if err != nil {
			return nil, err
		}
		err = createUnfinishedSite(block.Branch, ret.stateSiteDefs[i], numInputs+numBodySites+numOutputSites+i, abra.SITE_STATE)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func finalizeBranches(codeUnit *abra.CodeUnit, branchDefs []*branchSiteDefinitions) error {
	if codeUnit.Code.NumBranches != len(branchDefs) {
		return fmt.Errorf("inconsistency with number of branches")
	}
	var err error
	for i := 0; i < codeUnit.Code.NumBranches; i++ {
		err = finalizeBranch(codeUnit, codeUnit.Code.Blocks[codeUnit.Code.NumLUTs+i].Branch, branchDefs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func finalizeBranch(codeUnit *abra.CodeUnit, branch *abra.Branch, branchDef *branchSiteDefinitions) error {
	if branch.NumInputs != len(branchDef.inputLengths) {
		return fmt.Errorf("inconsistency with inputLengths")
	}
	for i := 0; i < branch.NumInputs; i++ {
		branch.AllSites[i] = construct.NewInputSite(branchDef.inputLengths[i], i)
	}
	var err error
	offs := branch.NumInputs
	for i := 0; i < branch.NumBodySites; i++ {
		err = finalizeSite(codeUnit, branch, branch.AllSites[offs+i], branchDef.bodySiteDefs[i])
		if err != nil {
			return err
		}
	}
	offs += branch.NumBodySites
	for i := 0; i < branch.NumOutputs; i++ {
		err = finalizeSite(codeUnit, branch, branch.AllSites[offs+i], branchDef.outputSiteDefs[i])
		if err != nil {
			return err
		}
	}
	offs += branch.NumOutputs
	for i := 0; i < branch.NumStateSites; i++ {
		err = finalizeSite(codeUnit, branch, branch.AllSites[offs+i], branchDef.stateSiteDefs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func finalizeSite(codeUnit *abra.CodeUnit, branch *abra.Branch, site *abra.Site, siteDef *siteDefinition) error {
	if site.IsKnot != siteDef.isKnot {
		return fmt.Errorf("inconsistency: site.IsKnot != siteDef.isKnot")
	}
	if site.IsKnot {
		if siteDef.blockIndex < 0 || siteDef.blockIndex >= len(codeUnit.Code.Blocks) {
			return fmt.Errorf("wrong block index in knot")
		}
		site.Knot.Block = codeUnit.Code.Blocks[siteDef.blockIndex]

		if len(site.Knot.Sites) != len(siteDef.inputSiteIndices) {
			return fmt.Errorf("inconsistency in inpout sites knot")
		}
		for i, idx := range siteDef.inputSiteIndices {
			if idx < 0 || idx >= branch.NumSites {
				return fmt.Errorf("site index out of range in knot")
			}
			if idx >= site.Index {
				if site.SiteType != abra.SITE_STATE && branch.AllSites[idx].SiteType != abra.SITE_STATE {
					return fmt.Errorf("wrong site index in knot: must be pointing backwards")
				}
			}
			site.Knot.Sites[i] = branch.AllSites[idx]
		}
	} else {
		if len(site.Merge.Sites) != len(siteDef.inputSiteIndices) {
			return fmt.Errorf("inconsistency in input sites in merge")
		}
		for i, idx := range siteDef.inputSiteIndices {
			if idx < 0 || idx >= branch.NumSites {
				return fmt.Errorf("site index out of range in merge")
			}
			if idx >= site.Index {
				if site.SiteType != abra.SITE_OUTPUT && branch.AllSites[idx].SiteType != abra.SITE_OUTPUT {
					return fmt.Errorf("wrong site index in merge: must be pointing backwards")
				}
			}
			site.Merge.Sites[i] = branch.AllSites[idx]
		}
	}
	return nil
}

// creates site with unresolved links to input sites and branch
func createUnfinishedSite(branch *abra.Branch, siteDef *siteDefinition, index int, siteType abra.SiteType) error {
	sites := make([]*abra.Site, len(siteDef.inputSiteIndices)) // empty array
	if siteDef.isKnot {
		branch.AllSites[index] = construct.NewKnotSite(0, "", nil, sites...)
	} else {
		branch.AllSites[index] = construct.NewMergeSite(0, "", sites...)
	}
	branch.AllSites[index].SiteType = siteType
	branch.AllSites[index].Index = index
	return nil
}

func ParseExternalBlock(tReader *tritReader, codeUnit *abra.CodeUnit) error {
	panic("implement me")
}

type siteDefinition struct {
	isKnot           bool
	inputSiteIndices []int
	blockIndex       int
}

type branchSiteDefinitions struct {
	inputLengths   []int
	bodySiteDefs   []*siteDefinition
	outputSiteDefs []*siteDefinition
	stateSiteDefs  []*siteDefinition
}

func ReadSiteDefinition(tReader *tritReader) (*siteDefinition, error) {
	ret := &siteDefinition{}
	siteType, err := readNTrits(tReader, 1)
	if siteType[0] != -1 && siteType[0] != 1 {
		return nil, fmt.Errorf("wrong site type")
	}
	ret.isKnot = siteType[0] == -1
	var numInputSites int
	numInputSites, err = ParsePosInt(tReader)
	if err != nil {
		return nil, err
	}
	ret.inputSiteIndices = make([]int, numInputSites)
	for i := 0; i < numInputSites; i++ {
		ret.inputSiteIndices[i], err = ParsePosInt(tReader)
		if err != nil {
			return nil, err
		}
	}
	if ret.isKnot {
		ret.blockIndex, err = ParsePosInt(tReader)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func ParsePosInt(tReader *tritReader) (int, error) {
	buf := make(Trits, 0, 31)
	exit := false
	for !exit {
		if len(buf) >= 31 {
			return -1, fmt.Errorf("ParsePosInt: wrong PosInt: longer that 31 bit at position %d", tReader.curPos-1)
		}
		t, eof := tReader.readTrit()
		switch {
		case eof:
			return -1, fmt.Errorf("ParsePosInt: unexpected EOF at position %d", tReader.curPos)
		case t == 0:
			exit = true
		case t == -1:
			buf = append(buf, 0)
		case t == 1:
			buf = append(buf, 1)
		default:
			return -1, fmt.Errorf("ParsePosInt: wrong trit at position %d", tReader.curPos-1)
		}
	}
	ret := 0
	for i := len(buf) - 1; i >= 0; i-- {
		ret <<= 1
		if buf[i] == 1 {
			ret |= 0x1
		}
	}
	return ret, nil
}

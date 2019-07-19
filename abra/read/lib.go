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
	for i := 0; i < codeUnit.Code.NumBranches; i++ {
		err = ParseBranchBlock(tReader, codeUnit)
		if err != nil {
			return err
		}
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

	return nil
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

func ParseBranchBlock(tReader *tritReader, codeUnit *abra.CodeUnit) error {
	blen, err := ParsePosInt(tReader)
	if err != nil {
		return err
	}
	var trits Trits
	trits, err = readNTrits(tReader, blen)
	if err != nil {
		return err
	}
	bReader := &tritReader{trits: trits}
	var numInputs, numBodySites, numOutputSites, numStateSites int

	numInputs, err = ParsePosInt(bReader)
	if err != nil {
		return err
	}
	inputLengths := make([]int, numInputs)
	for i := 0; i < numInputs; i++ {
		inputLengths[i], err = ParsePosInt(bReader)
		if err != nil {
			return err
		}
	}
	numBodySites, err = ParsePosInt(bReader)
	if err != nil {
		return err
	}
	numOutputSites, err = ParsePosInt(bReader)
	if err != nil {
		return err
	}
	numStateSites, err = ParsePosInt(bReader)
	if err != nil {
		return err
	}
	block := construct.AddNewBranchBlock(codeUnit, numInputs, numBodySites, numOutputSites, numStateSites)

	for i := 0; i < numInputs; i++ {
		block.Branch.AllSites[i] = construct.NewInputSite(inputLengths[i], i)
	}
	bodySiteDefs := make([]*siteDefinition, numBodySites)
	for i := 0; i < numBodySites; i++ {
		bodySiteDefs[i], err = ReadSiteDefinition(bReader)
		if err != nil {
			return err
		}
	}
	outputSiteDefs := make([]*siteDefinition, numOutputSites)
	for i := 0; i < numOutputSites; i++ {
		outputSiteDefs[i], err = ReadSiteDefinition(bReader)
		if err != nil {
			return err
		}
	}
	stateSiteDefs := make([]*siteDefinition, numStateSites)
	for i := 0; i < numStateSites; i++ {
		stateSiteDefs[i], err = ReadSiteDefinition(bReader)
		if err != nil {
			return err
		}
	}
	err = ParseSites(block.Branch, bodySiteDefs, abra.SITE_BODY)
	if err != nil {
		return err
	}
	err = ParseSites(block.Branch, outputSiteDefs, abra.SITE_OUTPUT)
	if err != nil {
		return err
	}
	err = ParseSites(block.Branch, stateSiteDefs, abra.SITE_STATE)
	if err != nil {
		return err
	}
	return nil
}

func ParseSites(branch *abra.Branch, siteDefs []*siteDefinition, siteType abra.SiteType) error {
	var err error
	switch siteType {
	case abra.SITE_BODY:
		c := 0
		for i := 0; i < branch.NumBodySites; i++ {
			branch.AllSites[branch.NumInputs+i], err = ParseSite(siteDefs[c], abra.SITE_BODY)
			if err != nil {
				return err
			}
		}
	case abra.SITE_OUTPUT:
	case abra.SITE_STATE:
	default:
		return fmt.Errorf("internal: wrong site type")
	}
	return nil
}

func ParseSite(siteDef *siteDefinition, siteType abra.SiteType) (*abra.Site, error) {
	return nil, nil
}

func ParseExternalBlock(tReader *tritReader, codeUnit *abra.CodeUnit) error {
	panic("implement me")
}

type siteDefinition struct {
	isKnot           bool
	numInputSites    int
	inputSiteIndices []int
	blockIndex       int
}

func ReadSiteDefinition(tReader *tritReader) (*siteDefinition, error) {
	ret := &siteDefinition{}
	siteType, err := readNTrits(tReader, 1)
	if siteType[0] != -1 && siteType[0] != 1 {
		return nil, fmt.Errorf("wrong site type")
	}
	ret.isKnot = siteType[0] == -1
	ret.numInputSites, err = ParsePosInt(tReader)
	if err != nil {
		return nil, err
	}
	ret.inputSiteIndices = make([]int, ret.numInputSites)
	for i := 0; i < ret.numInputSites; i++ {
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

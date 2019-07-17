package validate

import (
	. "github.com/lunfardo314/goq/abra"
	"sort"
	"strings"
)

type SortedBlocks []*Block

func (sb SortedBlocks) Len() int {
	return len(sb)
}

func (sb SortedBlocks) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}

func (sb SortedBlocks) Less(i, j int) bool {
	if sb[i].BlockType == sb[j].BlockType {
		return strings.Compare(sb[i].LookupName, sb[j].LookupName) < 0
	}
	switch {
	case sb[i].BlockType == BLOCK_LUT && sb[j].BlockType == BLOCK_BRANCH:
		return true
	case sb[i].BlockType == BLOCK_LUT && sb[j].BlockType == BLOCK_EXTERNAL:
		return true
	case sb[i].BlockType == BLOCK_BRANCH && sb[j].BlockType == BLOCK_EXTERNAL:
		return true

	case sb[i].BlockType == BLOCK_BRANCH && sb[j].BlockType == BLOCK_LUT:
		return false
	case sb[i].BlockType == BLOCK_EXTERNAL && sb[j].BlockType == BLOCK_LUT:
		return false
	case sb[i].BlockType == BLOCK_EXTERNAL && sb[j].BlockType == BLOCK_BRANCH:
		return false
	}
	panic("wrong block type")
}

func SortAndEnumerateBlocks(codeUnit *CodeUnit) (int, int, int) {
	var numBranch, numLUTs, numExternal int
	sort.Sort(SortedBlocks(codeUnit.Code.Blocks))
	for i, block := range codeUnit.Code.Blocks {
		block.Index = i
		switch block.BlockType {
		case BLOCK_LUT:
			numLUTs++
		case BLOCK_BRANCH:
			numBranch++
		case BLOCK_EXTERNAL:
			numExternal++
		default:
			panic("wrong block type")
		}
	}
	return numLUTs, numBranch, numExternal
}

type SortedSites []*Site

func (ss SortedSites) Len() int {
	return len(ss)
}

func (ss SortedSites) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

type stronglyOrderedTypes struct{ lhs, rhs SiteType }

var stronglylLessPairs = []stronglyOrderedTypes{
	{SITE_INPUT, SITE_BODY},
	{SITE_INPUT, SITE_STATE},
	{SITE_INPUT, SITE_OUTPUT},
	{SITE_BODY, SITE_STATE},
	{SITE_BODY, SITE_OUTPUT},
	{SITE_STATE, SITE_OUTPUT},
}

func (ss SortedSites) Less(i, j int) bool {
	for _, pair := range stronglylLessPairs {
		switch {
		case ss[i].SiteType == pair.lhs && ss[j].SiteType == pair.rhs:
			return true
		case ss[i].SiteType == pair.rhs && ss[j].SiteType == pair.lhs:
			return false
		}
	}
	if ss[i].SiteType != ss[j].SiteType {
		panic("inconsistency with conditions")
	}
	return ss[i].Index < ss[j].Index
}

func SortAndEnumerateSites(codeUnit *CodeUnit) {
	for _, block := range codeUnit.Code.Blocks {
		if block.BlockType == BLOCK_BRANCH {
			sort.Sort(SortedSites(block.Branch.AllSites))
			for i, site := range block.Branch.AllSites {
				site.Index = i
			}
		}
	}
}

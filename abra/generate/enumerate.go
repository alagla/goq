package generate

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

func SortAndEnumerateBocks(codeUnit *CodeUnit) (int, int) {
	var numBranch, numLUTs int
	sort.Sort(SortedBlocks(codeUnit.Code.Blocks))
	for i, block := range codeUnit.Code.Blocks {
		block.Index = i
		switch block.BlockType {
		case BLOCK_LUT:
			numLUTs++
		case BLOCK_BRANCH:
			numBranch++
		default:
			panic("wrong block type")
		}
	}
	return numLUTs, numBranch
}

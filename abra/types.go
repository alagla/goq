package abra

import . "github.com/iotaledger/iota.go/trinary"

// package contains code which allows loading, saving and interpretation of teh Abra tritcode
// it is independent from Qupla definitions

// https://github.com/iotaledger/omega-docs/blob/master/qbc/abra/Spec.md

type CodeUnit struct {
	EntityAttachments []*EntityAttachment
	Code              *Code
}

//Entity attachment:
//[ code hash (243 trits)
//, number of attachments (positive integer)
//, attachments...
//]

type EntityAttachment struct {
	CodeHash    Hash // 243 trits = 81 trytes
	Attachments []*Attachment
}

//Attachment:
//[ branch block index (positive integer)
//, maximum recursion depth (positive integer)
//, number of input environments (positive integer)
//, input environment data...
//, number of output environments (positive integer)
//, output environment data...
//]

type Attachment struct {
	Branch                *Branch
	MaximumRecursionDepth int
	InputEnvironments     []*InputEnvironmentData
	OutputEnvironments    []*OutputEnvironmentData
	// compile time
	InputEnvironmentsDict  map[Hash]*InputEnvironmentData
	OutputEnvironmentsDict map[Hash]*OutputEnvironmentData
}

//input environment data:
//[ environment hash
//, limit (positive integer)
//, first branch input index (positive integer)
//, last branch input index (positive integer)
//]

type InputEnvironmentData struct {
	EnvironmentHash Hash
	Limit           int
	//FirstBranchInputIndex int   //???
	//LastBranchInputIndex  int   //???
}

//output environment data:
//[ environment hash
//, delay (positive integer)
//, first branch output index (positive integer)
//, last branch output index (positive integer)
//]

type OutputEnvironmentData struct {
	EnvironmentHash Hash
	Delay           int
	//FirstBranchInputIndex int   //??
	//LastBranchInputIndex  int   //??
}

//code:
//[ tritcode version (positive integer [0])
//, number of lookup table blocks (positive integer)
//, 35-trit lookup tables (27 nullable trits in bct)...
//, number of branch blocks (positive integer)
//, branch block definitions ...
//, number of external blocks (positive integer)
//, external block definitions...
//]

type Code struct {
	TritcodeVersion int
	Blocks          []*Block
}

//LUT definition
//The lookup table is encoded as 27 nullable trits, which fits in a 35 -trit number as 27 binary-coded trits.
// A lookup table which returns 0 for any input would look, in binary, like 3F_FF_FF_FF_FF.
//
//Since this value only covers for any non-null possible inputs, we start encoding by starting at
// all negatives (first input as lowest-endian), ---, and continuing to increment: 0--, 1--, -0-, 00-, 10-, ..., 111.
//
//Thus, the most-significant pair of bits (binary-coded trits) corresponds to 111, and the
// least significant pair of bits corresponds to ---.
//
//This final value is treated as a binary number, and encoded within a 35-trit vector.

type LUT int64

//block (whether external, lut, or branch):
//[ number of trits in block definition (positive integer)
//, value...
//]

//branch:
//[ number of inputs (positive integer)
//, input lengths (positive integers)...
//, number of body sites (positive integer)
//, number of output sites (positive integer)
//, number of memory latch sites (positive integer)
//, body site definitions...
//, output site definitions...
//, memory latch site definitions...
//]

type Branch struct {
	inputSites  []*Site
	bodySites   []*Site
	outputSites []*Site
	stateSites  []*Site
	// compile time
	AllSites []*Site
	Size     int
}

//site:
//[ Merge / Knot? 1 trit (1/-)
//, value...
//]
type SiteType int

const (
	SITE_INPUT  = SiteType(0)
	SITE_BODY   = SiteType(1)
	SITE_STATE  = SiteType(2)
	SITE_OUTPUT = SiteType(3)
)

type Site struct {
	Index  int // index within branch
	IsKnot bool
	Merge  *Merge // SITE_MERGE
	Knot   *Knot  // SITE_KNOT
	Size   int    // SITE_INPUT
	// lookup name, compile time only
	LookupName string
	SiteType   SiteType
}

//Merge:
//[ number of input sites (positive integer)
//, input site indices (positive integers)...
//]

type Merge struct {
	//NumberOfInputSites int
	Sites []*Site
}

//Knot:
//[ number of input sites (positive integer)
//, input site indices (positive integers)...
//, block index
//]

type Knot struct {
	//NumberOfInputSites int
	Sites []*Site
	Block *Block
}

//external block:
//[ code hash
//, number of blocks to import (positive integer)
//, block indices (positive integers)...
//]

type ExternalBlock struct {
	CodeHash     Trits
	BlockIndices []int
}

type BlockType int

const (
	BLOCK_LUT      = BlockType(0)
	BLOCK_BRANCH   = BlockType(1)
	BLOCK_EXTERNAL = BlockType(2)
)

type Block struct {
	Index         int // block index, one for LUTs and branches
	BlockType     BlockType
	Branch        *Branch
	LUT           LUT
	ExternalBlock *ExternalBlock
	// lookup name, compile time only
	LookupName string
}

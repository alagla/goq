package abra

import . "github.com/iotaledger/iota.go/trinary"

// package contains code which allows loading, saving and interpretation of teh Abra tritcode
// it is independent from Qupla definitions

// https://github.com/iotaledger/omega-docs/blob/master/qbc/abra/Spec.md

type Tritcode interface {
	GetTritcode() Trits
}

type CodeUnit struct {
	EntityAttachment *EntityAttachment
	Code             *CodeStruct
}

//Entity attachment:
//[ code hash (243 trits)
//, number of attachments (positive integer)
//, attachments...
//]

type EntityAttachment struct {
	CodeHash    Trits // 243 trit
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
	BranchBlockIndex      int
	MaximumRecursionDepth int
	InputEnvironments     []*InputEnvironmentData
	OutputEnvironments    []*OutputEnvironmentData
}

//input environment data:
//[ environment hash
//, limit (positive integer)
//, first branch input index (positive integer)
//, last branch input index (positive integer)
//]

type InputEnvironmentData struct {
	EnvironmentHash       Trits
	Limit                 int
	FirstBranchInputIndex int
	LastBranchInputIndex  int
}

//output environment data:
//[ environment hash
//, delay (positive integer)
//, first branch output index (positive integer)
//, last branch output index (positive integer)
//]

type OutputEnvironmentData struct {
	EnvironmentHash       Trits
	Delay                 int
	FirstBranchInputIndex int
	LastBranchInputIndex  int
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

type CodeStruct struct {
	TritcodeVersion int
	LUTs            []*LUT
	Branches        []*Branch
	ExternalBlocks  []*ExternalBlock
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

type LUT uint64

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
	//NumberOfInputs   int
	InputLengths     []int
	BodySites        []*Site
	OutputSites      []*Site
	MemoryLatchSites []*Site
}

//site:
//[ merge / knot? 1 trit (1/-)
//, value...
//]

//merge:
//[ number of input sites (positive integer)
//, input site indices (positive integers)...
//]

type Merge struct {
	//NumberOfInputSites int
	InputSiteIndices []int
}

//knot:
//[ number of input sites (positive integer)
//, input site indices (positive integers)...
//, block index
//]

type Knot struct {
	//NumberOfInputSites int
	InputSiteIndices []int
	BlockIndex       int
}

type Site struct {
	isMerge bool
	merge   *Merge // isMerge == true
	knot    *Knot  // isMerge == false
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

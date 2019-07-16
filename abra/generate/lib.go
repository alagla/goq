package generate

import (
	"bytes"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abra"
)

type Tritcode bytes.Buffer

func (code *Tritcode) WriteTrits(trits Trits) error {
	b := make([]byte, len(trits))
	for i, v := range trits {
		b[i] = byte(v)
	}
	buf := bytes.Buffer(*code)
	_, err := buf.Write(b)
	return err
}

func (code *Tritcode) WriteBlock(block *Block) error {
	return code.WriteTrits(GetBlockTritcode(block))
}

//
func GetBlockTritcode(block *Block) Trits {
	switch block.BlockType {
	case BLOCK_LUT:
		return TritEncodeLUT(block.LUT)
	case BLOCK_BRANCH:
		panic("Implement me")
	case BLOCK_EXTERNAL:
		panic("Implement me")
	}
	panic("wrong block type")
}

func TritEncodeLUT(lut *LUT) Trits {
	ret := IntToTrits(int64(lut.Binary))
	ret = PadTrits(ret, 35)
	if len(ret) != 35 {
		panic("wrong LUT tritcode")
	}
	return ret
}

package generate

import (
	"bufio"
	"fmt"
	"github.com/lunfardo314/goq/abra"
	"os"
)

func SaveReadable(codeUnit *abra.CodeUnit, fname string) error {
	fout, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer fout.Close()
	w := bufio.NewWriter(fout)
	_, err = fmt.Fprintf(w, "Tritcode version: %d\n", codeUnit.Code.TritcodeVersion)
	if err != nil {
		return err
	}
	for _, block := range codeUnit.Code.Blocks {
		err = writeReadableBlock(w, block)
	}
	w.Flush()
	return nil
}

func blockTypeName(bt abra.BlockType) string {
	switch bt {
	case abra.BLOCK_LUT:
		return "LUT"
	case abra.BLOCK_BRANCH:
		return "BRANCH"
	case abra.BLOCK_EXTERNAL:
		return "EXTERNAL"
	}
	panic("bad")
}

func writeReadableBlock(w *bufio.Writer, b *abra.Block) error {
	var ln, qn string
	var err error

	qn = b.QuplaFunName
	if qn == "" {
		qn = "?"
	}
	ln = b.LookupName
	if ln == "" {
		ln = "?"
	}
	_, err = fmt.Fprintf(w, "block #%d %s, %s / %s, %d -> %d\n",
		b.Index, blockTypeName(b.BlockType), b.LookupName, b.QuplaFunName, b.SizeIn, b.SizeOut)
	if err != nil {
		return err
	}
	if b.BlockType != abra.BLOCK_BRANCH {
		return nil
	}
	for _, s := range b.Branch.AllSites {
		n := ""
		switch s.SiteType {
		case abra.SITE_INPUT:
			_, err = fmt.Fprintf(w, "     site #%d INPUT: size: %d   // %s\n", s.Index, s.Size, s.LookupName)
			continue
		case abra.SITE_BODY:
			n = "BODY  "
		case abra.SITE_OUTPUT:
			n = "OUTPUT"
		case abra.SITE_STATE:
			n = "STATE "
		}
		if s.IsKnot {
			_, err = fmt.Fprintf(w, "     site #%d %s - KNOT:  size: %d inp: %v ref_block: #%d (%s / %s)   // %s\n",
				s.Index, n, s.Size, getIndices(s.Knot.Sites), s.Knot.Block.Index,
				s.Knot.Block.LookupName, s.Knot.Block.QuplaFunName, s.LookupName)
		} else {
			_, err = fmt.Fprintf(w, "     site #%d %s - MERGE: size: %d inp: %v  // %s\n",
				s.Index, n, s.Size, getIndices(s.Merge.Sites), s.LookupName)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func getIndices(sites []*abra.Site) []int {
	ret := make([]int, 0, 10)
	for _, s := range sites {
		ret = append(ret, s.Index)
	}
	return ret
}

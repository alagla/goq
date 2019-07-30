package generate

import (
	"bufio"
	"fmt"
)

type AbraTest struct {
	BlockIndex int
	Input      string
	Expected   string
	IsFloat    bool
	Comment    string
}

func WriteAbraTest(w *bufio.Writer, t *AbraTest, index int) error {
	isFloat := "0"
	if t.IsFloat {
		isFloat = "1"
	}
	_, err := fmt.Fprintf(w, "test %3d %3d %s %s %s // %s\n",
		index, t.BlockIndex, t.Input, t.Expected, isFloat, t.Comment)
	return err
}

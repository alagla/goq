package qupla

import (
	"bufio"
	"fmt"
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
	"github.com/lunfardo314/goq/abra/generate"
	vabra "github.com/lunfardo314/goq/abra/validate"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/utils"
	"os"
)

func (module *QuplaModule) GetAbra(codeUnit *abra.CodeUnit) {
	// TODO environments etc
	Logf(2, "---- generating LUT blocks")
	count := 0
	for _, lut := range module.Luts {
		strRepr := lut.GetStringRepr()
		if cabra.FindLUTBlock(codeUnit, strRepr) != nil {
			continue
		}
		cabra.MustAddNewLUTBlock(codeUnit, strRepr, lut.Name)
		count++
	}

	Logf(2, "---- generating branch blocks")
	for _, fun := range module.Functions {
		fun.GetAbraBranchBlock(codeUnit)
	}

	vabra.SortAndEnumerateBlocks(codeUnit)
	vabra.SortAndEnumerateSites(codeUnit)

	Logf(0, "total %d LUTs, %d branches, %d external blocks",
		codeUnit.Code.NumLUTs, codeUnit.Code.NumBranches, codeUnit.Code.NumExternalBlocks)
}

func (module *QuplaModule) WriteAbraTests(codeUnit *abra.CodeUnit, fname string) error {
	fout, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer fout.Close()
	w := bufio.NewWriter(fout)

	// writing block lookup names and
	var ln, qn string
	var stateful string

	for idx, block := range codeUnit.Code.Blocks {
		qn = block.QuplaFunName
		if qn == "" {
			qn = "?"
		}
		ln = block.LookupName
		if ln == "" {
			ln = "?"
		}

		stateful = "?"
		switch block.BlockType {
		case abra.BLOCK_LUT:
			stateful = "0"
		case abra.BLOCK_BRANCH:
			fun := module.FindFuncDef(block.QuplaFunName)
			if fun != nil {
				stateful = boolStr(fun.HasState())
			}
		}

		_, err = fmt.Fprintf(w, "block %4d %s %30s %30s\n", idx, stateful, ln, qn)
		if err != nil {
			return err
		}
	}

	for idx, exec := range module.Execs {
		if !exec.isTest {
			continue
		}

		funExpr, ok := exec.Expr.(*FunctionExpr)
		if !ok {
			continue
		}
		abra_idx, _ := cabra.FindBlockByQuplaName(codeUnit, funExpr.FuncDef.Name)
		input := ""
		for _, subExpr := range funExpr.GetSubexpressions() {
			if valExpr, ok := subExpr.(*ValueExpr); ok {
				input += utils.TritsToString(valExpr.TritValue)
			} else {
				input += "?"
			}
		}
		tst := &generate.AbraTest{
			BlockIndex: abra_idx,
			Input:      input,
			Expected:   utils.TritsToString(exec.expected),
			IsFloat:    exec.isFloat,
			Comment:    exec.GetSource(),
		}
		err = generate.WriteAbraTest(w, tst, idx)
		if err != nil {
			return err
		}
	}
	return nil
}

func boolStr(c bool) string {
	if c {
		return "1"
	}
	return "0"
}

func concatArgs(exec *ExecStmt) string {
	ret := ""
	expr, ok := exec.Expr.(*FunctionExpr)
	if !ok {
		return "?"
	}
	for _, subExpr := range expr.GetSubexpressions() {
		valExpr, ok := subExpr.(*ValueExpr)
		if ok {
			ret += utils.TritsToString(valExpr.TritValue)
		} else {
			ret += "?"
		}
	}
	return ret
}

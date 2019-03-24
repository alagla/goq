package analyzeyaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeExecStmt(execStmtYAML *QuplaExecStmtYAML, module *QuplaModule) error {
	var err error
	var expr ExpressionInterface
	var ok bool
	expr, err = AnalyzeExpression(execStmtYAML.Expr, module, nil)
	if err != nil {
		return err
	}
	var funcExpr *FunctionExpr
	if funcExpr, ok = expr.(*FunctionExpr); !ok {
		return fmt.Errorf("top expression must be call to a function: '%v'", execStmtYAML.Source)
	}
	isTest := execStmtYAML.Expected != nil
	isFloat := execStmtYAML.IsFloat
	var expected Trits

	if isTest {
		exprExpected, err := AnalyzeExpression(execStmtYAML.Expected, module, nil)
		if err != nil {
			return err
		}
		// check sizes
		if err = MatchSizes(funcExpr, exprExpected); err != nil {
			return err
		}
		ve, ok := exprExpected.(*ValueExpr)
		if !ok {
			return fmt.Errorf("test '%v': left hand side must be ValueExpr", execStmtYAML.Source)
		}
		expected = ve.TritValue
		module.IncStat("numTest")
	} else {
		module.IncStat("numEval")
	}
	module.AddExec(NewExecStmt(execStmtYAML.Source, funcExpr, isTest, isFloat, expected, module))
	return nil
}

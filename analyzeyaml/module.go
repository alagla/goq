package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/utils"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeQuplaModule(name string, moduleYAML *QuplaModuleYAML) (*QuplaModule, bool) {
	ret := NewQuplaModule(name)

	retSucc := true
	logf(0, "Analyzing LUTs..")
	numLuts := 0
	for name, lutDefYAML := range moduleYAML.Luts {
		numLuts++
		if err := AnalyzeLutDef(name, lutDefYAML, ret); err != nil {
			ret.IncStat("numErr")
			logf(0, "%v", err)
			retSucc = false
		}
	}
	logf(1, "Analyzed %v LUTs", numLuts)

	retSucc = true

	logf(0, "Analyzing functions..")
	for funName, funDefYAML := range moduleYAML.Functions {
		if err := AnalyzePreliminaryFuncDef(funName, funDefYAML, ret); err != nil {
			ret.IncStat("numErr")
			errorf("%v", err)
			retSucc = false
		}
	}
	for funName, funDefYAML := range moduleYAML.Functions {
		if err := AnalyzeFuncDef(funName, funDefYAML, ret); err != nil {
			ret.IncStat("numErr")
			errorf("%v", err)
			retSucc = false
		}
	}
	logf(1, "Analyzed %v functions", len(ret.Functions))

	numWithStateVars, numStateful := ret.MarkStateful()
	logf(1, "Functions with state variables: %v", numWithStateVars)
	logf(1, "Functions with state (which references functions with state variables): %v", numStateful)

	//if n, ok := ret.stats["numEnvFundef"]; ok {
	//	logf(0, "Functions joins/affects environments: %v", n)
	//} else {
	//	logf(0, "No function joins/affects environments")
	//}

	for funname, fundef := range ret.Functions {
		if fundef.HasEnvStmt() {
			joins := StringSet{}
			for e, p := range fundef.GetJoinEnv() {
				joins.Append(fmt.Sprintf("%v(%v)", e, p))
			}
			affects := StringSet{}
			for e, p := range fundef.GetAffectEnv() {
				affects.Append(fmt.Sprintf("%v(%v)", e, p))
			}
			ret.Environments.AppendAll(affects)
			ret.Environments.AppendAll(joins)
			logf(1, "    Function '%v' joins: '%v', affects: '%v'",
				funname, joins.Join(","), affects.Join(","))
		}
	}
	if len(ret.Environments) > 0 {
		logf(1, "Environments: '%v'", ret.Environments.Join(", "))
	} else {
		logf(1, "Environments: none")
	}

	logf(0, "Analyzing execs (tests and evals)..")
	numExec := 0
	for _, execYAML := range moduleYAML.Execs {
		err := AnalyzeExecStmt(execYAML, ret)
		numExec++
		if err != nil {
			ret.IncStat("numErr")
			errorf("%v", err)
			retSucc = false
		}
	}
	logf(1, "Executable statements found: %v", numExec)

	return ret, retSucc
}

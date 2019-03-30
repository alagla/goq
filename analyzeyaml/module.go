package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	. "github.com/lunfardo314/goq/utils"
)

func AnalyzeQuplaModule(name string, moduleYAML *QuplaModuleYAML) (*QuplaModule, bool) {
	ret := NewQuplaModule(name)

	retSucc := true
	logf(1, "Analyzing LUTs..")
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

	logf(1, "Analyzing functions..")
	for funName, funDefYAML := range moduleYAML.Functions {
		if err := AnalyzeFunctionPreliminary(funName, funDefYAML, ret); err != nil {
			ret.IncStat("numErr")
			logf(0, "%v", err)
			retSucc = false
		}
	}
	for funName, funDefYAML := range moduleYAML.Functions {
		if err := AnalyzeFunction(funName, funDefYAML, ret); err != nil {
			ret.IncStat("numErr")
			logf(0, "%v", err)
			retSucc = false
		}
	}
	logf(1, "Analyzed %v functions", len(ret.Functions))

	// scans all functions, collects those, which has state vars themselves or dircetly/indirectly
	// references those with state vars
	numWithStateVars, numStateful := ret.MarkStateful()
	logf(1, "Functions with state variables: %v", numWithStateVars)
	logf(1, "Functions with state (which references functions with state variables): %v", numStateful)

	for funname, fundef := range ret.Functions {
		if fundef.HasEnvStmt() {
			joins := StringSet{}
			for e := range fundef.GetJoinEnv() {
				joins.Append(fmt.Sprintf("%v", e))
			}
			affects := StringSet{}
			for e := range fundef.GetAffectEnv() {
				affects.Append(fmt.Sprintf("%v", e))
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

	logf(1, "Analyzing execs (tests and evals)..")
	numExec := 0
	for _, execYAML := range moduleYAML.Execs {
		err := AnalyzeExecStmt(execYAML, ret)
		numExec++
		if err != nil {
			ret.IncStat("numErr")
			logf(0, "%v", err)
			retSucc = false
		}
	}
	logf(1, "Executable statements found: %v", numExec)

	return ret, retSucc
}

package qupla

import (
	"fmt"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
	"time"
)

type QuplaModule struct {
	yamlSource *QuplaModuleYAML
	factory    ExpressionFactory
	types      map[string]*QuplaTypeDef
	luts       map[string]*QuplaLutDef
	functions  map[string]*QuplaFuncDef
	execs      []*QuplaExecStmt
	stats      map[string]int
	processor  ProcessorInterface
}

type QuplaTypeField struct {
	offset int64
	size   int64
}

type QuplaTypeDef struct {
	size   int64
	fields map[string]*QuplaTypeField
}

func AnalyzeQuplaModule(moduleYAML *QuplaModuleYAML, factory ExpressionFactory) (*QuplaModule, bool) {
	ret := &QuplaModule{
		yamlSource: moduleYAML,
		factory:    factory,
		types:      make(map[string]*QuplaTypeDef),
		luts:       make(map[string]*QuplaLutDef),
		functions:  make(map[string]*QuplaFuncDef),
		execs:      make([]*QuplaExecStmt, 0, len(moduleYAML.Execs)),
		stats:      make(map[string]int),
		processor:  NewStackProcessor(),
	}
	//infof("Analyzing types..")
	//for name, td := range moduleYAML.Types {
	//	_ = ret.AnalyzeType(name, td)
	//}

	retSucc := true
	infof("Analyzing tests and evals..")
	for _, execYAML := range moduleYAML.Execs {
		err := AnalyzeExecStmt(execYAML, ret)
		if err != nil {
			ret.IncStat("numErr")
			errorf("%v", err)
			retSucc = false
		}
	}
	analyzedFunDefs := len(ret.functions)
	infof("Number of function definitions directly or indirectly referenced by execs: %v", analyzedFunDefs)

	infof("Analyzing all function definitions, which were not analyzed yet..")
	for funName := range moduleYAML.Functions {
		if _, err := ret.FindFuncDef(funName); err != nil {
			ret.IncStat("numErr")
			errorf("%v", err)
			retSucc = false
		}
	}
	infof("Additionally, were analyzed %v function definitions", len(ret.functions)-analyzedFunDefs)

	return ret, retSucc
}

//func (module *QuplaModule) AnalyzeType(name string, src *QuplaTypeDefYAML) bool {
//	if _, ok := module.types[name]; ok {
//		errorf("duplicate type name %v", name)
//		module.IncStat("numErr")
//		return false
//	}
//	ret := &QuplaTypeDef{
//		fields: make(map[string]*QuplaTypeField),
//	}
//	if src.Size != "*" {
//		if sz, err := strconv.Atoi(src.Size); err != nil {
//			errorf("wrong size '%v' in type '%v'", src.Size, name)
//			module.IncStat("numErr")
//			return false
//		} else {
//			ret.size = int64(sz)
//			module.types[name] = ret
//			return true
//		}
//	}
//
//	var offset int64
//	for fldname, fld := range src.Fields {
//		if sz, err := strconv.Atoi(fld.Size); err != nil {
//			errorf("wrong size '%v' in field '%v' of type '%v'", fld.Size, fldname, name)
//			module.IncStat("numErr")
//			return false
//
//		} else {
//			ret.fields[fldname] = &QuplaTypeField{
//				offset: offset,
//				size:   int64(sz),
//			}
//			offset += ret.size
//		}
//	}
//	module.types[name] = ret
//	return true
//}

func (module *QuplaModule) GetTypeFieldInfo(typeName, fldName string) (int64, int64, error) {
	if _, ok := module.types[typeName]; !ok {
		return 0, 0, fmt.Errorf("can't find type '%v", typeName)
	}
	fi, ok := module.types[typeName].fields[fldName]
	if !ok {
		return 0, 0, fmt.Errorf("can't find field '%v' in type '%v", fldName, typeName)
	}
	return fi.offset, fi.size, nil
}

func (module *QuplaModule) AnalyzeExpression(data interface{}, scope FuncDefInterface) (ExpressionInterface, error) {
	return module.factory.AnalyzeExpression(data, module, scope)
}

func (module *QuplaModule) AddExec(exec *QuplaExecStmt) {
	exec.num = len(module.execs)
	module.execs = append(module.execs, exec)
}

func (module *QuplaModule) AddFuncDef(name string, funcDef FuncDefInterface) {
	module.functions[name] = funcDef.(*QuplaFuncDef)
}

func (module *QuplaModule) AddLutDef(name string, lutDef LUTInterface) {
	module.luts[name] = lutDef.(*QuplaLutDef)
}

func (module *QuplaModule) FindFuncDef(name string) (FuncDefInterface, error) {
	var err error
	ret, ok := module.functions[name]
	if ok {
		return ret, nil
	}
	src, ok := module.yamlSource.Functions[name]
	if !ok {
		return nil, fmt.Errorf("can't find function definition '%v'", name)
	}
	err = AnalyzeFuncDef(name, src, module)
	if err != nil {
		return nil, fmt.Errorf("error while anlyzing fun def '%v': %v", name, err)
	}
	return module.functions[name], nil
}

func (module *QuplaModule) FindLUTDef(name string) (LUTInterface, error) {
	var err error
	ret, ok := module.luts[name]
	if ok {
		return ret, nil
	}
	src, ok := module.yamlSource.Luts[name]
	if !ok {
		return nil, fmt.Errorf("can't find LUT dfinition '%v'", name)
	}
	err = AnalyzeLutDef(name, src, module)
	if err != nil {
		return nil, err
	}

	if ret, ok = module.luts[name]; !ok {
		return nil, fmt.Errorf("inconsistency while analyzing LUT '%v'", name)
	}
	return ret, nil
}

func (module *QuplaModule) Execute(test bool) {
	if test {
		infof("Executing tests..")
	} else {
		infof("Executing evals and tests (total %v)..", len(module.execs))
	}

	testsPassed := 0
	testsSkipped := 0
	totalTests := 0
	start := time.Now()

	for _, exec := range module.execs {
		if exec.num == 7 {
			fmt.Printf("kuku\n")
		}
		infof("-----------------------")
		infof("Check exec statement: '%v'", exec.GetSource())

		if exec.HasState() {
			infof("SKIP because it has state: '%v'", exec.GetSource())
			testsSkipped++
			continue
		}
		if duration, passed, err := exec.Execute(); err != nil {
			errorf("Error: %v", err)
		} else {
			if exec.isTest {
				totalTests++
				if passed {
					testsPassed++
					infof("Test PASSED. Duration %v", duration)
				} else {
					infof("Test FAILED. Duration %v", duration)
				}
			} else {
				infof("Duration %v", duration)
			}
		}
		module.processor.Reset()
	}
	infof("Total tests and evals: %v", len(module.execs))
	var p string
	if testsPassed == 0 {
		p = "n/a"
	} else {
		p = fmt.Sprintf("%v%%", (testsPassed*100)/totalTests)
	}
	infof("---------------------")
	infof("---------------------")
	infof("Skipped: %v out of total %v executables", testsSkipped, len(module.execs))
	infof("Tests PASSED: %v out of %v (%v)", testsPassed, totalTests, p)
	infof("Total duration: %v ", time.Since(start))
}

func (module *QuplaModule) IncStat(key string) {
	if _, ok := module.stats[key]; !ok {
		module.stats[key] = 0
	}
	module.stats[key]++
}

func (module *QuplaModule) PrintStats() {
	fmt.Printf("Analyzed: \n")
	for k, v := range module.stats {
		fmt.Printf("  %v : %v\n", k, v)
	}
}

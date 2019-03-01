package qupla

import (
	"fmt"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaModule struct {
	yamlSource *QuplaModuleYAML
	factory    ExpressionFactory
	luts       map[string]*QuplaLutDef
	functions  map[string]*QuplaFuncDef
	execs      []*QuplaExecStmt
	stats      map[string]int
	processor  ProcessorInterface
}

func AnalyzeQuplaModule(moduleYAML *QuplaModuleYAML, factory ExpressionFactory) (*QuplaModule, bool) {
	ret := &QuplaModule{
		yamlSource: moduleYAML,
		factory:    factory,
		luts:       make(map[string]*QuplaLutDef),
		functions:  make(map[string]*QuplaFuncDef),
		execs:      make([]*QuplaExecStmt, 0, len(moduleYAML.Execs)),
		stats:      make(map[string]int),
		processor:  NewStackProcessor(),
	}
	infof("Analyzing...")
	retSucc := true
	for _, execYAML := range moduleYAML.Execs {
		err := AnalyzeExecStmt(execYAML, ret)
		if err != nil {
			ret.IncStat("numErr")
			errorf("%v", err)
			retSucc = false
			continue
		}
	}
	return ret, retSucc
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

func (module *QuplaModule) Execute() {
	for _, exec := range module.execs {
		if err := exec.Execute(); err != nil {
			errorf("execute error: %v", err)
		}
	}
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

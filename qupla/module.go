package qupla

import (
	"fmt"
	"github.com/lunfardo314/goq/quplayaml"
)

type QuplaModule struct {
	yamlSource *quplayaml.QuplaModuleYAML
	factory    ExpressionFactory
	types      map[string]*QuplaTypeDef
	luts       map[string]*QuplaLutDef
	functions  map[string]*QuplaFuncDef
	execs      []*QuplaExecStmt
	stats      map[string]int
}

func AnalyzeQuplaModule(moduleYAML *quplayaml.QuplaModuleYAML, factory ExpressionFactory) (*QuplaModule, bool) {
	ret := &QuplaModule{
		yamlSource: moduleYAML,
		factory:    factory,
		types:      make(map[string]*QuplaTypeDef),
		luts:       make(map[string]*QuplaLutDef),
		functions:  make(map[string]*QuplaFuncDef),
		execs:      make([]*QuplaExecStmt, 0, len(moduleYAML.Execs)),
		stats:      make(map[string]int),
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

func (module *QuplaModule) AnalyzeExpression(data interface{}, scope *QuplaFuncDef) (ExpressionInterface, error) {
	return module.factory.AnalyzeExpression(data, module, scope)
}

func (module *QuplaModule) AddExec(exec *QuplaExecStmt) {
	module.execs = append(module.execs, exec)
}

func (module *QuplaModule) AddFuncDef(name string, funcDef *QuplaFuncDef) {
	module.functions[name] = funcDef
}

func (module *QuplaModule) FindFuncDef(name string) (*QuplaFuncDef, error) {
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
		return nil, fmt.Errorf("error while anlyzing function definioton '%v': %v", name, err)
	}
	return module.functions[name], nil
}

func (module *QuplaModule) FindLUTDef(name string) (*QuplaLutDef, error) {
	var err error
	ret, ok := module.luts[name]
	if ok {
		return ret, nil
	}
	src, ok := module.yamlSource.Luts[name]
	if !ok {
		return nil, fmt.Errorf("can't find LUT dfinition '%v'", name)
	}
	ret, err = AnalyzeLutDef(name, src, module)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (module *QuplaModule) FindTypeDef(name string) *QuplaTypeDef {
	ret, ok := module.types[name]
	if ok {
		return ret
	}
	return nil
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

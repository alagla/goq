package qupla

import (
	"fmt"
	. "github.com/lunfardo314/goq/abstract"
	"github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/entities"
	. "github.com/lunfardo314/goq/utils"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"strings"
)

type QuplaModule struct {
	yamlSource   *QuplaModuleYAML
	name         string
	factory      ExpressionFactory
	types        map[string]*QuplaTypeDef
	luts         map[string]*QuplaLutDef
	functions    map[string]*QuplaFuncDef
	execs        []*QuplaExecStmt
	stats        map[string]int
	processor    ProcessorInterface
	environments StringSet
}

type QuplaTypeField struct {
	offset int64
	size   int64
}

type QuplaTypeDef struct {
	size   int64
	fields map[string]*QuplaTypeField
}

func AnalyzeQuplaModule(name string, moduleYAML *QuplaModuleYAML, factory ExpressionFactory) (*QuplaModule, bool) {
	ret := &QuplaModule{
		yamlSource:   moduleYAML,
		name:         name,
		factory:      factory,
		types:        make(map[string]*QuplaTypeDef),
		luts:         make(map[string]*QuplaLutDef),
		functions:    make(map[string]*QuplaFuncDef),
		execs:        make([]*QuplaExecStmt, 0, len(moduleYAML.Execs)),
		stats:        make(map[string]int),
		processor:    NewStackProcessor(),
		environments: make(StringSet),
	}
	//infof("Analyzing types..")
	//for name, td := range moduleYAML.Types {
	//	_ = ret.AnalyzeType(name, td)
	//}

	retSucc := true
	logf(0, "Analyzing execs (tests and evals)..")
	for _, execYAML := range moduleYAML.Execs {
		err := AnalyzeExecStmt(execYAML, ret)
		if err != nil {
			ret.IncStat("numErr")
			errorf("%v", err)
			retSucc = false
		}
	}
	analyzedFunDefs := len(ret.functions)
	logf(0, "Number of functions directly or indirectly referenced by execs: %v", analyzedFunDefs)

	logf(0, "Analyzing all functions, which were not analyzed yet")
	for funName := range moduleYAML.Functions {
		if _, err := ret.FindFuncDef(funName); err != nil {
			ret.IncStat("numErr")
			errorf("%v", err)
			retSucc = false
		}
	}
	logf(0, "Additionally were analyzed %v functions", len(ret.functions)-analyzedFunDefs)

	logf(0, "Determining stateful functions")
	numWithStateVars, numStateful := ret.markStateful()
	logf(0, "Found %v func def with state vars and %v stateful functions (which references them)",
		numWithStateVars, numStateful)

	if n, ok := ret.stats["numEnvFundef"]; ok {
		logf(0, "Functions joins/affects environments: %v", n)
	} else {
		logf(0, "No function joins/affects environments")
	}
	for funname, fundef := range ret.functions {
		if fundef.HasEnvStmt() {
			joins := StringSet{}
			for e, p := range fundef.GetJoinEnv() {
				joins.Append(fmt.Sprintf("%v(%v)", e, p))
			}
			affects := StringSet{}
			for e, p := range fundef.GetAffectEnv() {
				affects.Append(fmt.Sprintf("%v(%v)", e, p))
			}
			ret.environments.AppendAll(affects)
			ret.environments.AppendAll(joins)
			logf(1, "    Function '%v' joins: '%v', affects: '%v'",
				funname, joins.Join(","), affects.Join(","))
		}
	}
	if len(ret.environments) > 0 {
		logf(0, "Environments detected: '%v'", ret.environments.Join(", "))
	} else {
		logf(0, "Environments detected: none")
	}
	return ret, retSucc
}

func (module *QuplaModule) GetName() string {
	return module.name
}

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
	exec.idx = len(module.execs)
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

func (module *QuplaModule) FindExecs(substr string) []*QuplaExecStmt {
	ret := make([]*QuplaExecStmt, 0)
	for _, ex := range module.execs {
		if strings.Contains(ex.GetName(), substr) {
			ret = append(ret, ex)
		}
	}
	return ret
}
func (module *QuplaModule) IncStat(key string) {
	if _, ok := module.stats[key]; !ok {
		module.stats[key] = 0
	}
	module.stats[key]++
}

func (module *QuplaModule) PrintStats() {
	logf(0, "Module stats:")
	for k, v := range module.stats {
		logf(0, "  %v : %v", k, v)
	}
}

func (module *QuplaModule) markStateful() (int, int) {
	stateful := make(StringSet)
	for name, fd := range module.functions {
		if fd.hasStateVariables {
			stateful.Append(name)
		}
	}
	hasStateVars := len(stateful)
	newCollected := module.collectReferencingFuncs(stateful)
	for ; newCollected > 0; newCollected = module.collectReferencingFuncs(stateful) {
	}

	for name := range stateful {
		module.functions[name].hasState = true
	}
	return hasStateVars, len(stateful)
}

func (module *QuplaModule) collectReferencingFuncs(nameSet StringSet) int {
	tmpList := make([]string, 0)
	for name := range nameSet {
		tmpList = append(tmpList, name)
	}
	ret := 0
	for _, statefulName := range tmpList {
		for name, fd := range module.functions {
			if fd.References(statefulName) {
				if nameSet.Append(name) {
					ret++
				}
			}
		}
	}
	return ret
}

func (module *QuplaModule) AttachToDispatcher(disp *dispatcher.Dispatcher) bool {
	ret := true
	for _, funcdef := range module.functions {
		if !funcdef.HasEnvStmt() {
			continue
		}
		entity := entities.NewFunctionEntity(disp, funcdef, NewStackProcessor())
		if err := disp.Attach(entity, funcdef.GetJoinEnv(), funcdef.GetAffectEnv()); err != nil {
			logf(0, "error while attaching entity to dispatcher: %v", err)
			ret = false
		}
	}
	return ret
}

func (module *QuplaModule) ExecByIdx(idx int) *QuplaExecStmt {
	if idx < 0 || idx >= len(module.execs) {
		return nil
	}
	return module.execs[idx]
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

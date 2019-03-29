package qupla

import (
	"fmt"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/supervisor"
	. "github.com/lunfardo314/goq/utils"
	"strings"
)

type QuplaModule struct {
	name         string
	types        map[string]*QuplaTypeDef
	luts         map[string]*LutDef
	Functions    map[string]*Function
	execs        []*ExecStmt
	stats        map[string]int
	Environments StringSet
	stateHashMap *StateHashMap
}

type QuplaTypeField struct {
	offset int64
	size   int64
}

type QuplaTypeDef struct {
	size   int64
	fields map[string]*QuplaTypeField
}

func NewQuplaModule(name string) *QuplaModule {
	return &QuplaModule{
		name:         name,
		types:        make(map[string]*QuplaTypeDef),
		luts:         make(map[string]*LutDef),
		Functions:    make(map[string]*Function),
		execs:        make([]*ExecStmt, 0, 10),
		stats:        make(map[string]int),
		Environments: make(StringSet),
		stateHashMap: newStateHashMap(),
	}
}

func (module *QuplaModule) GetName() string {
	return module.name
}

func (module *QuplaModule) GetStateHashMap() *StateHashMap {
	return module.stateHashMap
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

func (module *QuplaModule) AddExec(exec *ExecStmt) {
	exec.idx = len(module.execs)
	module.execs = append(module.execs, exec)
}

func (module *QuplaModule) AddFuncDef(name string, funcDef *Function) error {
	if _, ok := module.Functions[name]; ok {
		return fmt.Errorf("duplicate function degfinition '%v'", name)
	}
	module.Functions[name] = funcDef
	return nil
}

func (module *QuplaModule) AddLutDef(name string, lutDef *LutDef) {
	module.luts[name] = lutDef
}

func (module *QuplaModule) FindFuncDef(name string) *Function {
	if ret, ok := module.Functions[name]; ok {
		return ret
	}
	return nil
}

func (module *QuplaModule) FindLUTDef(name string) (*LutDef, error) {
	ret, ok := module.luts[name]
	if !ok {
		return nil, fmt.Errorf("can't find LUT definition '%v'", name)
	}
	return ret, nil
}

func (module *QuplaModule) FindExecs(filterSubstr string) []*ExecStmt {
	ret := make([]*ExecStmt, 0)
	for _, ex := range module.execs {
		if strings.Contains(ex.GetName(), filterSubstr) {
			ret = append(ret, ex)
		}
	}
	return ret
}

func (module *QuplaModule) FindFuncs(filterSubstr string) []*Function {
	ret := make([]*Function, 0)
	for _, fun := range module.Functions {
		if strings.Contains(fun.Name, filterSubstr) {
			ret = append(ret, fun)
		}
	}
	return ret
}

func (module *QuplaModule) SetTraceLevel(traceLevel int, filterSubstr string) []*Function {
	funcs := module.FindFuncs(filterSubstr)
	for _, f := range funcs {
		f.SetTraceLevel(traceLevel)
	}
	return funcs
}

func (module *QuplaModule) IncStat(key string) {
	if _, ok := module.stats[key]; !ok {
		module.stats[key] = 0
	}
	module.stats[key]++
}

func (module *QuplaModule) PrintStats() {
	Logf(2, "Module statistics:")
	Logf(2, "  module name: '%v'", module.name)
	for k, v := range module.stats {
		Logf(2, "  %v : %v", k, v)
	}
}

func (module *QuplaModule) MarkStateful() (int, int) {
	stateful := make(StringSet)
	for name, fd := range module.Functions {
		if fd.HasStateVariables {
			stateful.Append(name)
		}
	}
	hasStateVars := len(stateful)
	newCollected := module.collectReferencingFuncs(stateful)
	for ; newCollected > 0; newCollected = module.collectReferencingFuncs(stateful) {
	}

	for name := range stateful {
		module.Functions[name].hasState = true
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
		for name, fd := range module.Functions {
			if fd.References(statefulName) {
				if nameSet.Append(name) {
					ret++
				}
			}
		}
	}
	return ret
}

func (module *QuplaModule) AttachToSupervisor(sv *supervisor.Supervisor) bool {
	ret := true
	for _, funcdef := range module.Functions {
		if !funcdef.HasEnvStmt() {
			continue
		}
		entity, err := NewFunctionEntity(sv, funcdef)
		if err != nil {
			Logf(0, "can't create entity: %v", err)
			ret = false
			continue
		}
		if err = sv.Attach(entity, funcdef.GetJoinEnv(), funcdef.GetAffectEnv()); err != nil {
			Logf(0, "error while attaching entity to dispatcher: %v", err)
			ret = false
		}
	}
	return ret
}

func (module *QuplaModule) ExecByIdx(idx int) *ExecStmt {
	if idx < 0 || idx >= len(module.execs) {
		return nil
	}
	return module.execs[idx]
}

//func (module *QuplaModule) AnalyzeType(Name string, src *QuplaTypeDefYAML) bool {
//	if _, ok := module.types[Name]; ok {
//		errorf("duplicate type Name %v", Name)
//		module.IncStat("numErr")
//		return false
//	}
//	ret := &QuplaTypeDef{
//		Fields: make(map[string]*QuplaTypeField),
//	}
//	if src.Size != "*" {
//		if sz, err := strconv.Atoi(src.Size); err != nil {
//			errorf("wrong Size '%v' in type '%v'", src.Size, Name)
//			module.IncStat("numErr")
//			return false
//		} else {
//			ret.Size = int64(sz)
//			module.types[Name] = ret
//			return true
//		}
//	}
//
//	var Offset int64
//	for fldname, fld := range src.Fields {
//		if sz, err := strconv.Atoi(fld.Size); err != nil {
//			errorf("wrong Size '%v' in field '%v' of type '%v'", fld.Size, fldname, Name)
//			module.IncStat("numErr")
//			return false
//
//		} else {
//			ret.Fields[fldname] = &QuplaTypeField{
//				Offset: Offset,
//				Size:   int64(sz),
//			}
//			Offset += ret.Size
//		}
//	}
//	module.types[Name] = ret
//	return true
//}

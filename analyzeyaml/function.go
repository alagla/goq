package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	"strconv"
)

// analyzes return size, creates variables and create entry in module's table
// not finished analysis, but is is ok to analyze correctly recursive calls
func AnalyzeFunctionPreliminary(name string, defYAML *QuplaFuncDefYAML, module *QuplaModule) error {
	ce, err := AnalyzeExpression(defYAML.ReturnType, module, nil)
	if err != nil {
		return err
	}
	var sz int
	if sz, err = GetConstValue(ce); err != nil {
		return err
	}
	def := NewFunction(name, sz)

	if err = createVarScope(defYAML, def, module); err != nil {
		return err
	}
	return module.AddFuncDef(name, def)
}

func AnalyzeFunction(name string, defYAML *QuplaFuncDefYAML, module *QuplaModule) error {
	var err error

	module.IncStat("numFuncDef")

	def := module.FindFuncDef(name)
	if def == nil {
		return fmt.Errorf("inconsistency: function must be preliminary analyzed first: '%v'", name)
	}
	if def.Analyzed {
		return fmt.Errorf("attempt to analyze function '%v' twice", name)
	}
	if err = analyzeEnvironmentStatements(defYAML, def, module); err != nil {
		return err
	}
	if def.HasEnvStmt() {
		module.IncStat("numEnvFundef")
	}
	if err = analyzeAssigns(defYAML, def, module); err != nil {
		return err
	}
	if err = finalizeLocalVars(def, module); err != nil {
		return err
	}
	// return expression
	if def.RetExpr, err = AnalyzeExpression(defYAML.ReturnExpr, module, def); err != nil {
		return err
	}
	if def.RetExpr == nil {
		return fmt.Errorf("in funcdef '%v': return expression can't be nil", def.Name)
	}
	def.Analyzed = true
	return nil
}
func analyzeEnvironmentStatements(defYAML *QuplaFuncDefYAML, def *Function, module *QuplaModule) error {
	for _, envYAML := range defYAML.Env {
		switch envYAML.Type {
		case "join":
			p := 1
			if envYAML.Limit != "" {
				if val, err := strconv.Atoi(envYAML.Limit); err != nil {
					return fmt.Errorf("join in '%v': %v", def.Name, err)
				} else {
					p = val
				}
			}
			def.Joins[envYAML.Name] = p
			module.IncStat("numEnvJoin")
		case "affect":
			p := 0
			if envYAML.Delay != "" {
				if val, err := strconv.Atoi(envYAML.Delay); err != nil {
					return fmt.Errorf("affect in '%v': %v", def.Name, err)
				} else {
					p = val
				}
			}
			def.Affects[envYAML.Name] = p
			module.IncStat("numEnvAffect")
		default:
			return fmt.Errorf("bad typeof environment statement in '%v': %v", def.Name, envYAML.Type)
		}
	}
	return nil
}

func AnalyzeVar(vi *VarInfo, defYAML *QuplaFuncDefYAML, def *Function, module *QuplaModule) error {
	if vi.Analyzed {
		return nil
		//panic(fmt.Errorf("attempt to analyze variable '%v' twice in '%v'", vi.Name, def.Name))
	}
	vi.Analyzed = true

	if vi.IsParam {
		vi.Assign = nil
		return nil
	}
	e, ok := defYAML.Assigns[vi.Name]
	if !ok {
		return fmt.Errorf("inconsistency with vars")
	}
	var err error
	if vi.Assign, err = AnalyzeExpression(e, module, def); err != nil {
		return err
	}
	if vi.IsState {
		if vi.Size != vi.Assign.Size() {
			return fmt.Errorf("expression and state variable has different sizes in the assign")
		}
	} else {
		vi.Size = vi.Assign.Size()
	}
	return nil
}

func createVarScope(src *QuplaFuncDefYAML, def *Function, module *QuplaModule) error {
	// function parameters (first numParams)
	def.NumParams = len(src.Params)
	for idx, arg := range src.Params {
		if def.GetVarIdx(arg.ArgName) >= 0 {
			return fmt.Errorf("duplicate arg Name '%v'", arg.ArgName)
		}
		def.LocalVars = append(def.LocalVars, &VarInfo{
			Idx:      idx,
			Name:     arg.ArgName,
			Size:     arg.Size,
			Analyzed: true,
			IsParam:  true,
			IsState:  false,
		})
	}
	// the rest of indices belong to local vars (incl state)
	if len(src.State) > 0 {
		def.HasStateVariables = true
		def.StateHashMap = module.GetStateHashMap()
	}

	var idx int
	for name, s := range src.State {
		idx = def.GetVarIdx(name)
		if idx >= 0 {
			return fmt.Errorf("wrong declared state variable: '%v' in '%v'", name, def.Name)
		} else {
			// for old value
			def.LocalVars = append(def.LocalVars, &VarInfo{
				Idx:     len(def.LocalVars),
				Name:    name,
				Size:    s.Size,
				IsState: true,
			})
		}
		module.IncStat("numStateVars")
	}
	// variables defined by assigns
	var vi *VarInfo
	for name := range src.Assigns {
		vi, _ = def.VarByName(name)
		if vi != nil {
			if vi.IsParam {
				return fmt.Errorf("cannot assign to function parameter: '%v' in '%v'", name, def.Name)
			}
			if !vi.IsState {
				return fmt.Errorf("several assignment to the same var '%v' in '%v' is not allowed", name, def.Name)
			}
		} else {
			def.LocalVars = append(def.LocalVars, &VarInfo{
				Idx:     len(def.LocalVars),
				Name:    name,
				Size:    0, // unknown yet
				IsState: false,
				IsParam: false,
			})
		}
	}
	return nil
}

func analyzeAssigns(defYAML *QuplaFuncDefYAML, def *Function, module *QuplaModule) error {
	for name := range defYAML.Assigns {
		// GetVarInfo analyzes expression if necessary
		vi, err := def.VarByName(name)
		if err != nil {
			panic(fmt.Errorf("'%v' : %v", def.Name, err))
		}
		if err = AnalyzeVar(vi, defYAML, def, module); err != nil {
			return err
		}
	}
	return nil
}

func finalizeLocalVars(def *Function, module *QuplaModule) error {
	var curOffset int
	def.InSize = 0
	for _, v := range def.LocalVars {
		if v.Size == 0 {
			v.Size = v.Assign.Size()
		}
		if v.Size == 0 {
			return fmt.Errorf("can't determine var size '%v': '%v'", v.Name, def.Name)
		}
		v.Offset = curOffset
		v.SliceEnd = v.Offset + v.Size
		curOffset += v.Size
		if !v.IsParam {
			if v.Assign == nil {
				return fmt.Errorf("variable '%v' in '%v' is not assigned", v.Name, def.Name)
			}
		} else {
			def.InSize += v.Size
			def.ParamSizes = append(def.ParamSizes, v.Size)
		}
		if v.Assign != nil && v.Assign.Size() != v.Size {
			return fmt.Errorf("sizes doesn't match for var '%v' in '%v'", v.Name, def.Name)
		}
	}
	def.BufLen = curOffset

	if def.HasStateVariables {
		module.IncStat("numStatefulFunDef")
	}
	return nil
}

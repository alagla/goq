package types

type QuplaModule struct {
	Types     map[string]*QuplaTypeDef `yaml:"types"`
	Luts      map[string]*QuplaLutDef  `yaml:"luts"`
	Functions map[string]*QuplaFuncDef `yaml:"functions"`
	Execs     []*QuplaExecStmt         `yaml:"execs"`
}

func (module *QuplaModule) Analyze() bool {
	al := module.AnalyzeLuts()
	ae := module.AnalyzeExecs()
	af := module.AnalyzeFuncDefs()
	return al && ae && af
}

func (module *QuplaModule) AnalyzeFuncDefs() bool {
	infof("Analyzing function definitions...")
	var numErr int
	for name, fd := range module.Functions {
		fd.SetName(name)
		if err := fd.Analyze(module); err != nil {
			numErr++
			errorf("Error in function '%v': %v", name, err)
		}
	}
	infof("Number of function definitons found: %v", len(module.Functions))
	if numErr == 0 {
		infof("Done analyzing function definitions. No errors.")
	} else {
		errorf("Failed analyzing function definitions. Errors found: %v", numErr)
	}
	return numErr == 0
}

func (module *QuplaModule) AnalyzeLuts() bool {
	infof("Analyzing luts...")
	var numErr int
	for _, lutDef := range module.Luts {
		if err := lutDef.Analyze(module); err != nil {
			numErr++
			errorf("Error in lut '%v': %v", lutDef.LutTable, err)
		}
	}
	infof("Number of LUTs found: %v", len(module.Luts))
	if numErr == 0 {
		infof("Done analyzing LUTs. No errors.")
	} else {
		errorf("Failed analyzing LUTs. Errors found: %v", numErr)
	}
	return numErr == 0
}

func (module *QuplaModule) AnalyzeExecs() bool {
	infof("Analyzing execs...")
	var numTest, numEval, numErr int
	var err error
	for _, exec := range module.Execs {
		exec.expr, err = exec.ExprWrap.Unwarp()
		if err != nil {
			numErr++
			errorf("%v", err)
			continue
		}
		err = exec.expr.Analyze(module)
		if err != nil {
			numErr++
			errorf("%v", err)
			continue
		}
		exec.isTest = exec.ExpectedWrap != nil
		if exec.isTest {
			exec.exprExpected, err = exec.ExpectedWrap.Unwarp()
			if err != nil {
				numErr++
				errorf("%v", err)
				continue
			}
			err = exec.exprExpected.Analyze(module)
			if err != nil {
				numErr++
				errorf("%v", err)
				continue
			}
			numTest++
		} else {
			exec.exprExpected = nil
			numEval++
		}
	}
	infof("Found tests: %v, evals: %v", numTest, numEval)
	if numErr == 0 {
		infof("Done analyzing execs. No errors.")
	} else {
		errorf("Failed analyzing execs. Errors found: %v", numErr)
	}
	return numErr == 0
}

func (module *QuplaModule) FindFuncDef(name string) *QuplaFuncDef {
	ret, ok := module.Functions[name]
	if ok {
		return ret
	}
	return nil
}

func (module *QuplaModule) FindLUTDef(name string) *QuplaLutDef {
	ret, ok := module.Luts[name]
	if ok {
		return ret
	}
	return nil
}

func (module *QuplaModule) FindTypeDef(name string) *QuplaTypeDef {
	ret, ok := module.Types[name]
	if ok {
		return ret
	}
	return nil
}

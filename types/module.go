package types

import "fmt"

type QuplaModule struct {
	Types     map[string]*QuplaTypeDef `yaml:"types"`
	Luts      map[string]*QuplaLutDef  `yaml:"luts"`
	Functions map[string]*QuplaFuncDef `yaml:"functions"`
	Execs     []*QuplaExecStmt         `yaml:"execs"`
}

func (module *QuplaModule) Analyze() bool {
	return module.AnalyzeExecs()
}

func (module *QuplaModule) AnalyzeExecs() bool {
	infof("Analyzing execs...")
	var numTest, numEval, numErr int
	var err error
	for _, exec := range module.Execs {
		exec.expr, err = exec.ExprWrap.Analyze(module, nil)
		if err != nil {
			numErr++
			errorf("%v", err)
			continue
		}
		exec.isTest = exec.ExpectedWrap != nil
		if exec.isTest {
			exec.exprExpected, err = exec.ExpectedWrap.Analyze(module, nil)
			if err != nil {
				numErr++
				errorf("%v", err)
				continue
			}
			// check sizes
			if err = MatchSizes(exec.expr, exec.exprExpected); err != nil {
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

func (module *QuplaModule) FindFuncDef(name string) (*QuplaFuncDef, error) {
	var err error
	var fd *QuplaFuncDef
	ret, ok := module.Functions[name]
	if ok {
		ret.SetName(name)
		//if name == "fullAdd_3" {
		//	fmt.Printf("kuku")
		//}
		fd, err = ret.Analyze(module)
		if err != nil {
			return nil, err
		}
		module.Functions[name] = fd
		return module.Functions[name], nil
	}
	return nil, fmt.Errorf("can't find function definition '%v'", name)
}

func (module *QuplaModule) FindLUTDef(name string) (*QuplaLutDef, error) {
	var err error
	ret, ok := module.Luts[name]
	if ok {
		ret.SetName(name)
		module.Luts[name], err = ret.Analyze(module)
		if err != nil {
			return nil, err
		}
		return module.Luts[name], nil
	}
	return nil, fmt.Errorf("can't find LUT definition '%v'", name)
}

func (module *QuplaModule) FindTypeDef(name string) *QuplaTypeDef {
	ret, ok := module.Types[name]
	if ok {
		return ret
	}
	return nil
}

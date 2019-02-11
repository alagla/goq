package main

type QuplaModule struct {
	Types     map[string]*QuplaTypeDef `yaml:"types"`
	Luts      map[string][]string      `yaml:"luts"`
	Functions map[string]*QuplaFuncDef `yaml:"functions"`
	Execs     []*QuplaExecStmt         `yaml:"execs"`
}

func (module *QuplaModule) analyze() error {
	infof("Analysing Qupla module...")
	var numTest, numEval, numErr int
	var err error
	for _, exec := range module.Execs {
		err = exec.Expr.Analyze()
		if err != nil {
			numErr++
			errorf("%v", err)
		}
		exec.isTest = exec.Expected != nil
		if exec.isTest {
			err = exec.Expected.Analyze()
			if err != nil {
				numErr++
				errorf("%v", err)
			}
			numTest++
		} else {
			numEval++
		}
	}
	infof("Found tests: %v, evals: %v, errors: %v", numTest, numEval, numErr)
	infof("Done analyzing")
	return nil
}

package qupla

import (
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
)

type Function struct {
	Analyzed          bool // finished analysis
	Joins             map[string]int
	Affects           map[string]int
	Name              string
	retSize           int
	RetExpr           ExpressionInterface
	LocalVars         []*VarInfo
	NumParams         int // idx < NumParams represents parameter, idx >= represents local var (assign)
	BufLen            int // total length of the local var buffer
	HasStateVariables bool
	hasState          bool
	InSize            int
	ParamSizes        []int
}

func (def *Function) HasState() bool {
	return def.HasStateVariables || def.hasState
}

func (def *Function) References(funName string) bool {
	for _, vi := range def.LocalVars {
		if vi.Assign != nil && vi.Assign.References(funName) {
			return true
		}
	}
	return def.RetExpr.References(funName)
}

func NewQuplaFuncDef(name string, size int) *Function {
	return &Function{
		Name:       name,
		retSize:    size,
		LocalVars:  make([]*VarInfo, 0, 10),
		Joins:      make(map[string]int),
		Affects:    make(map[string]int),
		ParamSizes: make([]int, 0, 5),
	}
}

func (def *Function) Size() int {
	return def.retSize
}

func (def *Function) ArgSize() int {
	return def.InSize
}

func (def *Function) HasEnvStmt() bool {
	return len(def.Joins) > 0 || len(def.Affects) > 0
}

func (def *Function) GetJoinEnv() map[string]int {
	return def.Joins
}

func (def *Function) GetAffectEnv() map[string]int {
	return def.Affects
}

func (def *Function) GetVarIdx(name string) int {
	for i, lv := range def.LocalVars {
		if lv.Name == name {
			return i
		}
	}
	return -1
}

func (def *Function) VarByIdx(idx int) (*VarInfo, error) {
	if idx < 0 || idx >= len(def.LocalVars) {
		return nil, fmt.Errorf("worng var idx %v", idx)
	}
	return def.LocalVars[idx], nil
}

func (def *Function) VarByName(name string) (*VarInfo, error) {
	idx := def.GetVarIdx(name)
	if idx < 0 {
		return nil, fmt.Errorf("can't finc variabe with name '%v'", name)
	}
	return def.VarByIdx(idx)
}

func (def *Function) CheckArgSizes(args []ExpressionInterface) error {
	for i := range args {
		if i >= def.NumParams || args[i].Size() != def.LocalVars[i].Size {
			return fmt.Errorf("param and arg # %v mismach in %v", i, def.Name)
		}
	}
	return nil
}

// mock expression with all null arguments
func (def *Function) NewFuncExpressionWithNulls() *FunctionExpr {
	ret := NewFunctionExpr("", def)

	offset := 0
	for _, sz := range def.ParamSizes {
		ret.AppendSubExpr(NewNullExpr(sz))
		offset += sz
	}
	return ret
}

func (def *Function) Eval(frame *EvalFrame, result trinary.Trits) bool {
	return def.RetExpr.Eval(frame, result)
}

package qupla

import (
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type Function struct {
	Analyzed          bool // finished analysis
	Joins             map[string]int
	Affects           map[string]int
	Name              string
	retSize           int64
	RetExpr           ExpressionInterface
	LocalVars         []*VarInfo
	NumParams         int64 // idx < NumParams represents parameter, idx >= represents local var (assign)
	BufLen            int64 // total length of the local var buffer
	HasStateVariables bool
	hasState          bool
	InSize            int64
	ParamSizes        []int64
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

func NewQuplaFuncDef(name string, size int64) *Function {
	return &Function{
		Name:       name,
		retSize:    size,
		LocalVars:  make([]*VarInfo, 0, 10),
		Joins:      make(map[string]int),
		Affects:    make(map[string]int),
		ParamSizes: make([]int64, 0, 5),
	}
}

func (def *Function) Size() int64 {
	return def.retSize
}

func (def *Function) ArgSize() int64 {
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

func (def *Function) GetVarIdx(name string) int64 {
	for i, lv := range def.LocalVars {
		if lv.Name == name {
			return int64(i)
		}
	}
	return -1
}

func (def *Function) VarByIdx(idx int64) *VarInfo {
	if idx < 0 || idx >= int64(len(def.LocalVars)) {
		return nil
	}
	return def.LocalVars[idx]
}

func (def *Function) VarByName(name string) *VarInfo {
	return def.VarByIdx(def.GetVarIdx(name))
}

func (def *Function) GetVarInfo(name string) (*VarInfo, error) {
	ret := def.VarByName(name)
	if !ret.Analyzed {
		// can only be called after analysis is completed
		panic(fmt.Errorf("var '%v' is not analyzed in '%v'", name, def.Name))
	}
	return ret, nil
}

func (def *Function) CheckArgSizes(args []ExpressionInterface) error {
	for i := range args {
		if int64(i) >= def.NumParams || args[i].Size() != def.LocalVars[i].Size {
			return fmt.Errorf("param and arg # %v mismach in %v", i, def.Name)
		}
	}
	return nil
}

func (def *Function) NewExpressionWithArgs(args trinary.Trits) (ExpressionInterface, error) {
	if def.InSize != int64(len(args)) {
		return nil, fmt.Errorf("Size mismatch: fundef '%v' has arg Size %v, trit vector's Size = %v",
			def.Name, def.ArgSize(), len(args))
	}
	ret := NewQuplaFuncExpr("", def)

	offset := int64(0)
	for _, sz := range def.ParamSizes {
		e := NewQuplaValueExpr(args[offset : offset+sz])
		ret.AppendSubExpr(e)
		offset += sz
	}
	return ret, nil
}

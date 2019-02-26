package abstract

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type ExpressionFactory interface {
	AnalyzeExpression(interface{}, ModuleInterface, FuncDefInterface) (ExpressionInterface, error)
}

type ModuleInterface interface {
	AnalyzeExpression(interface{}, FuncDefInterface) (ExpressionInterface, error)
	AddFuncDef(string, FuncDefInterface)
	FindFuncDef(string) (FuncDefInterface, error)
	AddLutDef(string, LUTInterface)
	FindLUTDef(string) (LUTInterface, error)
	IncStat(string)
}

type VarInfo struct {
	Name     string
	Analyzed bool
	Idx      int64
	Offset   int64
	Size     int64
	IsState  bool
	IsParam  bool
	Expr     ExpressionInterface
}
type FuncDefInterface interface {
	GetName() string
	Size() int64
	GetVarInfo(string, ModuleInterface) (*VarInfo, error)
}

type LUTInterface interface {
	Size() int64
	Lookup(Trits, Trits) bool
}

type ExpressionInterface interface {
	Size() int64
	Eval(ProcessorInterface, Trits) bool
}

type ProcessorInterface interface {
	Push(ExpressionInterface)
	Pull()
	Eval(ExpressionInterface, Trits) bool
	EvalVar(int64) bool
	Slice(int64, int64) Trits
}

func MatchSizes(e1, e2 ExpressionInterface) error {
	s1 := e1.Size()
	s2 := e2.Size()

	if s1 != s2 {
		return fmt.Errorf("sizes doesn't match: %v != %v", s1, s2)
	}
	return nil
}

func RequireSize(e ExpressionInterface, size int64) error {
	s := e.Size()

	if s != size {
		return fmt.Errorf("sizes doesn't match: required %v != %v", size, s)
	}
	return nil
}

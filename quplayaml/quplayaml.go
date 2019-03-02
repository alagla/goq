package quplayaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type QuplaModuleYAML struct {
	Types     map[string]*QuplaTypeDefYAML `yaml:"types"`
	Luts      map[string]*QuplaLutDefYAML  `yaml:"luts"`
	Functions map[string]*QuplaFuncDefYAML `yaml:"functions"`
	Execs     []*QuplaExecStmtYAML         `yaml:"execs"`
}

type QuplaExecStmtYAML struct {
	Source   string               `yaml:"source"`
	Expr     *QuplaExpressionYAML `yaml:"expr"`
	Expected *QuplaExpressionYAML `yaml:"expected,omitempty"`
}

type QuplaTypeFieldYAML struct {
	Vector string `yaml:"vector"`
	Size   string `yaml:"size"`
}

type QuplaTypeDefYAML struct {
	Size   string                         `yaml:"size"`
	Fields map[string]*QuplaTypeFieldYAML `yaml:"fields,omitempty"`
}

type QuplaLutDefYAML struct {
	LutTable []string `yaml:"lutTable"`
}

type QuplaEnvStmtYAML struct {
	Name string `yaml:"name"`
	Join bool   `yaml:"join"`
}

type QuplaFuncArgYAML struct {
	ArgName string               `yaml:"argName"`
	Size    int64                `yaml:"size"`
	Type    *QuplaExpressionYAML `yaml:"type"` // not used
}

type QuplaStateVar struct {
	Size int64  `yaml:"size"`
	Type string `yaml:"type"`
}

type QuplaFuncDefYAML struct {
	Source     string                          `yaml:"source"`
	ReturnType *QuplaExpressionYAML            `yaml:"returnType"` // only size is necessary
	Params     []*QuplaFuncArgYAML             `yaml:"params"`
	State      map[string]*QuplaStateVar       `yaml:"state"`
	Env        []*QuplaEnvStmtYAML             `yaml:"env,omitempty"`
	Assigns    map[string]*QuplaExpressionYAML `yaml:"assigns,omitempty"`
	ReturnExpr *QuplaExpressionYAML            `yaml:"return"`
}

type QuplaExpressionYAML struct {
	CondExpr      *QuplaCondExprYAML      `yaml:"CondExpr,omitempty"`
	LutExpr       *QuplaLutExprYAML       `yaml:"LutExpr,omitempty"`
	SliceExpr     *QuplaSliceExprYAML     `yaml:"SliceExpr,omitempty"`
	ValueExpr     *QuplaValueExprYAML     `yaml:"ValueExpr,omitempty"`
	SizeofExpr    *QuplaSizeofExprYAML    `yaml:"SizeofExpr,omitempty"`
	FuncExpr      *QuplaFuncExprYAML      `yaml:"FuncExpr,omitempty"`
	FieldExpr     *QuplaFieldExprYAML     `yaml:"FieldExpr,omitempty"`
	ConstNumber   *QuplaConstNumberYAML   `yaml:"ConstNumber,omitempty"`
	ConstTypeName *QuplaConstTypeNameYAML `yaml:"ConstTypeName,omitempty"`
	ConstTerm     *QuplaConstTermYAML     `yaml:"ConstTerm,omitempty"`
	ConstExpr     *QuplaConstExprYAML     `yaml:"ConstExpr,omitempty"`
	ConcatExpr    *QuplaConcatExprYAML    `yaml:"ConcatExpr,omitempty"`
	MergeExpr     *QuplaMergeExprYAML     `yaml:"MergeExpr,omitempty"`
	TypeExpr      *QuplaTypeExprYAML      `yaml:"TypeExpr,omitempty"`
	NullExpr      *QuplaNullExprYAML      `yaml:"NullExpr,omitempty"`
}

type QuplaNullExprYAML string

type QuplaConcatExprYAML struct {
	Source string               `yaml:"source"`
	Lhs    *QuplaExpressionYAML `yaml:"lhs"`
	Rhs    *QuplaExpressionYAML `yaml:"rhs"`
}

type QuplaCondExprYAML struct {
	Source string               `yaml:"source"`
	If     *QuplaExpressionYAML `yaml:"if"`
	Then   *QuplaExpressionYAML `yaml:"then"`
	Else   *QuplaExpressionYAML `yaml:"else"`
}

type QuplaConstExprYAML struct {
	Operator string               `yaml:"operator"`
	Lhs      *QuplaExpressionYAML `yaml:"lhs"`
	Rhs      *QuplaExpressionYAML `yaml:"rhs"`
}

type QuplaConstTermYAML struct {
	Operator string               `yaml:"operator"`
	Lhs      *QuplaExpressionYAML `yaml:"lhs"`
	Rhs      *QuplaExpressionYAML `yaml:"rhs"`
}

type QuplaConstTypeNameYAML struct {
	TypeName   string `yaml:"typeName"` // not used
	SizeString string `yaml:"size"`
}

type QuplaConstNumberYAML struct {
	Value string `yaml:"value"`
}

type QuplaFieldExprYAML struct {
	FieldName string               `yaml:"fieldName"`
	CondExpr  *QuplaExpressionYAML `yaml:"condExpr"`
}

type QuplaLutExprYAML struct {
	Source string                 `yaml:"source"`
	Name   string                 `yaml:"name"`
	Args   []*QuplaExpressionYAML `yaml:"args"`
}

type QuplaSliceExprYAML struct {
	Source    string               `yaml:"source"`
	Var       string               `yaml:"var"`
	Offset    int64                `yaml:"offset"`
	SliceSize int64                `yaml:"size"`
	StartExpr *QuplaExpressionYAML `yaml:"start,omitempty"` // not used
	EndExpr   *QuplaExpressionYAML `yaml:"end,omitempty"`   // not used
}

type QuplaValueExprYAML struct {
	Value  string `yaml:"value"`
	Trits  string `yaml:"trits"`
	Trytes string `yaml:"trytes"`
}

type QuplaSizeofExprYAML struct {
	Value  string `yaml:"value"`
	Trits  string `yaml:"trits"`
	Trytes string `yaml:"trytes"`
}

type QuplaFuncExprYAML struct {
	Source string                 `yaml:"source"`
	Name   string                 `yaml:"name"`
	Args   []*QuplaExpressionYAML `yaml:"args"`
}

type QuplaMergeExprYAML struct {
	Source string               `yaml:"source"`
	Lhs    *QuplaExpressionYAML `yaml:"lhs"`
	Rhs    *QuplaExpressionYAML `yaml:"rhs"`
}

// ----- ?????? do we need it?
type QuplaTypeExprYAML struct {
	Source        string                          `yaml:"source"`
	TypeNameConst *QuplaExpressionYAML            `yaml:"type"`
	Fields        map[string]*QuplaExpressionYAML `yaml:"fields"`
}

func NewQuplaModuleFromYAML(fname string) (*QuplaModuleYAML, error) {
	yamlFile, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer yamlFile.Close()

	yamlbytes, _ := ioutil.ReadAll(yamlFile)

	ret := &QuplaModuleYAML{}
	err = yaml.Unmarshal(yamlbytes, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (module *QuplaModuleYAML) WriteToFile(fname string) error {
	outData, err := yaml.Marshal(&module)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fname, outData, 0644)
}

func (e *QuplaExpressionYAML) Unwrap() (interface{}, error) {
	var ret interface{}
	var numCases int

	if e.CondExpr != nil {
		ret = e.CondExpr
		numCases++
	}
	if e.LutExpr != nil {
		ret = e.LutExpr
		numCases++
	}
	if e.SliceExpr != nil {
		ret = e.SliceExpr
		numCases++
	}
	if e.ValueExpr != nil {
		ret = e.ValueExpr
		numCases++
	}
	if e.SizeofExpr != nil {
		ret = e.SizeofExpr
		numCases++
	}
	if e.FuncExpr != nil {
		ret = e.FuncExpr
		numCases++
	}
	if e.FieldExpr != nil {
		ret = e.FieldExpr
		numCases++
	}
	if e.ConstNumber != nil {
		ret = e.ConstNumber
		numCases++
	}
	if e.ConstTypeName != nil {
		ret = e.ConstTypeName
		numCases++
	}
	if e.ConstTerm != nil {
		ret = e.ConstTerm
		numCases++
	}
	if e.ConstExpr != nil {
		ret = e.ConstExpr
		numCases++
	}
	if e.ConcatExpr != nil {
		ret = e.ConcatExpr
		numCases++
	}
	if e.MergeExpr != nil {
		ret = e.MergeExpr
		numCases++
	}
	if e.TypeExpr != nil {
		ret = e.TypeExpr
		numCases++
	}
	if e.NullExpr != nil {
		ret = e.NullExpr
		numCases++
	}
	if numCases == 0 {
		return nil, nil // null
	}
	if numCases != 1 {
		return nil, fmt.Errorf("internal error: must be no more than one expression case")
	}
	return ret, nil
}

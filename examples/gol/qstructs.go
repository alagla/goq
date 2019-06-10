package main

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type QField struct {
	name string
	size int
}

type QStruct []QField

func (qs QStruct) GetIdx(name string) int {
	for i, s := range qs {
		if s.name == name {
			return i
		}
	}
	return -1
}

func (qs QStruct) Size() int {
	var ret int
	for _, s := range qs {
		ret += s.size
	}
	return ret
}

func (qs QStruct) ToTrits(fields map[string]Trits) Trits {
	ret := make(Trits, qs.Size())
	var cnt int
	for _, s := range qs {
		fld := fields[s.name] // can panic
		copy(ret[cnt:cnt+s.size], fld)
		cnt += s.size
	}
	return ret
}

func (qs QStruct) Parse(effect Trits) (map[string]Trits, error) {
	if len(effect) != qs.Size() {
		return nil, fmt.Errorf("size of the effect must be equal to %v. Got %", qs.Size(), len(effect))
	}
	ret := make(map[string]Trits)
	var cnt int
	for _, s := range qs {
		ret[s.name] = effect[cnt : cnt+s.size]
		cnt += s.size
	}
	return ret, nil
}

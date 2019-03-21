package dispatcher_tests

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"math"
	"testing"
)

// test of iota.go trinary arithmetics

var pairs = []struct {
	a int64
	b int64
}{
	{1, 1},
	{-1, 1},
	{0, 1},
	{0, -1},
	{-1, 0},
	{1, 0},
	{1000, 1},
	{1000, -1},
	{math.MaxInt64, -1},
	{math.MinInt64, math.MinInt64},
}

func TestAddTrinary(t *testing.T) {
	fmt.Printf("Testing trinary arithmetics..\n")
	for idx, pair := range pairs {
		ta := IntToTrits(pair.a)
		tb := IntToTrits(pair.b)
		texpected := IntToTrits(pair.a + pair.b)
		length := len(ta)
		if length < len(tb) {
			length = len(tb)
		}
		if length < len(texpected) {
			length = len(texpected)
		}
		ta = PadTrits(ta, length)
		tb = PadTrits(tb, length)
		texpected = PadTrits(texpected, length)

		tsum := AddTrits(ta, tb)
		eq, err := TritsEqual(tsum, texpected)
		switch {
		case err != nil:
			t.Errorf("failed #%v: %v + %v: %v", idx, pair.a, pair.b, err)
		case !eq:
			t.Errorf("failed #%v: %v + %v. Result '%v' != expected '%v'",
				idx, pair.a, pair.b, utils.TritsToString(tsum), utils.TritsToString(texpected))
		}
	}
}

const (
	startUpper    = int64(1000000000000)
	startLower    = int64(-100000)
	numRunInt     = int64(1000000000)
	printStepInt  = numRunInt / 10
	numRunTrit    = int64(1000000)
	printStepTrit = numRunTrit / 10
)

func TestAddIntSpeed(t *testing.T) {
	fmt.Printf("Testing int64 arithmetics..\n")
	var a, b, s int64
	start := utils.UnixMsNow()
	for i := int64(0); i < numRunInt; i++ {
		a = startLower + i
		b = startUpper - i
		s = a + b

		if i%printStepInt == 0 {
			fmt.Println(a, b, s)
		}
	}

	durationMs := int64(utils.UnixMsNow() - start)
	fmt.Printf("Duration: %v msec, %v cycles per second\n", durationMs, (numRunInt*1000)/durationMs)
}

func TestAddTritToInt(t *testing.T) {
	fmt.Printf("Testing trits -> int(sum(trits, trits)) == sum(int, int)..\n")
	var a, b, s int64
	var tsum Trits
	start := utils.UnixMsNow()
	for i := int64(0); i < numRunTrit; i++ {
		a = startLower + i
		b = startUpper - i
		ta := IntToTrits(a)
		tb := IntToTrits(b)
		if len(tb) > len(ta) {
			ta = PadTrits(ta, len(tb))
		} else {
			tb = PadTrits(tb, len(ta))
		}
		tsum = AddTrits(ta, tb)
		backToInt := TritsToInt(tsum)
		s = a + b
		if backToInt != s {
			t.Errorf("failed %v + %v: %v != %v", a, b, backToInt, s)
		}
		if i%printStepTrit == 0 {
			fmt.Printf("%v + %v -> '%v'\n", a, b, utils.TritsToString(tsum))
		}
	}
	durationMs := int64(utils.UnixMsNow() - start)
	fmt.Printf("Duration: %v msec, %v cycles per second\n", durationMs, (numRunTrit*1000)/durationMs)
}

func TestAddTritSpeed(t *testing.T) {
	fmt.Printf("Benchmarking sum(trits(int), trits(int)) -> trits ..\n")
	var a, b int64
	var tsum Trits
	start := utils.UnixMsNow()
	for i := int64(0); i < numRunTrit; i++ {
		a = startLower + i
		b = startUpper - i
		ta := IntToTrits(a)
		tb := IntToTrits(b)
		if len(tb) > len(ta) {
			ta = PadTrits(ta, len(tb))
		} else {
			tb = PadTrits(tb, len(ta))
		}
		tsum = AddTrits(ta, tb)
		if i%printStepTrit == 0 {
			fmt.Printf("%v + %v -> '%v'\n", a, b, utils.TritsToString(tsum))
		}
	}

	durationMs := int64(utils.UnixMsNow() - start)
	fmt.Printf("Duration: %v msec, %v cycles per second\n", durationMs, (numRunTrit*1000)/durationMs)
}

func TestAddTritToInt2(t *testing.T) {
	fmt.Printf("benchmarking trits(sum(int(trits), int(trits))) -> trits ..\n")
	var a, b, s int64
	var tsum Trits
	start := utils.UnixMsNow()
	for i := int64(0); i < numRunTrit; i++ {
		a = startLower + i
		b = startUpper - i
		ta := IntToTrits(a)
		tb := IntToTrits(b)
		if len(tb) > len(ta) {
			ta = PadTrits(ta, len(tb))
		} else {
			tb = PadTrits(tb, len(ta))
		}
		s = a + b
		tsum = IntToTrits(s)
		if i%printStepTrit == 0 {
			fmt.Printf("%v + %v -> '%v'\n", a, b, utils.TritsToString(tsum))
		}
	}
	durationMs := int64(utils.UnixMsNow() - start)
	fmt.Printf("Duration: %v msec, %v cycles per second\n", durationMs, (numRunTrit*1000)/durationMs)
}

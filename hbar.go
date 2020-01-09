package hedera

import "fmt"

type Hbar struct {
	tinybar int64
	unit    HbarUnit
}

const max = int64(^uint(0) >> 1)
const min = -max - 1

var MaxHbar = Hbar{max, HBar}
var MinHbar = Hbar{min, HBar}
var ZeroHbar = Hbar{0, HBar}

func HbarFromDecimal(bars float64, unit HbarUnit) Hbar {
	return Hbar{int64(bars * float64(unit.toTinybarCount())), unit}
}

func HbarFromTinybar(tinybar int64) Hbar {
	return Hbar{tinybar: tinybar, unit: Tinybar}
}

func HbarOf(hbar float64) Hbar {
	return Hbar{tinybar: int64(hbar * 100_000_000), unit: HBar}
}

func (hbar Hbar) AsTinybar() uint64 {
	return uint64(hbar.tinybar)
}

func (hbar Hbar) As(unit HbarUnit) int64 {
	return hbar.tinybar * hbar.unit.toTinybarCount()
}

func (hbar Hbar) String() string {
	if hbar.unit == Tinybar {
		return fmt.Sprintf("%v %v", hbar.tinybar, hbar.unit.String())
	} else {
		return fmt.Sprintf("%v %v (%v tinybar)", hbar.tinybar, hbar.unit.String(), hbar.As(Tinybar))
	}
}

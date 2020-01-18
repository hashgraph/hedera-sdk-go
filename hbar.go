package hedera

import "fmt"

type Hbar struct {
	tinybar int64
}

const max = int64(^uint(0) >> 1)
const min = -max - 1

var MaxHbar = Hbar{max}

var MinHbar = Hbar{min}

var ZeroHbar = Hbar{0}

// todo: change this behavior to wrap around?
// HbarFrom creates a representation of Hbar in tinybar on the unit provided
// note: if the value of bars is out of range it will return 0
func HbarFrom(bars float64, unit HbarUnit) Hbar {
	return HbarFromTinybar(int64(bars * float64(unit.numberOfTinybar())))
}

// todo: change this behavior to wrap around?
// HbarFromTinybar creates a representation of Hbar in tinybars
// note: if the value of tinybar is out of range it will return 0
func HbarFromTinybar(tinybar int64) Hbar {
	if tinybar > max {
		return ZeroHbar
	}

	if tinybar < min {
		return ZeroHbar
	}

	return Hbar{tinybar}
}

func NewHbar(hbar float64) Hbar {
	return HbarFromTinybar(int64(hbar * float64(HbarUnits.Tinybar.numberOfTinybar())))
}

func (hbar Hbar) AsTinybar() int64 {
	return hbar.tinybar
}

func (hbar Hbar) As(unit HbarUnit) int64 {
	return hbar.tinybar * unit.numberOfTinybar()
}

// todo: handle different unit types?
func (hbar Hbar) String() string {
	return fmt.Sprintf("%v tinybar", hbar.tinybar)
}

func (hbar Hbar) negated() Hbar {
	return Hbar{
		tinybar: -hbar.tinybar,
	}
}

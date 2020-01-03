package hedera

type Hbar struct {
	tinybar int64
}

const max = int64(^uint(0) >> 1)
const min = -max - 1

var HbarMAX = Hbar{max}
var HbarMIN = Hbar{min}
var HbarZERO = Hbar{}

func HbarFromTinybar(tinybar uint64) Hbar {
	return Hbar{tinybar: int64(tinybar)}
}

func HbarOf(tinybar float64) Hbar {
	return Hbar{tinybar: int64(tinybar * 100_000_000)}
}

func (hbar Hbar) AsTinybar() uint64 {
	return uint64(hbar.tinybar)
}

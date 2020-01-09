package hedera

type Hbar struct {
	tinybar int64
}

const max = int64(^uint(0) >> 1)
const min = -max - 1

var MaxHbar = Hbar{max}
var MinHbar = Hbar{min}
var ZeroHbar = Hbar{}

func HbarFromTinybar(tinybar int64) Hbar {
	return Hbar{tinybar: tinybar}
}

func HbarOf(hbar float64) Hbar {
	return Hbar{tinybar: int64(hbar * 100_000_000)}
}

func (hbar Hbar) AsTinybar() uint64 {
	return uint64(hbar.tinybar)
}

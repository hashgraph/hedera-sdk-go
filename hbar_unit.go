package hedera

type HbarUnit string

const (
	Tinybar  HbarUnit = "tinybar"
	Microbar HbarUnit = "microbar"
	Millibar HbarUnit = "millibar"
	HBar     HbarUnit = "hbar"
	Kilobar  HbarUnit = "kilobar"
	Megabar  HbarUnit = "megabar"
	Gigabar  HbarUnit = "gigabar"
)

func (unit HbarUnit) Symbol() string {
	switch unit {
	case Tinybar:
		return "tℏ"
	case Microbar:
		return "μℏ"
	case Millibar:
		return "mℏ"
	case HBar:
		return "ℏ"
	case Kilobar:
		return "kℏ"
	case Megabar:
		return "Mℏ"
	case Gigabar:
		return "Gℏ"
	}

	panic("HbarUnit.Symbol() switch statement is non-exhaustive")
}

func (unit HbarUnit) String() string {
	return string(unit)
}

func (unit HbarUnit) toTinybarCount() int64 {
	switch unit {
	case Tinybar:
		return 1
	case Microbar:
		return 100
	case Millibar:
		return 100_000
	case HBar:
		return 100_000_000
	case Kilobar:
		return 100_000_000_000
	case Megabar:
		return 100_000_000_000_000
	case Gigabar:
		return 100_000_000_000_000_000
	}

	panic("HbarUnit.Symbol() switch statement is non-exhaustive")
}

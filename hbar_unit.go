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

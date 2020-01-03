package hedera

type HbarUnit string

const (
	Tinybar  HbarUnit = "tinybar"
	Microbar HbarUnit = "microbar"
	Millibar HbarUnit = "millibar"
	// Hbar     HbarUnit = "hbar"
	Kilobar HbarUnit = "kilobar"
	Megabar HbarUnit = "megabar"
	Gigabar HbarUnit = "gigabar"
)

func (unit HbarUnit) GetSymbol() string {
	switch unit {
	case Tinybar:
		return "tℏ"
	case Microbar:
		return "μℏ"
	case Millibar:
		return "mℏ"
	// case Hbar:
	// 	return "ℏ"
	case Kilobar:
		return "kℏ"
	case Megabar:
		return "Mℏ"
	case Gigabar:
		return "Gℏ"
	}

	// Unreachable
	return ""
}

func (unit HbarUnit) String() string {
	return string(unit)
}

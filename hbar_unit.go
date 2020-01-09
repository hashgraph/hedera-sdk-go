package hedera

type HbarUnit string

var HbarUnits = struct {
	Tinybar  HbarUnit
	Microbar HbarUnit
	Millibar HbarUnit
	Hbar     HbarUnit
	Kilobar  HbarUnit
	Megabar  HbarUnit
	Gigabar  HbarUnit
}{
	Tinybar:  HbarUnit("tinybar"),
	Microbar: HbarUnit("microbar"),
	Millibar: HbarUnit("millibar"),
	Hbar:     HbarUnit("hbar"),
	Kilobar:  HbarUnit("kilobar"),
	Megabar:  HbarUnit("megabar"),
	Gigabar:  HbarUnit("gigabar"),
}

func (unit HbarUnit) Symbol() string {
	switch unit {
	case HbarUnits.Tinybar:
		return "tℏ"
	case HbarUnits.Microbar:
		return "μℏ"
	case HbarUnits.Millibar:
		return "mℏ"
	case HbarUnits.Hbar:
		return "ℏ"
	case HbarUnits.Kilobar:
		return "kℏ"
	case HbarUnits.Megabar:
		return "Mℏ"
	case HbarUnits.Gigabar:
		return "Gℏ"
	}

	panic("unreacahble: HbarUnit.Symbol() switch statement is non-exhaustive")
}

func (unit HbarUnit) String() string {
	return string(unit)
}

func (unit HbarUnit) numberOfTinybar() int64 {
	switch unit {
	case HbarUnits.Tinybar:
		return 1
	case HbarUnits.Microbar:
		return 100
	case HbarUnits.Millibar:
		return 100_000
	case HbarUnits.Hbar:
		return 100_000_000
	case HbarUnits.Kilobar:
		return 100_000_000_000
	case HbarUnits.Megabar:
		return 100_000_000_000_000
	case HbarUnits.Gigabar:
		return 100_000_000_000_000_000
	}

	panic("unreacahble: HbarUnit.Symbol() switch statement is non-exhaustive")
}

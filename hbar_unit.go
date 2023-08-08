package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

type HbarUnit string

// HbarUnits is a set of HbarUnit
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

// Symbol returns the symbol representation of the HbarUnit
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

	panic("unreachable: HbarUnit.Symbol() switch statement is non-exhaustive")
}

// String returns a string representation of the HbarUnit
func (unit HbarUnit) String() string {
	return string(unit)
}

func (unit HbarUnit) _NumberOfTinybar() int64 {
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

	panic("unreachable: HbarUnit.Symbol() switch statement is non-exhaustive")
}

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

// Hbar is a typesafe wrapper around values of HBAR providing foolproof conversions to other denominations.
type Hbar struct {
	tinybar int64
}

// MaxHbar is the maximum amount the Hbar type can wrap.
var MaxHbar = Hbar{math.MaxInt64}

// MinHbar is the minimum amount the Hbar type can wrap.
var MinHbar = Hbar{math.MinInt64}

// ZeroHbar wraps a 0 value of Hbar.
var ZeroHbar = Hbar{0}

// HbarFrom creates a representation of Hbar in tinybar on the unit provided
func HbarFrom(bars float64, unit HbarUnit) Hbar {
	return HbarFromTinybar(int64(bars * float64(unit._NumberOfTinybar())))
}

// HbarFromTinybar creates a representation of Hbar in tinybars
func HbarFromTinybar(tinybar int64) Hbar {
	return Hbar{tinybar}
}

// NewHbar constructs a new Hbar from a possibly fractional amount of hbar.
func NewHbar(hbar float64) Hbar {
	return HbarFrom(hbar, HbarUnits.Hbar)
}

// AsTinybar returns the equivalent tinybar amount.
func (hbar Hbar) AsTinybar() int64 {
	return hbar.tinybar
}

func (hbar Hbar) As(unit HbarUnit) float64 {
	return float64(hbar.tinybar) / float64(unit._NumberOfTinybar())
}

func (hbar Hbar) String() string {
	// Format the string as tinybar if the value is 1000 tinybar or less
	if -10000 <= hbar.tinybar && hbar.tinybar <= 10000 {
		return fmt.Sprintf("%v %s", hbar.tinybar, HbarUnits.Tinybar.Symbol())
	}

	return fmt.Sprintf("%v %s", float64(hbar.tinybar)/float64(HbarUnits.Hbar._NumberOfTinybar()), HbarUnits.Hbar.Symbol())
}

func HbarFromString(hbar string) (Hbar, error) {
	var err error
	match := regexp.MustCompile(`^((?:\+|\-)?\d+(?:\.\d+)?)(?: (tℏ|μℏ|mℏ|ℏ|kℏ|Mℏ|Gℏ))?$`)

	matchArray := match.FindStringSubmatch(hbar)
	if len(matchArray) == 0 {
		return Hbar{}, errors.New("invalid number and/or symbol")
	}

	a, err := strconv.ParseFloat(matchArray[1], 64)
	if err != nil {
		return Hbar{}, err
	}

	return HbarFrom(a, _HbarUnitFromString(matchArray[2])), nil
}

func _HbarUnitFromString(symbol string) HbarUnit {
	switch symbol {
	case HbarUnits.Tinybar.Symbol():
		return HbarUnits.Tinybar
	case HbarUnits.Microbar.Symbol():
		return HbarUnits.Microbar
	case HbarUnits.Millibar.Symbol():
		return HbarUnits.Millibar
	case HbarUnits.Kilobar.Symbol():
		return HbarUnits.Kilobar
	case HbarUnits.Megabar.Symbol():
		return HbarUnits.Megabar
	case HbarUnits.Gigabar.Symbol():
		return HbarUnits.Gigabar
	default:
		return HbarUnits.Hbar
	}
}

func (hbar Hbar) ToString(unit HbarUnit) string {
	return fmt.Sprintf("%v %v", float64(hbar.tinybar)/float64(unit._NumberOfTinybar()), unit.Symbol())
}

func (hbar Hbar) Negated() Hbar {
	return Hbar{
		tinybar: -hbar.tinybar,
	}
}

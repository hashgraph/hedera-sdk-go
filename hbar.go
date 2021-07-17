package hedera

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"regexp"
	"strconv"
	"strings"
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
	return HbarFromTinybar(int64(bars * float64(unit.numberOfTinybar())))
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
	return float64(hbar.tinybar) / float64(unit.numberOfTinybar())
}

func (hbar Hbar) String() string {
	// Format the string as tinybar if the value is 1000 tinybar or less
	if -1000 <= hbar.tinybar && hbar.tinybar <= 1000 {
		return fmt.Sprintf("%v tℏ", hbar.tinybar)
	}

	return fmt.Sprintf("%v ℏ", float64(hbar.tinybar)/float64(HbarUnits.Hbar.numberOfTinybar()))
}

func HbarFromString(hbar string) (Hbar, error) {
	match, err := regexp.Compile(`^((?:\+|\-)?\d+(?:\.\d+)?)(?: (tℏ|μℏ|mℏ|ℏ|kℏ|Mℏ|Gℏ))?$`)
	if err != nil {
		return Hbar{}, err
	}

	matchArray := match.FindStringSubmatch(hbar)
	if len(matchArray)== 0{
		return Hbar{}, errors.New("Invalid number and/or symbol.")
	}

	a, err := strconv.ParseFloat(matchArray[1], 64)
	if err != nil {
		return Hbar{}, err
	}

	if strings.Contains(hbar, "tℏ") {
		return HbarFrom(a, HbarUnits.Tinybar), nil
	} else if strings.Contains(hbar, "μℏ") {
		return HbarFrom(a, HbarUnits.Microbar), nil
	} else if strings.Contains(hbar, "mℏ") {
		return HbarFrom(a, HbarUnits.Millibar), nil
	} else if strings.Contains(hbar, "kℏ") {
		return HbarFrom(a, HbarUnits.Kilobar), nil
	} else if strings.Contains(hbar, "Mℏ") {
		return HbarFrom(a, HbarUnits.Megabar), nil
	} else if strings.Contains(hbar, "Gℏ") {
		return HbarFrom(a, HbarUnits.Gigabar), nil
	} else {
		return HbarFrom(a, HbarUnits.Hbar), nil
	}
}

func (hbar Hbar) ToString(unit HbarUnit) string {
	return fmt.Sprintf("%v %v", float64(hbar.tinybar)/float64(unit.numberOfTinybar()), unit.Symbol())
}

func (hbar Hbar) Negated() Hbar {
	return Hbar{
		tinybar: -hbar.tinybar,
	}
}

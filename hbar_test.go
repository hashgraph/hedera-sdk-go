package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHbarFromTinybar(t *testing.T) {
	tinybar := HbarUnits.Hbar.numberOfTinybar()

	hbar := HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)

	tinybar = MaxHbar.tinybar

	hbar = HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)

	tinybar = MinHbar.tinybar

	hbar = HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)
}

func TestHbarUnit(t *testing.T) {
	tinybar := HbarUnits.Kilobar.numberOfTinybar()

	hbar := HbarFromTinybar(tinybar)

	hbar2, err := HbarFromString(hbar.ToString(HbarUnits.Kilobar))
	assert.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	tinybar = HbarUnits.Gigabar.numberOfTinybar()

	hbar = HbarFromTinybar(tinybar)

	hbar2, err = HbarFromString(hbar.ToString(HbarUnits.Gigabar))
	assert.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	tinybar = HbarUnits.Microbar.numberOfTinybar()

	hbar = HbarFromTinybar(tinybar)

	hbar2, err = HbarFromString(hbar.ToString(HbarUnits.Microbar))
	assert.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	hbar2, err = HbarFromString("-5.123 Gℏ")
	assert.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, int64(-512300000000000000))

	hbar2, err = HbarFromString("5")
	assert.NoError(t, err)
	assert.Equal(t, hbar2.ToString(HbarUnits.Hbar), "5 ℏ")

	hbar2, err = HbarFromString("+5.123 ℏ")
	assert.NoError(t, err)
	assert.Equal(t, hbar2.ToString(HbarUnits.Millibar), "5123 mℏ")

	hbar2, err = HbarFromString("1.151 uℏ")
	assert.Error(t, err)

	hbar2, err = HbarFromString("1.151.")
	assert.Error(t, err)
}

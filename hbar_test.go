package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHbarFromTinybar(t *testing.T) {
	tinybar := HbarUnits.Hbar.numberOfTinybar();

	hbar := HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)

	tinybar = MaxHbar.tinybar

	hbar = HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)

	tinybar = MinHbar.tinybar

	hbar = HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)
}

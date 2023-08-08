//go:build all || unit
// +build all unit

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

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitHbarFromTinybar(t *testing.T) {
	t.Parallel()

	tinybar := HbarUnits.Hbar._NumberOfTinybar()

	hbar := HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)

	tinybar = MaxHbar.tinybar

	hbar = HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)

	tinybar = MinHbar.tinybar

	hbar = HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)
}

func TestUnitHbarUnit(t *testing.T) {
	t.Parallel()

	tinybar := HbarUnits.Kilobar._NumberOfTinybar()

	hbar := HbarFromTinybar(tinybar)

	hbar2, err := HbarFromString(hbar.ToString(HbarUnits.Kilobar))
	require.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	tinybar = HbarUnits.Gigabar._NumberOfTinybar()

	hbar = HbarFromTinybar(tinybar)

	hbar2, err = HbarFromString(hbar.ToString(HbarUnits.Gigabar))
	require.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	tinybar = HbarUnits.Microbar._NumberOfTinybar()

	hbar = HbarFromTinybar(tinybar)

	hbar2, err = HbarFromString(hbar.ToString(HbarUnits.Microbar))
	require.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	hbar2, err = HbarFromString("-5.123 Gℏ")
	require.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, int64(-512300000000000000))

	hbar2, err = HbarFromString("5")
	require.NoError(t, err)
	assert.Equal(t, hbar2.ToString(HbarUnits.Hbar), "5 ℏ")

	hbar2, err = HbarFromString("+5.123 ℏ")
	require.NoError(t, err)
	assert.Equal(t, hbar2.ToString(HbarUnits.Millibar), "5123 mℏ")

	hbar2, err = HbarFromString("1.151 uℏ")
	assert.Error(t, err)

	hbar2, err = HbarFromString("1.151.")
	assert.Error(t, err)
}

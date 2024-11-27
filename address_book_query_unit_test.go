//go:build all || unit
// +build all unit

package hiero

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestUnitAddressBookQueryCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	file := FileID{File: 3, checksum: &checksum}

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewAddressBookQuery().
		SetFileID(file).
		SetLimit(3).
		SetMaxAttempts(4)

	err = query.validateNetworkOnIDs(client)

	require.NoError(t, err)
	query.GetFileID()
	query.GetLimit()
	query.GetFileID()
}

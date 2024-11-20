//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// testinf func (transfer TokenNftTransfer) ToBytes() and func NftTransferFromBytes(data []byte)
func TestUnitTokenNftTransferToBytes(t *testing.T) {
	t.Parallel()

	transfer := _TokenNftTransfer{
		SenderAccountID:   AccountID{Account: 3},
		ReceiverAccountID: AccountID{Account: 4},
		SerialNumber:      5,
		IsApproved:        true,
	}

	transferBytes := transfer.ToBytes()
	transferFromBytes, err := NftTransferFromBytes(transferBytes)

	assert.NoError(t, err)
	assert.Equal(t, transfer, transferFromBytes)

	// test invalid data from and to bytes
	_, err = NftTransferFromBytes([]byte{1, 2, 3})
	assert.Error(t, err)

	// test nil data from bytes and to bytes
	_, err = NftTransferFromBytes(nil)
	assert.Error(t, err)
}

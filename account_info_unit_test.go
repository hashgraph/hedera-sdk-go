//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The function checks the conversation methods on the AccountInfo struct. We check wether it is correctly converted to protobuf and back.
func TestUnitAccountInfoToBytes(t *testing.T) {
	t.Parallel()

	accInfoOriginal := *_MockAccountInfo()
	accInfoBytes := accInfoOriginal.ToBytes()

	accInfoFromBytes, err := AccountInfoFromBytes(accInfoBytes)

	require.NoError(t, err)
	require.Equal(t, accInfoOriginal.AccountID, accInfoFromBytes.AccountID)
	require.Equal(t, accInfoOriginal.ContractAccountID, accInfoFromBytes.ContractAccountID)
	require.Equal(t, accInfoOriginal.Key, accInfoFromBytes.Key)
	require.Equal(t, accInfoOriginal.LedgerID, accInfoFromBytes.LedgerID)
}
func _MockAccountInfo() *AccountInfo {
	privateKey, _ := PrivateKeyFromString(mockPrivateKey)
	accountID, _ := AccountIDFromString("0.0.123-esxsf")
	accountID.checksum = nil

	return &AccountInfo{
		AccountID:                      accountID,
		ContractAccountID:              "",
		IsDeleted:                      false,
		ProxyReceived:                  Hbar{},
		Key:                            privateKey.PublicKey(),
		Balance:                        HbarFromTinybar(100000000000),
		GenerateSendRecordThreshold:    Hbar{},
		GenerateReceiveRecordThreshold: Hbar{},
		ReceiverSigRequired:            false,
		ExpirationTime:                 time.Date(2222, 2, 2, 2, 2, 2, 2, time.Now().UTC().Location()),
		AutoRenewPeriod:                time.Duration(time.Duration(5).Seconds()),
		LiveHashes:                     nil,
		AccountMemo:                    "",
		OwnedNfts:                      0,
		MaxAutomaticTokenAssociations:  0,
		AliasKey:                       nil,
		LedgerID:                       *NewLedgerIDTestnet(),
		EthereumNonce:                  0,
		StakingInfo:                    nil,
	}

}

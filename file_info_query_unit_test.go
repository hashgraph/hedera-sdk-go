//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitFileInfoQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	fileInfo := NewFileInfoQuery().
		SetFileID(fileID)

	err = fileInfo.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitFileInfoQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	fileInfo := NewFileInfoQuery().
		SetFileID(fileID)

	err = fileInfo.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitFileInfoQueryMock(t *testing.T) {
	t.Parallel()

	newKey, err := PrivateKeyFromStringEd25519("302e020100300506032b657004220420a869f4c6191b9c8c99933e7f6b6611711737e4b1a1a5a4cb5370e719a1f6df98")
	require.NoError(t, err)
	key := newKey.PublicKey().BytesRaw()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_FileGetInfo{
				FileGetInfo: &services.FileGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_FileGetInfo{
				FileGetInfo: &services.FileGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_FileGetInfo{
				FileGetInfo: &services.FileGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					FileInfo: &services.FileGetInfoResponse_FileInfo{
						FileID:         &services.FileID{FileNum: 3},
						Size:           10,
						ExpirationTime: nil,
						Deleted:        false,
						Keys: &services.KeyList{
							Keys: []*services.Key{
								{
									Key: &services.Key_Ed25519{
										Ed25519: key,
									},
								},
							},
						},
						Memo:     "no memo",
						LedgerId: []byte{0},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewFileInfoQuery().
		SetFileID(FileID{File: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, cost, HbarFromTinybar(2))

	result, err := query.Execute(client)
	require.NoError(t, err)

	require.Equal(t, result.Keys.keys[0].String(), newKey.PublicKey().String())
	require.Equal(t, result.FileMemo, "no memo")
	require.Equal(t, result.IsDeleted, false)
	require.True(t, result.LedgerID.IsMainnet())
}

func TestUnitFileInfoQueryGet(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	fileID := FileID{File: 3, checksum: &checksum}
	deadline := time.Duration(time.Minute)
	accountId := AccountID{Account: 123}
	transactionID := TransactionIDGenerate(accountId)
	query := NewFileInfoQuery().
		SetFileID(fileID).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}}).
		SetGrpcDeadline(&deadline).
		SetMaxBackoff(1 * time.Minute).
		SetMinBackoff(500 * time.Millisecond).
		SetMaxRetry(5).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(500))
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	err = query.validateNetworkOnIDs(client)
	require.NoError(t, err)
	require.Equal(t, fileID, query.GetFileID())
	require.Equal(t, []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}, query.GetNodeAccountIDs())
	require.Equal(t, 500*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 1*time.Minute, query.GetMaxBackoff())
	require.Equal(t, 5, query.GetMaxRetryCount())
	require.Equal(t, transactionID, query.GetPaymentTransactionID())
	require.Equal(t, HbarFromTinybar(25), query.GetQueryPayment())
	require.Equal(t, NewHbar(500), query.GetMaxQueryPayment())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitFileInfoQuerySetNothing(t *testing.T) {
	t.Parallel()

	balance := NewFileInfoQuery()

	require.Equal(t, FileID{}, balance.GetFileID())
	require.Equal(t, []AccountID{}, balance.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, balance.GetMinBackoff())
	require.Equal(t, 8*time.Second, balance.GetMaxBackoff())
	require.Equal(t, 10, balance.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, balance.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, balance.GetQueryPayment())
	require.Equal(t, Hbar{}, balance.GetMaxQueryPayment())
}

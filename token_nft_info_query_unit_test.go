//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenNftGetInfoByNftIDValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	nftID, err := NftIDFromString("2@0.0.123-esxsf")
	require.NoError(t, err)

	nftInfo := NewTokenNftInfoQuery().
		SetNftID(nftID)

	err = nftInfo.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenNftGetInfoByNftIDValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	nftID, err := NftIDFromString("2@0.0.123-rmkykd")
	require.NoError(t, err)

	nftInfo := NewTokenNftInfoQuery().
		SetNftID(nftID)

	err = nftInfo.validateNetworkOnIDs(client)
	require.Error(t, err)
	if err != nil {
		require.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenNftInfoQueryNothingSet(t *testing.T) {
	t.Parallel()

	query := NewTokenNftInfoQuery()

	require.Equal(t, NftID{TokenID: TokenID{Shard: 0x0, Realm: 0x0, Token: 0x0, checksum: (*string)(nil)}, SerialNumber: 0}, query.GetNftID())
	require.Equal(t, []AccountID{}, query.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 8*time.Second, query.GetMaxBackoff())
	require.Equal(t, 10, query.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, query.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, query.GetQueryPayment())
	require.Equal(t, Hbar{}, query.GetMaxQueryPayment())
}

func TestUnitTokenNftInfoQueryGet(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	deadline := time.Second * 3
	token := TokenID{Token: 3, checksum: &checksum}
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(account)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewTokenNftInfoQuery().
		SetTokenID(token).
		SetNftID(token.Nft(334)).
		SetAccountID(account).
		SetEnd(5).
		SetStart(4).
		ByAccountID(account).
		ByTokenID(token).
		ByNftID(token.Nft(334)).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetNodeAccountIDs(nodeAccountID).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(23)).
		SetQueryPayment(NewHbar(3)).
		SetGrpcDeadline(&deadline)

	err = query.validateNetworkOnIDs(client)
	require.NoError(t, err)

	// Some assertions like SetStart, SetEnd, etc. are missing, because those fucntions are deprecated and empty
	require.Equal(t, token.Nft(334).String(), query.GetNftID().String())
	require.Equal(t, token.Nft(334).String(), query.GetNftID().String())
	require.Equal(t, nodeAccountID, nodeAccountID, query.GetNodeAccountIDs())
	require.Equal(t, time.Second*30, query.GetMaxBackoff())
	require.Equal(t, time.Second*10, query.GetMinBackoff())
	require.Equal(t, NewHbar(3), query.GetQueryPayment())
	require.Equal(t, NewHbar(23), query.GetMaxQueryPayment())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitTokenNftInfoQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_TokenGetNftInfo{
				TokenGetNftInfo: &services.TokenGetNftInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TokenGetNftInfo{
				TokenGetNftInfo: &services.TokenGetNftInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TokenGetNftInfo{
				TokenGetNftInfo: &services.TokenGetNftInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					Nft: &services.TokenNftInfo{
						NftID:        nil,
						AccountID:    nil,
						CreationTime: nil,
						Metadata:     nil,
						LedgerId:     nil,
						SpenderId:    nil,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	checksum := "dmqui"
	token := TokenID{Token: 3, checksum: &checksum}

	query := NewTokenNftInfoQuery().
		SetNftID(token.Nft(43)).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, HbarFromTinybar(2), cost)
	_, err = query.Execute(client)
	require.NoError(t, err)
}

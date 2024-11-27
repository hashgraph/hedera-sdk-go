//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenInfoQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tokenInfo := NewTokenInfoQuery().
		SetTokenID(tokenID)

	err = tokenInfo.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenInfoQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenInfo := NewTokenInfoQuery().
		SetTokenID(tokenID)

	err = tokenInfo.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenInfoFromBytesBadBytes(t *testing.T) {
	t.Parallel()

	bytes, err := base64.StdEncoding.DecodeString("tfhyY++/Q4BycortAgD4cmMKACB/")
	require.NoError(t, err)

	_, err = TokenInfoFromBytes(bytes)
	require.NoError(t, err)
}

func TestUnitTokenInfoFromBytesEmptyBytes(t *testing.T) {
	t.Parallel()

	_, err := TokenInfoFromBytes([]byte{})
	require.NoError(t, err)
}

func TestUnitTokenInfoQueryGet(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 7}
	deadline := time.Duration(time.Minute)
	accountId := AccountID{Account: 123}
	validStart := time.Now().Add(10 * time.Minute)
	balance := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}}).
		SetMaxRetry(3).
		SetMinBackoff(300 * time.Millisecond).
		SetMaxBackoff(10 * time.Second).
		SetPaymentTransactionID(TransactionID{AccountID: &accountId, ValidStart: &validStart}).
		SetMaxQueryPayment(NewHbar(500)).
		SetGrpcDeadline(&deadline)

	require.Equal(t, tokenID, balance.GetTokenID())
	require.Equal(t, []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}, balance.GetNodeAccountIDs())
	require.Equal(t, 300*time.Millisecond, balance.GetMinBackoff())
	require.Equal(t, 10*time.Second, balance.GetMaxBackoff())
	require.Equal(t, 3, balance.GetMaxRetryCount())
	require.Equal(t, TransactionID{AccountID: &accountId, ValidStart: &validStart}, balance.GetPaymentTransactionID())
	require.Equal(t, HbarFromTinybar(25), balance.GetQueryPayment())
	require.Equal(t, NewHbar(500), balance.GetMaxQueryPayment())
}

func TestUnitTokenInfoQueryNothingSet(t *testing.T) {
	t.Parallel()

	balance := NewTokenInfoQuery()

	require.Equal(t, TokenID{}, balance.GetTokenID())
	require.Equal(t, []AccountID{}, balance.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, balance.GetMinBackoff())
	require.Equal(t, 8*time.Second, balance.GetMaxBackoff())
	require.Equal(t, 10, balance.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, balance.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, balance.GetQueryPayment())
	require.Equal(t, Hbar{}, balance.GetMaxQueryPayment())
}

func TestUnitTokenInfoQueryCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	deadline := time.Second * 3
	token := TokenID{Token: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewTokenInfoQuery().
		SetTokenID(token).
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

	require.Equal(t, nodeAccountID, query.GetNodeAccountIDs())
	require.Equal(t, time.Second*30, query.GetMaxBackoff())
	require.Equal(t, time.Second*10, query.GetMinBackoff())
	require.Equal(t, token, query.GetTokenID())
	require.Equal(t, NewHbar(3), query.GetQueryPayment())
	require.Equal(t, NewHbar(23), query.GetMaxQueryPayment())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitTokenInfoQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_TokenGetInfo{
				TokenGetInfo: &services.TokenGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TokenGetInfo{
				TokenGetInfo: &services.TokenGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TokenGetInfo{
				TokenGetInfo: &services.TokenGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					TokenInfo: &services.TokenInfo{
						TokenId:             nil,
						Name:                "",
						Symbol:              "",
						Decimals:            0,
						TotalSupply:         0,
						Treasury:            nil,
						AdminKey:            nil,
						KycKey:              nil,
						FreezeKey:           nil,
						WipeKey:             nil,
						SupplyKey:           nil,
						DefaultFreezeStatus: 0,
						DefaultKycStatus:    0,
						Deleted:             false,
						AutoRenewAccount:    nil,
						AutoRenewPeriod:     nil,
						Expiry:              nil,
						Memo:                "",
						TokenType:           0,
						SupplyType:          0,
						MaxSupply:           0,
						FeeScheduleKey:      nil,
						CustomFees:          nil,
						PauseKey:            nil,
						PauseStatus:         0,
						MetadataKey:         nil,
						LedgerId:            nil,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewTokenInfoQuery().
		SetTokenID(TokenID{Token: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, HbarFromTinybar(2), cost)
	_, err = query.Execute(client)
	require.NoError(t, err)
}

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

func TestUnitAccountRecordQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	recordQuery := NewAccountRecordsQuery().
		SetAccountID(accountID)

	err = recordQuery.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountRecordQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	recordQuery := NewAccountRecordsQuery().
		SetAccountID(accountID)

	err = recordQuery.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountRecordsQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_CryptoGetAccountRecords{
				CryptoGetAccountRecords: &services.CryptoGetAccountRecordsResponse{
					Header:    &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
					AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 1800}},
				},
			},
		},
		&services.Response{
			Response: &services.Response_CryptoGetAccountRecords{
				CryptoGetAccountRecords: &services.CryptoGetAccountRecordsResponse{
					Header:    &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
					AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 1800}},
				},
			},
		},
		&services.Response{
			Response: &services.Response_CryptoGetAccountRecords{
				CryptoGetAccountRecords: &services.CryptoGetAccountRecordsResponse{
					Header:    &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 1},
					AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 1800}},
					Records: []*services.TransactionRecord{
						{
							TransactionHash:    []byte{1},
							ConsensusTimestamp: &services.Timestamp{Nanos: 12313123, Seconds: 2313},
							TransactionID: &services.TransactionID{
								TransactionValidStart: &services.Timestamp{Nanos: 12313123, Seconds: 2313},
								AccountID:             &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 1800}},
								Scheduled:             false,
								Nonce:                 0,
							},
							Memo:           "",
							TransactionFee: 0,
						},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAccountID(AccountID{Account: 1800}).
		SetMaxQueryPayment(NewHbar(1))

	_, err := query.GetCost(client)
	require.NoError(t, err)
	recordsQuery, err := query.Execute(client)
	require.NoError(t, err)

	require.Equal(t, len(recordsQuery), 1)
	require.Equal(t, recordsQuery[0].TransactionID.AccountID.Account, uint64(1800))
}

func TestUnitAccountRecordsQueryGet(t *testing.T) {
	t.Parallel()

	spenderAccountID1 := AccountID{Account: 7}

	balance := NewAccountRecordsQuery().
		SetAccountID(spenderAccountID1).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(10)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetAccountID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitAccountRecordsQuerySetNothing(t *testing.T) {
	t.Parallel()

	balance := NewAccountRecordsQuery()

	balance.GetAccountID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitAccountRecordsQueryCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 3
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewAccountRecordsQuery().
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetAccountID(account).
		SetNodeAccountIDs(nodeAccountID).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(23)).
		SetQueryPayment(NewHbar(3)).
		SetGrpcDeadline(&grpc)

	err = query.validateNetworkOnIDs(client)

	require.NoError(t, err)
	query.GetNodeAccountIDs()
	query.GetMaxBackoff()
	query.GetMinBackoff()
	query.getName()
	query.GetAccountID()
}

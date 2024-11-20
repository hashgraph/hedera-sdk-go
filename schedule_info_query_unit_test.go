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

func TestUnitScheduleInfoQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	scheduleInfo := NewScheduleInfoQuery().
		SetScheduleID(scheduleID)

	err = scheduleInfo.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleInfoQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	scheduleInfo := NewScheduleInfoQuery().
		SetScheduleID(scheduleID)

	err = scheduleInfo.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitScheduleInfoQueryGet(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	accountId := AccountID{Account: 123}
	deadline := time.Duration(time.Minute)
	validStart := time.Now().Add(10 * time.Minute)
	scheduleID := ScheduleID{Schedule: 3, checksum: &checksum}

	query := NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
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
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	err = query.validateNetworkOnIDs(client)
	require.NoError(t, err)
	require.Equal(t, scheduleID, query.GetScheduleID())
	require.Equal(t, []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}, query.GetNodeAccountIDs())
	require.Equal(t, 300*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 10*time.Second, query.GetMaxBackoff())
	require.Equal(t, 3, query.GetMaxRetryCount())
	require.Equal(t, TransactionID{AccountID: &AccountID{Account: 123}, ValidStart: &validStart}, query.GetPaymentTransactionID())
	require.Equal(t, HbarFromTinybar(25), query.GetQueryPayment())
	require.Equal(t, NewHbar(500), query.GetMaxQueryPayment())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitScheduleInfoQuerySetNothing(t *testing.T) {
	t.Parallel()

	info := NewScheduleInfoQuery()

	require.Equal(t, ScheduleID{}, info.GetScheduleID())
	require.Equal(t, []AccountID{}, info.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, info.GetMinBackoff())
	require.Equal(t, 8*time.Second, info.GetMaxBackoff())
	require.Equal(t, 10, info.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, info.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, info.GetQueryPayment())
	require.Equal(t, Hbar{}, info.GetMaxQueryPayment())
}

func TestUnitScheduleInfoQueryCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	deadline := time.Second * 3
	schedule := ScheduleID{Schedule: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewScheduleInfoQuery().
		SetScheduleID(schedule).
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
	require.Equal(t, 30*time.Second, query.GetMaxBackoff())
	require.Equal(t, 10*time.Second, query.GetMinBackoff())
	require.NotEmpty(t, query.getName())
	require.Equal(t, schedule, query.GetScheduleID())
	require.Equal(t, NewHbar(3), query.GetQueryPayment())
	require.Equal(t, NewHbar(23), query.GetMaxQueryPayment())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitScheduleInfoQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_ScheduleGetInfo{
				ScheduleGetInfo: &services.ScheduleGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ScheduleGetInfo{
				ScheduleGetInfo: &services.ScheduleGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ScheduleGetInfo{
				ScheduleGetInfo: &services.ScheduleGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					ScheduleInfo: &services.ScheduleInfo{
						ScheduleID:               nil,
						Data:                     nil,
						ExpirationTime:           nil,
						ScheduledTransactionBody: nil,
						Memo:                     "",
						AdminKey:                 nil,
						Signers:                  nil,
						CreatorAccountID:         nil,
						PayerAccountID:           nil,
						ScheduledTransactionID:   nil,
						LedgerId:                 nil,
						WaitForExpiry:            false,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewScheduleInfoQuery().
		SetScheduleID(ScheduleID{Schedule: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, HbarFromTinybar(2), cost)
	_, err = query.Execute(client)
	require.NoError(t, err)
}

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

func TestUnitTopicInfoQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	topicID, err := TopicIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	err = topicID.Validate(client)
	require.NoError(t, err)
	topicInfo := NewTopicInfoQuery().
		SetTopicID(topicID)

	err = topicInfo.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTopicInfoQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	topicID, err := TopicIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	err = topicID.Validate(client)
	require.Error(t, err)
	topicInfo := NewTopicInfoQuery().
		SetTopicID(topicID)

	err = topicInfo.validateNetworkOnIDs(client)
	require.Error(t, err)
	if err != nil {
		require.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTopicInfoQueryGet(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	topicID := TopicID{Topic: 3, checksum: &checksum}
	deadline := time.Duration(time.Minute)
	accountId := AccountID{Account: 123}
	transactionID := TransactionIDGenerate(accountId)
	query := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}}).
		SetMaxRetry(3).
		SetMinBackoff(300 * time.Millisecond).
		SetMaxBackoff(10 * time.Second).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(500)).
		SetGrpcDeadline(&deadline)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	err = query.validateNetworkOnIDs(client)
	require.NoError(t, err)
	require.Equal(t, topicID, query.GetTopicID())
	require.Equal(t, []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}, query.GetNodeAccountIDs())
	require.Equal(t, 300*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 10*time.Second, query.GetMaxBackoff())
	require.Equal(t, 3, query.GetMaxRetryCount())
	require.Equal(t, transactionID, query.GetPaymentTransactionID())
	require.Equal(t, HbarFromTinybar(25), query.GetQueryPayment())
	require.Equal(t, NewHbar(500), query.GetMaxQueryPayment())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitTopicInfoQueryNothingSet(t *testing.T) {
	t.Parallel()

	query := NewTopicInfoQuery()

	require.Equal(t, TopicID{}, query.GetTopicID())
	require.Equal(t, []AccountID{}, query.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 8*time.Second, query.GetMaxBackoff())
	require.Equal(t, 10, query.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, query.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, query.GetQueryPayment())
	require.Equal(t, Hbar{}, query.GetMaxQueryPayment())
}

func TestUnitTopicInfoQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_ConsensusGetTopicInfo{
				ConsensusGetTopicInfo: &services.ConsensusGetTopicInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ConsensusGetTopicInfo{
				ConsensusGetTopicInfo: &services.ConsensusGetTopicInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ConsensusGetTopicInfo{
				ConsensusGetTopicInfo: &services.ConsensusGetTopicInfoResponse{
					Header:  &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					TopicID: nil,
					TopicInfo: &services.ConsensusTopicInfo{
						Memo:             "",
						RunningHash:      nil,
						SequenceNumber:   0,
						ExpirationTime:   nil,
						AdminKey:         nil,
						SubmitKey:        nil,
						AutoRenewPeriod:  nil,
						AutoRenewAccount: nil,
						LedgerId:         nil,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	checksum := "dmqui"
	topic := TopicID{Topic: 3, checksum: &checksum}

	query := NewTopicInfoQuery().
		SetTopicID(topic).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, HbarFromTinybar(2), cost)
	_, err = query.Execute(client)
	require.NoError(t, err)
}

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

func TestUnitNetworkVersionInfoQuerySetNothing(t *testing.T) {
	t.Parallel()

	query := NewNetworkVersionQuery()

	require.Equal(t, []AccountID{}, query.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 8*time.Second, query.GetMaxBackoff())
	require.Equal(t, 10, query.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, query.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, query.GetQueryPayment())
	require.Equal(t, Hbar{}, query.GetMaxQueryPayment())
}

func TestNetworkVersionInfoQuery_Get(t *testing.T) {
	t.Parallel()

	deadline := time.Duration(time.Minute)
	transactionID := TransactionIDGenerate(AccountID{Account: 324})
	query := NewNetworkVersionQuery().
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(10)).
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetMaxRetry(5).
		SetMaxBackoff(10 * time.Second).
		SetMinBackoff(1 * time.Second).
		SetPaymentTransactionID(transactionID).
		SetGrpcDeadline(&deadline)

	require.Equal(t, NewHbar(2), query.GetQueryPayment())
	require.Equal(t, NewHbar(10), query.GetMaxQueryPayment())
	require.Equal(t, []AccountID{{Account: 3}, {Account: 4}}, query.GetNodeAccountIDs())
	require.Equal(t, 5, query.GetMaxRetryCount())
	require.Equal(t, 10*time.Second, query.GetMaxBackoff())
	require.Equal(t, 1*time.Second, query.GetMinBackoff())
	require.Equal(t, transactionID, query.GetPaymentTransactionID())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitNetworkVersionInfoQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_NetworkGetVersionInfo{
				NetworkGetVersionInfo: &services.NetworkGetVersionInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewNetworkVersionQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, HbarFromTinybar(2), cost)
}

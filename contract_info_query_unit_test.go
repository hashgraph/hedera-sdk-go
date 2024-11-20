//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitContractInfoQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	contractInfoQuery := NewContractInfoQuery().
		SetContractID(contractID)

	err = contractInfoQuery.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitContractInfoQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	contractInfoQuery := NewContractInfoQuery().
		SetContractID(contractID)

	err = contractInfoQuery.validateNetworkOnIDs(client)
	require.Error(t, err)
	if err != nil {
		require.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitContractInfoQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_ContractGetInfo{
				ContractGetInfo: &services.ContractGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ContractGetInfo{
				ContractGetInfo: &services.ContractGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ContractGetInfo{
				ContractGetInfo: &services.ContractGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					ContractInfo: &services.ContractGetInfoResponse_ContractInfo{
						ContractID:         &services.ContractID{Contract: &services.ContractID_ContractNum{ContractNum: 3}},
						AccountID:          &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 4}},
						ContractAccountID:  "",
						AdminKey:           nil,
						ExpirationTime:     nil,
						AutoRenewPeriod:    nil,
						Storage:            0,
						Memo:               "yes",
						Balance:            0,
						Deleted:            false,
						TokenRelationships: nil,
						LedgerId:           nil,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewContractInfoQuery().
		SetContractID(ContractID{Contract: 3}).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{{Account: 3}})

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, cost, HbarFromTinybar(2))

	result, err := query.Execute(client)
	require.NoError(t, err)

	require.Equal(t, result.ContractID.Contract, uint64(3))
	require.Equal(t, result.AccountID.Account, uint64(4))
	require.Equal(t, result.ContractMemo, "yes")
}

func TestUnitContractInfoQueryGetTransactionIDMock(t *testing.T) {
	t.Skip("Skipping test as it is currently broken with the addition of generating new payment transactions for queries")
	t.Parallel()

	transactionID := TransactionIDGenerate(AccountID{Account: 123})
	call := func(request *services.Query) *services.Response {
		if query, ok := request.Query.(*services.Query_ContractGetInfo); ok {
			paymentTransacction := query.ContractGetInfo.Header.Payment

			require.NotEmpty(t, paymentTransacction.BodyBytes)
			transactionBody := services.TransactionBody{}
			_ = protobuf.Unmarshal(paymentTransacction.BodyBytes, &transactionBody)

			require.NotNil(t, transactionBody.TransactionID)
			tempTransactionID := _TransactionIDFromProtobuf(transactionBody.TransactionID)
			require.Equal(t, transactionID.String(), tempTransactionID.String())
		}

		return &services.Response{
			Response: &services.Response_ContractGetInfo{
				ContractGetInfo: &services.ContractGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					ContractInfo: &services.ContractGetInfoResponse_ContractInfo{
						ContractID:         &services.ContractID{Contract: &services.ContractID_ContractNum{ContractNum: 3}},
						AccountID:          &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 4}},
						ContractAccountID:  "",
						AdminKey:           nil,
						ExpirationTime:     nil,
						AutoRenewPeriod:    nil,
						Storage:            0,
						Memo:               "yes",
						Balance:            0,
						Deleted:            false,
						TokenRelationships: nil,
						LedgerId:           nil,
					},
				},
			},
		}
	}
	responses := [][]interface{}{{
		call,
	}}

	client, server := NewMockClientAndServer(responses)

	result, err := NewContractInfoQuery().
		SetContractID(ContractID{Contract: 3}).
		SetMaxQueryPayment(NewHbar(1)).
		SetPaymentTransactionID(transactionID).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, result.ContractID.Contract, uint64(3))
	require.Equal(t, result.AccountID.Account, uint64(4))
	require.Equal(t, result.ContractMemo, "yes")

	server.Close()
}

func TestUnitContractInfoQueryGet(t *testing.T) {
	t.Parallel()

	spenderContractID := ContractID{Contract: 7}

	balance := NewContractInfoQuery().
		SetContractID(spenderContractID).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	require.Equal(t, spenderContractID, balance.GetContractID())
	require.Equal(t, []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}, balance.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, balance.GetMinBackoff())
	require.Equal(t, 8*time.Second, balance.GetMaxBackoff())
	require.Equal(t, 10, balance.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, balance.GetPaymentTransactionID())
	require.Equal(t, HbarFromTinybar(25), balance.GetQueryPayment())
	require.Equal(t, NewHbar(1), balance.GetMaxQueryPayment())

}

func TestUnitContractInfoQueryCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	deadline := time.Second * 3
	contract := ContractID{Contract: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewContractInfoQuery().
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetContractID(contract).
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
	require.Equal(t, contract, query.GetContractID())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

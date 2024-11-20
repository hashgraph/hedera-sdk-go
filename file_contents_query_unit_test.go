//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitFileContentsQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(fileID)

	err = fileContents.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitFileContentsQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(fileID)

	err = fileContents.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitFileContentsQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_FileGetContents{
				FileGetContents: &services.FileGetContentsResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 3},
				},
			},
		},
		&services.Response{
			Response: &services.Response_FileGetContents{
				FileGetContents: &services.FileGetContentsResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 3},
				},
			},
		},
		&services.Response{
			Response: &services.Response_FileGetContents{
				FileGetContents: &services.FileGetContentsResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					FileContents: &services.FileGetContentsResponse_FileContents{
						FileID:   &services.FileID{FileNum: 3},
						Contents: []byte{123},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewFileContentsQuery().
		SetFileID(FileID{File: 3}).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{{Account: 3}})

	_, err := query.GetCost(client)
	require.NoError(t, err)

	result, err := query.Execute(client)
	require.NoError(t, err)

	require.Equal(t, bytes.Compare(result, []byte{123}), 0)
}

func TestUnitFileContentsQueryGet(t *testing.T) {
	t.Parallel()

	fileID := FileID{File: 7}

	balance := NewFileContentsQuery().
		SetFileID(fileID).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetFileID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitFileContentsQuerySetNothing(t *testing.T) {
	t.Parallel()

	balance := NewFileContentsQuery()

	balance.GetFileID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitFileContentsQueryCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 3
	file := FileID{File: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewFileContentsQuery().
		SetFileID(file).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
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
	query.GetFileID()
	query.GetQueryPayment()
	query.GetMaxQueryPayment()
}

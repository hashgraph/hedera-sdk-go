//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitNodeCreateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tx := NewNodeCreateTransaction().
		SetAccountID(accountID)

	err = tx.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitNodeCreateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tx := NewNodeCreateTransaction().
		SetAccountID(accountID)

	err = tx.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func endpoints(offsets ...uint) []Endpoint {
	endpoints := make([]Endpoint, 0)

	for _, offset := range offsets {
		endpoints = append(endpoints, Endpoint{
			address: []byte{byte(offset), byte(offset), byte(offset), byte(offset)},
		})
	}

	return endpoints
}

func TestUnitNodeCreateTransactionMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		status.New(codes.Unavailable, "node is UNAVAILABLE").Err(),
		status.New(codes.Internal, "Received RST_STREAM with code 0").Err(),
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_COST_ANSWER,
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status: services.ResponseCodeEnum_RECEIPT_NOT_FOUND,
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status: services.ResponseCodeEnum_SUCCESS,
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
							AccountNum: 234,
						}},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	resp, err := NewNodeCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetAdminKey(newKey).
		SetDescription("test").
		SetGossipEndpoints(endpoints(0, 1, 2)).
		SetServiceEndpoints(endpoints(3, 4, 5)).
		SetGossipCaCertificate([]byte{111}).
		SetGrpcCertificateHash([]byte{222}).
		SetTransactionID(tran).
		Execute(client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, receipt.AccountID, &AccountID{Account: 234})
}

func TestUnitNodeCreateTransactionGet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	key, err := PrivateKeyGenerateEd25519()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(key).
		SetTransactionMemo("").
		SetDescription("test").
		SetGossipEndpoints(endpoints(0, 1, 2)).
		SetServiceEndpoints(endpoints(3, 4, 5)).
		SetGossipCaCertificate([]byte{111}).
		SetGrpcCertificateHash([]byte{222}).
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetAccountID()
	transaction.GetDescription()
	transaction.GetGossipEndpoints()
	transaction.GetServiceEndpoints()
	transaction.GetGossipCaCertificate()
	transaction.GetGrpcCertificateHash()
	transaction.GetAdminKey()
}

func TestUnitNodeCreateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetAccountID()
	transaction.GetDescription()
	transaction.GetGossipEndpoints()
	transaction.GetServiceEndpoints()
	transaction.GetGossipCaCertificate()
	transaction.GetGrpcCertificateHash()
	transaction.GetAdminKey()
}

func TestUnitNodeCreateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	stackedAccountID := AccountID{Account: 5}

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	gossipEndpoints := endpoints(1, 2, 3)
	serviceEndpoints := endpoints(3, 4, 5)
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAccountID(stackedAccountID).
		SetAdminKey(key).
		SetTransactionMemo("").
		SetDescription("test").
		SetGossipEndpoints(gossipEndpoints).
		SetServiceEndpoints(serviceEndpoints).
		SetGossipCaCertificate([]byte{111}).
		SetGrpcCertificateHash([]byte{222}).
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetNodeCreate()
	require.Equal(t, proto.AccountId.String(), stackedAccountID._ToProtobuf().String())
	require.Equal(t, proto.Description, "test")
	require.Equal(t, proto.GossipEndpoint[0], gossipEndpoints[0]._ToProtobuf())
	require.Equal(t, proto.ServiceEndpoint[0], serviceEndpoints[0]._ToProtobuf())
	require.Equal(t, proto.GossipCaCertificate, []byte{111})
	require.Equal(t, proto.GrpcCertificateHash, []byte{222})
	require.Equal(t, proto.AdminKey, key._ToProtoKey())
}

func TestUnitNodeCreateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	trx, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(key).
		SetAccountID(account).
		SetMaxTransactionFee(NewHbar(3)).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetTransactionMemo("no").
		SetTransactionValidDuration(time.Second * 30).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	trx.validateNetworkOnIDs(client)
	_, err = trx.Schedule()
	require.NoError(t, err)
	trx.GetTransactionID()
	trx.GetNodeAccountIDs()
	trx.GetMaxRetry()
	trx.GetMaxTransactionFee()
	trx.GetMaxBackoff()
	trx.GetMinBackoff()
	trx.GetRegenerateTransactionID()
	byt, err := trx.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)
	sig, err := key.SignTransaction(trx)
	require.NoError(t, err)

	_, err = trx.GetTransactionHash()
	require.NoError(t, err)
	trx.GetMaxTransactionFee()
	trx.GetTransactionMemo()
	trx.GetRegenerateTransactionID()
	trx.GetAccountID()
	trx.GetDescription()
	trx.GetGossipEndpoints()
	trx.GetServiceEndpoints()
	trx.GetGossipCaCertificate()
	trx.GetGrpcCertificateHash()
	trx.GetAdminKey()
	_, err = trx.GetSignatures()
	require.NoError(t, err)
	trx.getName()
	switch b := txFromBytes.(type) {
	case NodeCreateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

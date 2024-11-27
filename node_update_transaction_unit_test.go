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

func TestUnitNodeUpdateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tx := NewNodeUpdateTransaction().
		SetAccountID(accountID)

	err = tx.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitNodeUpdateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tx := NewNodeUpdateTransaction().
		SetAccountID(accountID)

	err = tx.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitNodeUpdateTransactionMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
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
						Status: services.ResponseCodeEnum_SUCCESS,
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
							AccountNum: 234,
						}},
						NodeId: 1,
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

	resp, err := NewNodeUpdateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetAdminKey(newKey).
		SetNodeID(1).
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
	require.Equal(t, uint64(1), receipt.NodeID)
}

func TestUnitNodeUpdateTransactionGet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	key, err := PrivateKeyGenerateEd25519()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetNodeID(1).
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
	transaction.GetNodeID()
}

func TestUnitNodeUpdateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeUpdateTransaction().
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
	transaction.GetNodeID()
}

func TestUnitNodeUpdateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	stackedAccountID := AccountID{Account: 5}

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	gossipEndpoints := endpoints(1, 2, 3)
	serviceEndpoints := endpoints(3, 4, 5)
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeID(1).
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

	proto := transaction.build().GetNodeUpdate()
	require.Equal(t, proto.AccountId.String(), stackedAccountID._ToProtobuf().String())
	require.Equal(t, proto.Description.Value, "test")
	require.Equal(t, proto.GossipEndpoint[0], gossipEndpoints[0]._ToProtobuf())
	require.Equal(t, proto.ServiceEndpoint[0], serviceEndpoints[0]._ToProtobuf())
	require.Equal(t, proto.GossipCaCertificate.Value, []byte{111})
	require.Equal(t, proto.GrpcCertificateHash.Value, []byte{222})
	require.Equal(t, proto.AdminKey, key._ToProtoKey())
	require.Equal(t, proto.NodeId, uint64(1))
}

func TestUnitNodeUpdateTransactionCoverage(t *testing.T) {
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

	trx, err := NewNodeUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(key).
		SetNodeID(1).
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
	trx.GetNodeID()
	_, err = trx.GetSignatures()
	require.NoError(t, err)
	trx.getName()
	switch b := txFromBytes.(type) {
	case NodeUpdateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

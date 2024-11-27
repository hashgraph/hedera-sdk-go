//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountDeleteTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	deleteAccount := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(accountID)

	err = deleteAccount.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountDeleteTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	deleteAccount := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(accountID)

	err = deleteAccount.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountDeleteTransactionGet(t *testing.T) {
	t.Parallel()

	spenderAccountID1 := AccountID{Account: 7}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTransferAccountID(spenderAccountID1).
		SetAccountID(spenderAccountID1).
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
	transaction.GetTransferAccountID()
	transaction.GetAccountID()
}

func TestUnitAccountDeleteTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountDeleteTransaction().
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
	transaction.GetTransferAccountID()
	transaction.GetAccountID()
}

func TestUnitAccountDeleteTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	spenderAccountID1 := AccountID{Account: 7}
	transferAccountID := AccountID{Account: 8}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTransferAccountID(transferAccountID).
		SetAccountID(spenderAccountID1).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetCryptoDelete()
	require.Equal(t, proto.TransferAccountID.String(), transferAccountID._ToProtobuf().String())
	require.Equal(t, proto.DeleteAccountID.String(), spenderAccountID1._ToProtobuf().String())
}

func TestUnitAccountDeleteTransactionTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 3
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewAccountDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAccountID(account).
		SetTransferAccountID(account).
		SetGrpcDeadline(&grpc).
		SetMaxTransactionFee(NewHbar(3)).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetTransactionMemo("no").
		SetTransactionValidDuration(time.Second * 30).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.validateNetworkOnIDs(client)

	_, err = transaction.Schedule()
	require.NoError(t, err)
	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()
	transaction.GetMaxRetry()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxBackoff()
	transaction.GetMinBackoff()
	transaction.GetRegenerateTransactionID()
	byt, err := transaction.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)
	sig, err := newKey.SignTransaction(transaction)
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)
	transaction.GetAccountID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetTransferAccountID()
	transaction.getName()
	switch b := txFromBytes.(type) {
	case AccountDeleteTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitAccountDeleteTransactionTransactionMock(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 3}}
	transactionID := TransactionIDGenerate(AccountID{Account: 3})

	call := func(request *services.Transaction) *services.TransactionResponse {
		require.NotEmpty(t, request.SignedTransactionBytes)
		signedTransaction := services.SignedTransaction{}
		_ = protobuf.Unmarshal(request.SignedTransactionBytes, &signedTransaction)

		require.NotEmpty(t, signedTransaction.BodyBytes)
		transactionBody := services.TransactionBody{}
		_ = protobuf.Unmarshal(signedTransaction.BodyBytes, &transactionBody)

		require.NotNil(t, transactionBody.TransactionID)
		transactionId := transactionBody.TransactionID.String()
		require.NotEqual(t, "", transactionId)

		sigMap := signedTransaction.GetSigMap()
		require.NotNil(t, sigMap)

		return &services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		}
	}
	responses := [][]interface{}{{
		call,
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	_, err := NewAccountDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAccountID(AccountID{Account: 3}).
		Execute(client)
	require.NoError(t, err)
}

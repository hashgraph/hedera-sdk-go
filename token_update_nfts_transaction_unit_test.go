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
	protobuf "google.golang.org/protobuf/proto"
)

func TestUnitTokenUpdateNftsTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tokenUpdate := NewTokenUpdateNftsTransaction().
		SetTokenID(tokenID).
		SetSerialNumbers([]int64{1, 2, 3}).
		SetMetadata([]byte("metadata"))

	err = tokenUpdate.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenUpdateNftsTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenUpdate := NewTokenUpdateNftsTransaction().
		SetTokenID(tokenID).
		SetSerialNumbers([]int64{1, 2, 3}).
		SetMetadata([]byte("metadata"))

	err = tokenUpdate.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenUpdateNftsTransactionGet(t *testing.T) {
	t.Parallel()

	grpc := time.Second * 30
	checksum := "dmqui"
	tokenID := TokenID{Token: 3, checksum: &checksum}

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewTokenUpdateNftsTransaction().
		SetTokenID(tokenID).
		SetSerialNumbers([]int64{1, 2, 3}).
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetMetadata([]byte("metadata")).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		SetGrpcDeadline(&grpc).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetMaxTransactionFee(NewHbar(10)).
		SetLogLevel(LoggerLevelInfo).
		Freeze()

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
	_, err = transaction.Schedule()
	require.NoError(t, err)

	require.NotNil(t, transaction.GetTokenID())
	require.NotNil(t, transaction.GetSerialNumbers())
	require.NotNil(t, transaction.GetMetadata())
	require.NotNil(t, transaction.GetMaxTransactionFee())
	require.NotNil(t, transaction.GetTransactionMemo())
	require.NotNil(t, transaction.GetRegenerateTransactionID())
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	require.NotNil(t, transaction.GetRegenerateTransactionID())
	require.NotNil(t, transaction.GetMaxTransactionFee())
	require.NotNil(t, transaction.GetRegenerateTransactionID())
	byt, err := transaction.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)
	sig, err := newKey.SignTransaction(transaction)
	require.NoError(t, err)
	require.NotNil(t, transaction.getName())
	require.NotNil(t, transaction.GetMaxRetry())
	require.NotNil(t, transaction.GetMaxBackoff())
	require.NotNil(t, transaction.GetMinBackoff())
	switch b := txFromBytes.(type) {
	case TokenUpdateNfts:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitTokenUpdateNftsTransactionNothingSet(t *testing.T) {
	t.Parallel()
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenUpdateNftsTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).Freeze()

	require.NoError(t, err)
	require.Nil(t, transaction.GetTokenID())
	require.Nil(t, transaction.GetSerialNumbers())
	require.Nil(t, transaction.GetMetadata())
	require.NotNil(t, transaction.GetMaxTransactionFee())
	require.NotNil(t, transaction.GetTransactionMemo())
	require.NotNil(t, transaction.GetRegenerateTransactionID())
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	require.NotNil(t, transaction.GetRegenerateTransactionID())
	require.NotNil(t, transaction.GetMaxTransactionFee())
	require.NotNil(t, transaction.GetRegenerateTransactionID())
	require.NotNil(t, transaction.getName())
	require.NotNil(t, transaction.GetMaxRetry())
	require.NotNil(t, transaction.GetMaxBackoff())
	require.NotNil(t, transaction.GetMinBackoff())
	proto := transaction.build().GetTokenUpdateNfts()
	require.Nil(t, proto.Metadata)
	require.Nil(t, proto.SerialNumbers)
	require.Nil(t, proto.Token)
}

func TestUnitTokenUpdateNftTransactionMock(t *testing.T) {
	t.Parallel()

	newKey, err := PrivateKeyFromStringEd25519("302e020100300506032b657004220420a869f4c6191b9c8c99933e7f6b6611711737e4b1a1a5a4cb5370e719a1f6df98")
	require.NoError(t, err)

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

	checksum := "dmqui"
	token := TokenID{Token: 3, checksum: &checksum}

	freeze, err := NewTokenUpdateNftsTransaction().
		SetTokenID(token).
		SetSerialNumbers([]int64{1, 2, 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMetadata([]byte("metadata")).FreezeWith(client)
	require.NoError(t, err)

	_, err = freeze.Sign(newKey).Execute(client)
	require.NoError(t, err)

}

func TestUnitTokenUpdateNftsTransactionSignWithOperator(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	tokenID, err := TokenIDFromString("0.0.123")
	require.NoError(t, err)

	tokenUpdate := NewTokenUpdateNftsTransaction().
		SetTokenID(tokenID).
		SetSerialNumbers([]int64{1, 2, 3}).
		SetMetadata([]byte("metadata"))

	_, err = tokenUpdate.SignWithOperator(client)
	require.NoError(t, err)
}

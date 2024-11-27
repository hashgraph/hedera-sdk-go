//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenMintTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tokenMint := NewTokenMintTransaction().
		SetTokenID(tokenID)

	err = tokenMint.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenMintTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenMint := NewTokenMintTransaction().
		SetTokenID(tokenID)

	err = tokenMint.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenMintTransactionGet(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 7}

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenMintTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenID(tokenID).
		SetAmount(123).
		SetMetadata([]byte{123, 123}).
		SetMetadatas([][]byte{{123, 123}, {13, 13}}).
		SetMaxTransactionFee(NewHbar(10)).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetTokenID()
	transaction.GetAmount()
	transaction.GetMetadatas()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTokenMintTransactionNothingSet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenMintTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetTokenID()
	transaction.GetAmount()
	transaction.GetMetadatas()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTokenMintTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 7}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenMintTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenID(tokenID).
		SetAmount(323).
		SetMetadatas([][]byte{{50}, {50}}).
		SetMaxTransactionFee(NewHbar(10)).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetTokenMint()
	require.Equal(t, proto.Token.String(), tokenID._ToProtobuf().String())
	require.Equal(t, proto.Amount, uint64(323))
	require.Equal(t, proto.Metadata, [][]byte{{50}, {50}})
}

func TestUnitTokenMintTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	token := TokenID{Token: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewTokenMintTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenID(token).
		SetMetadata([]byte{1}).
		SetMetadatas([][]byte{{1, 0}}).
		SetAmount(5).
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

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
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
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	transaction.GetTokenID()
	transaction.GetMetadatas()
	transaction.GetAmount()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case TokenMintTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitTokenMintTransactionMock(t *testing.T) {
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

	freez, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTokenID(token).
		SetMetadata([]byte{1}).
		SetAmount(5).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}

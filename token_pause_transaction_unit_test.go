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

func TestUnitTokenPause(t *testing.T) {
	t.Parallel()

	accountID, err := AccountIDFromString("0.0.5005")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.5005")
	require.NoError(t, err)

	tx, err := NewTokenPauseTransaction().
		SetNodeAccountIDs([]AccountID{accountID}).
		SetTransactionID(TransactionIDGenerate(accountID)).
		SetTokenID(tokenID).
		Freeze()
	require.NoError(t, err)

	pb := tx.build()
	require.Equal(t, pb.GetTokenPause().GetToken().String(), tokenID._ToProtobuf().String())
}

func TestUnitTokenUnpause(t *testing.T) {
	t.Parallel()

	accountID, err := AccountIDFromString("0.0.5005")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.5005")
	require.NoError(t, err)

	tx, err := NewTokenUnpauseTransaction().
		SetNodeAccountIDs([]AccountID{accountID}).
		SetTransactionID(TransactionIDGenerate(accountID)).
		SetTokenID(tokenID).
		Freeze()
	require.NoError(t, err)

	pb := tx.build()
	require.Equal(t, pb.GetTokenUnpause().GetToken().String(), tokenID._ToProtobuf().String())
}
func TestUnitTokenPauseSchedule(t *testing.T) {
	accountID, err := AccountIDFromString("0.0.5005")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.5005")
	require.NoError(t, err)

	tx, err := NewTokenPauseTransaction().
		SetNodeAccountIDs([]AccountID{accountID}).
		SetTransactionID(TransactionIDGenerate(accountID)).
		SetTokenID(tokenID).
		Schedule()
	require.NoError(t, err)

	scheduled := tx.schedulableBody.GetTokenPause()
	require.Equal(t, scheduled.Token.String(), tokenID._ToProtobuf().String())
}

func TestUnitTokenPauseTransactionGet(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 7}

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenPauseTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenID(tokenID).
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
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTokenPauseTransactionNothingSet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenPauseTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetTokenID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTokenUnpauseTransactionGet(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 7}

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenUnpauseTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenID(tokenID).
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
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTokenUnpauseTransactionNothingSet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenUnpauseTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetTokenID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTokenUnpauseTransactionCoverage(t *testing.T) {
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

	transaction, err := NewTokenUnpauseTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenID(token).
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
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case TokenUnpauseTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitTokenUnpauseTransactionMock(t *testing.T) {
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

	freez, err := NewTokenUnpauseTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTokenID(token).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}

func TestUnitTokenPauseTransaction_SetMaxRetry(t *testing.T) {
	t.Parallel()
	transaction := NewTokenPauseTransaction()
	transaction.SetMaxRetry(5)

	require.Equal(t, 5, transaction.GetMaxRetry())
}

func TestUnitTokenPauseTransaction_AddSignature(t *testing.T) {
	t.Parallel()
	client, _ := ClientFromConfig([]byte(testClientJSONWithoutMirrorNetwork))

	nodeAccountId, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)

	nodeIdList := []AccountID{nodeAccountId}

	transaction, err := NewTokenPauseTransaction().
		SetNodeAccountIDs(nodeIdList).
		FreezeWith(client)

	privateKey, _ := PrivateKeyGenerateEd25519()

	signature, err := privateKey.SignTransaction(transaction)
	require.NoError(t, err)

	signs, err := transaction.GetSignatures()
	for key := range signs[nodeAccountId] {
		require.Equal(t, signs[nodeAccountId][key], signature)
	}

	require.NoError(t, err)

	privateKey2, _ := PrivateKeyGenerateEd25519()
	publicKey2 := privateKey2.PublicKey()

	signedTransaction := transaction.AddSignature(publicKey2, signature)
	signs2, err := signedTransaction.GetSignatures()
	require.NoError(t, err)

	for key := range signs2[nodeAccountId] {
		require.Equal(t, signs2[nodeAccountId][key], signature)
	}
}

func TestUnitTokenPauseTransaction_SignWithOperator(t *testing.T) {
	t.Parallel()
	client, _ := ClientFromConfig([]byte(testClientJSONWithoutMirrorNetwork))
	privateKey, _ := PrivateKeyGenerateEd25519()
	publicKey := privateKey.PublicKey()
	operatorId, _ := AccountIDFromString("0.0.10")

	client.SetOperator(operatorId, privateKey)

	nodeAccountId, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	nodeIdList := []AccountID{nodeAccountId}

	transaction, err := NewTokenPauseTransaction().
		SetNodeAccountIDs(nodeIdList).
		SetTokenID(TokenID{Token: 3}).
		SignWithOperator(client)

	require.NoError(t, err)
	require.NotNil(t, transaction)

	privateKey2, _ := PrivateKeyGenerateEd25519()
	publicKey2 := privateKey2.PublicKey()
	client.SetOperator(operatorId, privateKey2)

	transactionSignedWithOp, err := transaction.SignWithOperator(client)
	require.NoError(t, err)
	require.NotNil(t, transactionSignedWithOp)

	assert.Contains(t, transactionSignedWithOp.Transaction.publicKeys, publicKey)
	assert.Contains(t, transactionSignedWithOp.Transaction.publicKeys, publicKey2)

	// test errors
	client.operator = nil
	tx, err := NewTokenPauseTransaction().
		SetNodeAccountIDs(nodeIdList).
		SetTokenID(TokenID{Token: 3}).
		SignWithOperator(client)

	require.Error(t, err)
	require.Nil(t, tx)

	client = nil
	tx, err = NewTokenPauseTransaction().
		SetNodeAccountIDs(nodeIdList).
		SetTokenID(TokenID{Token: 3}).
		SignWithOperator(client)

	require.Error(t, err)
	require.Nil(t, tx)
}

func TestUnitTokenPauseTransaction_SetMaxBackoff(t *testing.T) {
	t.Parallel()
	transaction := NewTokenPauseTransaction()
	maxBackoff := 10 * time.Second

	transaction.SetMaxBackoff(maxBackoff)

	require.Equal(t, maxBackoff, transaction.GetMaxBackoff())

	// test max.Nanoseconds() < 0
	transaction2 := NewTokenPauseTransaction()
	maxBackoff2 := -1 * time.Second

	require.Panics(t, func() { transaction2.SetMaxBackoff(maxBackoff2) })

	// test max.Nanoseconds() < min.Nanoseconds()
	transaction3 := NewTokenPauseTransaction()
	maxBackoff3 := 1 * time.Second
	minBackoff3 := 2 * time.Second

	transaction3.SetMinBackoff(minBackoff3)

	require.Panics(t, func() { transaction3.SetMaxBackoff(maxBackoff3) })
}

func TestUnitTokenPauseTransaction_GetMaxBackoff(t *testing.T) {
	t.Parallel()
	transaction := NewTokenPauseTransaction()

	require.Equal(t, 8*time.Second, transaction.GetMaxBackoff())
}

func TestUnitTokenPauseTransaction_SetMinBackoff(t *testing.T) {
	t.Parallel()
	transaction := NewTokenPauseTransaction()
	minBackoff := 1 * time.Second

	transaction.SetMinBackoff(minBackoff)

	require.Equal(t, minBackoff, transaction.GetMinBackoff())

	// test min.Nanoseconds() < 0
	transaction2 := NewTokenPauseTransaction()
	minBackoff2 := -1 * time.Second

	require.Panics(t, func() { transaction2.SetMinBackoff(minBackoff2) })

	// test transaction.maxBackoff.Nanoseconds() < min.Nanoseconds()
	transaction3 := NewTokenPauseTransaction()
	minBackoff3 := 10 * time.Second

	require.Panics(t, func() { transaction3.SetMinBackoff(minBackoff3) })
}

func TestUnitTokenPauseTransaction_GetMinBackoff(t *testing.T) {
	t.Parallel()
	transaction := NewTokenPauseTransaction()

	require.Equal(t, 250*time.Millisecond, transaction.GetMinBackoff())
}

func TestUnitTokenPauseTransaction_SetLogLevel(t *testing.T) {
	t.Parallel()
	transaction := NewTokenPauseTransaction()

	transaction.SetLogLevel(LoggerLevelDebug)

	level := transaction.GetLogLevel()
	require.Equal(t, LoggerLevelDebug, *level)
}

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

func TestUnitTopicCreateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	topicCreate := NewTopicCreateTransaction().
		SetAutoRenewAccountID(accountID)

	err = topicCreate.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTopicCreateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	topicCreate := NewTopicCreateTransaction().
		SetAutoRenewAccountID(accountID)

	err = topicCreate.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTopicCreateTransactionGet(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 3}

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTopicCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAutoRenewAccountID(accountID).
		SetAdminKey(newKey).
		SetSubmitKey(newKey).
		SetTopicMemo("ad").
		SetAutoRenewPeriod(60 * time.Second).
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

	transaction.GetAutoRenewAccountID()
	transaction.GetAdminKey()
	transaction.GetSubmitKey()
	transaction.GetTopicMemo()
	transaction.GetAutoRenewPeriod()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTopicCreateTransactionNothingSet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTopicCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetAutoRenewAccountID()
	transaction.GetAdminKey()
	transaction.GetSubmitKey()
	transaction.GetTopicMemo()
	transaction.GetAutoRenewPeriod()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTopicCreateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 23}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	newKey2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transaction, err := NewTopicCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(newKey).
		SetSubmitKey(newKey2).
		SetAutoRenewAccountID(accountID).
		SetTopicMemo("memo").
		SetAutoRenewPeriod(time.Second * 3).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetConsensusCreateTopic()
	require.Equal(t, proto.AdminKey.String(), newKey._ToProtoKey().String())
	require.Equal(t, proto.SubmitKey.String(), newKey2._ToProtoKey().String())
	require.Equal(t, proto.Memo, "memo")
	require.Equal(t, proto.AutoRenewPeriod.Seconds, _DurationToProtobuf(time.Second*3).Seconds)
	require.Equal(t, proto.AutoRenewAccount.String(), accountID._ToProtobuf().String())
}

func TestUnitTopicCreateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewTopicCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(newKey).
		SetTopicMemo("ad").
		SetSubmitKey(newKey).
		SetAutoRenewAccountID(account).
		SetAutoRenewPeriod(time.Second * 30).
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
	transaction.GetAdminKey()
	transaction.GetSubmitKey()
	transaction.GetTopicMemo()
	transaction.GetAutoRenewAccountID()
	transaction.GetAutoRenewPeriod()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case TopicCreateTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitTopicCreateTransactionMock(t *testing.T) {
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

	freez, err := NewTopicCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAdminKey(newKey).
		SetSubmitKey(newKey).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}

func TestUnitTopicCreateTransactionSerialization(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	topicCreate, err := NewTopicCreateTransaction().
		SetTransactionID(transactionID).
		SetAdminKey(newKey).
		SetNodeAccountIDs(nodeAccountID).
		SetSubmitKey(newKey).
		SetTopicMemo("ad").
		SetAutoRenewPeriod(time.Second * 30).
		Freeze()
	require.NoError(t, err)

	transactionBytes, err := topicCreate.ToBytes()
	require.NoError(t, err)

	txParsed, err := TransactionFromBytes(transactionBytes)
	require.NoError(t, err)

	result, ok := txParsed.(TopicCreateTransaction)
	require.True(t, ok)

	require.Equal(t, topicCreate.GetTopicMemo(), result.GetTopicMemo())
	require.Equal(t, topicCreate.GetAutoRenewPeriod(), result.GetAutoRenewPeriod())
	adminKey, _ := result.GetAdminKey()
	require.Equal(t, newKey.PublicKey(), adminKey)
	submitKey, _ := result.GetSubmitKey()
	require.Equal(t, newKey.PublicKey(), submitKey)
}

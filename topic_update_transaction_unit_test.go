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

func TestUnitTopicUpdateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	topicID, err := TopicIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	topicUpdate := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAutoRenewAccountID(accountID)

	err = topicUpdate.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTopicUpdateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	topicID, err := TopicIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	topicUpdate := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAutoRenewAccountID(accountID)

	err = topicUpdate.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTopicUpdateTransactionGet(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 3}

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTopicUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAutoRenewAccountID(accountID).
		SetTopicID(TopicID{Topic: 7}).
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

	transaction.GetTopicID()
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

func TestUnitTopicUpdateTransactionNothingSet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTopicUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetTopicID()
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

func TestUnitTopicUpdateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	topicID := TopicID{Topic: 5}
	accountID := AccountID{Account: 23}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	newKey2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transaction, err := NewTopicUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTopicID(topicID).
		SetAutoRenewAccountID(accountID).
		SetAdminKey(newKey).
		SetSubmitKey(newKey2).
		SetTopicMemo("memo").
		SetAutoRenewPeriod(time.Second * 3).
		SetExpirationTime(time.Unix(34, 12)).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetConsensusUpdateTopic()
	require.Equal(t, proto.AdminKey.String(), newKey._ToProtoKey().String())
	require.Equal(t, proto.TopicID.String(), topicID._ToProtobuf().String())
	require.Equal(t, proto.AutoRenewAccount.String(), accountID._ToProtobuf().String())
	require.Equal(t, proto.SubmitKey.String(), newKey2._ToProtoKey().String())
	require.Equal(t, proto.Memo.Value, "memo")
	require.Equal(t, proto.AutoRenewPeriod.Seconds, _DurationToProtobuf(time.Second*3).Seconds)
	require.Equal(t, proto.ExpirationTime.String(), _TimeToProtobuf(time.Unix(34, 12)).String())
}

func TestUnitTopicUpdateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	account := AccountID{Account: 3, checksum: &checksum}
	topic := TopicID{Topic: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewTopicUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTopicID(topic).
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

func TestUnitTopicUpdateTransactionMock(t *testing.T) {
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
	topic := TopicID{Topic: 3, checksum: &checksum}

	freez, err := NewTopicUpdateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTopicID(topic).
		SetAdminKey(newKey).
		SetSubmitKey(newKey).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}

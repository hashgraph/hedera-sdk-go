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

func TestUnitScheduleCreateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	scheduleCreate := NewScheduleCreateTransaction().
		SetPayerAccountID(accountID)

	err = scheduleCreate.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleCreateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	scheduleCreate := NewScheduleCreateTransaction().
		SetPayerAccountID(accountID)

	err = scheduleCreate.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitScheduleSignTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	err = scheduleID.Validate(client)
	require.NoError(t, err)
	scheduleSign := NewScheduleSignTransaction().
		SetScheduleID(scheduleID)

	err = scheduleSign.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleSignTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	err = scheduleID.Validate(client)
	require.Error(t, err)
	scheduleSign := NewScheduleSignTransaction().
		SetScheduleID(scheduleID)

	err = scheduleSign.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitScheduleDeleteTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	scheduleDelete := NewScheduleDeleteTransaction().
		SetScheduleID(scheduleID)

	err = scheduleDelete.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleDeleteTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	scheduleDelete := NewScheduleDeleteTransaction().
		SetScheduleID(scheduleID)

	err = scheduleDelete.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitScheduleCreateTransactionGet(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 7}

	newKey, err := PrivateKeyGenerateEd25519()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewScheduleCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetScheduledTransaction(NewTransferTransaction())

	transaction, err = transaction.
		SetPayerAccountID(accountID).
		SetAdminKey(newKey).
		SetScheduleMemo("").
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

	transaction.GetAdminKey()
	transaction.GetPayerAccountID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetScheduleMemo()
}

func TestUnitScheduleCreateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewScheduleCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetAdminKey()
	transaction.GetPayerAccountID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetScheduleMemo()
}

func TestUnitScheduleDeleteTransactionGet(t *testing.T) {
	t.Parallel()

	scheduleID := ScheduleID{Schedule: 7}

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewScheduleDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetScheduleID(scheduleID).
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

	transaction.GetScheduleID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitScheduleDeleteTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewScheduleDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetScheduleID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitScheduleCreateTransactionCoverage(t *testing.T) {
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

	accountCreate, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		Freeze()
	require.NoError(t, err)

	transaction, err := NewScheduleCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(newKey).
		SetScheduleMemo("no").
		SetPayerAccountID(account).
		SetExpirationTime(time.Unix(3, 23)).
		SetWaitForExpiry(true).
		SetScheduledTransaction(accountCreate)
	require.NoError(t, err)

	transaction, err = transaction.
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
	transaction.GetScheduleMemo()
	transaction.GetPayerAccountID()
	transaction.GetExpirationTime()
	transaction.GetWaitForExpiry()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case ScheduleCreateTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitScheduleCreateTransactionMock(t *testing.T) {
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
	account := AccountID{Account: 3, checksum: &checksum}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	accountCreate, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		Freeze()
	require.NoError(t, err)

	freez, err := NewScheduleCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAdminKey(newKey).
		SetScheduleMemo("no").
		SetPayerAccountID(account).
		SetScheduledTransaction(accountCreate)
	require.NoError(t, err)

	_, err = freez.Execute(client)
	require.NoError(t, err)
}

func TestUnitScheduleDeleteTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	schedule := ScheduleID{Schedule: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewScheduleDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetScheduleID(schedule).
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
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	transaction.GetScheduleID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case ScheduleDeleteTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitScheduleDeleteTransactionMock(t *testing.T) {
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
	schedule := ScheduleID{Schedule: 3, checksum: &checksum}

	freez, err := NewScheduleDeleteTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetScheduleID(schedule).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}

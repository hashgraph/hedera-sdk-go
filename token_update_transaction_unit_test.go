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

func TestUnitTokenUpdateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tokenUpdate := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenUpdate.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenUpdateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenUpdate := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenUpdate.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenUpdateTransactionGet(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	accountID := AccountID{Account: 3, checksum: &checksum}
	token := TokenID{Token: 3, checksum: &checksum}

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewTokenUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetTokenMetadata([]byte{1, 2, 3}).
		SetTreasuryAccountID(accountID).
		SetTokenID(token).
		SetAdminKey(newKey).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetPauseKey(newKey).
		SetMetadataKey(newKey).
		SetKeyVerificationMode(NO_VALIDATION).
		SetExpirationTime(time.Now()).
		SetAutoRenewPeriod(60 * time.Second).
		SetAutoRenewAccount(accountID).
		SetMaxTransactionFee(NewHbar(10)).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		SetGrpcDeadline(&grpc).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
	_, err = transaction.Schedule()
	require.NoError(t, err)
	transaction.GetTokenName()
	transaction.GetTokenSymbol()
	transaction.GetTreasuryAccountID()
	transaction.GetAdminKey()
	transaction.GetFreezeKey()
	transaction.GetWipeKey()
	transaction.GetKycKey()
	transaction.GetSupplyKey()
	transaction.GetPauseKey()
	transaction.GetMetadataKey()
	transaction.GetKeyVerificationMode()
	transaction.GetExpirationTime()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	byt, err := transaction.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)
	sig, err := newKey.SignTransaction(transaction)
	require.NoError(t, err)
	transaction.getName()
	transaction.GetMaxRetry()
	transaction.GetMaxBackoff()
	transaction.GetMinBackoff()
	switch b := txFromBytes.(type) {
	case TokenUpdateTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitTokenUpdateTransactionNothingSet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetTokenName()
	transaction.GetTokenSymbol()
	transaction.GetTreasuryAccountID()
	transaction.GetAdminKey()
	transaction.GetFreezeKey()
	transaction.GetWipeKey()
	transaction.GetKycKey()
	transaction.GetSupplyKey()
	transaction.GetPauseKey()
	transaction.GetMetadataKey()
	transaction.GetKeyVerificationMode()
	transaction.GetExpirationTime()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	proto := transaction.build().GetTokenUpdate()
	require.Nil(t, proto.Token)
	require.Nil(t, proto.AutoRenewAccount)
	require.Nil(t, proto.AdminKey)
	require.Nil(t, proto.Expiry)
	require.Nil(t, proto.FeeScheduleKey)
	require.Nil(t, proto.FreezeKey)
	require.Nil(t, proto.KycKey)
	require.Nil(t, proto.FeeScheduleKey)
	require.Nil(t, proto.PauseKey)
	require.Nil(t, proto.SupplyKey)
	require.Nil(t, proto.MetadataKey)
	require.Nil(t, proto.Treasury)
}

func TestUnitTokenUpdateTransactionKeyCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	keys := make([]PrivateKey, 8)
	var err error

	for i := 0; i < len(keys); i++ {
		keys[i], err = PrivateKeyGenerateEd25519()
		require.NoError(t, err)
	}

	transaction, err := NewTokenUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(keys[0]).
		SetFreezeKey(keys[1]).
		SetWipeKey(keys[2]).
		SetKycKey(keys[3]).
		SetSupplyKey(keys[4]).
		SetPauseKey(keys[5]).
		SetFeeScheduleKey(keys[6]).
		SetMetadataKey(keys[7]).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetTokenUpdate()
	require.Equal(t, proto.AdminKey.String(), keys[0]._ToProtoKey().String())
	require.Equal(t, proto.FreezeKey.String(), keys[1]._ToProtoKey().String())
	require.Equal(t, proto.WipeKey.String(), keys[2]._ToProtoKey().String())
	require.Equal(t, proto.KycKey.String(), keys[3]._ToProtoKey().String())
	require.Equal(t, proto.SupplyKey.String(), keys[4]._ToProtoKey().String())
	require.Equal(t, proto.PauseKey.String(), keys[5]._ToProtoKey().String())
	require.Equal(t, proto.FeeScheduleKey.String(), keys[6]._ToProtoKey().String())
	require.Equal(t, proto.MetadataKey.String(), keys[7]._ToProtoKey().String())
}

func TestUnitTokenUpdateTransactionMock(t *testing.T) {
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
	accountID := AccountID{Account: 3, checksum: &checksum}
	token := TokenID{Token: 3, checksum: &checksum}

	freez, err := NewTokenUpdateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTokenID(token).
		SetTreasuryAccountID(accountID).
		SetAdminKey(newKey).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetPauseKey(newKey).
		SetMetadataKey(newKey).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}

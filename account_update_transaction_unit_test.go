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

func TestUnitAccountUpdateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	accountUpdate := NewAccountUpdateTransaction().
		SetProxyAccountID(accountID)

	err = accountUpdate.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountUpdateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	accountUpdate := NewAccountUpdateTransaction().
		SetProxyAccountID(accountID)

	err = accountUpdate.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountUpdateTransactionMock(t *testing.T) {
	t.Parallel()

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

		for _, sigPair := range sigMap.SigPair {
			verified := false

			switch k := sigPair.Signature.(type) {
			case *services.SignaturePair_Ed25519:
				pbTemp, _ := PublicKeyFromBytesEd25519(sigPair.PubKeyPrefix)
				verified = pbTemp.Verify(signedTransaction.BodyBytes, k.Ed25519)
			case *services.SignaturePair_ECDSASecp256K1:
				pbTemp, _ := PublicKeyFromBytesECDSA(sigPair.PubKeyPrefix)
				verified = pbTemp.Verify(signedTransaction.BodyBytes, k.ECDSASecp256K1)
			}
			require.True(t, verified)
		}

		if bod, ok := transactionBody.Data.(*services.TransactionBody_CryptoUpdateAccount); ok {
			require.Equal(t, bod.CryptoUpdateAccount.Memo.Value, "no")
			require.Equal(t, bod.CryptoUpdateAccount.AccountIDToUpdate.GetAccountNum(), int64(123))
			//alias := services.Key{}
			//_ = protobuf.Unmarshal(bod.CryptoUpdateAccount.Alias, &alias)
			//require.Equal(t, hex.EncodeToString(alias.GetEd25519()), "1480272863d39c42f902bc11601a968eaf30ad662694e3044c86d5df46fabfd2")
		}

		return &services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		}
	}
	responses := [][]interface{}{{
		call,
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()
	//302a300506032b65700321001480272863d39c42f902bc11601a968eaf30ad662694e3044c86d5df46fabfd2
	newKey, err := PrivateKeyFromStringEd25519("302e020100300506032b657004220420278184257eb568d0e5fcfc1df99828b039b4776da05855dc5af105996e6200d1")
	require.NoError(t, err)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	_, err = NewAccountUpdateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(tran).
		SetAccountMemo("no").
		SetAccountID(AccountID{Account: 123}).
		SetAliasKey(newKey.PublicKey()).
		Execute(client)
	require.NoError(t, err)
}

func TestUnitAccountUpdateTransactionGet(t *testing.T) {
	t.Parallel()

	spenderAccountID1 := AccountID{Account: 7}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	key, err := PrivateKeyGenerateEd25519()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAccountID(spenderAccountID1).
		SetKey(key).
		SetProxyAccountID(spenderAccountID1).
		SetAccountMemo("").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetAccountID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxAutomaticTokenAssociations()
	transaction.GetProxyAccountID()
	transaction.GetRegenerateTransactionID()
	transaction.GetKey()
	transaction.GetAutoRenewPeriod()
	transaction.GetReceiverSignatureRequired()
}

func TestUnitAccountUpdateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetAccountID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxAutomaticTokenAssociations()
	transaction.GetProxyAccountID()
	transaction.GetRegenerateTransactionID()
	transaction.GetKey()
	transaction.GetAutoRenewPeriod()
	transaction.GetReceiverSignatureRequired()
}

func TestUnitAccountUpdateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	stackedAccountID := AccountID{Account: 5}
	accountID := AccountID{Account: 6}

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKey(key).
		SetAccountID(accountID).
		SetAccountMemo("ty").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetStakedAccountID(stackedAccountID).
		SetDeclineStakingReward(true).
		SetAutoRenewPeriod(60 * time.Second).
		SetExpirationTime(time.Unix(34, 56)).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetCryptoUpdateAccount()
	require.Equal(t, proto.AccountIDToUpdate.String(), accountID._ToProtobuf().String())
	require.Equal(t, proto.Key.String(), key._ToProtoKey().String())
	require.Equal(t, proto.Memo.Value, "ty")
	require.Equal(t, proto.ReceiverSigRequiredField.(*services.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper).ReceiverSigRequiredWrapper.Value, true)
	require.Equal(t, proto.MaxAutomaticTokenAssociations.GetValue(), int32(2))
	require.Equal(t, proto.StakedId.(*services.CryptoUpdateTransactionBody_StakedAccountId).StakedAccountId.String(),
		stackedAccountID._ToProtobuf().String())
	require.Equal(t, proto.DeclineReward.Value, true)
	require.Equal(t, proto.AutoRenewPeriod.String(), _DurationToProtobuf(60*time.Second).String())
	require.Equal(t, proto.ExpirationTime.String(), _TimeToProtobuf(time.Unix(34, 56)).String())
}

func TestUnitAccountUpdateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewAccountUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKey(key).
		SetAccountID(account).
		SetAccountMemo("ty").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetStakedAccountID(account).
		SetStakedNodeID(4).
		SetDeclineStakingReward(true).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetMaxTransactionFee(NewHbar(3)).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetTransactionMemo("no").
		SetTransactionValidDuration(time.Second * 30).
		SetRegenerateTransactionID(false).
		SetGrpcDeadline(&grpc).
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
	sig, err := key.SignTransaction(transaction)
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	transaction.GetStakedAccountID()
	transaction.GetStakedNodeID()
	transaction.ClearStakedAccountID()
	transaction.ClearStakedNodeID()
	transaction.GetDeclineStakingReward()
	transaction.GetExpirationTime()
	transaction.GetAccountMemo()

	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case *AccountUpdateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

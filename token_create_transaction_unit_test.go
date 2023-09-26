//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"testing"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenCreateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	supple := TokenSupplyTypeFinite

	println("token supply type", supple.String())

	typ := TokenTypeFungibleCommon

	println("token type", typ.String())

	tokenCreate := NewTokenCreateTransaction().
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenCreate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenCreateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenCreate := NewTokenCreateTransaction().
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenCreate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenCreateTransactionGet(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	accountID := AccountID{Account: 3, checksum: &checksum}

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewTokenCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetDecimals(3).
		SetCustomFees([]Fee{}).
		SetTokenType(TokenTypeFungibleCommon).
		SetSupplyType(TokenSupplyTypeFinite).
		SetMaxSupply(12312123123123).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(accountID).
		SetAdminKey(newKey).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetPauseKey(newKey).
		SetExpirationTime(time.Now()).
		SetFreezeDefault(false).
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

	err = transaction._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
	_, err = transaction.Schedule()
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetTokenName()
	transaction.GetTokenSymbol()
	transaction.GetTokenMemo()
	transaction.GetDecimals()
	transaction.GetCustomFees()
	transaction.GetTokenType()
	transaction.GetSupplyType()
	transaction.GetMaxSupply()
	transaction.GetInitialSupply()
	transaction.GetTreasuryAccountID()
	transaction.GetAdminKey()
	transaction.GetFreezeKey()
	transaction.GetWipeKey()
	transaction.GetKycKey()
	transaction.GetSupplyKey()
	transaction.GetPauseKey()
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
	sig, err := newKey.SignTransaction(&transaction.Transaction)
	require.NoError(t, err)
	transaction._GetLogID()
	switch b := txFromBytes.(type) {
	case TokenCreateTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitTokenCreateTransactionNothingSet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenCreateTransaction().
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
	transaction.GetTokenMemo()
	transaction.GetDecimals()
	transaction.GetCustomFees()
	transaction.GetTokenType()
	transaction.GetSupplyType()
	transaction.GetMaxSupply()
	transaction.GetInitialSupply()
	transaction.GetTreasuryAccountID()
	transaction.GetAdminKey()
	transaction.GetFreezeKey()
	transaction.GetWipeKey()
	transaction.GetKycKey()
	transaction.GetSupplyKey()
	transaction.GetPauseKey()
	transaction.GetExpirationTime()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	proto := transaction._Build().GetTokenCreation()
	require.Nil(t, proto.AutoRenewAccount)
	require.Nil(t, proto.AdminKey)
	require.Nil(t, proto.Expiry)
	require.Nil(t, proto.FeeScheduleKey)
	require.Nil(t, proto.FreezeKey)
	require.Nil(t, proto.KycKey)
	require.Nil(t, proto.FeeScheduleKey)
	require.Nil(t, proto.PauseKey)
	require.Nil(t, proto.SupplyKey)
	require.Nil(t, proto.Treasury)
}

func TestUnitTokenCreateTransactionKeyCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	keys := make([]PrivateKey, 7)
	var err error

	for i := 0; i < 7; i++ {
		keys[i], err = PrivateKeyGenerateEd25519()
		require.NoError(t, err)
	}

	transaction, err := NewTokenCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(keys[0]).
		SetFreezeKey(keys[1]).
		SetWipeKey(keys[2]).
		SetKycKey(keys[3]).
		SetSupplyKey(keys[4]).
		SetPauseKey(keys[5]).
		SetFeeScheduleKey(keys[6]).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction._Build().GetTokenCreation()
	require.Equal(t, proto.AdminKey.String(), keys[0]._ToProtoKey().String())
	require.Equal(t, proto.FreezeKey.String(), keys[1]._ToProtoKey().String())
	require.Equal(t, proto.WipeKey.String(), keys[2]._ToProtoKey().String())
	require.Equal(t, proto.KycKey.String(), keys[3]._ToProtoKey().String())
	require.Equal(t, proto.SupplyKey.String(), keys[4]._ToProtoKey().String())
	require.Equal(t, proto.PauseKey.String(), keys[5]._ToProtoKey().String())
	require.Equal(t, proto.FeeScheduleKey.String(), keys[6]._ToProtoKey().String())
}

func TestUnitTokenCreateTransactionMock(t *testing.T) {
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

	freez, err := NewTokenCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTreasuryAccountID(accountID).
		SetAdminKey(newKey).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetPauseKey(newKey).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}

func TestUnitTokenCreateAutoRenewAccountSetting(t *testing.T) {
	client, err := _NewMockClient()
	autoRenewAccount := AccountID{Account: 3}
	frozenTx, err := NewTokenCreateTransaction().
		SetAutoRenewAccount(autoRenewAccount).
		SetAutoRenewPeriod(24 * time.Hour).
		FreezeWith(client)
	require.NoError(t, err)
	require.Equal(t, autoRenewAccount, frozenTx.GetAutoRenewAccount())
}

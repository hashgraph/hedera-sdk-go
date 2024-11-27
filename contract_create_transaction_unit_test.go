//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitContractCreateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	contractCreate := NewContractCreateTransaction().
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractCreate.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitContractCreateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	contractCreate := NewContractCreateTransaction().
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractCreate.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitContractCreateTransactionMock(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyFromStringEd25519("302e020100300506032b657004220420d45e1557156908c967804615af59a000be88c7aa7058bfcbe0f46b16c28f887d")
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

		if bod, ok := transactionBody.Data.(*services.TransactionBody_ContractCreateInstance); ok {
			require.Equal(t, bod.ContractCreateInstance.InitcodeSource.(*services.ContractCreateTransactionBody_FileID).FileID.FileNum, int64(123))
			params := NewContractFunctionParameters().AddString("hello from hiero")
			require.Equal(t, bytes.Compare(bod.ContractCreateInstance.ConstructorParameters, params._Build(nil)), 0)
			require.Equal(t, bod.ContractCreateInstance.Memo, "hiero-sdk-go::TestContractCreateTransaction_Execute")
			require.Equal(t, bod.ContractCreateInstance.Gas, int64(100000))
			require.Equal(t, hex.EncodeToString(bod.ContractCreateInstance.AdminKey.GetEd25519()), key.PublicKey().StringRaw())
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

	_, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetGas(100000).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hiero")).
		SetBytecodeFileID(FileID{File: 123}).
		SetContractMemo("hiero-sdk-go::TestContractCreateTransaction_Execute").
		Execute(client)
	require.NoError(t, err)
}

func TestUnitContractCreateTransactionGet(t *testing.T) {
	t.Parallel()

	spenderAccountID1 := AccountID{Account: 7}
	fileID := FileID{File: 7}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewContractCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetGas(21341).
		SetProxyAccountID(spenderAccountID1).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetContractMemo("yes").
		SetBytecodeFileID(fileID).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetContractMemo()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetProxyAccountID()
	transaction.GetRegenerateTransactionID()
	transaction.GetAutoRenewPeriod()
	transaction.GetBytecodeFileID()
	transaction.GetContractMemo()
	transaction.GetGas()
}

func TestUnitContractCreateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewContractCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetContractMemo()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetProxyAccountID()
	transaction.GetRegenerateTransactionID()
	transaction.GetAutoRenewPeriod()
	transaction.GetBytecodeFileID()
	transaction.GetContractMemo()
	transaction.GetGas()
}

func TestUnitContractCreateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 7}
	fileID := FileID{File: 7}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transaction, err := NewContractCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(newKey.PublicKey()).
		SetBytecodeFileID(fileID).
		SetGas(500).
		SetInitialBalance(NewHbar(50)).
		SetConstructorParametersRaw([]byte{34}).
		SetContractMemo("yes").
		SetStakedAccountID(accountID).
		SetMaxAutomaticTokenAssociations(3).
		SetAutoRenewPeriod(time.Second * 3).
		SetAutoRenewAccountID(accountID).
		SetDeclineStakingReward(true).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetContractCreateInstance()
	require.Equal(t, proto.AdminKey.String(), newKey._ToProtoKey().String())
	require.Equal(t, proto.GetFileID().String(), fileID._ToProtobuf().String())
	require.Equal(t, proto.Memo, "yes")
	require.Equal(t, proto.StakedId.(*services.ContractCreateTransactionBody_StakedAccountId).StakedAccountId.String(),
		accountID._ToProtobuf().String())
	require.Equal(t, proto.MaxAutomaticTokenAssociations, int32(3))
	require.Equal(t, proto.AutoRenewPeriod.String(), _DurationToProtobuf(time.Second*3).String())
	require.Equal(t, proto.AutoRenewAccountId.String(), accountID._ToProtobuf().String())
	require.Equal(t, proto.DeclineReward, true)
	require.Equal(t, proto.ConstructorParameters, []byte{34})
	require.Equal(t, proto.InitialBalance, NewHbar(50).AsTinybar())
	require.Equal(t, proto.Gas, int64(500))
}

func TestUnitContractCreateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	file := FileID{File: 3, checksum: &checksum}
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewContractCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(newKey.PublicKey()).
		SetBytecodeFileID(file).
		SetGas(500).
		SetInitialBalance(NewHbar(50)).
		SetConstructorParametersRaw([]byte{34}).
		SetContractMemo("yes").
		SetStakedAccountID(account).
		SetStakedNodeID(3).
		SetBytecode([]byte{0}).
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
	transaction.GetStakedAccountID()
	transaction.GetStakedNodeID()
	transaction.GetConstructorParameters()
	transaction.GetBytecode()
	transaction.GetDeclineStakingReward()
	transaction.GetInitialBalance()
	transaction.GetGas()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case ContractCreateTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

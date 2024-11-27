//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitContractUpdateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	contractInfoQuery := NewContractUpdateTransaction().
		SetContractID(contractID).
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractInfoQuery.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitContractUpdateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	contractInfoQuery := NewContractUpdateTransaction().
		SetContractID(contractID).
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractInfoQuery.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitContractUpdateTransactionMock(t *testing.T) {
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

		if bod, ok := transactionBody.Data.(*services.TransactionBody_ContractUpdateInstance); ok {
			require.Equal(t, bod.ContractUpdateInstance.ContractID.GetContractNum(), int64(3))
			if mem, ok2 := bod.ContractUpdateInstance.MemoField.(*services.ContractUpdateTransactionBody_MemoWrapper); ok2 {
				require.Equal(t, mem.MemoWrapper.GetValue(), "yes")
			}
			require.Equal(t, hex.EncodeToString(bod.ContractUpdateInstance.GetAdminKey().GetEd25519()), "1480272863d39c42f902bc11601a968eaf30ad662694e3044c86d5df46fabfd2")
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

	_, err = NewContractUpdateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(tran).
		SetAdminKey(newKey.PublicKey()).
		SetContractMemo("yes").
		SetContractID(ContractID{Contract: 3}).
		Execute(client)
	require.NoError(t, err)
}

func TestUnitContractUpdateTransactionGet(t *testing.T) {
	t.Parallel()

	contractID := ContractID{Contract: 7}

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transaction, err := NewContractUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetContractID(contractID).
		SetAdminKey(newKey.PublicKey()).
		SetContractMemo("yes").
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetContractID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetAdminKey()
	transaction.GetRegenerateTransactionID()
	transaction.GetContractMemo()
}

func TestUnitContractUpdateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewContractUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetContractID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetAdminKey()
	transaction.GetRegenerateTransactionID()
	transaction.GetContractMemo()
}

func TestUnitContractUpdateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	contractID := ContractID{Contract: 7}
	accountID := AccountID{Account: 7}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transaction, err := NewContractUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetContractID(contractID).
		SetAdminKey(newKey.PublicKey()).
		SetContractMemo("yes").
		SetStakedAccountID(accountID).
		SetMaxAutomaticTokenAssociations(3).
		SetAutoRenewPeriod(time.Second * 3).
		SetExpirationTime(time.Unix(34, 3)).
		SetAutoRenewAccountID(accountID).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetContractUpdateInstance()
	require.Equal(t, proto.AdminKey.String(), newKey._ToProtoKey().String())
	require.Equal(t, proto.ContractID.String(), contractID._ToProtobuf().String())
	require.Equal(t, proto.MemoField.(*services.ContractUpdateTransactionBody_MemoWrapper).MemoWrapper.Value, "yes")
	require.Equal(t, proto.StakedId.(*services.ContractUpdateTransactionBody_StakedAccountId).StakedAccountId.String(),
		accountID._ToProtobuf().String())
	require.Equal(t, proto.MaxAutomaticTokenAssociations.Value, int32(3))
	require.Equal(t, proto.AutoRenewPeriod.String(), _DurationToProtobuf(time.Second*3).String())
	require.Equal(t, proto.ExpirationTime.String(), _TimeToProtobuf(time.Unix(34, 3)).String())
	require.Equal(t, proto.AutoRenewAccountId.String(), accountID._ToProtobuf().String())
}

func TestUnitContractUpdateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	file := FileID{File: 3, checksum: &checksum}
	account := AccountID{Account: 3, checksum: &checksum}
	contract := ContractID{Contract: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewContractUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(newKey.PublicKey()).
		SetBytecodeFileID(file).
		SetContractMemo("yes").
		SetStakedAccountID(account).
		SetStakedNodeID(3).
		SetContractID(contract).
		SetProxyAccountID(account).
		SetDeclineStakingReward(true).
		SetAutoRenewPeriod(time.Second * 30).
		SetExpirationTime(time.Unix(345, 566)).
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
	transaction.GetProxyAccountID()
	transaction.GetBytecodeFileID()
	transaction.GetAutoRenewPeriod()
	transaction.GetAutoRenewAccountID()
	transaction.GetExpirationTime()
	transaction.GetDeclineStakingReward()
	transaction.ClearStakedAccountID()
	transaction.ClearStakedNodeID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case ContractUpdateTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

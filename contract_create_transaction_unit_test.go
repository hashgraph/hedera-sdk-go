//go:build all || unit
// +build all unit

package hedera

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitContractCreateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	contractCreate := NewContractCreateTransaction().
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractCreate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitContractCreateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	contractCreate := NewContractCreateTransaction().
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractCreate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitMockContractCreateTransaction(t *testing.T) {
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
			require.Equal(t, bod.ContractCreateInstance.FileID.FileNum, int64(123))
			params := NewContractFunctionParameters().AddString("hello from hedera")
			require.Equal(t, bytes.Compare(bod.ContractCreateInstance.ConstructorParameters, params._Build(nil)), 0)
			require.Equal(t, bod.ContractCreateInstance.Memo, "hedera-sdk-go::TestContractCreateTransaction_Execute")
			require.Equal(t, bod.ContractCreateInstance.Gas, int64(75000))
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
		SetGas(75000).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(FileID{File: 123}).
		SetContractMemo("hedera-sdk-go::TestContractCreateTransaction_Execute").
		Execute(client)
	require.NoError(t, err)
}

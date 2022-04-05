//go:build all || unit
// +build all unit

package hedera

import (
	"bytes"
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitFileUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	fileUpdate := NewFileUpdateTransaction().
		SetFileID(fileID)

	err = fileUpdate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitFileUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	fileUpdate := NewFileUpdateTransaction().
		SetFileID(fileID)

	err = fileUpdate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitMockFileUpdateTransaction(t *testing.T) {
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

		key, _ := PrivateKeyFromStringEd25519("302e020100300506032b657004220420d45e1557156908c967804615af59a000be88c7aa7058bfcbe0f46b16c28f887d")
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

		require.Equal(t, bytes.Compare(sigMap.SigPair[1].PubKeyPrefix, key.PublicKey().BytesRaw()), 0)
		require.Equal(t, bytes.Compare(sigMap.SigPair[0].PubKeyPrefix, newKey.PublicKey().BytesRaw()), 0)

		if bod, ok := transactionBody.Data.(*services.TransactionBody_FileUpdate); ok {
			require.Equal(t, bod.FileUpdate.FileID.FileNum, int64(3))
			require.Equal(t, bytes.Compare(bod.FileUpdate.Contents, []byte{123}), 0)
			require.Equal(t, bod.FileUpdate.Memo, &wrapperspb.StringValue{Value: "no memo"})
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

	freez, err := NewFileUpdateTransaction().
		SetFileID(FileID{File: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetFileMemo("no memo").
		SetKeys(newKey).
		SetContents([]byte{123}).
		FreezeWith(client)
	require.NoError(t, err)

	_, err = freez.Sign(newKey).Execute(client)
	require.NoError(t, err)
}

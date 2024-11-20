//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitFileAppendTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	fileAppend := NewFileAppendTransaction().
		SetFileID(fileID)

	err = fileAppend.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitFileAppendTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	fileAppend := NewFileAppendTransaction().
		SetFileID(fileID)

	err = fileAppend.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitFileAppendTransactionMock(t *testing.T) {
	t.Parallel()

	fil := []byte(" world!")
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

		if bod, ok := transactionBody.Data.(*services.TransactionBody_FileAppend); ok {
			require.Equal(t, bod.FileAppend.FileID.FileNum, int64(3))
			require.Equal(t, bytes.Compare(bod.FileAppend.Contents, fil), 0)
		}

		return &services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		}
	}
	responses := [][]interface{}{{
		call, &services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status: services.ResponseCodeEnum_SUCCESS,
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
							AccountNum: 234,
						}},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	_, err := NewFileAppendTransaction().
		SetFileID(FileID{File: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetContents(fil).
		Execute(client)
	require.NoError(t, err)
}

func TestUnitFileAppendTransactionGet(t *testing.T) {
	t.Parallel()

	fileID := FileID{File: 7}

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewFileAppendTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetFileID(fileID).
		SetContents([]byte("Hello, World")).
		SetMaxTransactionFee(NewHbar(10)).
		SetTransactionMemo("").
		SetMaxChunkSize(12).
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetFileID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetContents()
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxChunkSize()
}

//func TestUnitFileAppendTransactionNothingSet(t *testing.T) {
//	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
//	transactionID := TransactionIDGenerate(AccountID{Account: 324})
//
//	transaction, err := NewFileAppendTransaction().
//		SetTransactionID(transactionID).
//		SetNodeAccountIDs(nodeAccountID).
//		Freeze()
//	require.NoError(t, err)
//
//	transaction.GetTransactionID()
//	transaction.GetNodeAccountIDs()
//
//	_, err = transaction.GetTransactionHash()
//	require.NoError(t, err)
//
//	transaction.GetFileID()
//	transaction.GetMaxTransactionFee()
//	transaction.GetTransactionMemo()
//	transaction.GetRegenerateTransactionID()
//	_, err = transaction.GetSignatures()
//	require.NoError(t, err)
//	transaction.GetRegenerateTransactionID()
//	transaction.GetMaxTransactionFee()
//	transaction.GetContents()
//	transaction.GetRegenerateTransactionID()
//	transaction.GetMaxChunkSize()
//}

func TestUnitFileAppendTransactionBigContentsMock(t *testing.T) {
	t.Parallel()

	var previousTransactionID string

	receipt := &services.Response{
		Response: &services.Response_TransactionGetReceipt{
			TransactionGetReceipt: &services.TransactionGetReceiptResponse{
				Header: &services.ResponseHeader{
					NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
					ResponseType:                services.ResponseType_ANSWER_ONLY,
				},
				Receipt: &services.TransactionReceipt{
					Status: services.ResponseCodeEnum_SUCCESS,
					FileID: &services.FileID{FileNum: 3},
				},
			},
		},
	}

	contents := ""

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
		if previousTransactionID == "" {
			previousTransactionID = transactionId
		} else {
			require.NotEqual(t, transactionId, previousTransactionID)
			previousTransactionID = transactionId
		}

		contents += string(transactionBody.Data.(*services.TransactionBody_FileAppend).FileAppend.Contents)

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

		return &services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		}
	}
	responses := [][]interface{}{{
		call, receipt, call, receipt, call, receipt, call, receipt, call, receipt, call, receipt, call,
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	_, err := NewFileAppendTransaction().
		SetFileID(FileID{File: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetContents([]byte(bigContents2)).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, bigContents2, contents)
}

func TestUnitFileAppendTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	file := FileID{File: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewFileAppendTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetFileID(file).
		SetContents([]byte{1}).
		SetMaxChunkSize(5).
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
	transaction.GetContents()
	transaction.GetMaxChunkSize()
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
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case FileAppendTransaction:
		b.AddSignature(newKey.PublicKey(), sig)
	}
}

func TestUnitFileAppendTransactionSerialization(t *testing.T) {
	t.Parallel()

	fileID := FileID{File: 7}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewFileAppendTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetFileID(fileID).
		SetContents([]byte("Hello, World")).
		Freeze()
	require.NoError(t, err)

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	txParsed, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	result, ok := txParsed.(FileAppendTransaction)
	require.True(t, ok)

	require.Equal(t, transactionID.AccountID, result.GetTransactionID().AccountID)
	require.Equal(t, fileID, result.GetFileID())
	require.Equal(t, transaction.GetContents(), result.GetContents())
}

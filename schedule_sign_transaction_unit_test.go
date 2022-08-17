package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/require"
	protobuf "google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestUnitScheduleSignTransactionCoverage(t *testing.T) {
	checksum := "dmqui"
	grpc := time.Second * 30
	schedule := ScheduleID{Schedule: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)

	transaction, err := NewScheduleSignTransaction().
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

	transaction.Sign(newKey)

	transaction._ValidateNetworkOnIDs(client)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()
	transaction.GetMaxRetry()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxBackoff()
	transaction.GetMinBackoff()
	transaction.GetRegenerateTransactionID()
	byt, err := transaction.ToBytes()
	require.NoError(t, err)
	_, err = TransactionFromBytes(byt)
	require.NoError(t, err)
	_, err = newKey.SignTransaction(&transaction.Transaction)
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	transaction.GetScheduleID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction._GetLogID()
	//switch b := txFromBytes.(type) {
	//case ScheduleSignTransaction:
	//	b.AddSignature(newKey.PublicKey(), sig)
	//}
}

func TestUnitScheduleSignTransactionMock(t *testing.T) {
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
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewScheduleSignTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetScheduleID(schedule).
		Freeze()
	require.NoError(t, err)

	transaction, err = transaction.SignWithOperator(client)
	require.NoError(t, err)
	
	_, err = transaction.Execute(client)
	require.NoError(t, err)
}

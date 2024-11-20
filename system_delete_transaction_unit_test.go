//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

var testContractId = ContractID{Contract: 5}
var testExpirationTime = time.Now().Add(24 * time.Hour)
var testFileId = FileID{File: 3}
var testTrxValidDuration = 24 * time.Hour

func TestUnitSystemDeleteTransactionFromProtobuf(t *testing.T) {
	t.Parallel()

	trx, trxBody := _CreateProtoBufTrxBody()
	sysDeleteTrx := _SystemDeleteTransactionFromProtobuf(trx, trxBody)
	require.NotNil(t, sysDeleteTrx)
	require.Equal(t, "memo", sysDeleteTrx.memo)
	require.Equal(t, uint64(5), sysDeleteTrx.transactionFee)
	require.Equal(t, uint64(10), sysDeleteTrx.defaultMaxTransactionFee)
}

func TestUnitSystemDeleteTrxGettersAndSetters(t *testing.T) {
	t.Parallel()
	deleteTrx := _SetupSystemDeleteTrx()

	require.Equal(t, testContractId, deleteTrx.GetContractID())
	require.Equal(t, testExpirationTime.Unix(), deleteTrx.GetExpirationTime())
	require.Equal(t, testFileId, deleteTrx.GetFileID())
	require.Equal(t, testTrxValidDuration, deleteTrx.GetTransactionValidDuration())
}

func TestUnitSystemDeleteTrxValidateNetworkOnIDs(t *testing.T) {
	t.Parallel()
	deleteTrx := _SetupSystemDeleteTrx()
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)

	error := deleteTrx.validateNetworkOnIDs(client)
	require.NoError(t, error)
}

func TestUnitSystemDeleteTrxBuild(t *testing.T) {
	t.Parallel()
	deleteTrx := _SetupSystemDeleteTrx()

	trxBody := deleteTrx.build()

	require.NotNil(t, trxBody)
	require.Equal(t, "memo", trxBody.Memo)
	require.Equal(t, uint64(0), trxBody.TransactionFee)
	require.Equal(t, int64(testTrxValidDuration.Seconds()), trxBody.TransactionValidDuration.Seconds)
}

func TestUnitSystemDeleteTrxExecute(t *testing.T) {
	t.Parallel()
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	require.NoError(t, err)
	deleteTrx := _SetupSystemDeleteTrx()

	contractId, _ := ContractIDFromString("0.0.123-esxsf")
	deleteTrx.SetContractID(contractId)

	fileId, _ := FileIDFromString("0.0.123-esxsf")
	deleteTrx.SetFileID(fileId)

	_, err = deleteTrx.FreezeWith(client)

	deleteTrx.Sign(*client.operator.privateKey)
	response, _ := deleteTrx.Execute(client)
	require.Equal(t, deleteTrx.transactionID, response.TransactionID)

}

func TestUnitSystemConstructNewScheduleDeleteTransactionProtobuf(t *testing.T) {
	t.Parallel()
	deleteTrx := _SetupSystemUndeleteTrx()

	protoBody, err := deleteTrx.buildScheduled()
	require.NoError(t, err)
	require.NotNil(t, protoBody)
	require.Equal(t, "memo", protoBody.Memo)
	require.Equal(t, uint64(0), protoBody.TransactionFee)
}

func _CreateProtoBufTrxBody() (Transaction[*SystemDeleteTransaction], *services.TransactionBody) {
	transaction := Transaction[*SystemDeleteTransaction]{BaseTransaction: &BaseTransaction{transactionFee: 5, memo: "memo", defaultMaxTransactionFee: 10}}
	transactionBody := &services.TransactionBody{
		Data: &services.TransactionBody_SystemDelete{SystemDelete: &services.SystemDeleteTransactionBody{ExpirationTime: &services.TimestampSeconds{Seconds: 100}}}}

	return transaction, transactionBody
}

func _SetupSystemDeleteTrx() *SystemDeleteTransaction {

	return NewSystemDeleteTransaction().
		SetContractID(testContractId).
		SetExpirationTime(testExpirationTime).
		SetFileID(testFileId).
		SetTransactionValidDuration(testTrxValidDuration).
		SetTransactionMemo("memo")
}

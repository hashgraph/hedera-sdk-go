//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitSystemUndeleteTransactionFromProtobuf(t *testing.T) {
	t.Parallel()

	trx, trxBody := _CreateProtoBufUndeleteTrxBody()
	sysUndeleteTrx := _SystemUndeleteTransactionFromProtobuf(trx, trxBody)
	require.NotNil(t, sysUndeleteTrx)
	require.Equal(t, "memo", sysUndeleteTrx.memo)
	require.Equal(t, uint64(5), sysUndeleteTrx.transactionFee)
	require.Equal(t, uint64(10), sysUndeleteTrx.defaultMaxTransactionFee)
}

func TestUnitSystemUndeleteTrxGettersAndSetters(t *testing.T) {
	t.Parallel()
	undeleteTrx := _SetupSystemUndeleteTrx()

	require.Equal(t, testContractId, undeleteTrx.GetContractID())
	require.Equal(t, undeleteTrx.GetNodeAccountIDs(), []AccountID{AccountID{Account: 3}})
	require.Equal(t, testFileId, undeleteTrx.GetFileID())
	require.Equal(t, testTrxValidDuration, undeleteTrx.GetTransactionValidDuration())
	require.Equal(t, testTrxValidDuration, *undeleteTrx.GetGrpcDeadline())
}

func TestUnitSystemUndeleteTrxValidateNetworkOnIDs(t *testing.T) {
	t.Parallel()
	undeleteTrx := _SetupSystemUndeleteTrx()
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)

	error := undeleteTrx.validateNetworkOnIDs(client)
	require.NoError(t, error)
}

func TestUnitSystemUndeleteTrxBuild(t *testing.T) {
	t.Parallel()
	deleteTrx := _SetupSystemUndeleteTrx()

	trxBody := deleteTrx.build()
	require.NotNil(t, trxBody)
	require.Equal(t, "memo", trxBody.Memo)
	require.Equal(t, uint64(0), trxBody.TransactionFee)
	require.Equal(t, int64(testTrxValidDuration.Seconds()), trxBody.TransactionValidDuration.Seconds)
}

func TestUnitSystemUndeleteTrxExecute(t *testing.T) {
	t.Parallel()
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	require.NoError(t, err)
	undeleteTrx := _SetupSystemUndeleteTrx()

	contractId, _ := ContractIDFromString("0.0.123-esxsf")
	undeleteTrx.SetContractID(contractId)

	fileId, _ := FileIDFromString("0.0.123-esxsf")
	undeleteTrx.SetFileID(fileId)

	_, err = undeleteTrx.FreezeWith(client)
	undeleteTrx.Sign(*client.operator.privateKey)
	response, _ := undeleteTrx.Execute(client)

	require.Equal(t, undeleteTrx.transactionID, response.TransactionID)
}

func TestUnitSystemConstructNewScheduleUndeleteTransactionProtobuf(t *testing.T) {
	t.Parallel()
	undeleteTrx := _SetupSystemUndeleteTrx()

	protoBody, err := undeleteTrx.buildScheduled()
	require.NoError(t, err)
	require.NotNil(t, protoBody)
	require.Equal(t, "memo", protoBody.Memo)
	require.Equal(t, uint64(0), protoBody.TransactionFee)
}

func _CreateProtoBufUndeleteTrxBody() (Transaction[*SystemUndeleteTransaction], *services.TransactionBody) {
	transaction := Transaction[*SystemUndeleteTransaction]{BaseTransaction: &BaseTransaction{transactionFee: 5, memo: "memo", defaultMaxTransactionFee: 10}}
	transactionBody := &services.TransactionBody{
		Data: &services.TransactionBody_SystemUndelete{SystemUndelete: &services.SystemUndeleteTransactionBody{}}}

	return transaction, transactionBody
}

func _SetupSystemUndeleteTrx() *SystemUndeleteTransaction {
	testAccountID := AccountID{Account: 3}

	return NewSystemUndeleteTransaction().
		SetContractID(testContractId).
		SetFileID(testFileId).
		SetTransactionValidDuration(testTrxValidDuration).
		SetTransactionMemo("memo").
		SetGrpcDeadline(&testTrxValidDuration).
		SetNodeAccountIDs([]AccountID{testAccountID})
}

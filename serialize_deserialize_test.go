//go:build all || e2e
// +build all e2e

package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrationSerializeTransactionDeserializeAndAgainSerializeHasTheSameBytesFreezeBeforeSer(t *testing.T) {
	t.Parallel()

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	env := NewIntegrationTestEnv(t)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	transactionOriginal, _ := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).SignWithOperator(env.Client)
	transactionOriginal, _ = transactionOriginal.FreezeWith(env.Client)

	require.NoError(t, err)
	firstBytes, _ := transactionOriginal.ToBytes()

	txFromBytes, err := TransactionFromBytes(firstBytes)
	require.NoError(t, err)

	transaction := txFromBytes.(AccountCreateTransaction)
	secondBytes, err := transaction.ToBytes()
	fmt.Println(len(secondBytes))
	fmt.Println(secondBytes)
	require.NoError(t, err)

	assert.Equal(t, firstBytes, secondBytes)
}

func TestIntegrationSerializeTransactionDeserializeAndAgainSerializeHasTheSameBytesDontFreeze(t *testing.T) {
	t.Parallel()

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)
	originalTransaction := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance)
	firstBytes, err := originalTransaction.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(firstBytes)
	require.NoError(t, err)
	transaction := txFromBytes.(AccountCreateTransaction)

	secondBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	assert.Equal(t, firstBytes, secondBytes)
}

func TestIntegrationSerializeTransactionWithoutNodeAccountIdDeserialiseAndExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	transactionOriginal := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance)

	require.NoError(t, err)
	resp, _ := transactionOriginal.ToBytes()

	txFromBytes, err := TransactionFromBytes(resp)
	require.NoError(t, err)

	transaction := txFromBytes.(AccountCreateTransaction)
	_, err = transaction.
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)

	require.NoError(t, err)
}

func TestIntegrationAddSignatureSerializeDeserializeAddAnotherSignatureExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// Generate new key to use with new account
	newKey, err := GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	resp, err := NewAccountCreateTransaction().SetKey(newKey).Execute(env.Client)
	receipt, err := resp.GetReceipt(env.Client)
	newAccountId := *receipt.AccountID

	// Prepare and sign the tx and send it to be signed by another actor
	txBefore := NewTransferTransaction().SetTransactionMemo("Serialize/Deserialize transaction test").AddHbarTransfer(env.OperatorID, NewHbar(-1)).AddHbarTransfer(newAccountId, NewHbar(1)).
		Sign(env.OperatorKey)

	bytes, err := txBefore.ToBytes()

	FromBytes, err := TransactionFromBytes(bytes)
	if err != nil {
		panic(err)
	}
	txFromBytes := FromBytes.(TransferTransaction)
	// Assert the fields are the same:
	assert.Equal(t, txFromBytes.signedTransactions._Length(), txBefore.signedTransactions._Length())
	assert.Equal(t, txFromBytes.memo, txBefore.memo)

	executed, err := txFromBytes.Sign(newKey).Execute(env.Client)
	if err != nil {
		panic(err)
	}
	receipt, err = executed.GetReceipt(env.Client)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tx successfully executed. Here is receipt:", receipt)
}

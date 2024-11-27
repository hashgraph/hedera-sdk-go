//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	frozenTx, err := txFromBytes.FreezeWith(env.Client)
	require.NoError(t, err)

	executed, err := frozenTx.Sign(newKey).Execute(env.Client)
	if err != nil {
		panic(err)
	}
	receipt, err = executed.GetReceipt(env.Client)
	assert.Equal(t, receipt.Status, StatusSuccess)
	if err != nil {
		panic(err)
	}
}

func TestIntegrationTransactionShouldReturnFailedReceiptWhenFieldsAreNotSet(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	// Prepare and sign the tx and send it to be signed by another actor
	txBefore := NewTransferTransaction().SetTransactionMemo("Serialize/Deserialize transaction test").AddHbarTransfer(env.OperatorID, NewHbar(-1)).
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

	_, err = txFromBytes.Execute(env.Client)
	assert.Error(t, err)
}

func TestIntegrationAddSignatureSerializeDeserialiseExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)
	updateBytes, err := tx.ToBytes()
	require.NoError(t, err)

	sig1, err := newKey.SignTransaction(tx)
	require.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	require.NoError(t, err)

	if newTx, ok := tx2.(AccountDeleteTransaction); ok {
		assert.True(t, newTx.IsFrozen())
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).Execute(env.Client)
		require.NoError(t, err)
	}

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)

}

func TestIntegrationTopicCreateTransactionAfterSerialization(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tx := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetSubmitKey(env.Client.GetOperatorPublicKey()).
		SetTopicMemo(topicMemo)

	// Serialize unfinished transaction
	bytes, err := tx.ToBytes()

	fromBytes, err := TransactionFromBytes(bytes)
	require.NoError(t, err)
	// Deserialize and add node accounts transaction
	transaction := fromBytes.(TopicCreateTransaction)
	resp, err := transaction.SetNodeAccountIDs(env.NodeAccountIDs).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationTopicSubmitTransactionSerializationDeserialization(t *testing.T) {
	const bigContents2 = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur aliquam augue sem, ut mattis dui laoreet a. Curabitur consequat est euismod, scelerisque metus et, tristique dui. Nulla commodo mauris ut faucibus ultricies. Quisque venenatis nisl nec augue tempus, at efficitur elit eleifend. Duis pharetra felis metus, sed dapibus urna vehicula id. Duis non venenatis turpis, sit amet ornare orci. Donec non interdum quam. Sed finibus nunc et risus finibus, non sagittis lorem cursus. Proin pellentesque tempor aliquam. Sed congue nisl in enim bibendum, condimentum vehicula nisi feugiat.

Suspendisse non sodales arcu. Suspendisse sodales, lorem ac mollis blandit, ipsum neque porttitor nulla, et sodales arcu ante fermentum tellus. Integer sagittis dolor sed augue fringilla accumsan. Cras vitae finibus arcu, sit amet varius dolor. Etiam id finibus dolor, vitae luctus velit. Proin efficitur augue nec pharetra accumsan. Aliquam lobortis nisl diam, vel fermentum purus finibus id. Etiam at finibus orci, et tincidunt turpis. Aliquam imperdiet congue lacus vel facilisis. Phasellus id magna vitae enim dapibus vestibulum vitae quis augue. Morbi eu consequat enim. Maecenas neque nulla, pulvinar sit amet consequat sed, tempor sed magna. Mauris lacinia sem feugiat faucibus aliquet. Etiam congue non turpis at commodo. Nulla facilisi.

Nunc velit turpis, cursus ornare fringilla eu, lacinia posuere leo. Mauris rutrum ultricies dui et suscipit. Curabitur in euismod ligula. Curabitur vitae faucibus orci. Phasellus volutpat vestibulum diam sit amet vestibulum. In vel purus leo. Nulla condimentum lectus vestibulum lectus faucibus, id lobortis eros consequat. Proin mollis libero elit, vel aliquet nisi imperdiet et. Morbi ornare est velit, in vehicula nunc malesuada quis. Donec vehicula convallis interdum.
`

	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tx := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetSubmitKey(env.Client.GetOperatorPublicKey()).
		SetTopicMemo(topicMemo)

	// Serialize unfinished transaction
	bytes, err := tx.ToBytes()

	fromBytes, err := TransactionFromBytes(bytes)
	require.NoError(t, err)
	// Deserialize and add node accounts transaction
	transaction := fromBytes.(TopicCreateTransaction)
	resp, err := transaction.SetNodeAccountIDs(env.NodeAccountIDs).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	submitBytes, err := NewTopicMessageSubmitTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMessage([]byte(bigContents2)).
		SetTopicID(topicID).ToBytes()
	require.NoError(t, err)

	fromBytes, err = TransactionFromBytes(submitBytes)
	require.NoError(t, err)

	topicSubmitTx := fromBytes.(TopicMessageSubmitTransaction)
	_, err = topicSubmitTx.Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

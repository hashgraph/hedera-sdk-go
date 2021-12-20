//+build all e2e

package hedera

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTransferTransactionCanTransferHbar(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationTransferTransactionTransferHbarNothingSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationTransferTransactionTransferHbarPositiveFlippedAmount(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(10)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	frozen, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(10)).
		AddHbarTransfer(accountID, NewHbar(-10)).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	frozen.Sign(newKey)

	resp, err = frozen.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func DisabledTestIntegrationTransferTransactionTransferHbarLoadOf1000(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)
	var err error
	tx := make([]*TransferTransaction, 500)
	response := make([]TransactionResponse, len(tx))
	receipt := make([]TransactionReceipt, len(tx))

	for i := 0; i < len(tx); i++ {
		tx[i], err = NewTransferTransaction().
			AddHbarTransfer(env.Client.GetOperatorAccountID(), HbarFromTinybar(-10)).
			AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(10)).
			FreezeWith(env.Client)
		if err != nil {
			panic(err)
		}

		_, err = tx[i].SignWithOperator(env.Client)
		if err != nil {
			panic(err)
		}

		response[i], err = tx[i].Execute(env.Client)
		if err != nil {
			panic(err)
		}

		receipt[i], err = response[i].GetReceipt(env.Client)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\r%v", i)
	}
}

func TestIntegrationTransferTransactionCanTransferFromBytes(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	newBalance := NewHbar(10)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(accountID, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	transferTx.Sign(newKey)

	transferTxBytes, err := transferTx.ToBytes()
	require.NoError(t, err)

	transactionInterface, err := TransactionFromBytes(transferTxBytes)
	require.NoError(t, err)

	resp, err = TransactionExecute(transactionInterface, env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationTransferTransactionCanTransferFromBytesAfter(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	newBalance := NewHbar(10)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(accountID, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	transferTxBytes, err := transferTx.ToBytes()
	require.NoError(t, err)

	transactionInterface, err := TransactionFromBytes(transferTxBytes)
	require.NoError(t, err)

	signedTx, err := TransactionSign(transactionInterface, newKey)
	require.NoError(t, err)

	resp, err = TransactionExecute(signedTx, env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationTransferTransactionCanTransferSignature(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	newBalance := NewHbar(10)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	transferTx, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(accountID, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	transferTxBytes, err := transferTx.ToBytes()
	require.NoError(t, err)

	signature, err := newKey.SignTransaction(&transferTx.Transaction)

	transactionInterface, err := TransactionFromBytes(transferTxBytes)
	require.NoError(t, err)

	signedTx, err := TransactionAddSignature(transactionInterface, newKey.PublicKey(), signature)
	require.NoError(t, err)

	resp, err = TransactionExecute(signedTx, env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationTransferTransactionCanTransferHbarWithAliasID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	key, err := GeneratePrivateKey()
	require.NoError(t, err)
	aliasAccountID := key.ToAccountID(0, 0)

	resp, err := NewTransferTransaction().
		AddHbarTransfer(env.OperatorID, NewHbar(1).Negated()).
		AddHbarTransfer(*aliasAccountID, NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

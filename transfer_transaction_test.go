package hedera

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationTransferTransactionCanTransferHbar(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTransferTransactionTransferHbarNothingSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTransferTransactionTransferHbarPositiveFlippedAmount(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(10)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	frozen, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(10)).
		AddHbarTransfer(accountID, NewHbar(-10)).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	frozen.Sign(newKey)

	resp, err = frozen.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
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

package hedera

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationAccountCreateTransactionCanFreezeModify(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	tx = tx.SetAccountID(accountID)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "transaction is immutable; it has at least one signature or has been explicitly frozen", err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionNoKey(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status KEY_REQUIRED received for transaction %s", resp.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionAddSignature(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	updateBytes, err := tx.ToBytes()
	assert.NoError(t, err)

	sig1, err := newKey.SignTransaction(&tx.Transaction)
	assert.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	assert.NoError(t, err)

	if newTx, ok := tx2.(AccountDeleteTransaction); ok {
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).Execute(env.Client)
		assert.NoError(t, err)
	}

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountCreateTransactionSetProxyAccountID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		SetProxyAccountID(accountID).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID2 := *receipt.AccountID

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID2).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, accountID.String(), info.ProxyAccountID.String())

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestUnitAccountCreateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	createAccount := NewAccountCreateTransaction().
		SetProxyAccountID(accountID)

	err = createAccount._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitAccountCreateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	createAccount := NewAccountCreateTransaction().
		SetProxyAccountID(accountID)

	err = createAccount._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch; some IDs have different networks set", err.Error())
	}
}

func TestIntegrationAccountCreateTransactionNetwork(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	env.Client.SetAutoValidateChecksums(true)
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	newClient := Client{}
	networkName := NetworkNameMainnet
	newClient.network.networkName = &networkName
	accountID._SetNetworkWithClient(&newClient)

	_, err = NewAccountDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	assert.Error(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

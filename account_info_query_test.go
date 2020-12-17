package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeAccountInfoQuery(t *testing.T) {
	query := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 3}).
		Query

	assert.Equal(t, `cryptoGetInfo:{header:{}accountID:{accountNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestAccountInfoQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)
	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, accountID, info.AccountID)
	assert.Equal(t, false, info.IsDeleted)
	assert.Equal(t, newKey.PublicKey(), info.Key)
	assert.Equal(t, newBalance.tinybar, info.Balance.tinybar)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_AccountInfo_NoAccountID(t *testing.T) {
	client := newTestClient(t)

	_, err := NewAccountInfoQuery().
		Execute(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_ACCOUNT_ID"), err.Error())
}

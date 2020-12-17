package hedera

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeLiveHashQuery(t *testing.T) {
	query := NewLiveHashQuery().
		SetAccountID(AccountID{Account: 3}).
		Query

	assert.Equal(t, `cryptoGetLiveHash:{header:{}accountID:{accountNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestLiveHashQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	_hash, err := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {

	}

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp2, err := NewLiveHashAddTransaction().
		SetAccountID(accountID).
		SetDuration(24 * 30 * time.Hour).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash).
		SetKeys(newKey.PublicKey()).
		Execute(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED received for transaction %s", resp2.TransactionID), err.Error())

	_, err = resp2.GetReceipt(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("Invalid node AccountID was set for transaction: %s", resp2.NodeID), err.Error())

	_, err = NewLiveHashQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash).
		Execute(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED"), err.Error())

	resp2, err = NewLiveHashDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash).
		Execute(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED received for transaction %s", resp2.TransactionID), err.Error())

	_, err = resp2.GetReceipt(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("Invalid node AccountID was set for transaction: %s", resp2.NodeID), err.Error())

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

package hedera

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLiveHashAddTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	_hash, err := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {

	}

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
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
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED received for transaction %s", resp2.TransactionID), err.Error())
	}

	resp2, err = NewLiveHashDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED received for transaction %s", resp2.TransactionID), err.Error())
	}

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

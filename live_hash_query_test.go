package hedera

import (
	"encoding/hex"
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

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	resp, err = NewLiveHashAddTransaction().
		SetAccountID(accountID).
		SetDuration(24 * 30 * time.Hour).
		SetNodeAccountIDs(nodeIDs).
		SetHash(_hash).
		SetKeys(newKey.PublicKey()).
		Execute(client)

	assert.Error(t, err)

	_, err = NewLiveHashDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetHash(_hash).
		Execute(client)
	assert.Error(t, err)

	_, err = NewLiveHashQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetHash(_hash).
		Execute(client)
	assert.Error(t, err)

	resp, err = NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetTransferAccountID(client.GetOperatorAccountID()).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

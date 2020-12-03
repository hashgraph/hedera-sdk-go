package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeAccountRecordsQuery(t *testing.T) {
	query := NewAccountRecordsQuery().
		SetAccountID(AccountID{Account: 3}).
		Query

	assert.Equal(t, `cryptoGetAccountRecords:{header:{}accountID:{accountNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestAccountRecordQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs(nodeIDs).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(client)
	assert.NoError(t, err)

	recordsQuery, err := NewAccountRecordsQuery().
		SetNodeAccountIDs(nodeIDs).
		SetAccountID(client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(recordsQuery))
}

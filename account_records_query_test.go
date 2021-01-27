package hedera

import (
	"fmt"
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
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(client)
	assert.NoError(t, err)

	recordsQuery, err := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(recordsQuery))
}

func TestAccountRecordQueryCost_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(client)
	assert.NoError(t, err)

	records := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(client.GetOperatorAccountID())

	cost, err := records.GetCost(client)
	assert.NoError(t, err)

	recordsQuery, err := records.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(recordsQuery))
}

func TestAccountRecordQueryCost_BigMax_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(client)
	assert.NoError(t, err)

	records := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(100000))

	_, err = records.GetCost(client)
	assert.NoError(t, err)

	recordsQuery, err := records.Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(recordsQuery))
}

func TestAccountRecordQueryCost_SmallMax_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(client)
	assert.NoError(t, err)

	records := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(client.GetOperatorAccountID()).
		SetMaxQueryPayment(HbarFromTinybar(1))

	cost, err := records.GetCost(client)
	assert.NoError(t, err)

	recordsQuery, err := records.Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("cost of AccountRecordsQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tħ"), err.Error())
	}

	assert.Equal(t, 0, len(recordsQuery))
}

func TestAccountRecordQueryCost_InsufficientFee_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(client)
	assert.NoError(t, err)

	records := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(client.GetOperatorAccountID())

	_, err = records.GetCost(client)
	assert.NoError(t, err)

	recordsQuery, err := records.SetQueryPayment(HbarFromTinybar(1)).Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INSUFFICIENT_TX_FEE"), err.Error())
	}

	assert.Equal(t, 0, len(recordsQuery))
}

func Test_AccountRecord_NoAccountID(t *testing.T) {
	client := newTestClient(t)

	_, err := NewAccountRecordsQuery().
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_ACCOUNT_ID"), err.Error())
	}
}

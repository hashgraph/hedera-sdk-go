package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenBurnTransaction_Execute(t *testing.T) {
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

	resp, err = NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFreezeKey(client.GetOperatorPublicKey()).
		SetWipeKey(client.GetOperatorPublicKey()).
		SetKycKey(client.GetOperatorPublicKey()).
		SetSupplyKey(client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenBurnTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAmount(10).
		SetTokenID(tokenID).
		Execute(client)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

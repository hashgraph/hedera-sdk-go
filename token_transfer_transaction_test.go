package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTokenTransferTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	operatorKey := client.operator.privateKey
	operatorAccountID := client.operator.accountID

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	txID, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	accountID := receipt.GetAccountID()
	assert.NoError(t, err)

	txID, err = NewTokenCreateTransaction().
		SetName("ffff").
		SetSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasury(operatorAccountID).
		SetAdminKey(operatorKey.PublicKey()).
		SetFreezeKey(operatorKey.PublicKey()).
		SetWipeKey(operatorKey.PublicKey()).
		SetKycKey(operatorKey.PublicKey()).
		SetSupplyKey(operatorKey.PublicKey()).
		SetFreezeDefault(false).
		SetExpirationTime(uint64(time.Now().Add(7890000 * time.Second).Unix())).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := receipt.GetTokenID()
	assert.NoError(t, err)

	txID, err = NewTokenMintTransaction().
		SetTokenID(tokenID).
		SetAmount(10).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	txID, err = NewTokenBurnTransaction().
		SetTokenID(tokenID).
		SetAmount(10).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenAssociateTransaction().
		SetAccountID(accountID).
		AddTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenGrantKycTransaction().
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenTransferTransaction().
		AddSender(tokenID, operatorAccountID, 10).
		AddRecipient(tokenID, accountID, 10).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenFreezeTransaction().
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenUnfreezeTransaction().
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenWipeTransaction().
		SetAccountID(accountID).
		SetTokenID(tokenID).
		SetAmount(10).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenRevokeKycTransaction().
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenDissociateTransaction().
		SetAccountID(accountID).
		AddTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenDeleteTransaction().
		SetTokenID(tokenID).
		Execute(client)

	tx, err := NewAccountDeleteTransaction().
		SetDeleteAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(NewTransactionID(accountID)).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}

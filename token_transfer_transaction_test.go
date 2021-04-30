package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTokenTransferTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	txID, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(120)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	accountID1 := receipt.GetAccountID()

	txID, err = NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	accountID2 := receipt.GetAccountID()

	tx, err := NewTokenCreateTransaction().
		SetName("ffff").
		SetSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasury(accountID1).
		SetAdminKey(newKey.PublicKey()).
		SetFreezeKey(newKey.PublicKey()).
		SetWipeKey(newKey.PublicKey()).
		SetKycKey(newKey.PublicKey()).
		SetSupplyKey(newKey.PublicKey()).
		SetFreezeDefault(false).
		SetExpirationTime(uint64(time.Now().Add(7890000 * time.Second).Unix())).
		SetMaxTransactionFee(NewHbar(100)).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	receipt, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := receipt.GetTokenID()
	assert.NoError(t, err)

	tx, err = NewTokenMintTransaction().
		SetTokenID(tokenID).
		SetAmount(10).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	receipt, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	tx, err = NewTokenBurnTransaction().
		SetTokenID(tokenID).
		SetAmount(10).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	tx, err = NewTokenAssociateTransaction().
		SetAccountID(accountID2).
		AddTokenID(tokenID).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	tx, err = NewTokenGrantKycTransaction().
		SetAccountID(accountID2).
		SetTokenID(tokenID).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	tx, err = NewTokenTransferTransaction().
		AddSender(tokenID, accountID1, 10).
		AddRecipient(tokenID, accountID2, 10).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	tx, err = NewTokenFreezeTransaction().
		SetAccountID(accountID2).
		SetTokenID(tokenID).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	tx, err = NewTokenUnfreezeTransaction().
		SetAccountID(accountID2).
		SetTokenID(tokenID).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	tx, err = NewTokenWipeTransaction().
		SetAccountID(accountID2).
		SetTokenID(tokenID).
		SetAmount(10).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	tx, err = NewTokenRevokeKycTransaction().
		SetAccountID(accountID2).
		SetTokenID(tokenID).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	tx, err = NewTokenDissociateTransaction().
		SetAccountID(accountID2).
		AddTokenID(tokenID).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)
}

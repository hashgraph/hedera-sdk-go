package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestTokenTransferTransaction_Execute(t *testing.T) {
	client, err := ClientFromFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	var operatorAccountID AccountID
	var operatorKey Ed25519PrivateKey

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err = AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err = Ed25519PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

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

	txID, err = NewTokenBurnTransaction().
		SetTokenID(tokenID).
		SetAmount(10).
		Execute(client)
	assert.NoError(t, err)

	txID, err = NewTokenAssociateTransaction().
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
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
		SetTokenIDs(tokenID).
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

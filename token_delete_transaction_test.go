package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenDeleteTransaction_Execute(t *testing.T) {
	client := newTestClient(t, true)

	resp, err := NewTokenCreateTransaction().
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

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_TokenDelete_NoKeys(t *testing.T) {
	client := newTestClient(t, true)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	tokenCreate, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(newKey.PublicKey()).
		SetFreezeDefault(false).
		FreezeWith(client)
	assert.NoError(t, err)

	tokenCreate, err = tokenCreate.
		SignWithOperator(client)
	assert.NoError(t, err)

	resp, err := tokenCreate.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_TokenDelete_NoTokenID(t *testing.T) {
	client := newTestClient(t, true)

	resp, err := NewTokenDeleteTransaction().
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOKEN_ID received for transaction %s", resp.TransactionID), err.Error())
	}
}

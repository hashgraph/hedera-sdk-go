package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestTokenInfoQuery_Execute(t *testing.T) {
	var client *Client

	network := os.Getenv("HEDERA_NETWORK")

	if network == "previewnet" {
		client = ClientForPreviewnet()
	}

	var err error
	client, err = ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	client.SetMirrorNetwork([]string{"hcs.previewnet.mirrornode.hedera.com:5600"})

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")
	var operatorKey PrivateKey
	var operatorAccountID AccountID

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err = AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err = PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	resp, err := NewTokenCreateTransaction().
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
		SetExpirationTime(uint64(time.Now().Unix() + 86400*90)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	info, err := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(2)).
		SetTokenID(tokenID).
		Execute(client)

	assert.Equal(t, info.TokenID, tokenID)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(3))
	assert.Equal(t, info.Treasury, operatorAccountID)
	assert.Equal(t, info.AdminKey, operatorKey.PublicKey())
	assert.Equal(t, info.KycKey, operatorKey.PublicKey())
	assert.Equal(t, info.FreezeKey, operatorKey.PublicKey())
	assert.Equal(t, info.WipeKey, operatorKey.PublicKey())
	assert.Equal(t, info.SupplyKey, operatorKey.PublicKey())
	assert.False(t, *info.DefaultFreezeStatus)
	assert.False(t, *info.DefaultKycStatus)

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

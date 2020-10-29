package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestTokenUpdateTransaction_Execute(t *testing.T) {
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

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSymbol("A").
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

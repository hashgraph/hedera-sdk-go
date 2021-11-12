//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitClientFromConfig(t *testing.T) {
	client, err := ClientFromConfig([]byte(testClientJSON))
	assert.NoError(t, err)

	assert.NotNil(t, client)
	assert.Equal(t, 10, len(client.network.network))
	assert.Nil(t, client.operator)
}

func TestUnitClientSetNetwork(t *testing.T) {
	client := ClientForPreviewnet()

	assert.NotNil(t, client)
	assert.Equal(t, 5, len(client.network.network))
	assert.Nil(t, client.operator)

	network := make(map[string]AccountID)
	network["35.237.200.180:50211"] = AccountID{0, 0, 3, nil}
	network["35.186.191.247:50211"] = AccountID{0, 0, 4, nil}
	network["35.192.2.25:50211"] = AccountID{0, 0, 5, nil}
	network["35.199.161.108:50211"] = AccountID{0, 0, 6, nil}
	network["35.203.82.240:50211"] = AccountID{0, 0, 7, nil}

	err := client.SetNetwork(network)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(client.network.network))
}

func TestUnitClientFromConfigWithOperator(t *testing.T) {
	client, err := ClientFromConfig([]byte(testClientJSONWithOperator))
	assert.NoError(t, err)

	assert.NotNil(t, client)

	testOperatorKey, err := PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	assert.NoError(t, err)

	assert.Equal(t, 10, len(client.network.network))
	assert.NotNil(t, client.operator)
	assert.Equal(t, testOperatorKey.keyData, client.operator.privateKey.keyData)
	assert.Equal(t, AccountID{Account: 3}.Account, client.operator.accountID.Account)
}

func TestUnitClientFromConfigWrongType(t *testing.T) {
	_, err := ClientFromConfig([]byte(testClientJSONWrongTypeMirror))
	if err != nil {
		assert.Equal(t, "mirrorNetwork is expected to be either string or an array of strings", err.Error())
	}

	_, err = ClientFromConfig([]byte(testClientJSONWrongTypeNetwork))
	if err != nil {
		assert.Equal(t, "network is expected to be map of string to string, or string", err.Error())
	}
}

package hedera

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testClientJSON string = `{
    "network": {
		"35.237.200.180:50211": "0.0.3",
		"35.186.191.247:50211": "0.0.4",
		"35.192.2.25:50211": "0.0.5",
		"35.199.161.108:50211": "0.0.6",
		"35.203.82.240:50211": "0.0.7",
		"35.236.5.219:50211": "0.0.8",
		"35.197.192.225:50211": "0.0.9",
		"35.242.233.154:50211": "0.0.10",
		"35.240.118.96:50211": "0.0.11",
		"35.204.86.32:50211": "0.0.12"
    },
    "mirrorNetwork": "testnet"
}`

const testClientJSONWithOperator string = `{
    "network": {
		"35.237.200.180:50211": "0.0.3",
		"35.186.191.247:50211": "0.0.4",
		"35.192.2.25:50211": "0.0.5",
		"35.199.161.108:50211": "0.0.6",
		"35.203.82.240:50211": "0.0.7",
		"35.236.5.219:50211": "0.0.8",
		"35.197.192.225:50211": "0.0.9",
		"35.242.233.154:50211": "0.0.10",
		"35.240.118.96:50211": "0.0.11",
		"35.204.86.32:50211": "0.0.12"
    },
    "operator": {
        "accountId": "0.0.3",
        "privateKey": "302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"
    },
    "mirrorNetwork": "testnet"
}`

const testClientJSONWrongTypeMirror string = `{
    "network": "testnet",
    "operator": {
        "accountId": "0.0.3",
        "privateKey": "302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"
    },
 	"mirrorNetwork": 5
}`

const testClientJSONWrongTypeNetwork string = `{
    "network": 1,
    "operator": {
        "accountId": "0.0.3",
        "privateKey": "302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"
    },
 	"mirrorNetwork": ["hcs.testnet.mirrornode.hedera.com:5600"]
}`

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

func TestClientFromConfigWithOperator(t *testing.T) {
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

func TestClientFromConfigWrongType(t *testing.T) {
	_, err := ClientFromConfig([]byte(testClientJSONWrongTypeMirror))
	if err != nil {
		assert.Equal(t, fmt.Sprintf("mirrorNetwork is expected to be either string or an array of strings"), err.Error())
	}

	_, err = ClientFromConfig([]byte(testClientJSONWrongTypeNetwork))
	if err != nil {
		assert.Equal(t, fmt.Sprintf("network is expected to be map of string to string, or string"), err.Error())
	}
}

func TestIntegrationClientPingAllGoodNetwork(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	env.Client.SetMaxNodeAttempts(1)
	env.Client.PingAll()

	net := env.Client.GetNetwork()

	keys := make([]string, len(net))
	val := make([]AccountID, len(net))
	i := 0
	for st, n := range net {
		keys[i] = st
		val[i] = n
		i++
	}

	_, err := NewAccountBalanceQuery().
		SetAccountID(val[0]).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationClientPingAllBadNetwork(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	tempClient := newClient(env.Client.GetNetwork(), env.Client.GetMirrorNetwork(), *env.Client.networkName)
	tempClient.SetOperator(env.OperatorID, env.OperatorKey)

	tempClient.SetMaxNodeAttempts(1)
	tempClient.SetMaxNodesPerTransaction(2)
	tempClient.SetMaxAttempts(3)
	net := tempClient.GetNetwork()
	assert.True(t, len(net) > 1)

	keys := make([]string, len(net))
	val := make([]AccountID, len(net))
	i := 0
	for st, n := range net {
		keys[i] = st
		val[i] = n
		i++
	}

	tempNet := make(map[string]AccountID, 2)
	tempNet["in-process:name"] = val[0]
	tempNet[keys[1]] = val[1]

	err := tempClient.SetNetwork(tempNet)
	assert.NoError(t, err)

	tempClient.PingAll()

	net = tempClient.GetNetwork()
	i = 0
	for st, n := range net {
		keys[i] = st
		val[i] = n
		i++
	}

	_, err = NewAccountBalanceQuery().
		SetAccountID(val[0]).
		Execute(tempClient)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(tempClient.GetNetwork()))

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

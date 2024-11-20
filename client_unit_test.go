//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitClientFromConfig(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfig([]byte(testClientJSON))
	require.NoError(t, err)

	assert.NotNil(t, client)
	assert.True(t, len(client.network.network) > 0)
	assert.Nil(t, client.operator)
}

func TestUnitClientFromConfigWithOperator(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfig([]byte(testClientJSONWithOperator))
	require.NoError(t, err)

	assert.NotNil(t, client)

	testOperatorKey, err := PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	require.NoError(t, err)

	assert.True(t, len(client.network.network) > 0)
	assert.NotNil(t, client.operator)
	assert.Equal(t, testOperatorKey.ed25519PrivateKey.keyData, client.operator.privateKey.ed25519PrivateKey.keyData)
	assert.Equal(t, AccountID{Account: 3}.Account, client.operator.accountID.Account)
}

func TestUnitClientFromConfigWithoutMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfig([]byte(testClientJSONWithoutMirrorNetwork))
	require.NoError(t, err)
	assert.NotNil(t, client)

	assert.True(t, len(client.network.network) > 0)
	assert.True(t, len(client.GetMirrorNetwork()) == 0)
}

func TestUnitClientFromConfigWrongMirrorNetworkType(t *testing.T) {
	t.Parallel()

	_, err := ClientFromConfig([]byte(testClientJSONWrongTypeMirror))
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "mirrorNetwork is expected to be a string, an array of strings or nil", err.Error())
	}
}

func TestUnitClientFromConfigWrongNetworkType(t *testing.T) {
	t.Parallel()

	_, err := ClientFromConfig([]byte(testClientJSONWrongTypeNetwork))
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network is expected to be map of string to string, or string", err.Error())
	}
}

func TestUnitClientFromConfigWrongAccountIDNetworkType(t *testing.T) {
	_, err := ClientFromConfig([]byte(testClientJSONWrongAccountIDNetwork))
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "expected {shard}.{realm}.{num}", err.Error())
	}
}

func TestUnitClientFromCorrectConfigFile(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfigFile("client-config-with-operator.json")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.operator)
	assert.Equal(t, AccountID{Account: 3}.Account, client.operator.accountID.Account)
	assert.Equal(t, "a608e2130a0a3cb34f86e757303c862bee353d9ab77ba4387ec084f881d420d4", client.operator.privateKey.StringRaw())
}

func TestUnitClientFromMissingConfigFile(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfigFile("missing.json")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestUnitClientSetNetworkExtensive(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	nodes := make(map[string]AccountID, 2)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	network := client.GetNetwork()
	assert.Equal(t, 2, len(network))
	assert.Equal(t, network["0.testnet.hedera.com:50211"], AccountID{0, 0, 3, nil, nil, nil})
	assert.Equal(t, network["1.testnet.hedera.com:50211"], AccountID{0, 0, 4, nil, nil, nil})

	nodes = make(map[string]AccountID, 2)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["2.testnet.hedera.com:50211"] = AccountID{0, 0, 5, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	network = client.GetNetwork()
	assert.Equal(t, 3, len(network))
	assert.Equal(t, network["0.testnet.hedera.com:50211"], AccountID{0, 0, 3, nil, nil, nil})
	assert.Equal(t, network["1.testnet.hedera.com:50211"], AccountID{0, 0, 4, nil, nil, nil})
	assert.Equal(t, network["2.testnet.hedera.com:50211"], AccountID{0, 0, 5, nil, nil, nil})

	nodes = make(map[string]AccountID, 1)
	nodes["2.testnet.hedera.com:50211"] = AccountID{0, 0, 5, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	network = client.GetNetwork()
	networkMirror := client.GetMirrorNetwork()
	assert.Equal(t, 1, len(network))
	assert.Equal(t, network["2.testnet.hedera.com:50211"], AccountID{0, 0, 5, nil, nil, nil})
	// There is only one mirror address, no matter the transport security
	assert.Equal(t, "nonexistent-mirror-testnet:443", networkMirror[0])

	client.SetTransportSecurity(true)
	client.SetCertificateVerification(true)
	network = client.GetNetwork()
	networkMirror = client.GetMirrorNetwork()
	assert.Equal(t, network["2.testnet.hedera.com:50212"], AccountID{0, 0, 5, nil, nil, nil})
	assert.Equal(t, "nonexistent-mirror-testnet:443", networkMirror[0])

	err = client.Close()
	require.NoError(t, err)
}

func TestUnitClientSetMirrorNetwork(t *testing.T) {
	t.Parallel()

	mirrorNetworkString := "testnet.mirrornode.hedera.com:443"
	mirrorNetwork1String := "testnet1.mirrornode.hedera.com:443"
	defaultNetwork := make([]string, 0)
	defaultNetwork = append(defaultNetwork, mirrorNetworkString)
	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork(defaultNetwork)

	mirrorNetwork := client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])

	defaultNetworkWithExtraNode := make([]string, 0)
	defaultNetworkWithExtraNode = append(defaultNetworkWithExtraNode, mirrorNetworkString)
	defaultNetworkWithExtraNode = append(defaultNetworkWithExtraNode, mirrorNetwork1String)

	client.SetMirrorNetwork(defaultNetworkWithExtraNode)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 2, len(mirrorNetwork))
	require.True(t, contains(mirrorNetwork, mirrorNetworkString))
	require.True(t, contains(mirrorNetwork, mirrorNetwork1String))

	defaultNetwork = make([]string, 0)
	defaultNetwork = append(defaultNetwork, mirrorNetwork1String)

	client.SetMirrorNetwork(defaultNetwork)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetwork1String, mirrorNetwork[0])

	defaultNetwork = make([]string, 0)
	defaultNetwork = append(defaultNetwork, mirrorNetworkString)

	client.SetMirrorNetwork(defaultNetwork)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])

	client.SetTransportSecurity(true)
	mirrorNetwork = client.GetMirrorNetwork()
	// SetTransportSecurity is deprecated, so the mirror node should not be updated
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])

	err = client.Close()
	require.NoError(t, err)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestUnitClientSetMultipleNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	nodes := make(map[string]AccountID, 8)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["34.94.106.61:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["50.18.132.211:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["138.91.142.219:50211"] = AccountID{0, 0, 3, nil, nil, nil}

	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["35.237.119.55:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["3.212.6.13:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["52.168.76.241:50211"] = AccountID{0, 0, 4, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	net := client.GetNetwork()

	if val, ok := net["0.testnet.hedera.com:50211"]; ok {
		require.Equal(t, val.String(), "0.0.3")
	}

	if val, ok := net["1.testnet.hedera.com:50211"]; ok {
		require.Equal(t, val.String(), "0.0.4")
	}

	if val, ok := net["50.18.132.211:50211"]; ok {
		require.Equal(t, val.String(), "0.0.3")
	}

	if val, ok := net["3.212.6.13:50211"]; ok {
		require.Equal(t, val.String(), "0.0.4")
	}

}

func TestUnitClientLogger(t *testing.T) {
	client := ClientForTestnet()

	var buf bytes.Buffer
	writer := zerolog.ConsoleWriter{Out: &buf, TimeFormat: time.RFC3339}

	hederaLoger := NewLogger("test", LoggerLevelTrace)

	l := zerolog.New(&writer)
	hederaLoger.logger = &l

	client.SetLogger(hederaLoger)
	client.SetLogLevel(LoggerLevelInfo)

	client.logger.Trace("trace message", "traceKey", "traceValue")
	client.logger.Debug("debug message", "debugKey", "debugValue")
	client.logger.Info("info message", "infoKey", "infoValue")
	client.logger.Warn("warn message", "warnKey", "warnValue")
	client.logger.Error("error message", "errorKey", "errorValue")

	assert.NotContains(t, buf.String(), "trace message")
	assert.NotContains(t, buf.String(), "debug message")
	assert.Contains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "error message")

	buf.Reset()
	client.SetLogLevel(LoggerLevelWarn)
	client.logger.Trace("trace message", "traceKey", "traceValue")
	client.logger.Debug("debug message", "debugKey", "debugValue")
	client.logger.Info("info message", "infoKey", "infoValue")
	client.logger.Warn("warn message", "warnKey", "warnValue")
	client.logger.Error("error message", "errorKey", "errorValue")

	assert.NotContains(t, buf.String(), "trace message")
	assert.NotContains(t, buf.String(), "debug message")
	assert.NotContains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "error message")

	buf.Reset()
	client.SetLogLevel(LoggerLevelTrace)
	client.logger.Trace("trace message", "traceKey", "traceValue")
	client.logger.Debug("debug message", "debugKey", "debugValue")
	client.logger.Info("info message", "infoKey", "infoValue")
	client.logger.Warn("warn message", "warnKey", "warnValue")
	client.logger.Error("error message", "errorKey", "errorValue")

	assert.Contains(t, buf.String(), "trace message")
	assert.Contains(t, buf.String(), "debug message")
	assert.Contains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "error message")

	hl := client.GetLogger()
	assert.Equal(t, hl, hederaLoger)
}

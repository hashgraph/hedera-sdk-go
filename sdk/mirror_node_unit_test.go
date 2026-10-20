//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestUnitMirrorNodeGetSchemeErrors(t *testing.T) {
	t.Parallel()

	node := &_MirrorNode{
		_ManagedNode: &_ManagedNode{
			address: &_ManagedNodeAddress{
				address: nil,
				port:    443,
			},
		},
	}

	_, err := node.getScheme()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mirror node address is not set")
}

func TestUnitMirrorNodeGetBaseRestUrl(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		address     string
		expectedURL string
		shouldError bool
	}{
		{
			name:        "HTTPS with port 443",
			address:     "mirror.example.com:443",
			expectedURL: "https://mirror.example.com:443/api/v1",
		},
		{
			name:        "HTTP with port 80",
			address:     "mirror.example.com:80",
			expectedURL: "http://mirror.example.com:80/api/v1",
		},
		{
			name:        "HTTPS with custom port",
			address:     "mirror.example.com:8443",
			expectedURL: "https://mirror.example.com:8443/api/v1",
		},
		{
			name:        "localhost gets special handling",
			address:     "localhost:8080",
			expectedURL: "http://localhost:38081/api/v1",
		},
		{
			name:        "127.0.0.1 gets special handling",
			address:     "127.0.0.1:9999",
			expectedURL: "http://127.0.0.1:38081/api/v1",
		},
		{
			name:        "testnet mirror",
			address:     "testnet.mirrornode.hedera.com:443",
			expectedURL: "https://testnet.mirrornode.hedera.com:443/api/v1",
		},
		{
			name:        "mainnet mirror",
			address:     "mainnet-public.mirrornode.hedera.com:443",
			expectedURL: "https://mainnet-public.mirrornode.hedera.com:443/api/v1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			node, err := _NewMirrorNode(test.address)
			require.NoError(t, err)

			url, err := node.getBaseRestUrl()

			if test.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if test.address == "localhost:8080" {
					assert.Equal(t, "http://localhost:38081/api/v1", url)
				} else {
					assert.Equal(t, test.expectedURL, url)
				}
			}
		})
	}
}

func TestUnitMirrorNodeGetBaseRestUrlErrors(t *testing.T) {
	t.Parallel()

	node := &_MirrorNode{
		_ManagedNode: &_ManagedNode{
			address: &_ManagedNodeAddress{
				address: nil,
				port:    443,
			},
		},
	}

	_, err := node.getBaseRestUrl()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mirror node address is not set")
}

func TestUnitMirrorNodeBackoffMethods(t *testing.T) {
	t.Parallel()

	node, err := _NewMirrorNode("mirror.example.com:443")
	require.NoError(t, err)

	// Test default min backoff
	defaultMinBackoff := node._GetMinBackoff()
	assert.Equal(t, 250*time.Millisecond, defaultMinBackoff)

	// Test setting min backoff
	newMinBackoff := 500 * time.Millisecond
	node._SetMinBackoff(newMinBackoff)
	assert.Equal(t, newMinBackoff, node._GetMinBackoff())

	// Test setting max backoff
	newMaxBackoff := 10 * time.Second
	node._SetMaxBackoff(newMaxBackoff)
	assert.Equal(t, newMaxBackoff, node._GetMaxBackoff())
}

func TestUnitMirrorNodeHealthAndUsage(t *testing.T) {
	t.Parallel()

	node, err := _NewMirrorNode("mirror.example.com:443")
	require.NoError(t, err)

	// Test initial state
	assert.True(t, node._IsHealthy())
	assert.Equal(t, int64(0), node._GetUseCount())
	assert.Equal(t, int64(0), node._GetAttempts())

	// Test using the node
	node._InUse()
	assert.Equal(t, int64(1), node._GetUseCount())

	// Test backoff operations
	node._IncreaseBackoff()
	node._DecreaseBackoff()

	// Test wait time (should not panic)
	waitTime := node._Wait()
	assert.GreaterOrEqual(t, waitTime, time.Duration(0))
}

func TestUnitMirrorNodeAddressOperations(t *testing.T) {
	t.Parallel()

	node, err := _NewMirrorNode("mirror.example.com:443")
	require.NoError(t, err)

	// Test getting address
	address := node._GetAddress()
	assert.Equal(t, "mirror.example.com:443", address)

	// Test secure/insecure conversion
	secureNode := node._ToSecure()
	require.NotNil(t, secureNode)

	insecureNode := node._ToInsecure()
	require.NotNil(t, insecureNode)

	// Test getting managed node
	managedNode := node._GetManagedNode()
	require.NotNil(t, managedNode)
}

func TestUnitMirrorNodeCertificateSettings(t *testing.T) {
	t.Parallel()

	node, err := _NewMirrorNode("mirror.example.com:443")
	require.NoError(t, err)

	// Test certificate verification (always returns false for mirror nodes)
	assert.False(t, node._GetVerifyCertificate())

	// Test setting certificate verification (should not error)
	node._SetVerifyCertificate(true)
	assert.False(t, node._GetVerifyCertificate()) // Still false after setting

	node._SetVerifyCertificate(false)
	assert.False(t, node._GetVerifyCertificate())
}

func TestUnitMirrorNodeClose(t *testing.T) {
	t.Parallel()

	node, err := _NewMirrorNode("mirror.example.com:443")
	require.NoError(t, err)

	// Test closing node (should not error when no client is set)
	err = node._Close()
	assert.NoError(t, err)
}

func TestUnitMirrorNodeStreamBackoff(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	require.False(t, waitMirrorStreamBackoff(ctx, 10))
	require.Less(t, time.Since(start), time.Second)

	require.True(t, waitMirrorStreamBackoff(context.Background(), 0))
}

func setupMirrorServer(t *testing.T, register func(server *grpc.Server)) *Client {
	t.Helper()

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	server := grpc.NewServer()
	register(server)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{listener.Addr().String()})
	t.Cleanup(func() { _ = client.Close() })

	return client
}

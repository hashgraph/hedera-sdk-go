//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"os"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

var callCount int64

func incrementCallCount() {
	atomic.AddInt64(&callCount, 1)
}

func getCallCount() int64 {
	return atomic.LoadInt64(&callCount)
}

func TestIntegrationOneSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	client := ClientForNetwork(env.Client.GetNetwork()).SetOperatorWith(env.OriginalOperatorID, env.OriginalOperatorKey, signingServiceTwo)
	response, err := NewTransferTransaction().
		AddHbarTransfer(env.OriginalOperatorID, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(client)
	require.NoError(t, err)

	_, err = response.GetReceipt(client)
	require.NoError(t, err)

	require.Equal(t, int64(1), getCallCount())
	client.Close()

}

func signingServiceTwo(txBytes []byte) []byte {
	localOperatorPrivateKey, _ := PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	incrementCallCount()

	signature := localOperatorPrivateKey.Sign(txBytes)
	return signature
}

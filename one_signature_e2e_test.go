//go:build all || e2e
// +build all e2e

package hedera

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type Count struct {
	mx    *sync.Mutex
	count int64
}

func NewCount() *Count {
	return &Count{mx: new(sync.Mutex), count: 0}
}

func (c *Count) Incr() {
	c.mx.Lock()
	c.count++
	c.mx.Unlock()
}

func (c *Count) Count() int64 {
	c.mx.Lock()
	count := c.count
	c.mx.Unlock()
	return count
}

var fncCount = NewCount()

func fnc() {
	fncCount.Incr()
}

func TestIntegrationOneSignature(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	client := ClientForNetwork(env.Client.GetNetwork()).SetOperatorWith(env.OriginalOperatorID, env.OriginalOperatorKey, signingServiceTwo)
	response, err := NewTransferTransaction().
		AddHbarTransfer(env.OriginalOperatorID, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(client)
	require.NoError(t, err)

	_, err = response.GetReceipt(client)
	require.NoError(t, err)

	require.Equal(t, int64(1), fncCount.count)
	client.Close()
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func signingServiceTwo(txBytes []byte) []byte {
	localOperatorPrivateKey, _ := PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	go fnc()

	signature := localOperatorPrivateKey.Sign(txBytes)
	return signature
}

//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAddressBookQueryCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	result, err := NewAddressBookQuery().
		SetFileID(FileID{0, 0, 101, nil}).
		SetLimit(14).
		Execute(env.Client)
	require.NoError(t, err)

	//for _, k := range result.NodeAddresses {
	//	println(k.AccountID.String())
	//	for _, s := range k.Addresses {
	//		println(s.String())
	//	}
	//}

	require.NotEqual(t, len(result.NodeAddresses), 0)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

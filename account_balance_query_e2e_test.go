//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountBalanceQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewAccountBalanceQuery().
		SetAccountID(env.OriginalOperatorID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountBalanceQueryCanGetTokenBalance(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	balance, err := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, balance, balance)
	// TODO: assert.Equal(t, uint64(1000000), balance.Tokens.Get(*tokenID))
	// TODO: assert.Equal(t, uint64(3), balance.TokenDecimals.Get(*tokenID))
	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationAccountBalanceQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	balance := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.SetMaxQueryPayment(cost).
		Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountBalanceQuerySetBigMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountBalanceQuerySetSmallMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountBalanceQueryCanSetQueryPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetQueryPayment(NewHbar(0)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.OperatorID)

	cost, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountBalanceQueryCostCanSetPaymentOneTinybar(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetQueryPayment(NewHbar(0)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountBalanceQueryNoAccountIDError(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	assert.True(t, err.Error() == "exceptional precheck status INVALID_ACCOUNT_ID")

}
func TestIntegrationAccountBalanceQueryWorksWithHollowAccountAlias(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	// Extract the ECDSA public key public key
	publicKey := privateKey.PublicKey()
	// Extract the Ethereum public address
	aliasAccountId := *publicKey.ToAccountID(0, 0)
	evmAddress := publicKey.ToEvmAddress()

	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)

	// Transfer tokens using the `TransferTransaction` to the Etherum Account Address
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(4)).
		AddHbarTransfer(env.OperatorID, NewHbar(-4)).Execute(env.Client)
	require.NoError(t, err)

	// Get the child receipt or child record to return the Hiero Account ID for the new account that was created
	_, err = tx.SetIncludeChildren(true).SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = NewAccountBalanceQuery().SetAccountID(aliasAccountId).Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationAccountBalanceQueryCanConnectToMainnetTls(t *testing.T) {
	t.Skip("AccountBalanceQuery is throttled on mainnet")
	t.Parallel()
	client := ClientForMainnet()
	client.SetTransportSecurity(true)

	succeededOnce := false
	for address, accountID := range client.GetNetwork() {

		if !strings.HasSuffix(address, ":50212") {
			t.Errorf("Expected entry key to end with ':50212', but got %s", address)
		}

		accountIDs := []AccountID{accountID}
		_, err := NewAccountBalanceQuery().
			SetNodeAccountIDs(accountIDs).
			SetAccountID(accountID).
			Execute(client)
		if err == nil {
			succeededOnce = true
		}
	}
	assert.True(t, succeededOnce)
}

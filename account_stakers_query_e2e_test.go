//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountStakersQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountStakersQuery().
		SetAccountID(env.Client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(env.Client)
	assert.Error(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountStakersQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)

}

func TestIntegrationAccountStakersQuerySetBigMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(100000)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)

}

func TestIntegrationAccountStakersQuerySetSmallMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)

}

func TestIntegrationAccountStakersQueryInsufficientFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status NOT_SUPPORTED", err.Error())
	}

}

func TestIntegrationAccountStakersQueryNoAccountID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewAccountStakersQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status NOT_SUPPORTED", err.Error())
	}

}

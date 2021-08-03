package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationAccountStakersQueryCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountStakersQuery().
		SetAccountID(env.Client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(env.Client)
	assert.Error(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountStakersQueryGetCost(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountStakersQuerySetBigMaxPayment(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(100000)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountStakersQuerySetSmallMaxPayment(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountStakersQueryInsufficientFee(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED"), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationAccountStakersQueryNoAccountID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountStakersQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED"), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

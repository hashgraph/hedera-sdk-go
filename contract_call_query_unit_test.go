//go:build all || unit
// +build all unit

package hedera

import (
	"bytes"
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitContractCallQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	contractCall := NewContractCallQuery().
		SetContractID(contractID)

	err = contractCall._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitContractCallQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	contractCall := NewContractCallQuery().
		SetContractID(contractID)

	err = contractCall._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestUnitMockContractCallQuery(t *testing.T) {
	message := "getMessage"
	params := NewContractFunctionParameters()
	params._Build(&message)

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_ContractCallLocal{
				ContractCallLocal: &services.ContractCallLocalResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ContractCallLocal{
				ContractCallLocal: &services.ContractCallLocalResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					FunctionResult: &services.ContractFunctionResult{
						ContractID:         &services.ContractID{Contract: &services.ContractID_ContractNum{ContractNum: 123}},
						GasUsed:            75000,
						ContractCallResult: params._Build(&message),
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)

	result, err := NewContractCallQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetContractID(ContractID{Contract: 123}).
		SetQueryPayment(NewHbar(1)).
		SetGas(100000).
		SetFunction(message, nil).
		SetMaxQueryPayment(NewHbar(5)).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, bytes.Compare(result.ContractCallResult, params._Build(&message)), 0)
	require.Equal(t, result.GasUsed, uint64(75000))
	require.Equal(t, result.ContractID.Contract, uint64(123))

	server.Close()
}

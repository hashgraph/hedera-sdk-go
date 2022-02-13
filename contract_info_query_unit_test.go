//go:build all || unit
// +build all unit

package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitContractInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	contractInfoQuery := NewContractInfoQuery().
		SetContractID(contractID)

	err = contractInfoQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitContractInfoQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	contractInfoQuery := NewContractInfoQuery().
		SetContractID(contractID)

	err = contractInfoQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestUnitMockContractInfoQuery(t *testing.T) {
	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_ContractGetInfo{
				ContractGetInfo: &services.ContractGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ContractGetInfo{
				ContractGetInfo: &services.ContractGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					ContractInfo: &services.ContractGetInfoResponse_ContractInfo{
						ContractID:         &services.ContractID{Contract: &services.ContractID_ContractNum{ContractNum: 3}},
						AccountID:          &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 4}},
						ContractAccountID:  "",
						AdminKey:           nil,
						ExpirationTime:     nil,
						AutoRenewPeriod:    nil,
						Storage:            0,
						Memo:               "yes",
						Balance:            0,
						Deleted:            false,
						TokenRelationships: nil,
						LedgerId:           nil,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)

	result, err := NewContractInfoQuery().
		SetContractID(ContractID{Contract: 3}).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, result.ContractID.Contract, uint64(3))
	require.Equal(t, result.AccountID.Account, uint64(4))
	require.Equal(t, result.ContractMemo, "yes")

	server.Close()
}

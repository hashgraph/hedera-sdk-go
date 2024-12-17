//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMirrorNodeContractQuerySetAndGetContractID(t *testing.T) {
	contractID := ContractID{Shard: 0, Realm: 0, Contract: 1234}

	queries := []interface {
		setContractID(ContractID)
		GetContractID() ContractID
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setContractID(contractID)
		assert.Equal(t, contractID, query.GetContractID())
	}
}

func TestMirrorNodeContractQuerySetAndGetSenderEvmAddress(t *testing.T) {
	evmAddress := "0x1234567890abcdef1234567890abcdef12345678"

	queries := []interface {
		setSenderEvmAddress(string)
		GetSenderEvmAddress() string
		GetSender() AccountID
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setSenderEvmAddress(evmAddress)
		assert.Equal(t, evmAddress, query.GetSenderEvmAddress())
		assert.Equal(t, AccountID{}, query.GetSender())
	}
}

func TestMirrorNodeContractQuerySetAndGetContractEvmAddress(t *testing.T) {
	evmAddress := "0x1234567890abcdef1234567890abcdef12345678"

	queries := []interface {
		setContractEvmAddress(string)
		GetContractEvmAddress() string
		GetContractID() ContractID
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setContractEvmAddress(evmAddress)
		assert.Equal(t, evmAddress, query.GetContractEvmAddress())
		assert.Equal(t, ContractID{}, query.GetContractID())
	}
}

func TestMirrorNodeContractQuerySetAndGetCallData(t *testing.T) {
	params := []byte("test")

	queries := []interface {
		setFunctionParameters([]byte)
		GetCallData() []byte
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setFunctionParameters(params)
		assert.Equal(t, params, query.GetCallData())
	}
}

func TestMirrorNodeContractQuerySetFunctionWithoutParameters(t *testing.T) {
	queries := []interface {
		setFunction(string, *ContractFunctionParameters)
		GetCallData() []byte
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setFunction("myFunction", nil)
		assert.NotNil(t, query.GetCallData())
	}
}

func TestMirrorNodeContractQuerySetAndGetBlockNumber(t *testing.T) {
	blockNumber := int64(123456)

	queries := []interface {
		setBlockNumber(int64)
		GetBlockNumber() int64
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setBlockNumber(blockNumber)
		assert.Equal(t, blockNumber, query.GetBlockNumber())
	}
}

func TestMirrorNodeContractQuerySetAndGetValue(t *testing.T) {
	value := int64(1000)

	queries := []interface {
		setValue(int64)
		GetValue() int64
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setValue(value)
		assert.Equal(t, value, query.GetValue())
	}
}

func TestMirrorNodeContractQuerySetAndGetGasLimit(t *testing.T) {
	gas := int64(50000)

	queries := []interface {
		setGasLimit(int64)
		GetGasLimit() int64
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setGasLimit(gas)
		assert.Equal(t, gas, query.GetGasLimit())
	}
}

func TestMirrorNodeContractQuerySetAndGetGasPrice(t *testing.T) {
	gasPrice := int64(200)

	queries := []interface {
		setGasPrice(int64)
		GetGasPrice() int64
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
		NewMirrorNodeContractCallQuery(),
	}

	for _, query := range queries {
		query.setGasPrice(gasPrice)
		assert.Equal(t, gasPrice, query.GetGasPrice())
	}
}

func TestMirrorNodeContractQueryEstimateGasWithMissingContractIDOrEvmAddressThrowsException(t *testing.T) {
	queries := []interface {
		setFunction(string, *ContractFunctionParameters)
		estimateGas(*Client) (uint64, error)
	}{
		&mirrorNodeContractQuery{},
		NewMirrorNodeContractEstimateGasQuery(),
	}

	for _, query := range queries {
		query.setFunction("testFunction", NewContractFunctionParameters().AddString("params"))
		_, err := query.estimateGas(nil)
		require.Error(t, err)
	}
}

func TestMirrorNodeContractQueryCreateJSONPayloadAllFieldsSet(t *testing.T) {
	query := &mirrorNodeContractQuery{
		callData:           []byte("testData"),
		senderEvmAddress:   stringPtr("0x1234567890abcdef1234567890abcdef12345678"),
		contractEvmAddress: stringPtr("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
		gasLimit:           int64Ptr(50000),
		gasPrice:           int64Ptr(2000),
		value:              int64Ptr(1000),
		blockNumber:        int64Ptr(123456),
	}

	jsonPayload, err := query.createJSONPayload(true, "latest")
	assert.NoError(t, err)

	expectedJson := `{"data":"7465737444617461","to":"0xabcdefabcdefabcdefabcdefabcdefabcdef","estimate":true,"blockNumber":"latest","from":"0x1234567890abcdef1234567890abcdef12345678","gas":50000,"gasPrice":2000,"value":1000}`
	assert.JSONEq(t, expectedJson, jsonPayload)
}

func TestMirrorNodeContractQueryCreateJSONPayloadOnlyRequiredFieldsSet(t *testing.T) {
	query := &mirrorNodeContractQuery{
		callData:           []byte("testData"),
		contractEvmAddress: stringPtr("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
	}

	jsonPayload, err := query.createJSONPayload(true, "latest")
	assert.NoError(t, err)

	expectedJson := `{"data":"7465737444617461","to":"0xabcdefabcdefabcdefabcdefabcdefabcdef","estimate":true,"blockNumber":"latest"}`
	assert.JSONEq(t, expectedJson, jsonPayload)
}

func TestMirrorNodeContractQueryCreateJSONPayloadSomeOptionalFieldsSet(t *testing.T) {
	query := &mirrorNodeContractQuery{
		callData:           []byte("testData"),
		senderEvmAddress:   stringPtr("0x1234567890abcdef1234567890abcdef12345678"),
		contractEvmAddress: stringPtr("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
		gasLimit:           int64Ptr(50000),
		value:              int64Ptr(1000),
	}

	jsonPayload, err := query.createJSONPayload(false, "latest")
	assert.NoError(t, err)

	expectedJson := `{"data":"7465737444617461","to":"0xabcdefabcdefabcdefabcdefabcdefabcdef","estimate":false,"blockNumber":"latest","from":"0x1234567890abcdef1234567890abcdef12345678","gas":50000,"value":1000}`
	assert.JSONEq(t, expectedJson, jsonPayload)
}

func TestMirrorNodeContractQueryCreateJSONPayloadAllOptionalFieldsDefault(t *testing.T) {
	query := &mirrorNodeContractQuery{
		callData:           []byte("testData"),
		contractEvmAddress: stringPtr("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
	}

	jsonPayload, err := query.createJSONPayload(false, "latest")
	assert.NoError(t, err)

	expectedJson := `{"data":"7465737444617461","to":"0xabcdefabcdefabcdefabcdefabcdefabcdef","estimate":false,"blockNumber":"latest"}`
	assert.JSONEq(t, expectedJson, jsonPayload)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

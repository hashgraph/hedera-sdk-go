package hiero

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// SPDX-License-Identifier: Apache-2.0

// mirrorNodeContractQuery returns a result from EVM execution such as cost-free execution of read-only smart contract
// queries, gas estimation, and transient simulation of read-write operations.
type mirrorNodeContractQuery struct {
	// The contract we are sending the transaction to
	contractID         *ContractID
	contractEvmAddress *string
	// The account we are sending the transaction from
	sender           *AccountID
	senderEvmAddress *string
	// The transaction callData
	callData []byte
	// The amount we are sending to the contract
	value *int64
	// The gas limit
	gasLimit *int64
	// The gas price
	gasPrice *int64
	// The block number for the simulation
	blockNumber *int64
}

// setContractID sets the contract instance to call
func (query *mirrorNodeContractQuery) setContractID(contractID ContractID) {
	query.contractID = &contractID
}

// GetContractID returns the contract instance to call
func (query *mirrorNodeContractQuery) GetContractID() ContractID {
	if query.contractID == nil {
		return ContractID{}
	}

	return *query.contractID
}

// setContractEvmAddress Set the 20-byte EVM address of the contract to call.
func (query *mirrorNodeContractQuery) setContractEvmAddress(contractEvmAddress string) {
	query.contractEvmAddress = &contractEvmAddress
	query.contractID = nil
}

// GetContractEvmAddress returns the 20-byte EVM address of the contract to call.
func (query *mirrorNodeContractQuery) GetContractEvmAddress() string {
	if query.contractEvmAddress == nil {
		return ""
	}
	return *query.contractEvmAddress
}

// setSender sets the sender of the transaction simulation.
func (query *mirrorNodeContractQuery) setSender(sender AccountID) {
	query.sender = &sender
}

// GetSender returns the sender of the transaction simulation.
func (query *mirrorNodeContractQuery) GetSender() AccountID {
	if query.sender == nil {
		return AccountID{}
	}

	return *query.sender
}

// setSenderEvmAddress Set the 20-byte EVM address of the sender of the transaction simulation.
func (query *mirrorNodeContractQuery) setSenderEvmAddress(senderEvmAddress string) {
	query.senderEvmAddress = &senderEvmAddress
	query.sender = nil
}

// GetSenderEvmAddress returns the 20-byte EVM address of the sender of the transaction simulation.
func (query *mirrorNodeContractQuery) GetSenderEvmAddress() string {
	if query.senderEvmAddress == nil {
		return ""
	}

	return *query.senderEvmAddress
}

// setFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (query *mirrorNodeContractQuery) setFunction(name string, params *ContractFunctionParameters) {
	if params == nil {
		params = NewContractFunctionParameters()
	}

	query.callData = params._Build(&name)
}

// setFunctionParameters sets the function parameters as their raw bytes.
func (q *mirrorNodeContractQuery) setFunctionParameters(byteArray []byte) {
	q.callData = byteArray
}

// GetCallData returns the calldata
func (query *mirrorNodeContractQuery) GetCallData() []byte {
	return query.callData
}

// setValue Sets the amount of value (in tinybars or wei) to be sent to the contract in the transaction.
func (query *mirrorNodeContractQuery) setValue(value int64) {
	query.value = &value
}

// GetValue returns the amount of value (in tinybars or wei) to be sent to the contract in the transaction.
func (query *mirrorNodeContractQuery) GetValue() int64 {
	if query.value == nil {
		return 0
	}

	return *query.value
}

// setGasLimit Sets the gas limit for the contract call. This specifies the maximum amount of gas that the transaction can consume.
func (query *mirrorNodeContractQuery) setGasLimit(gasLimit int64) {
	query.gasLimit = &gasLimit
}

// GetGasLimit returns the gas limit for the contract call. This specifies the maximum amount of gas that the transaction can consume.
func (query *mirrorNodeContractQuery) GetGasLimit() int64 {
	if query.gasLimit == nil {
		return 0
	}

	return *query.gasLimit
}

// setGasPrice Sets the gas price to be used for the contract call. This specifies the price of each unit of gas used in the transaction.
func (query *mirrorNodeContractQuery) setGasPrice(gasPrice int64) {
	query.gasPrice = &gasPrice
}

// GetGasPrice returns the gas price to be used for the contract call. This specifies the price of each unit of gas used in the transaction.
func (query *mirrorNodeContractQuery) GetGasPrice() int64 {
	if query.gasPrice == nil {
		return 0
	}

	return *query.gasPrice
}

// setBlockNumber Sets the block number for the simulation of the contract call. The block number determines the context of the contract call simulation within the blockchain.
func (query *mirrorNodeContractQuery) setBlockNumber(blockNumber int64) {
	query.blockNumber = &blockNumber
}

// GetBlockNumber returns the block number for the simulation of the contract call. The block number determines the context of the contract call simulation within the blockchain.
func (query *mirrorNodeContractQuery) GetBlockNumber() int64 {
	if query.blockNumber == nil {
		return 0
	}

	return *query.blockNumber
}

// Returns gas estimation for the EVM execution
func (query *mirrorNodeContractQuery) estimateGas(client *Client) (uint64, error) {
	err := query.fillEvmAddresses()
	if err != nil {
		return 0, err
	}

	jsonPayload, err := query.createJSONPayload(true, "latest")
	if err != nil {
		return 0, err
	}

	result, err := query.performContractCallToMirrorNode(client, jsonPayload)
	if err != nil {
		return 0, err
	}

	hexString, ok := result["result"].(string)
	if !ok {
		return 0, fmt.Errorf("result is not a string")
	}
	hexString = strings.TrimPrefix(hexString, "0x")
	gas, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse the result: %w", err)
	}
	return gas, nil
}

// Does transient simulation of read-write operations and returns the result in hexadecimal string format. The result can be any solidity type.
func (query *mirrorNodeContractQuery) call(client *Client) (string, error) {
	err := query.fillEvmAddresses()
	if err != nil {
		return "", err
	}

	var blockNumber string
	if query.blockNumber == nil {
		blockNumber = "latest"
	} else {
		blockNumber = fmt.Sprintf("%d", *query.blockNumber)
	}
	jsonPayload, err := query.createJSONPayload(false, blockNumber)
	if err != nil {
		return "", err
	}

	result, err := query.performContractCallToMirrorNode(client, jsonPayload)
	if err != nil {
		return "", err
	}

	hexString, ok := result["result"].(string)
	if !ok {
		return "", fmt.Errorf("result is not a string")
	}
	return hexString, nil
}

// Retrieve and set the evm addresses if necessary
func (query *mirrorNodeContractQuery) fillEvmAddresses() error {
	// fill contractEvmAddress
	if query.contractEvmAddress == nil {
		if query.contractID == nil {
			return errors.New("contractID is not set")
		}
		address := query.contractID.ToSolidityAddress()
		query.contractEvmAddress = &address
	}

	// fill senderEvmAddress
	if query.senderEvmAddress == nil && query.sender != nil {
		address := query.sender.ToSolidityAddress()
		query.senderEvmAddress = &address
	}
	return nil
}

func (query *mirrorNodeContractQuery) performContractCallToMirrorNode(client *Client, jsonPayload string) (map[string]any, error) {
	if client.mirrorNetwork == nil || len(client.GetMirrorNetwork()) == 0 {
		return nil, errors.New("mirror node is not set")
	}
	mirrorUrl := client.GetMirrorNetwork()[0]
	index := strings.Index(mirrorUrl, ":")
	if index == -1 {
		return nil, errors.New("invalid mirrorUrl format")
	}
	mirrorUrl = mirrorUrl[:index]

	var url string
	protocol := httpsString
	port := ""

	if client.GetLedgerID().String() == "" {
		protocol = httpString
		port = ":8545"
	}
	url = fmt.Sprintf("%s://%s%s/api/v1/contracts/call", protocol, mirrorUrl, port)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(jsonPayload))) // #nosec
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 response from Mirror Node: %d, details: %s", resp.StatusCode, body)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (query *mirrorNodeContractQuery) createJSONPayload(estimate bool, blockNumber string) (string, error) {
	hexData := hex.EncodeToString(query.callData)

	payload := map[string]any{
		"data":        hexData,
		"to":          query.contractEvmAddress,
		"estimate":    estimate,
		"blockNumber": blockNumber,
	}

	// Conditionally add fields if they are set to non-default values
	if query.senderEvmAddress != nil {
		payload["from"] = query.senderEvmAddress
	}
	if query.gasLimit != nil {
		payload["gas"] = query.gasLimit
	}
	if query.gasPrice != nil {
		payload["gasPrice"] = query.gasPrice
	}
	if query.value != nil {
		payload["value"] = query.value
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonBytes), nil
}

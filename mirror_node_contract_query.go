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

// GetContractID returns the contract instance to call
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetContractID() ContractID {
	if mirrorNodeContractQuery.contractID == nil {
		return ContractID{}
	}

	return *mirrorNodeContractQuery.contractID
}

// GetContractEvmAddress returns the 20-byte EVM address of the contract to call.
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetContractEvmAddress() string {
	if mirrorNodeContractQuery.contractEvmAddress == nil {
		return ""
	}
	return *mirrorNodeContractQuery.contractEvmAddress
}

// GetSender returns the sender of the transaction simulation.
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetSender() AccountID {
	if mirrorNodeContractQuery.sender == nil {
		return AccountID{}
	}

	return *mirrorNodeContractQuery.sender
}

// GetSenderEvmAddress returns the 20-byte EVM address of the sender of the transaction simulation.
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetSenderEvmAddress() string {
	if mirrorNodeContractQuery.senderEvmAddress == nil {
		return ""
	}

	return *mirrorNodeContractQuery.senderEvmAddress
}

// setFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (mirrorNodeContractQuery *mirrorNodeContractQuery) setFunction(name string, params *ContractFunctionParameters) {
	if params == nil {
		params = NewContractFunctionParameters()
	}

	mirrorNodeContractQuery.callData = params._Build(&name)
}

// GetCallData returns the calldata
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetCallData() []byte {
	return mirrorNodeContractQuery.callData
}

// GetValue returns the amount of value (in tinybars or wei) to be sent to the contract in the transaction.
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetValue() int64 {
	if mirrorNodeContractQuery.value == nil {
		return 0
	}

	return *mirrorNodeContractQuery.value
}

// GetGasLimit returns the gas limit for the contract call. This specifies the maximum amount of gas that the transaction can consume.
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetGasLimit() int64 {
	if mirrorNodeContractQuery.gasLimit == nil {
		return 0
	}

	return *mirrorNodeContractQuery.gasLimit
}

// GetGasPrice returns the gas price to be used for the contract call. This specifies the price of each unit of gas used in the transaction.
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetGasPrice() int64 {
	if mirrorNodeContractQuery.gasPrice == nil {
		return 0
	}

	return *mirrorNodeContractQuery.gasPrice
}

// GetBlockNumber returns the block number for the simulation of the contract call. The block number determines the context of the contract call simulation within the blockchain.
func (mirrorNodeContractQuery *mirrorNodeContractQuery) GetBlockNumber() int64 {
	if mirrorNodeContractQuery.blockNumber == nil {
		return 0
	}

	return *mirrorNodeContractQuery.blockNumber
}

// Returns gas estimation for the EVM execution
func (mirrorNodeContractQuery *mirrorNodeContractQuery) estimateGas(client *Client) (uint64, error) {
	err := mirrorNodeContractQuery.fillEvmAddresses()
	if err != nil {
		return 0, err
	}

	jsonPayload, err := mirrorNodeContractQuery.createJSONPayload(true, "latest")
	if err != nil {
		return 0, err
	}

	result, err := mirrorNodeContractQuery.performContractCallToMirrorNode(client, jsonPayload)
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
func (mirrorNodeContractQuery *mirrorNodeContractQuery) call(client *Client) (string, error) {
	err := mirrorNodeContractQuery.fillEvmAddresses()
	if err != nil {
		return "", err
	}

	var blockNumber string
	if mirrorNodeContractQuery.blockNumber == nil {
		blockNumber = "latest"
	} else {
		blockNumber = fmt.Sprintf("%d", *mirrorNodeContractQuery.blockNumber)
	}
	jsonPayload, err := mirrorNodeContractQuery.createJSONPayload(false, blockNumber)
	if err != nil {
		return "", err
	}

	result, err := mirrorNodeContractQuery.performContractCallToMirrorNode(client, jsonPayload)
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
func (mirrorNodeContractQuery *mirrorNodeContractQuery) fillEvmAddresses() error {
	// fill contractEvmAddress
	if mirrorNodeContractQuery.contractEvmAddress == nil {
		if mirrorNodeContractQuery.contractID == nil {
			return errors.New("contractID is not set")
		}
		address := mirrorNodeContractQuery.contractID.ToSolidityAddress()
		mirrorNodeContractQuery.contractEvmAddress = &address
	}

	// fill senderEvmAddress
	if mirrorNodeContractQuery.senderEvmAddress == nil && mirrorNodeContractQuery.sender != nil {
		address := mirrorNodeContractQuery.sender.ToSolidityAddress()
		mirrorNodeContractQuery.senderEvmAddress = &address
	}
	return nil
}

func (mirrorNodeContractQuery *mirrorNodeContractQuery) performContractCallToMirrorNode(client *Client, jsonPayload string) (map[string]any, error) {
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
	protocol := "https"
	port := ""

	if client.GetLedgerID().String() == "" {
		protocol = "http"
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

func (mirrorNodeContractQuery *mirrorNodeContractQuery) createJSONPayload(estimate bool, blockNumber string) (string, error) {
	hexData := hex.EncodeToString(mirrorNodeContractQuery.callData)

	payload := map[string]any{
		"data":        hexData,
		"to":          mirrorNodeContractQuery.contractEvmAddress,
		"estimate":    estimate,
		"blockNumber": blockNumber,
	}

	// Conditionally add fields if they are set to non-default values
	if mirrorNodeContractQuery.senderEvmAddress != nil {
		payload["from"] = mirrorNodeContractQuery.senderEvmAddress
	}
	if mirrorNodeContractQuery.gasLimit != nil {
		payload["gas"] = mirrorNodeContractQuery.gasLimit
	}
	if mirrorNodeContractQuery.gasPrice != nil {
		payload["gasPrice"] = mirrorNodeContractQuery.gasPrice
	}
	if mirrorNodeContractQuery.value != nil {
		payload["value"] = mirrorNodeContractQuery.value
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonBytes), nil
}

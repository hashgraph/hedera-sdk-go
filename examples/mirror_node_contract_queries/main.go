package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

var contractBytecode = "60806040526040518060400160405280600581526020017f68656c6c6f0000000000000000000000000000000000000000000000000000008152505f90816100479190610293565b50348015610053575f80fd5b50610362565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806100d457607f821691505b6020821081036100e7576100e6610090565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026101497fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8261010e565b610153868361010e565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61019761019261018d8461016b565b610174565b61016b565b9050919050565b5f819050919050565b6101b08361017d565b6101c46101bc8261019e565b84845461011a565b825550505050565b5f90565b6101d86101cc565b6101e38184846101a7565b505050565b5b81811015610206576101fb5f826101d0565b6001810190506101e9565b5050565b601f82111561024b5761021c816100ed565b610225846100ff565b81016020851015610234578190505b610248610240856100ff565b8301826101e8565b50505b505050565b5f82821c905092915050565b5f61026b5f1984600802610250565b1980831691505092915050565b5f610283838361025c565b9150826002028217905092915050565b61029c82610059565b67ffffffffffffffff8111156102b5576102b4610063565b5b6102bf82546100bd565b6102ca82828561020a565b5f60209050601f8311600181146102fb575f84156102e9578287015190505b6102f38582610278565b86555061035a565b601f198416610309866100ed565b5f5b828110156103305784890151825560018201915060208501945060208101905061030b565b8683101561034d5784890151610349601f89168261025c565b8355505b6001600288020188555050505b505050505050565b6102178061036f5f395ff3fe608060405234801561000f575f80fd5b5060043610610029575f3560e01c8063ce6d41de1461002d575b5f80fd5b61003561004b565b6040516100429190610164565b60405180910390f35b60605f8054610059906101b1565b80601f0160208091040260200160405190810160405280929190818152602001828054610085906101b1565b80156100d05780601f106100a7576101008083540402835291602001916100d0565b820191905f5260205f20905b8154815290600101906020018083116100b357829003601f168201915b5050505050905090565b5f81519050919050565b5f82825260208201905092915050565b5f5b838110156101115780820151818401526020810190506100f6565b5f8484015250505050565b5f601f19601f8301169050919050565b5f610136826100da565b61014081856100e4565b93506101508185602086016100f4565b6101598161011c565b840191505092915050565b5f6020820190508181035f83015261017c818461012c565b905092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806101c857607f821691505b6020821081036101db576101da610184565b5b5091905056fea26469706673582212202a86c27939bfab6d4a2c61ebbf096d8424e17e22dfdd42320f6e2654863581e964736f6c634300081a0033"

/**
 * @summary Mirror Node Contract Queries for estimation or simulation of transactions
 */
func main() {
	var client *hiero.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	fmt.Println("Example Start!")

	/*
	 * Step 1: Create a contract
	 */

	hexBytecode, err := hex.DecodeString(contractBytecode)
	if err != nil {
		panic(fmt.Sprintf("%v : error decoding bytecode", err))
	}

	resp, err := hiero.NewContractCreateTransaction().
		SetBytecode(hexBytecode).
		SetGas(2000000).
		SetTransactionMemo("Create a simple contract instance").
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating contract", err))
	}

	receipt, err := resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving receipt", err))
	}
	contractId := receipt.ContractID
	fmt.Println("Created new contract with ID: ", contractId)

	/*
	 * Step 2: Wait for mirror node to import data
	 */
	time.Sleep(5 * time.Second)

	/*
	 * Step 3: Estimate the gas needed
	 */
	gas, err := hiero.NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractId).
		SetSender(operatorAccountID).
		SetGasLimit(30_000).
		SetGasPrice(1234).
		SetFunction("getMessage", nil).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error estimating gas", err))
	}

	fmt.Printf("Gas needed for this query: %v\n", gas)

	/*
	 * Step 4: Do the query against the consensus node using the estimated gas
	 */

	_, err = hiero.NewContractCallQuery().
		SetContractID(*contractId).
		SetGas(gas).
		SetFunction("getMessage", nil).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error executing contract call query", err))
	}

	/*
	 * Step 5: Simulate the transaction for free, using the mirror node
	 */

	result, err := hiero.NewMirrorNodeContractCallQuery().
		SetContractID(*contractId).
		SetSender(operatorAccountID).
		SetGasLimit(30_000).
		SetGasPrice(1234).
		SetBlockNumber(10000).
		SetFunction("getMessage", nil).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error during contract call", err))
	}

	decodedResult, err := DecodeABIHexString(result)
	if err != nil {
		panic(fmt.Sprintf("%v : error decoding result", err))
	}

	fmt.Println("Simulation result: ", decodedResult)
	fmt.Println("Contract call result: ", decodedResult)

	/*
	 * Clean up:
	 */
	client.Close()

	fmt.Println("Example Complete!")

}

// DecodeABIHexString decodes a hex-encoded ABI (Application Binary Interface) string into a UTF-8 string.
// It assumes the input follows the ABI encoding standard for dynamic data.
func DecodeABIHexString(hexStr string) (string, error) {
	// Remove the `0x` prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// Ensure the input has at least 128 characters (for metadata and length)
	if len(hexStr) < 128 {
		return "", fmt.Errorf("input hex string is too short")
	}

	// Parse the length of the dynamic data (64 to 128) as a hexadecimal integer
	lengthHex := hexStr[64:128]
	length, err := parseHexToInt(lengthHex)
	if err != nil {
		return "", fmt.Errorf("failed to parse length from hex string: %w", err)
	}

	// Calculate the start and end index for the dynamic data
	startIndex := 128
	endIndex := startIndex + length*2
	if len(hexStr) < endIndex {
		return "", fmt.Errorf("input hex string is too short for the expected dynamic data")
	}

	// Extract the dynamic data substring
	hexData := hexStr[startIndex:endIndex]

	// Convert the hex data to bytes
	bytes, err := hex.DecodeString(hexData)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex data: %w", err)
	}

	// Convert the bytes to a UTF-8 string and return
	return string(bytes), nil
}

// parseHexToInt parses a hexadecimal string into an integer
func parseHexToInt(hexStr string) (int, error) {
	if len(hexStr) == 0 {
		return 0, fmt.Errorf("empty hex string")
	}
	parsed, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, err
	}
	// Convert parsed bytes to an integer
	result := 0
	for _, b := range parsed {
		result = result*256 + int(b)
	}
	return result, nil
}

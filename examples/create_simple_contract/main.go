package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

// a simple contract struct
type contract struct {
	// ignore the link references since it is empty
	Object    string `json:"object"`
	OpCodes   string `json:"opcodes"`
	SourceMap string `json:"sourceMap"`
}

func main() {
	var client *hedera.Client
	var err error

	net := os.Getenv("HEDERA_NETWORK")
	client, err = hedera.ClientForName(net)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(fmt.Sprintf("%v : error converting string to AccountID", err))
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	defer func() {
		err = client.Close()
		if err != nil {
			panic(fmt.Sprintf("%v : error closing client", err))
		}
	}()

	// R contents from hello_world.json file
	rawContract, err := os.ReadFile("./hello_world.json")
	if err != nil {
		panic(fmt.Sprintf("%v : error reading hello_world.json", err))
	}

	// Initialize simple contract
	contract := contract{}

	// Unmarshal the json read from the file into the simple contract
	err = json.Unmarshal([]byte(rawContract), &contract)
	if err != nil {
		panic(fmt.Sprintf("%v : error unmarshaling the json file", err))
	}

	// Convert contract to bytes
	contractByteCode := []byte(contract.Object)

	fmt.Println("Simple contract example")
	fmt.Printf("Contract bytecode size: %v bytes\n", len(contractByteCode))

	// Upload a file containing the byte code
	byteCodeTransactionID, err := hedera.NewFileCreateTransaction().
		SetMaxTransactionFee(hedera.NewHbar(2)).
		// All keys at the top level of a key list must sign to create or modify the file
		SetKeys(client.GetOperatorPublicKey()).
		// Initial contents, in our case it's the contract object converted to bytes
		SetContents(contractByteCode).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error creating file", err))
	}

	// Get the record
	byteCodeTransactionRecord, err := byteCodeTransactionID.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting file creation record", err))
	}

	fmt.Printf("contract bytecode file upload fee: %v\n", byteCodeTransactionRecord.TransactionFee)

	// Get the file ID from the record we got
	byteCodeFileID := *byteCodeTransactionRecord.Receipt.FileID

	fmt.Printf("contract bytecode file: %v\n", byteCodeFileID)

	// Instantiate the contract instance
	contractTransactionResponse, err := hedera.NewContractCreateTransaction().
		// Failing to set this to a sufficient amount will result in "INSUFFICIENT_GAS" status
		SetGas(100000).
		// The file ID we got from the record of the file created previously
		SetBytecodeFileID(byteCodeFileID).
		// Setting an admin key allows you to delete the contract in the future
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error creating contract", err))
	}

	// get the record for the contract we created
	contractRecord, err := contractTransactionResponse.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving contract creation record", err))
	}

	contractCreateResult, err := contractRecord.GetContractCreateResult()
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving contract creation result", err))
	}

	// get the contract ID from the record
	newContractID := *contractRecord.Receipt.ContractID

	fmt.Printf("Contract create gas used: %v\n", contractCreateResult.GasUsed)
	fmt.Printf("Contract create transaction fee: %v\n", contractRecord.TransactionFee)
	fmt.Printf("Contract: %v\n", newContractID)

	// Call the contract to receive the greeting
	callResult, err := hedera.NewContractCallQuery().
		SetContractID(newContractID).
		// The amount of gas to use for the call
		// All of the gas offered will be used and charged a corresponding fee
		SetGas(100000).
		// This query requires payment, depends on gas used
		SetQueryPayment(hedera.NewHbar(1)).
		// Specified which function to call, and the parameters to pass to the function
		SetFunction("greet", nil).
		// This requires payment
		SetMaxQueryPayment(hedera.NewHbar(5)).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error executing contract call query", err))
	}

	fmt.Printf("Call gas used: %v\n", callResult.GasUsed)
	fmt.Printf("Message: %v\n", callResult.GetString(0))

	// Clean up, delete the transaction
	deleteTransactionResponse, err := hedera.NewContractDeleteTransaction().
		// Only thing required here is the contract ID
		SetContractID(newContractID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error deleting contract", err))
	}

	deleteTransactionReceipt, err := deleteTransactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving contract delete receipt", err))
	}

	fmt.Printf("Status of transaction deletion: %v\n", deleteTransactionReceipt.Status)
}

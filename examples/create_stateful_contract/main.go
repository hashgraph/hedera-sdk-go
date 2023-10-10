package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

type contract struct {
	Abi string `json:"abi"`
	Bin string `json:"bin"`
}

type contracts struct {
	Contracts map[string]contract `json:"contracts"`
	Version   string              `json:"version"`
}

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	// Make sure to close client after running
	defer func() {
		err = client.Close()
		if err != nil {
			panic(fmt.Sprintf("%v : error closing client", err))
		}
	}()

	// Read in the compiled contract from stateful.json
	rawSmartContract, err := os.ReadFile("./stateful.json")
	if err != nil {
		panic(fmt.Sprintf("%v : error reading stateful.json", err))
	}

	// Initialize contracts
	var smartContract contracts = contracts{}

	// Parse the rawSmartContract into smartContract
	err = json.Unmarshal([]byte(rawSmartContract), &smartContract)
	if err != nil {
		panic(fmt.Sprintf("%v : error unmarshaling", err))
	}

	// Retrieve the bytecode from the parsed smart contract
	smartContractByteCode := smartContract.Contracts["stateful.sol:StatefulContract"].Bin

	fmt.Println("Stateful contract example")
	fmt.Printf("Contract bytecode size: %v bytes\n", len(smartContractByteCode))

	// Upload a file containing the byte code
	byteCodeTransactionResponse, err := hedera.NewFileCreateTransaction().
		// A file is not implicitly owned by anyone, even the operator
		// But we do use operator's key for this one
		SetKeys(client.GetOperatorPublicKey()).
		// Set the stateful contract bytes for this
		SetContents([]byte(smartContractByteCode)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating file", err))
	}

	// Retrieve the receipt to make sure the transaction went through and to get bytecode file ID
	byteCodeTransactionReceipt, err := byteCodeTransactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting file create transaction receipt", err))
	}

	// Retrieve bytecode file ID from the receipt
	byteCodeFileID := *byteCodeTransactionReceipt.FileID

	fmt.Printf("contract bytecode file: %v\n", byteCodeFileID)

	// Set the parameters that should be passed to the contract constructor
	// In this case we are passing in a string with the value "hello from hedera!"
	// as the only parameter that is passed to the contract
	contractFunctionParams := hedera.NewContractFunctionParameters().
		AddString("hello from hedera")

	// Instantiate the contract instance
	contractTransactionID, err := hedera.NewContractCreateTransaction().
		// Set gas to create the contract
		// Failing to set this to a sufficient amount will result in "INSUFFICIENT_GAS" status
		SetGas(100000).
		// Failing to set parameters when required will result in "CONTRACT_REVERT_EXECUTED" status
		SetConstructorParameters(contractFunctionParams).
		// The contract bytecode must be set to the file ID containing the contract bytecode
		SetBytecodeFileID(byteCodeFileID).
		// Set the admin key on the contract in case the contract should be deleted or
		// updated in the future
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error creating contract", err))
	}

	// Get the new contract record to make sure the transaction ran successfully
	contractRecord, err := contractTransactionID.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving contract creation record", err))
	}

	// Get the contract create result from the record
	contractCreateResult, err := contractRecord.GetContractCreateResult()
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving contract creation result", err))
	}

	// Get the new contract ID from the receipt contained in the record
	newContractID := *contractRecord.Receipt.ContractID

	fmt.Printf("Contract create gas used: %v\n", contractCreateResult.GasUsed)
	fmt.Printf("Contract create transaction fee: %v\n", contractRecord.TransactionFee)
	fmt.Printf("contract: %v\n", newContractID)

	// Ask for the current message (set on creation)
	callResult, err := hedera.NewContractCallQuery().
		// Set which contract
		SetContractID(newContractID).
		// The amount of gas to use for the call
		// All of the gas offered will be used and charged a corresponding fee
		SetGas(100000).
		// This query requires payment, depends on gas used
		SetQueryPayment(hedera.NewHbar(1)).
		// nil -> no parameters
		// Specified which function to call, and the parameters to pass to the function
		SetFunction("getMessage", nil).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing contract call query", err))
	}

	fmt.Printf("Call gas used: %v\n", callResult.GasUsed)
	// Get the message from the result
	// The `0` is the index to fetch a particular type from
	//
	// e.g. type of `getMessage` was `(uint32, string)`
	// then you'd need to get each field separately using:
	//      uint32 := callResult.getUint32(0);
	//      string := callResult.getString(1);
	fmt.Printf("Message: %v\n", callResult.GetString(0))

	// In this case we are passing in a string with the value "Hello from Hedera again!"
	// as the only parameter that is passed to the contract
	contractFunctionParams = hedera.NewContractFunctionParameters().
		AddString("Hello from Hedera again!")

	// Update the message
	contractExecuteID, err := hedera.NewContractExecuteTransaction().
		// Set which contract
		SetContractID(newContractID).
		// Set the gas to execute the contract call
		SetGas(100000).
		// Set the function to call and the parameters to send
		// in this case we're calling function "set_message" with a single
		// string parameter of value "Hello from Hedera again!"
		// If instead the "setMessage" method were to require "uint32, string"
		// parameters then you must do:
		//     contractFunctionParams := hedera.NewContractFunctionParameters().
		//          .addUint32(1)
		//          .addString("string 3")
		SetFunction("setMessage", contractFunctionParams).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error executing contract", err))
	}

	// Retrieve the record to make sure the execute transaction ran
	contractExecuteRecord, err := contractExecuteID.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving contract execution record", err))
	}

	// Get the contract execute result, that contains gas used
	contractExecuteResult, err := contractExecuteRecord.GetContractExecuteResult()
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving contract exe", err))
	}

	// Print gas used
	fmt.Printf("Execute gas used: %v\n", contractExecuteResult.GasUsed)

	// Call a method on a contract that exists on Hedera
	secondCallResult, err := hedera.NewContractCallQuery().
		// Set which contract
		SetContractID(newContractID).
		// Set gas to use
		SetGas(100000).
		// Set the query payment explicitly since sometimes automatic payment calculated
		// is too low
		SetQueryPayment(hedera.NewHbar(1)).
		// Set the function to call on the contract
		SetFunction("getMessage", nil).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error executing contract call query", err))
	}

	// Get gas used
	fmt.Printf("Call gas used: %v\n", secondCallResult.GasUsed)
	// Get a string from the result at index 0
	fmt.Printf("Message: %v\n", secondCallResult.GetString(0))
}

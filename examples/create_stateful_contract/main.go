package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	client.SetOperator(operatorAccountID, operatorKey)

	defer func() {
		err = client.Close()
		if err != nil {
			println(err.Error(), ": error closing client")
			return
		}
	}()

	rawSmartContract, err := ioutil.ReadFile("./stateful.json")
	if err != nil {
		println(err.Error(), ": error reading stateful.json")
		return
	}

	var smartContract contracts = contracts{}

	err = json.Unmarshal([]byte(rawSmartContract), &smartContract)
	if err != nil {
		println(err.Error(), ": error unmarshaling")
		return
	}

	smartContractByteCode := smartContract.Contracts["stateful.sol:StatefulContract"].Bin

	fmt.Println("Stateful contract example")
	fmt.Printf("Contract bytecode size: %v bytes\n", len(smartContractByteCode))

	// Upload a file containing the byte code
	byteCodeTransactionResponse, err := hedera.NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte(smartContractByteCode)).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating file")
		return
	}

	byteCodeTransactionReceipt, err := byteCodeTransactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting file create transaction receipt")
		return
	}

	byteCodeFileID := *byteCodeTransactionReceipt.FileID

	fmt.Printf("contract bytecode file: %v\n", byteCodeFileID)

	contractFunctionParams := hedera.NewContractFunctionParameters().
		AddString("hello from hedera")

	// Instantiate the contract instance
	contractTransactionID, err := hedera.NewContractCreateTransaction().
		// Failing to set this to a sufficient amount will result in "INSUFFICIENT_GAS" status
		SetGas(75000).
		// Failing to set parameters when required will result in "CONTRACT_REVERT_EXECUTED" status
		SetConstructorParameters(contractFunctionParams).
		SetBytecodeFileID(byteCodeFileID).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating contract")
		return
	}

	contractRecord, err := contractTransactionID.GetRecord(client)
	if err != nil {
		println(err.Error(), ": error retrieving contract creation record")
		return
	}

	contractCreateResult, err := contractRecord.GetContractCreateResult()
	if err != nil {
		println(err.Error(), ": error retrieving contract creation result")
		return
	}

	newContractID := *contractRecord.Receipt.ContractID

	fmt.Printf("Contract create gas used: %v\n", contractCreateResult.GasUsed)
	fmt.Printf("Contract create transaction fee: %v\n", contractRecord.TransactionFee)
	fmt.Printf("contract: %v\n", newContractID)

	// Ask for the current message (set on creation)
	callResult, err := hedera.NewContractCallQuery().
		SetContractID(newContractID).
		SetGas(75000).
		SetQueryPayment(hedera.NewHbar(1)).
		// nil -> no parameters
		SetFunction("getMessage", nil).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing contract call query")
		return
	}

	fmt.Printf("Call gas used: %v\n", callResult.GasUsed)
	fmt.Printf("Message: %v\n", callResult.GetString(0))

	contractFunctionParams = hedera.NewContractFunctionParameters().
		AddString("Hello from Hedera again!")

	// Update the message
	contractExecuteID, err := hedera.NewContractExecuteTransaction().
		SetContractID(newContractID).
		SetGas(75000).
		SetFunction("setMessage", contractFunctionParams).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing contract")
		return
	}

	contractExecuteRecord, err := contractExecuteID.GetRecord(client)
	if err != nil {
		println(err.Error(), ": error retrieving contract execution record")
		return
	}

	contractExecuteResult, err := contractExecuteRecord.GetContractExecuteResult()
	if err != nil {
		println(err.Error(), ": error retrieving contract exe")
		return
	}

	fmt.Printf("Execute gas used: %v\n", contractExecuteResult.GasUsed)

	secondCallResult, err := hedera.NewContractCallQuery().
		SetContractID(newContractID).
		SetGas(75000).
		SetQueryPayment(hedera.NewHbar(1)).
		SetFunction("getMessage", nil).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing contract call query")
		return
	}

	fmt.Printf("Call gas used: %v\n", secondCallResult.GasUsed)
	fmt.Printf("Message: %v\n", secondCallResult.GetString(0))
}

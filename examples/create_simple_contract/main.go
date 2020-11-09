package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
)

type contract struct {
	// ignore the link references since it is empty
	Object    string `json:"object"`
	OpCodes   string `json:"opcodes"`
	SourceMap string `json:"sourceMap"`
}

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			println("not error", err.Error())
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")
	var operatorKey hedera.PrivateKey

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err = hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	defer func() {
		err = client.Close()
		if err != nil {
			panic(err)
		}
	}()

	rawContract, err := ioutil.ReadFile("./hello_world.json")
	if err != nil {
		panic(err)
	}

	var contract contract = contract{}

	err = json.Unmarshal([]byte(rawContract), &contract)
	if err != nil {
		panic(err)
	}

	contractByteCode := []byte(contract.Object)

	fmt.Println("Simple contract example")
	fmt.Printf("Contract bytecode size: %v bytes\n", len(contractByteCode))

	// Upload a file containing the byte code
	byteCodeTransactionID, err := hedera.NewFileCreateTransaction().
		SetMaxTransactionFee(hedera.NewHbar(2)).
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(contractByteCode).
		Execute(client)

	if err != nil {
		panic(err)
	}

	byteCodeTransactionRecord, err := byteCodeTransactionID.GetRecord(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("contract bytecode file upload fee: %v\n", byteCodeTransactionRecord.TransactionFee)

	byteCodeFileID := *byteCodeTransactionRecord.Receipt.FileID

	fmt.Printf("contract bytecode file: %v\n", byteCodeFileID)

	// Instantiate the contract instance
	contractTransactionResponse, err := hedera.NewContractCreateTransaction().
		SetMaxTransactionFee(hedera.NewHbar(15)).
		// Failing to set this to a sufficient amount will result in "INSUFFICIENT_GAS" status
		SetGas(2000).
		SetBytecodeFileID(byteCodeFileID).
		// Setting an admin key allows you to delete the contract in the future
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)

	if err != nil {
		panic(err)
	}

	contractRecord, err := contractTransactionResponse.GetRecord(client)
	if err != nil {
		panic(err)
	}

	contractCreateResult, err := contractRecord.GetContractCreateResult()
	if err != nil {
		panic(err)
	}

	newContractID := *contractRecord.Receipt.ContractID

	fmt.Printf("Contract create gas used: %v\n", contractCreateResult.GasUsed)
	fmt.Printf("Contract create transaction fee: %v\n", contractRecord.TransactionFee)
	fmt.Printf("Contract: %v\n", newContractID)

	// Call the contract to receive the greeting
	callResult, err := hedera.NewContractCallQuery().
		SetContractID(newContractID).
		SetGas(30000).
		SetFunction("greet", nil).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Call gas used: %v\n", callResult.GasUsed)
	fmt.Printf("Message: %v\n", callResult.GetString(0))

	// delete the transaction
	deleteTransactionResponse, err := hedera.NewContractDeleteTransaction().
		SetContractID(newContractID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	deleteTransactionReceipt, err := deleteTransactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status of transaction deletion: %v\n", deleteTransactionReceipt.Status)
}

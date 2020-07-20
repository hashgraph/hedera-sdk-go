package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
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
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	client := hedera.ClientForTestnet().
		SetOperator(operatorAccountID, operatorPrivateKey)

	defer func() {
		err = client.Close()
		if err != nil {
			panic(err)
		}
	}()

	rawSmartContract, err := ioutil.ReadFile("./stateful.json")
	if err != nil {
		panic(err)
	}

	var smartContract contracts = contracts{}

	err = json.Unmarshal([]byte(rawSmartContract), &smartContract)
	if err != nil {
		panic(err)
	}

	smartContractByteCode := smartContract.Contracts["stateful.sol:StatefulContract"].Bin

	fmt.Println("Stateful contract example")
	fmt.Printf("Contract bytecode size: %v bytes\n", len(smartContractByteCode))

	// Upload a file containing the byte code
	byteCodeTransactionID, err := hedera.NewFileCreateTransaction().
		SetMaxTransactionFee(hedera.NewHbar(2)).
		AddKey(operatorPrivateKey.PublicKey()).
		SetContents([]byte(smartContractByteCode)).
		Execute(client)

	if err != nil {
		panic(err)
	}

	byteCodeTransactionReceipt, err := byteCodeTransactionID.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	byteCodeFileID := byteCodeTransactionReceipt.GetFileID()

	fmt.Printf("contract bytecode file: %v\n", byteCodeFileID)

	contractFunctionParams := hedera.NewContractFunctionParams().
		AddString("hello from hedera")

	// Instantiate the contract instance
	contractTransactionID, err := hedera.NewContractCreateTransaction().
		SetMaxTransactionFee(hedera.NewHbar(15)).
		// Failing to set this to a sufficient amount will result in "INSUFFICIENT_GAS" status
		SetGas(2000).
		// Failing to set parameters when required will result in "CONTRACT_REVERT_EXECUTED" status
		SetConstructorParams(contractFunctionParams).
		SetBytecodeFileID(byteCodeFileID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	contractRecord, err := contractTransactionID.GetRecord(client)
	if err != nil {
		panic(err)
	}

	contractCreateResult, err := contractRecord.GetContractCreateResult()
	if err != nil {
		panic(err)
	}

	newContractID := contractRecord.Receipt.GetContractID()

	fmt.Printf("Contract create gas used: %v\n", contractCreateResult.GasUsed)
	fmt.Printf("Contract create transaction fee: %v\n", contractRecord.TransactionFee)
	fmt.Printf("contract: %v\n", newContractID)

	// Ask for the current message (set on creation)
	callResult, err := hedera.NewContractCallQuery().
		SetContractID(newContractID).
		SetGas(1000).
		// nil -> no parameters
		SetFunction("getMessage", nil).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Call gas used: %v\n", callResult.GasUsed)
	fmt.Printf("Message: %v\n", callResult.GetString(0))

	contractFunctionParams = hedera.NewContractFunctionParams().
		AddString("Hello from Hedera again!")

	// Update the message
	contractExecuteID, err := hedera.NewContractExecuteTransaction().
		SetContractID(newContractID).
		SetGas(7000).
		SetFunction("setMessage", contractFunctionParams).
		SetMaxTransactionFee(hedera.HbarFrom(8, hedera.HbarUnits.Hbar)).
		Execute(client)

	if err != nil {
		panic(err)
	}

	contractExecuteRecord, err := contractExecuteID.GetRecord(client)
	if err != nil {
		panic(err)
	}

	contractExecuteResult, err := contractExecuteRecord.GetContractExecuteResult()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Execute gas used: %v\n", contractExecuteResult.GasUsed)

	secondCallResult, err := hedera.NewContractCallQuery().
		SetContractID(newContractID).
		SetGas(1000).
		SetFunction("getMessage", nil).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Call gas used: %v\n", secondCallResult.GasUsed)
	fmt.Printf("Message: %v\n", secondCallResult.GetString(0))
}

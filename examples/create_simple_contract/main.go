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
		// SetMaxTransactionFee(hedera.HbarFrom(50, hedera.HbarUnits.Hbar)).
		// SetMaxQueryPayment(hedera.HbarFrom(50, hedera.HbarUnits.Hbar))

	defer func() {
		err = client.Close()
		if err != nil {
			panic(err)
		}
	}()

	// This path assumes you are running it from the sdk root
	rawContract, err := ioutil.ReadFile("./examples/create_simple_contract/hello_world.json")
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
		AddKey(operatorPrivateKey.PublicKey()).
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

	byteCodeFileID := byteCodeTransactionRecord.Receipt.GetFileID()

	fmt.Printf("contract bytecode file: %v\n", byteCodeFileID)

	// Instantiate the contract instance
	contractTransactionID, err := hedera.NewContractCreateTransaction().
		SetMaxTransactionFee(hedera.NewHbar(15)).
		// Failing to set this to a sufficient amount will result in "INSUFFICIENT_GAS" status
		SetGas(2000).
		SetBytecodeFileID(byteCodeFileID).
		// Setting an admin key allows you to delete the contract in the future
		SetAdminKey(operatorPrivateKey.PublicKey()).
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
	deleteTransactionID, err := hedera.NewContractDeleteTransaction().
		SetContractID(newContractID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	deleteTransactionReceipt, err := deleteTransactionID.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status of transaction deletion: %v\n", deleteTransactionReceipt.Status)
}

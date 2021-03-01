package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
)

func main() {
	var client = hedera.ClientForPreviewnet()
	var err error

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != ""{
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			println(err.Error(), ": error converting string to AccountID")
			return
		}

		operatorKey, err := hedera.Ed25519PrivateKeyFromString(configOperatorKey)
		if err != nil {
			println(err.Error(), ": error converting string to PrivateKey")
			return
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	keys := make([]hedera.Ed25519PrivateKey, 3)
	pubKeys := make([]hedera.PublicKey, 3)

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	for i := range keys {
		newKey, err := hedera.GenerateEd25519PrivateKey()
		if err != nil {
			println(err.Error(), ": error generating PrivateKey}")
			return
		}

		fmt.Printf("Key %v:\n", i)
		fmt.Printf("private = %v\n", newKey)
		fmt.Printf("public = %v\n", newKey.PublicKey())

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	// A threshold key with a threshold of 2 and length of 3 requires
	// at least 2 of the 3 keys to sign anything modifying the account
	keyList := hedera.NewKeyList().
		AddAll(pubKeys)

	//fmt.Printf("threshold key %v\n", thresholdKey)

	createResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(keyList).
		SetNodeAccountID(hedera.AccountID{0,0,3}).
		SetInitialBalance(hedera.NewHbar(10)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing create account transaction")
		return
	}

	transactionReceipt, err := createResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt")
		return
	}

	transactionID := hedera.NewTransactionID(client.GetOperatorID())

	println("transactionId for scheduled transaction = ", transactionID.String())

	newAccountID := transactionReceipt.GetAccountID()

	fmt.Printf("account = %v\n", newAccountID)

	transferTx, err := hedera.NewTransferTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountID(hedera.AccountID{0,0,3}).
		AddHbarTransfer(newAccountID, hedera.HbarFrom(-1, hedera.HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorID(), hedera.HbarFrom(1, hedera.HbarUnits.Hbar)).
		Build(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer transaction")
		return
	}

	// Manually sign with 2 of the private keys provided in the threshold
	transferTx = transferTx.
		Sign(keys[0]).
		Sign(keys[1])

	scheduled := transferTx.Schedule()
	signatures1, err := scheduled.GetScheduleSignatures()

	if len(signatures1) != 2 {
		println("Scheduled transaction has incorrect number of signatures: ", len(signatures1))
		return
	}

	scheduleResponse, err := scheduled.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing schedule create")
		return
	}

	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting schedule create receipt")
		return
	}

	scheduleID := scheduleReceipt.GetScheduleID()

	info, err := hedera.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting schedule info")
		return
	}

	println("schedule info signatories = ", len(info.Signatories))

	transfer, err := info.GetTransaction()
	if err != nil {
		println(err.Error(), ": error getting transaction from schedule info")
		return
	}

	//var transfers map[hedera.AccountID]hedera.Hbar
	var key3Signature []byte

	switch tx := transfer.(type){
	case hedera.TransferTransaction:
		built, err := tx.Build(client)
		if err != nil {
			println(err.Error(), ": error building transfer transaction")
			return
		}
		key3Signature, err = keys[2].SignTransaction(&built)
		if err != nil {
			println(err.Error(), ": error signing transfer transaction")
			return
		}
	}

	//if len(transfers) != 2{
	//
	//	println("more transfers than expected ", len(transfers))
	//	return
	//}

	//if transfers[newAccountID].AsTinybar() != -hedera.NewHbar(1).AsTinybar(){
	//	println("transfer for ", newAccountID.String(), " is not whats is expected")
	//}
	//
	//if transfers[client.GetOperatorID()].AsTinybar() != hedera.NewHbar(1).AsTinybar(){
	//	println("transfer for ", client.GetOperatorID().String(), " is not whats is expected")
	//}

	println("sending schedule sign transaction")

	signTransaction := hedera.NewScheduleSignTransaction().
		SetNodeAccountID(hedera.AccountID{0,0,3}).
		SetScheduleID(scheduleID).
		AddScheduleSignature(keys[2].PublicKey(), key3Signature)

	signatures2, err := signTransaction.GetScheduleSignatures()
	if err != nil {
		println(err.Error(), ": error getting schedule sign transaction signatures")
		return
	}

	if len(signatures2) != 1 {
		println("Scheduled sign transaction has incorrect number of signatures: ", len(signatures2))
		return
	}

	resp, err := signTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing schedule sign transaction")
		return
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error executing schedule sign receipt")
		return
	}

	_, err = hedera.
		NewScheduleInfoQuery().
		Execute(client)
	if err != nil {
		println(err.Error(), ": error retrieving info query after sign transaction")
		return
	}
}

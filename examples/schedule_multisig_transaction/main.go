package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" && client.GetOperatorPublicKey().Bytes() == nil {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			println(err.Error(), ": error converting string to AccountID")
			return
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			println(err.Error(), ": error converting string to PrivateKey")
			return
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	keys := make([]hedera.PrivateKey, 3)
	pubKeys := make([]hedera.PublicKey, 3)

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	for i := range keys {
		newKey, err := hedera.GeneratePrivateKey()
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
		AddAllPublicKeys(pubKeys)

	//fmt.Printf("threshold key %v\n", thresholdKey)

	createResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(keyList).
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

	transactionID := hedera.TransactionIDGenerate(client.GetOperatorAccountID())

	println("transactionId for scheduled transaction = ", transactionID.String())

	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)

	transferTx, err := hedera.NewTransferTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
		AddHbarTransfer(newAccountID, hedera.HbarFrom(-1, hedera.HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorAccountID(), hedera.HbarFrom(1, hedera.HbarUnits.Hbar)).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer transaction")
		return
	}

	// Manually sign with 2 of the private keys provided in the threshold
	transferTx = transferTx.
		Sign(keys[0]).
		Sign(keys[1])

	scheduled, err := transferTx.Schedule()
	if err != nil {
		println(err.Error(), ": error scheduling Transfer Transaction")
		return
	}
	signatures1, err := scheduled.GetScheduledSignatures()
	if err != nil {
		println(err.Error(), ": error getting scheduled signatures")
		return
	}

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

	scheduleID := *scheduleReceipt.ScheduleID

	info, err := hedera.NewScheduleInfoQuery().
		SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting schedule info")
		return
	}

	println("schedule info signatories = ", info.Signatories.String())

	transfer, err := info.GetTransaction()
	if err != nil {
		println(err.Error(), ": error getting transaction from schedule info")
		return
	}

	var transfers map[hedera.AccountID]hedera.Hbar
	var key3Signature []byte
	switch tx := transfer.(type){
	case hedera.TransferTransaction:
		transfers = tx.GetHbarTransfers()
		//println(tx.Transaction.)
		key3Signature, err = keys[2].SignTransaction(&tx.Transaction)
		if err != nil {
			println(err.Error(), ": error signing transfer transaction")
			return
		}
	}

	if len(transfers) != 2{
		println("more transfers than expected")
		return
	}

	if transfers[newAccountID].AsTinybar() != -hedera.NewHbar(1).AsTinybar(){
		println("transfer for ", newAccountID.String(), " is not whats is expected")
	}

	if transfers[client.GetOperatorAccountID()].AsTinybar() != hedera.NewHbar(1).AsTinybar(){
		println("transfer for ", client.GetOperatorAccountID().String(), " is not whats is expected")
	}

	println("sending schedule sign transaction")

	signTransaction := hedera.NewScheduleSignTransaction().
		SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
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
		SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error retrieving info query after sign transaction")
		return
	}
}

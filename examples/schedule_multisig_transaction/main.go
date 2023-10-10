package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

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

	keys := make([]hedera.PrivateKey, 3)
	pubKeys := make([]hedera.PublicKey, 3)

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	// Loop to generate keys for the KeyList
	for i := range keys {
		newKey, err := hedera.GeneratePrivateKey()
		if err != nil {
			panic(fmt.Sprintf("%v : error generating PrivateKey}", err))
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

	// We are using all of these keys, so the scheduled transaction doesn't automatically go through
	// It works perfectly fine with just one key
	createResponse, err := hedera.NewAccountCreateTransaction().
		// The key that must sign each transfer out of the account. If receiverSigRequired is true, then
		// it must also sign any transfer into the account.
		SetKey(keyList).
		SetInitialBalance(hedera.NewHbar(10)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing create account transaction", err))
	}

	// Make sure the transaction succeeded
	transactionReceipt, err := createResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}

	// Pre-generating transaction id for the scheduled transaction so we can track it
	transactionID := hedera.TransactionIDGenerate(client.GetOperatorAccountID())

	println("transactionId for scheduled transaction = ", transactionID.String())

	// Not really necessary as its client.GetOperatorAccountID()
	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)

	// Creating a non frozen transaction for the scheduled transaction
	// In this case its TransferTransaction
	transferTx := hedera.NewTransferTransaction().
		SetTransactionID(transactionID).
		AddHbarTransfer(newAccountID, hedera.HbarFrom(-1, hedera.HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorAccountID(), hedera.HbarFrom(1, hedera.HbarUnits.Hbar))

	// Scheduling it, this gives us hedera.ScheduleCreateTransaction
	scheduled, err := transferTx.Schedule()
	if err != nil {
		panic(fmt.Sprintf("%v : error scheduling Transfer Transaction", err))
	}

	// Executing the scheduled transaction
	scheduleResponse, err := scheduled.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing schedule create", err))
	}

	// Make sure it executed successfully
	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule create receipt", err))
	}

	// Taking out the schedule ID
	scheduleID := *scheduleReceipt.ScheduleID

	// Using the schedule ID to get the schedule transaction info, which contains the whole scheduled transaction
	info, err := hedera.NewScheduleInfoQuery().
		SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule info", err))
	}

	// Taking out the TransferTransaction from earlier
	transfer, err := info.GetScheduledTransaction()
	if err != nil {
		panic(fmt.Sprintf("%v : error getting transaction from schedule info", err))
	}

	// Converting it from interface to hedera.TransferTransaction() and retrieving the amount of transfers
	// to check if we have the right one, and that it's not empty
	var transfers map[hedera.AccountID]hedera.Hbar
	switch tx := transfer.(type) {
	case *hedera.TransferTransaction:
		transfers = tx.GetHbarTransfers()
	}

	if len(transfers) != 2 {
		println("more transfers than expected")
	}

	// Checking if the Hbar values are correct
	if transfers[newAccountID].AsTinybar() != -hedera.NewHbar(1).AsTinybar() {
		println("transfer for ", newAccountID.String(), " is not whats is expected")
	}

	// Checking if the Hbar values are correct
	if transfers[client.GetOperatorAccountID()].AsTinybar() != hedera.NewHbar(1).AsTinybar() {
		println("transfer for ", client.GetOperatorAccountID().String(), " is not whats is expected")
	}

	println("sending schedule sign transaction")

	// Creating a scheduled sign transaction, we have to sign with all of the keys in the KeyList
	signTransaction, err := hedera.NewScheduleSignTransaction().
		SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
		SetScheduleID(scheduleID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing sign transaction", err))
	}

	// Signing the scheduled transaction
	signTransaction.Sign(keys[0])
	signTransaction.Sign(keys[1])
	signTransaction.Sign(keys[2])

	resp, err := signTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing schedule sign transaction", err))
	}

	// Getting the receipt to make sure the signing executed properly
	_, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing schedule sign receipt", err))
	}

	// Making sure the scheduled transaction executed properly with schedule info query
	info, err = hedera.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving schedule info after signing", err))
	}

	// Checking if the scheduled transaction was executed and signed, and retrieving the signatories
	if !info.ExecutedAt.IsZero() {
		println("Singing success, signed at: ", info.ExecutedAt.String())
		println("Signatories: ", info.Signatories.String())
	}else{
		panic("Signing failed")
	}
}

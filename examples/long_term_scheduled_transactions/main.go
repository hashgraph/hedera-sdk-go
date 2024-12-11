package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

/**
 * @summary HIP-423 https://hips.hedera.com/hip/hip-423
 * @description Long term scheduled transactions
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
		Step 1: Create key pairs
	*/
	privateKey1, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating private key", err))
	}
	publicKey1 := privateKey1.PublicKey()

	privateKey2, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating private key", err))
	}

	fmt.Println("Creating a Key List...")
	keyList := hiero.NewKeyList().
		AddAllPublicKeys([]hiero.PublicKey{publicKey1, privateKey2.PublicKey()}).
		SetThreshold(2)
	fmt.Println("Created a Key List: ", keyList)

	/*
		Step 2: Create the account
	*/
	fmt.Println("Creating new account...")
	createResponse, err := hiero.NewAccountCreateTransaction().
		SetKey(keyList).
		SetInitialBalance(hiero.NewHbar(2)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}
	transactionReceipt, err := createResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}
	alice := *transactionReceipt.AccountID
	fmt.Println("Created new account with ID: ", alice)

	/*
		Step 3:
		Schedule a transfer transaction of 1 Hbar from the created account to the
		operator account with an expirationTime of
		24 hours in the future and waitForExpiry=false

	*/
	fmt.Println("Creating new scheduled transaction with 1 day expiry")
	transfer := hiero.NewTransferTransaction().
		AddHbarTransfer(alice, hiero.NewHbar(1).Negated()).
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(1)).
		SetMaxTransactionFee(hiero.NewHbar(10))

	schedule, err := transfer.Schedule()
	if err != nil {
		panic(fmt.Sprintf("%v : error scheduling transaction", err))
	}
	scheduleResponse, err := schedule.
		SetWaitForExpiry(false).
		SetExpirationTime(time.Now().Add(24 * time.Hour)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error scheduling transaction", err))
	}
	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule receipt", err))
	}
	scheduleID := *scheduleReceipt.ScheduleID

	/*
		Step 4:
		Sign the transaction with one key and verify the transaction is not executed
	*/
	fmt.Println("Signing the new scheduled transaction with 1 key")
	frozenSign, err := hiero.NewScheduleSignTransaction().
		SetScheduleID(scheduleID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}

	signResponse, err := frozenSign.
		Sign(privateKey1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}
	_, err = signResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}

	info, err := hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule info", err))
	}
	fmt.Println("Scheduled transaction is not yet executed. Executed at: ", info.ExecutedAt)

	/*
		Step 5:
		Sign the transaction with the other key and verify the transaction executes successfully
	*/

	accountBalance, err := hiero.NewAccountBalanceQuery().
		SetAccountID(alice).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account balance", err))
	}
	fmt.Println("Alice's account balance before scheduled transfer", accountBalance.Hbars)

	fmt.Println("Signing the new scheduled transaction with the 2nd key")
	frozenSign, err = hiero.NewScheduleSignTransaction().
		SetScheduleID(scheduleID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}

	signResponse, err = frozenSign.
		Sign(privateKey2).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}
	_, err = signResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}

	info, err = hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule info", err))
	}
	accountBalance, err = hiero.NewAccountBalanceQuery().
		SetAccountID(alice).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account balance", err))
	}
	fmt.Println("Alice's account balance after scheduled transfer", accountBalance.Hbars)
	fmt.Println("Scheduled transaction is executed. Executed at: ", info.ExecutedAt)

	/*
		Step 6:
		Schedule another transfer transaction of 1 Hbar from the account to the operator account
		with an expirationTime of 10 seconds in the future and waitForExpiry=true .
	*/
	fmt.Println("Creating new scheduled transaction with 10 seconds expiry")
	transfer = hiero.NewTransferTransaction().
		AddHbarTransfer(alice, hiero.NewHbar(1).Negated()).
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(1)).
		SetMaxTransactionFee(hiero.NewHbar(10))

	schedule, err = transfer.Schedule()
	if err != nil {
		panic(fmt.Sprintf("%v : error scheduling transaction", err))
	}
	scheduleResponse, err = schedule.
		SetWaitForExpiry(true).
		SetExpirationTime(time.Now().Add(10 * time.Second)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error scheduling transaction", err))
	}
	scheduleReceipt, err = scheduleResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule receipt", err))
	}
	scheduleID2 := *scheduleReceipt.ScheduleID
	txId := *scheduleReceipt.ScheduledTransactionID

	/*
		Step 7:
		Sign the transaction with one key and verify the transaction is not executed
	*/
	fmt.Println("Signing the new scheduled transaction with 1 key")
	frozenSign, err = hiero.NewScheduleSignTransaction().
		SetScheduleID(scheduleID2).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}

	signResponse, err = frozenSign.
		Sign(privateKey1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}
	_, err = signResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing schedule transaction", err))
	}

	info, err = hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID2).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule info", err))
	}
	fmt.Println("Scheduled transaction is not yet executed. Executed at: ", info.ExecutedAt)

	/*
		Step 8:
		Update the account’s key to be only the one key
		that has already signed the scheduled transfer.
	*/
	fmt.Println("Updating Alice's key to be the 1st key")
	frozenAccountUpdate, err := hiero.NewAccountUpdateTransaction().
		SetAccountID(alice).
		SetKey(publicKey1).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating account key", err))
	}

	accountUpdateResp, err := frozenAccountUpdate.
		Sign(privateKey1).
		Sign(privateKey2).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating account key", err))
	}
	_, err = accountUpdateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating account key", err))
	}

	/*
		Step 9:
		Verify that the transfer successfully executes roughly at the time of its expiration.
	*/
	accountBalance, err = hiero.NewAccountBalanceQuery().
		SetAccountID(alice).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account balance", err))
	}
	fmt.Println("Alice's account balance before scheduled transfer", accountBalance.Hbars)

	startTime := time.Now()
	for time.Since(startTime) < 10*time.Second {
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Elapsed time: %.1f seconds\r", time.Since(startTime).Seconds())
	}

	accountBalance, err = hiero.NewAccountBalanceQuery().
		SetAccountID(alice).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account balance", err))
	}
	fmt.Println("Alice's account balance after scheduled transfer", accountBalance.Hbars)

	record, err := hiero.NewTransactionRecordQuery().
		SetTransactionID(txId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting transaction record", err))
	}
	fmt.Println("Transaction status: ", record.Receipt.Status)

	/*
	 * Clean up:
	 */
	client.Close()

	fmt.Println("Example Complete!")
}

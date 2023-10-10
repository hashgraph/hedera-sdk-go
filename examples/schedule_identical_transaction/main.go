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

	pubKeys := make([]hedera.PublicKey, 3)
	clients := make([]*hedera.Client, 3)
	accounts := make([]hedera.AccountID, 3)

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	var scheduleID *hedera.ScheduleID

	// Loop to generate keys, clients, and accounts
	for i := range pubKeys {
		newKey, err := hedera.GeneratePrivateKey()
		if err != nil {
			panic(fmt.Sprintf("%v : error generating PrivateKey", err))
		}

		fmt.Printf("Key %v:\n", i)
		fmt.Printf("private = %v\n", newKey)
		fmt.Printf("public = %v\n", newKey.PublicKey())

		pubKeys[i] = newKey.PublicKey()

		createResponse, err := hedera.NewAccountCreateTransaction().
			// The key that must sign each transfer out of the account. If receiverSigRequired is true, then
			// it must also sign any transfer into the account.
			SetKey(newKey).
			SetInitialBalance(hedera.NewHbar(1)).
			Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error creating account", err))
		}

		// Make sure the transaction succeeded
		transactionReceipt, err := createResponse.GetReceipt(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error getting receipt 1", err))
		}

		newClient, err := hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
		if err != nil {
			panic(fmt.Sprintf("%v : error creating client", err))
		}
		newClient = newClient.SetOperator(*transactionReceipt.AccountID, newKey)

		clients[i] = newClient
		accounts[i] = *transactionReceipt.AccountID

		fmt.Printf("account = %v\n", accounts[i])
	}

	// A threshold key with a threshold of 2 and length of 3 requires
	// at least 2 of the 3 keys to sign anything modifying the account
	keyList := hedera.KeyListWithThreshold(2).
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
		panic(fmt.Sprintf("%v : error getting receipt 2", err))
	}
	thresholdAccount := *transactionReceipt.AccountID

	fmt.Printf("threshold account = %v\n", thresholdAccount)

	for _, client := range clients {
		operator := client.GetOperatorAccountID().String()

		// Each client creates an identical transaction, sending 1 hbar to each of the created accounts,
		// sent from the threshold Account
		tx := hedera.NewTransferTransaction()
		for _, account := range accounts {
			tx.AddHbarTransfer(account, hedera.NewHbar(1))
		}
		tx.AddHbarTransfer(thresholdAccount, hedera.NewHbar(3).Negated())

		tx, err := tx.FreezeWith(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error while freezing transaction for client ", err))
		}

		signedTransaction, err := tx.SignWithOperator(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error while signing with operator client ", operator))
		}

		scheduledTx, err := hedera.NewScheduleCreateTransaction().
			SetScheduledTransaction(signedTransaction)
		if err != nil {
			panic(fmt.Sprintf("%v : error while setting scheduled transaction with operator client", operator))
		}

		scheduledTx = scheduledTx.
			SetPayerAccountID(thresholdAccount)

		response, err := scheduledTx.Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error while executing schedule create transaction with operator", operator))
		}

		receipt, err := hedera.NewTransactionReceiptQuery().
			SetTransactionID(response.TransactionID).
			SetNodeAccountIDs([]hedera.AccountID{response.NodeID}).
			Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error while getting schedule create receipt transaction with operator", operator))
		}

		fmt.Printf("operator [%s]: scheduleID = %v\n", operator, receipt.ScheduleID)

		// Save the schedule ID, so that it can be asserted for each client submission
		if scheduleID == nil {
			scheduleID = receipt.ScheduleID
		}

		if scheduleID.String() != receipt.ScheduleID.String() {
			panic("invalid generated schedule id, expected " + scheduleID.String() + ", got "+ receipt.ScheduleID.String())
		}

		// If the status return by the receipt is related to already created, execute a schedule sign transaction
		if receipt.Status == hedera.StatusIdenticalScheduleAlreadyCreated {
			signTransaction, err := hedera.NewScheduleSignTransaction().
				SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
				SetScheduleID(*receipt.ScheduleID).
				Execute(client)

			if err != nil {
				panic(fmt.Sprintf("%v : error while executing scheduled sign with operator", operator))
			}

			_, err = signTransaction.GetReceipt(client)
			if err != nil {
				if err.Error() != "exceptional receipt status: SCHEDULE_ALREADY_EXECUTED" {
					panic(fmt.Sprintf("%v : error while getting scheduled sign with operator ", operator))
				}
			}
		}
	}

	// Making sure the scheduled transaction executed properly with schedule info query
	info, err := hedera.NewScheduleInfoQuery().
		SetScheduleID(*scheduleID).
		SetNodeAccountIDs([]hedera.AccountID{createResponse.NodeID}).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving schedule info after signing", err))
	}

	// Checking if the scheduled transaction was executed and signed, and retrieving the signatories
	if !info.ExecutedAt.IsZero() {
		println("Signing success, signed at: ", info.ExecutedAt.String())
		println("Signatories: ", info.Signatories.String())
	}else{
		panic("Signing failed")
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

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

	// make the key arrays
	keys := make([]hiero.PrivateKey, 3)
	pubKeys := make([]hiero.PublicKey, 3)

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	// generate the keys and put them in their respective arrays
	for i := range keys {
		newKey, err := hiero.GeneratePrivateKey()
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
	thresholdPublicKeys := hiero.KeyListWithThreshold(2).
		AddAllPublicKeys(pubKeys)

	println()
	fmt.Printf("threshold keys: %v\n", thresholdPublicKeys)
	println()

	// setup account create transaction with the public threshold keys, then freeze it for singing
	transaction, err := hiero.NewAccountCreateTransaction().
		// Only thing required to create account is the key
		SetKey(thresholdPublicKeys).
		// Setting the initial balance to be 6 Hbars
		SetInitialBalance(hiero.NewHbar(6)).
		// Presetting transaction ID, this is not required
		SetTransactionID(hiero.TransactionIDGenerate(client.GetOperatorAccountID())).
		SetTransactionMemo("sdk example create_account_with_threshold_keys/main.go").
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing create account transaction", err))
	}

	// Sign with all the private keys
	for i := range keys {
		transaction = transaction.Sign(keys[i])
	}

	// Finally, execute the transaction getting the response
	transactionResponse, err := transaction.Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}

	// Get the receipt to see everything worked
	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account create receipt", err))
	}

	// Get the new account ID
	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)

	// Now we have to make sure everything worked with a transfer transaction using the new account ID
	transferTx, err := hiero.NewTransferTransaction().
		// Presetting transaction ID is not required
		SetTransactionID(hiero.TransactionIDGenerate(newAccountID)).
		// Setting node id is not required, but it guarantees the account will be available without waiting for propagation
		SetNodeAccountIDs([]hiero.AccountID{transactionResponse.NodeID}).
		// Negate the Hbar if its being taken out of the account
		AddHbarTransfer(newAccountID, hiero.HbarFrom(-5, hiero.HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.HbarFrom(5, hiero.HbarUnits.Hbar)).
		FreezeWith(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error freezing transfer transaction", err))
	}

	// Manually sign with 2 of the private keys provided in the threshold
	transactionResponse, err = transferTx.
		Sign(keys[0]).
		Sign(keys[1]).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing create account transaction", err))
	}

	// Make sure the transaction executes properly
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving transfer receipt", err))
	}

	fmt.Printf("status of transfer transaction: %v\n", transactionReceipt.Status)

	// This query is free
	// Here we check if transfer transaction actually succeeded
	balance, err := hiero.NewAccountBalanceQuery().
		// The account ID to check balance of
		SetAccountID(newAccountID).
		SetNodeAccountIDs([]hiero.AccountID{transactionResponse.NodeID}).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account balance query", err))
	}

	fmt.Printf("account balance after transfer: %v\n", balance.Hbars.String())
}

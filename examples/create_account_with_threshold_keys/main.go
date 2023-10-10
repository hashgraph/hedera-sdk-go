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

	// make the key arrays
	keys := make([]hedera.PrivateKey, 3)
	pubKeys := make([]hedera.PublicKey, 3)

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	// generate the keys and put them in their respective arrays
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
	thresholdPublicKeys := hedera.KeyListWithThreshold(2).
		AddAllPublicKeys(pubKeys)

	println()
	fmt.Printf("threshold keys: %v\n", thresholdPublicKeys)
	println()

	// setup account create transaction with the public threshold keys, then freeze it for singing
	transaction, err := hedera.NewAccountCreateTransaction().
		// Only thing required to create account is the key
		SetKey(thresholdPublicKeys).
		// Setting the initial balance to be 6 Hbars
		SetInitialBalance(hedera.NewHbar(6)).
		// Presetting transaction ID, this is not required
		SetTransactionID(hedera.TransactionIDGenerate(client.GetOperatorAccountID())).
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
	transferTx, err := hedera.NewTransferTransaction().
		// Presetting transaction ID is not required
		SetTransactionID(hedera.TransactionIDGenerate(newAccountID)).
		// Setting node id is not required, but it guarantees the account will be available without waiting for propagation
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// Negate the Hbar if its being taken out of the account
		AddHbarTransfer(newAccountID, hedera.HbarFrom(-5, hedera.HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorAccountID(), hedera.HbarFrom(5, hedera.HbarUnits.Hbar)).
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
	balance, err := hedera.NewAccountBalanceQuery().
		// The account ID to check balance of
		SetAccountID(newAccountID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account balance query", err))
	}

	fmt.Printf("account balance after transfer: %v\n", balance.Hbars.String())
}

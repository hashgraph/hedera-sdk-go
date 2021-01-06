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
			println(err.Error(), ": error setting up client from config file")
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
	thresholdPublicKeys := hedera.KeyListWithThreshold(2).
		AddAllPublicKeys(pubKeys)

	//fmt.Printf("threshold key %v\n", thresholdKey)

	transaction, err := hedera.NewAccountCreateTransaction().
		SetKey(thresholdPublicKeys).
		SetInitialBalance(hedera.NewHbar(6)).
		SetTransactionID(hedera.TransactionIDGenerate(client.GetOperatorAccountID())).
		SetTransactionMemo("sdk example create_account_with_threshold_keys/main.go").
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing create account transaction")
		return
	}

	for i := range keys {
		transaction = transaction.Sign(keys[i])
	}

	transactionResponse, err := transaction.Execute(client)

	if err != nil {
		println(err.Error(), ": error creating account")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		println(err.Error(), ": error retrieving account create receipt")
		return
	}

	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)

	transferTx, err := hedera.NewTransferTransaction().
		SetTransactionID(hedera.TransactionIDGenerate(newAccountID)).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		AddHbarSender(newAccountID, hedera.HbarFrom(5, hedera.HbarUnits.Hbar)).
		AddHbarRecipient(client.GetOperatorAccountID(), hedera.HbarFrom(5, hedera.HbarUnits.Hbar)).
		FreezeWith(client)

	if err != nil {
		println(err.Error(), ": error freezing transfer transaction")
		return
	}

	// Manually sign with 2 of the private keys provided in the threshold
	transactionResponse, err = transferTx.
		Sign(keys[0]).
		Sign(keys[1]).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error freezing create account transaction")
		return
	}

	// Must wait for the transaction to go to consensus
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving transfer receipt")
		return
	}

	fmt.Printf("status of transfer transaction: %v\n", transactionReceipt.Status)

	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(newAccountID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account balance query")
		return
	}

	fmt.Printf("account balance after transfer: %v\n", balance.Hbars.String())
}

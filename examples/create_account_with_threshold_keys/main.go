package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			println("not error", err.Error())
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")
	var operatorKey hedera.PrivateKey

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err = hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
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
			panic(err)
		}

		fmt.Printf("Key %v:\n", i)
		fmt.Printf("private = %v\n", newKey)
		fmt.Printf("public = %v\n", newKey.PublicKey())

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	// A threshold key with a threshold of 2 and length of 3 requires
	// at least 2 of the 3 keys to sign anything modifying the account
	thresholdKey := hedera.KeyListWithThreshold(2).
		AddAllPublicKeys(pubKeys)

	//fmt.Printf("threshold key %v\n", thresholdKey)

	transaction, err := hedera.NewAccountCreateTransaction().
		SetKey(thresholdKey).
		SetInitialBalance(hedera.NewHbar(6)).
		SetTransactionID(hedera.TransactionIDGenerate(client.GetOperatorID())).
		SetTransactionMemo("sdk example create_account_with_threshold_keys/main.go").
		FreezeWith(client)

	if err != nil {
		panic(err)
	}

	transactionResponse, err := transaction.Sign(operatorKey).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)

	transferTx, err := hedera.NewCryptoTransferTransaction().
		SetTransactionID(hedera.TransactionIDGenerate(newAccountID)).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		AddSender(newAccountID, hedera.HbarFrom(5, hedera.HbarUnits.Hbar)).
		AddRecipient(client.GetOperatorID(), hedera.HbarFrom(5, hedera.HbarUnits.Hbar)).
		FreezeWith(client)

	if err != nil {
		panic(err)
	}

	// Manually sign with 2 of the private keys provided in the threshold
	transactionResponse, err = transferTx.
		Sign(keys[0]).
		Sign(keys[1]).
		Execute(client)

	if err != nil {
		panic(err)
	}

	// Must wait for the transaction to go to consensus
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("status of transfer transaction: %v\n", transactionReceipt.Status)

	// Operator must be set
	client.SetOperator(client.GetOperatorID(), operatorKey)

	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(newAccountID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		Execute(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("account balance after transfer: %v\n", balance.Hbars.String())
}

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
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	transactionResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(hedera.NewHbar(2)).
		SetInitialBalance(hedera.NewHbar(1)).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	accountID := *transactionReceipt.AccountID
	if err != nil {
		panic(err)
	}

	cost, err := hedera.NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetMaxQueryPayment(hedera.NewHbar(1)).
		GetCost(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Estimated txCost to be applied is %v\n", cost)

	transaction, err := hedera.NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetMaxTransactionFee(hedera.NewHbar(1)).
		SetTransactionID(hedera.TransactionIDGenerate(accountID)).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	transactionResponse, err = transaction.
		Sign(newKey).
		Execute(client)
	if err != nil {
		panic(err)
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}
}

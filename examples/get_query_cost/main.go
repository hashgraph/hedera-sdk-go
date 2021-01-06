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

	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	transactionResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(hedera.NewHbar(2)).
		SetInitialBalance(hedera.NewHbar(1)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating account")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account creation receipt")
		return
	}

	accountID := *transactionReceipt.AccountID

	cost, err := hedera.NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetMaxQueryPayment(hedera.NewHbar(1)).
		GetCost(client)
	if err != nil {
		println(err.Error(), ": error retrieving account info query cost")
		return
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
		println(err.Error(), ": error freezing account delete transaction")
		return
	}

	transactionResponse, err = transaction.
		Sign(newKey).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting account")
		return
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account deletion receipt")
		return
	}
}

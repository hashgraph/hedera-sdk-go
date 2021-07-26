package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	client.SetOperator(operatorAccountID, operatorKey)

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

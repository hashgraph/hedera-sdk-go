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

	bobsKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating Bob's key", err))
	}

	bobsAccountCreate, err := hedera.NewAccountCreateTransaction().
		SetReceiverSignatureRequired(true).
		SetKey(bobsKey).
		SetInitialBalance(hedera.NewHbar(10)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account creation", err))
	}

	bobsAccountCreate.Sign(bobsKey)

	response, err := bobsAccountCreate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating Bob's account", err))
	}

	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}

	if transactionReceipt.AccountID == nil {
		panic(fmt.Sprintf("%v : missing Bob's AccountID", err))
	}

	bobsID := *transactionReceipt.AccountID

	println("Alice's ID:", client.GetOperatorAccountID().String())
	println("Bob's ID:", bobsID.String())

	bobsInitialBalance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(bobsID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting Bob's balance", err))
	}

	println("Bob's initial balance:", bobsInitialBalance.Hbars.String())

	transactionID := hedera.TransactionIDGenerate(bobsID)

	transferToSchedule := hedera.NewTransferTransaction().
		SetTransactionID(transactionID).
		AddHbarTransfer(client.GetOperatorAccountID(), hedera.HbarFrom(-2, hedera.HbarUnits.Hbar)).
		AddHbarTransfer(bobsID, hedera.HbarFrom(2, hedera.HbarUnits.Hbar))

	scheduleTransaction, err := transferToSchedule.Schedule()
	if err != nil {
		panic(fmt.Sprintf("%v : error setting schedule transaction", err))
	}

	frozenScheduleTransaction, err := scheduleTransaction.FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing scheduled transaction", err))
	}

	frozenScheduleTransaction.Sign(bobsKey)

	response, err = frozenScheduleTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing create scheduled transaction", err))
	}

	transactionReceipt, err = response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule create receipt", err))
	}

	if transactionReceipt.ScheduleID == nil {
		panic(fmt.Sprintf("%v : missing Bob's ScheduleID", err))
	}

	bobsBalanceAfterSchedule, err := hedera.NewAccountBalanceQuery().
		SetAccountID(bobsID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting Bob's balance", err))
	}

	println("Bob's balance after schedule:", bobsBalanceAfterSchedule.Hbars.String())

	//clean up

	deleteAccount, err := hedera.NewAccountDeleteTransaction().
		SetAccountID(bobsID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error cleaning up", err))
	}

	deleteAccount.Sign(bobsKey)

	response, err = deleteAccount.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error cleaning up", err))
	}

	_, err = response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error cleaning up", err))
	}
}

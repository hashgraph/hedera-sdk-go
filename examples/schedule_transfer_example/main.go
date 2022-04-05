package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	bobsKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		println(err.Error(), ": error generating Bob's key")
		return
	}

	bobsAccountCreate, err := hedera.NewAccountCreateTransaction().
		SetReceiverSignatureRequired(true).
		SetKey(bobsKey).
		SetInitialBalance(hedera.NewHbar(10)).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account creation")
		return
	}

	bobsAccountCreate.Sign(bobsKey)

	response, err := bobsAccountCreate.Execute(client)
	if err != nil {
		println(err.Error(), ": error creating Bob's account")
		return
	}

	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt")
		return
	}

	if transactionReceipt.AccountID == nil {
		println(err.Error(), ": missing Bob's AccountID")
		return
	}

	bobsID := *transactionReceipt.AccountID

	println("Alice's ID:", client.GetOperatorAccountID().String())
	println("Bob's ID:", bobsID.String())

	bobsInitialBalance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(bobsID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting Bob's balance")
		return
	}

	println("Bob's initial balance:", bobsInitialBalance.Hbars.String())

	transactionID := hedera.TransactionIDGenerate(bobsID)

	transferToSchedule := hedera.NewTransferTransaction().
		SetTransactionID(transactionID).
		AddHbarTransfer(client.GetOperatorAccountID(), hedera.HbarFrom(-2, hedera.HbarUnits.Hbar)).
		AddHbarTransfer(bobsID, hedera.HbarFrom(2, hedera.HbarUnits.Hbar))

	scheduleTransaction, err := transferToSchedule.Schedule()
	if err != nil {
		println(err.Error(), ": error setting schedule transaction")
		return
	}

	frozenScheduleTransaction, err := scheduleTransaction.FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing scheduled transaction")
		return
	}

	frozenScheduleTransaction.Sign(bobsKey)

	response, err = frozenScheduleTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing create scheduled transaction")
		return
	}

	transactionReceipt, err = response.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting schedule create receipt")
		return
	}

	if transactionReceipt.ScheduleID == nil {
		println(err.Error(), ": missing Bob's ScheduleID")
		return
	}

	bobsBalanceAfterSchedule, err := hedera.NewAccountBalanceQuery().
		SetAccountID(bobsID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting Bob's balance")
		return
	}

	println("Bob's balance after schedule:", bobsBalanceAfterSchedule.Hbars.String())

	//clean up

	deleteAccount, err := hedera.NewAccountDeleteTransaction().
		SetAccountID(bobsID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error cleaning up")
		return
	}

	deleteAccount.Sign(bobsKey)

	response, err = deleteAccount.Execute(client)
	if err != nil {
		println(err.Error(), ": error cleaning up")
		return
	}

	_, err = response.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error cleaning up")
		return
	}
}

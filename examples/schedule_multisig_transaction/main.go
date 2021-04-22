package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
)

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			println(err.Error(), ": error converting string to AccountID")
			return
		}

		operatorKey, err := hedera.Ed25519PrivateKeyFromString(configOperatorKey)
		if err != nil {
			println(err.Error(), ": error converting string to PrivateKey")
			return
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	keys := make([]hedera.Ed25519PrivateKey, 3)
	pubKeys := make([]hedera.PublicKey, 3)

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	for i := range keys {
		newKey, err := hedera.GenerateEd25519PrivateKey()
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
	keyList := hedera.NewKeyList().
		AddAllPublicKeys(pubKeys)

	//fmt.Printf("threshold key %v\n", thresholdKey)

	createResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(keyList).
		SetInitialBalance(hedera.NewHbar(10)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing create account transaction")
		return
	}

	transactionReceipt, err := createResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt")
		return
	}

	transactionID := hedera.NewTransactionID(client.GetOperatorID())

	println("transactionId for scheduled transaction = ", transactionID.String())

	newAccountID := transactionReceipt.GetAccountID()

	fmt.Printf("account = %v\n", newAccountID)

	transferTx := hedera.NewTransferTransaction().
		SetTransactionID(transactionID).
		AddHbarTransfer(newAccountID, hedera.HbarFrom(-1, hedera.HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorID(), hedera.HbarFrom(1, hedera.HbarUnits.Hbar))

	scheduled, err := transferTx.Schedule()
	if err != nil {
		println(err.Error(), ": error scheduling Transfer Transaction")
		return
	}

	scheduleResponse, err := scheduled.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing schedule create")
		return
	}

	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting schedule create receipt")
		return
	}

	scheduleID := scheduleReceipt.GetScheduleID()

	info, err := hedera.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting schedule info")
		return
	}

	_, err = info.GetScheduledTransaction()
	if err != nil {
		println(err.Error(), ": error getting transaction from schedule info")
		return
	}

	signTransaction, err := hedera.NewScheduleSignTransaction().
		SetScheduleID(scheduleID).
		Build(client)
	if err != nil {
		println(err.Error(), ": error freezing sign transaction")
		return
	}

	signTransaction.Sign(keys[0])
	signTransaction.Sign(keys[1])
	signTransaction.Sign(keys[2])

	resp, err := signTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing schedule sign transaction")
		return
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error executing schedule sign receipt")
		return
	}

	info, err = hedera.
		NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error retrieving schedule info after signing")
		return
	}
	if !info.ExecutedAt.IsZero() {
		println("Singing success, signed at: ", info.ExecutedAt.String())
		println("Signatories: ")
		for _, key := range info.Signers {
			switch edKey := key.(type){
			case hedera.Ed25519PublicKey:
				println(edKey.String())
			}
		}
		return
	}
}

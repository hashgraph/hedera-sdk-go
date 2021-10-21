package main

import (
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

	key1, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}
	key2, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	accountTxResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(key1.PublicKey()).
		SetInitialBalance(hedera.ZeroHbar).
		SetTransactionID(hedera.TransactionIDGenerate(client.GetOperatorAccountID())).
		SetTransactionMemo("sdk example create_account__with_manual_signing/main.go").
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating account")
		return
	}

	println("transaction ID:", accountTxResponse.TransactionID.String())

	accountTxReceipt, err := accountTxResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account creation receipt")
		return
	}

	accountID := *accountTxReceipt.AccountID
	println("account =", accountID.String())
	println("key =", key1.PublicKey().String())
	println(":: update public key of account", accountID.String())
	println("set key =", key2.PublicKey().String())

	accountUpdateTx, err := hedera.NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetKey(key2.PublicKey()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account update transaction")
		return
	}

	accountUpdateTx.Sign(key1)
	accountUpdateTx.Sign(key2)

	accountUpdateTxResponse, err := accountUpdateTx.Execute(client)
	if err != nil {
		println(err.Error(), ": error updating account")
		return
	}

	println("transaction ID:", accountUpdateTxResponse.TransactionID.String())

	_, err = accountUpdateTxResponse.GetReceipt(client)

	println(":: getAccount and check our current key")

	info, err := hedera.NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account info query")
		return
	}

	println("key =", info.Key.String())
}

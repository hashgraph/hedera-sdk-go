package main

import (
	"github.com/hashgraph/hedera-sdk-go"
	"os"
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

	if configOperatorID != "" && configOperatorKey != "" && client.GetOperatorPublicKey().Bytes() == nil {
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

	key1, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	key2, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	accountTxResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(key1.PublicKey()).
		SetInitialBalance(hedera.ZeroHbar).
		SetTransactionID(hedera.TransactionIDGenerate(client.GetOperatorAccountID())).
		SetTransactionMemo("sdk example create_account__with_manual_signing/main.go").
		Execute(client)
	if err != nil {
		panic(err)
	}

	println("transaction ID:", accountTxResponse.TransactionID.String())

	accountTxReceipt, err := accountTxResponse.GetReceipt(client)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	accountUpdateTx.Sign(key1)
	accountUpdateTx.Sign(key2)

	accountUpdateTxResponse, err := accountUpdateTx.Execute(client)
	if err != nil {
		panic(err)
	}

	println("transaction ID:", accountUpdateTxResponse.TransactionID.String())

	_, err = accountUpdateTxResponse.GetReceipt(client)

	println(":: getAccount and check our current key")

	info, err := hedera.NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		panic(err)
	}

	println("key =", info.Key.String())
}

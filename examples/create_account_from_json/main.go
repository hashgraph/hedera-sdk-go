package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	js := []byte(
		"{\n\"network\":{\n\"0.testnet.hedera.com:50211\":\"0.0.3\"," +
			"\n\"1.testnet.hedera.com:50211\":\"0.0.4\"," +
			"\n\"2.testnet.hedera.com:50211\":\"0.0.5\"," +
			"\n\"3.testnet.hedera.com:50211\":\"0.0.6\"\n}," +
			"\n\"operator\":{\n\"accountId\":\"0.0.56313\"," +
			"\n\"privateKey\":\"302e020100300506032b657004220420c581ebedb27097be2e22b4df5a2117fdc1c1e41ac7b43ece2eff5acfa6973739\"\n}\n}\n")

	client, err = hedera.ClientFromConfig(js)
	if err != nil {
		println(err.Error(), ": error creating client from JSON")
	}

	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating private key")
		return
	}

	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	transactionResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetTransactionMemo("go sdk example create_account/main.go").
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account create transaction")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt")
		return
	}

	if transactionReceipt.AccountID != nil {
		fmt.Printf("account = %v\n", transactionReceipt.AccountID)
	}
}

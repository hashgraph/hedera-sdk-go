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

	var testnetNodes = map[string]hedera.AccountID{
		"0.testnet.hedera.com:50211": {Account: 3},
		"1.testnet.hedera.com:50211": {Account: 4},
		"2.testnet.hedera.com:50211": {Account: 5},
		"3.testnet.hedera.com:50211": {Account: 6},
	}

	fmt.Println("pinging the testnet")

	for address, id := range testnetNodes {

		// client.Ping(AccountID) returns an error, if the error is nil than the
		// ping was successful, otherwise the error will contain information to
		// potentially help diagnose the failure of the ping
		status := client.Ping(id)

		if status == nil {
			fmt.Printf("Status of node at %v with ID %v ... ok\n", address, id)
			continue
		}

		fmt.Printf("Status of node at %v with ID %v ... %s\n", address, id, status)
	}
}

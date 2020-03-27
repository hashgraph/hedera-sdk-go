package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
)

func main() {
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	// note that this is essentially the same as calling hedera.ClientForTestnet,
	// the example constructs it from scratch for the sake of showing iterating over
	// the list of nodes and pinging them
	var testnetNodes = map[string]hedera.AccountID{
		"0.testnet.hedera.com:50211": {Account: 3},
		"1.testnet.hedera.com:50211": {Account: 4},
		"2.testnet.hedera.com:50211": {Account: 5},
		"3.testnet.hedera.com:50211": {Account: 6},
	}

	client := hedera.NewClient(testnetNodes)

	client.SetOperator(operatorAccountID, operatorPrivateKey)

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

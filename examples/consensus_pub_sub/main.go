package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hashgraph/hedera-sdk-go"
)

var nodeAddress = os.Getenv("NODE_ADDRESS")
var operatorKey = os.Getenv("OPERATOR_KEY")
var mirrorNodeAddress = os.Getenv("MIRROR_NODE_ADDRESS")

func main() {
	consensusClient, err := hedera.NewConsensusClient(mirrorNodeAddress)
	if err != nil {
		panic(err)
	}

	consensusClient.SetErrorHandler(func(err error) {
		fmt.Printf("Received error: %v", err)
	})

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(operatorKey)

	client := hedera.NewClient(map[string]hedera.AccountID{
		nodeAddress: {Account: 3},
	})

	if err != nil {
		panic(err)
	}

	client.SetOperator(
		// Operator Account ID
		hedera.AccountID{Account: 2},
		// Operator Private Key
		operatorPrivateKey,
	)

	transactionId, err := hedera.NewConsensusTopicCreateTransaction().
		SetTransactionMemo("sdk example create_pub_sub/main.go").
		SetMaxTransactionFee(1000000000).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionId.Receipt(client)

	if err != nil {
		panic(err)
	}

	topicID := transactionReceipt.ConsensusTopicID()

	fmt.Printf("topicID: %v\n", topicID)

	_, err = consensusClient.Subscribe(topicID, nil, func(message hedera.ConsensusMessage) {
		fmt.Printf("%v recived topic message %v\n", message.ConsensusTimestamp, message.String())
	})

	if err != nil {
		fmt.Printf("Failed to Subscribe to topic")
		panic(err)
	}

	for i := 0; true; i++ {
		id, err := hedera.NewConsensusSubmitMessageTransaction().
			SetTopicID(topicID).
			SetMessage([]byte(fmt.Sprintf("Hello, HCS! Message %v", i))).
			Execute(client)

		if err != nil {
			panic(err)
		}

		_, err = id.Receipt(client)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Sent Message %v\n", i)

		time.Sleep(2500 * time.Millisecond)
	}
}

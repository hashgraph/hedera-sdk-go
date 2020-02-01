package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	nodeAddress := os.Getenv("NODE_ADDRESS")

	nodeID, err := hedera.AccountIDFromString(os.Getenv("NODE_ID"))
	if err != nil {
		panic(err)
	}

	client := hedera.NewClient(map[string]hedera.AccountID{
		nodeAddress: nodeID,
	})

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))

	if err != nil {
		panic(err)
	}

	client.SetOperator(
		// Operator Account ID
		operatorAccountID,
		// Operator Private Key
		operatorPrivateKey,
	)

	transactionId, err := hedera.NewConsensusTopicCreateTransaction().
		SetTransactionMemo("sdk example create_pub_sub/main.go").
		// SetMaxTransactionFee(hedera.HbarFrom(8, hedera.HbarUnits.Hbar)).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionId.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	fmt.Println(transactionReceipt.Status)

	topicID := transactionReceipt.ConsensusTopicID()

	fmt.Printf("topicID: %v\n", topicID)

	mirrorNodeAddress := os.Getenv("MIRROR_NODE_ADDRESS")

	mirrorClient, err := hedera.NewMirrorClient(mirrorNodeAddress)
	if err != nil {
		panic(err)
	}

	topicQuery, err := hedera.NewMirrorConsensusTopicQuery().
		SetTopicID(topicID).
		Subscribe(
			mirrorClient,
			func(resp hedera.MirrorConsensusTopicResponse) {
				fmt.Println(string(resp.Message))
			},
			func(err error) {
				fmt.Println(err.Error())
			})

	if err != nil {
		panic(err)
	}

	for i := 0; true; i++ {
		id, err := hedera.NewConsensusMessageSubmitTransaction().
			SetTopicID(topicID).
			SetMessage([]byte(fmt.Sprintf("Hello, HCS! Message %v", i))).
			Execute(client)

		if err != nil {
			panic(err)
		}

		_, err = id.GetReceipt(client)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Sent Message %v\n", i)

		time.Sleep(2500 * time.Millisecond)
	}

	topicQuery.Unsubscribe()
}

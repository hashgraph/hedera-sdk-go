package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	client := hedera.ClientForTestnet()

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))

	if err != nil {
		panic(err)
	}

	client.SetOperator(operatorAccountID, operatorPrivateKey)

	transactionID, err := hedera.NewTopicCreateTransaction().
		SetTransactionMemo("go sdk example create_pub_sub/main.go").
		// SetMaxTransactionFee(hedera.HbarFrom(8, hedera.HbarUnits.Hbar)).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	topicID := *transactionReceipt.TopicID

	fmt.Printf("topicID: %v\n", topicID)

	mirrorNodeAddress := os.Getenv("MIRROR_NODE_ADDRESS")

	mirrorClient := hedera.ClientForTestnet()
	mirrorClient.SetMirrorNetwork([]string{mirrorNodeAddress})
	if err != nil {
		panic(err)
	}

	_, err = hedera.NewMirrorConsensusTopicQuery().
		SetTopicID(topicID).
		Subscribe(
			mirrorClient,
			func(resp hedera.MirrorConsensusTopicResponse) {
				fmt.Printf("received message: %v\n", string(resp.Message))
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
			SetMessage([]byte(fmt.Sprintf("Hello HCS from Go! Message %v", i))).
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
}

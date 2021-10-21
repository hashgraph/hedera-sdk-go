package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

const content = `Programming is the process of creating a set of instructions that tell a computer how to perform a task. Programming can be done using a variety of computer programming languages, such as JavaScript, Python, and C++`

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

	transactionResponse, err := hedera.NewTopicCreateTransaction().
		SetTransactionMemo("go sdk example create_pub_sub/main.go").
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating topic")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		println(err.Error(), ": error getting topic create receipt")
		return
	}

	topicID := *transactionReceipt.TopicID

	fmt.Printf("topicID: %v\n", topicID)

	start := time.Now()

	_, err = hedera.NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(client, func(message hedera.TopicMessage) {
			print("Received message ", message.SequenceNumber, "\r")
		})

	if err != nil {
		println(err.Error(), ": error subscribing to the topic")
		return
	}

	for {
		_, err = hedera.NewTopicMessageSubmitTransaction().
			SetMessage([]byte(content)).
			SetTopicID(topicID).
			Execute(client)

		if err != nil {
			println(err.Error(), ": error submitting topic")
			return
		}

		if uint64(time.Since(start).Seconds()) > 60*10 {
			break
		}

		time.Sleep(2000)
	}

	println()

	transactionResponse, err = hedera.NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetMaxTransactionFee(hedera.NewHbar(5)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting topic")
		return
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for topic deletion")
		return
	}
}

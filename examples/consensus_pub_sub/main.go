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

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	// Defaults the operator account ID and key such that all generated transactions will be paid for
	// by this account and be signed by this key
	client.SetOperator(operatorAccountID, operatorKey)

	// Make a new topic
	transactionResponse, err := hedera.NewTopicCreateTransaction().
		SetTransactionMemo("go sdk example create_pub_sub/main.go").
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating topic")
		return
	}

	// Get the receipt
	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		println(err.Error(), ": error getting topic create receipt")
		return
	}

	// get the topic id from receipt
	topicID := *transactionReceipt.TopicID

	fmt.Printf("topicID: %v\n", topicID)

	start := time.Now()

	// Setup a mirror client to print out messages as we receive them
	_, err = hedera.NewTopicMessageQuery().
		// For which topic ID
		SetTopicID(topicID).
		// When to start
		SetStartTime(time.Unix(0, 0)).
		Subscribe(client, func(message hedera.TopicMessage) {
			print("Received message ", message.SequenceNumber, "\r")
		})

	if err != nil {
		println(err.Error(), ": error subscribing to the topic")
		return
	}

	// Loop submit transaction with "content" as message, wait a bit to make sure it propagates
	for {
		_, err = hedera.NewTopicMessageSubmitTransaction().
			// The message we are submitting
			SetMessage([]byte(content)).
			// To which topic ID
			SetTopicID(topicID).
			Execute(client)

		if err != nil {
			println(err.Error(), ": error submitting topic")
			return
		}

		// Setting up how long the loop wil run
		if uint64(time.Since(start).Seconds()) > 60*10 {
			break
		}

		// Sleep to make sure everything propagates
		time.Sleep(2000)
	}

	println()

	// Clean up by deleting the topic, etc
	transactionResponse, err = hedera.NewTopicDeleteTransaction().
		// Which topic ID
		SetTopicID(topicID).
		// Making sure it works right away, without propagation, by setting the same node as topic create
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// Setting the max fee just in case
		SetMaxTransactionFee(hedera.NewHbar(5)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting topic")
		return
	}

	// Get the receipt to make sure everything went through
	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for topic deletion")
		return
	}
}

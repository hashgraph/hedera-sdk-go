package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

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

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	//generate new submit key
	submitKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	println("acc", client.GetOperatorAccountID().String())

	// Create new topic ID
	transactionResponse, err := hedera.NewTopicCreateTransaction().
		// You don't need any of this to create a topic
		// If key is not set all submissions are allowed
		SetTransactionMemo("HCS topic with submit key").
		// Access control for TopicSubmitMessage.
		// If unspecified, no access control is performed, all submissions are allowed.
		SetSubmitKey(submitKey.PublicKey()).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating topic")
		return
	}

	// Get receipt
	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		println(err.Error(), ": error retrieving topic create transaction receipt")
		return
	}

	// Get topic ID from receipt
	topicID := *transactionReceipt.TopicID

	println("Created new topic", topicID.String(), "with ED25519 submitKey of", submitKey.String())

	time.Sleep(5 * time.Second)

	// Setup a mirror client to print out messages as we receive them
	_, err = hedera.NewTopicMessageQuery().
		// Sets for which topic
		SetTopicID(topicID).
		// Set when the query starts
		SetStartTime(time.Unix(0, 0)).
		// What to do when messages are received
		Subscribe(client, func(message hedera.TopicMessage) {
			// Print out the timestamp and the message
			println(message.ConsensusTimestamp.String(), " received topic message:", string(message.Contents))
		})
	if err != nil {
		println(err.Error(), ": error subscribing")
		return
	}

	for i := 0; i < 3; i++ {
		message := "random message " + strconv.Itoa(rand.Int())

		println("Publishing message:", message)

		// Prepare a message send transaction that requires a submit key from "somewhere else"
		submitTx, err := hedera.NewTopicMessageSubmitTransaction().
			// Sets the topic ID we want to send to
			SetTopicID(topicID).
			// Sets the message
			SetMessage([]byte(message)).
			FreezeWith(client)
		if err != nil {
			println(err.Error(), ": error freezing topic message submit transaction")
			return
		}

		// Sign with that submit key we gave the topic
		submitTx.Sign(submitKey)

		// Now actually submit the transaction
		submitTxResponse, err := submitTx.Execute(client)
		if err != nil {
			println(err.Error(), ": error executing topic message submit transaction")
			return
		}

		// Get the receipt to ensure there were no errors
		_, err = submitTxResponse.GetReceipt(client)
		if err != nil {
			println(err.Error(), ": error retrieving topic message submit transaction receipt")
			return
		}

		// Wait a bit for it to propagate
		time.Sleep(2 * time.Second)
	}
}

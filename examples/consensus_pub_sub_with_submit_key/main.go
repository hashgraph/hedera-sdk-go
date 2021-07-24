package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	"math/rand"
	"os"
	"strconv"
	"time"
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

	submitKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	println("acc", client.GetOperatorAccountID().String())

	transactionResponse, err := hedera.NewTopicCreateTransaction().
		SetTransactionMemo("HCS topic with submit key").
		SetSubmitKey(submitKey.PublicKey()).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating topic")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		println(err.Error(), ": error retrieving topic create transaction receipt")
		return
	}

	topicID := *transactionReceipt.TopicID

	println("Created new topic", topicID.String(), "with ED25519 submitKey of", submitKey.String())

	time.Sleep(5 * time.Second)

	_, err = hedera.NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(client, func(message hedera.TopicMessage) {
			println(message.ConsensusTimestamp.String(), " received topic message:", string(message.Contents))
		})
	if err != nil {
		println(err.Error(), ": error subscribing")
		return
	}

	for i := 0; i < 3; i++ {
		message := "random message " + strconv.Itoa(rand.Int())

		println("Publishing message:", message)

		submitTx, err := hedera.NewTopicMessageSubmitTransaction().
			SetTopicID(topicID).
			SetMessage([]byte(message)).
			FreezeWith(client)
		if err != nil {
			println(err.Error(), ": error freezing topic message submit transaction")
			return
		}

		submitTx.Sign(submitKey)
		submitTxResponse, err := submitTx.Execute(client)
		if err != nil {
			println(err.Error(), ": error executing topic message submit transaction")
			return
		}

		_, err = submitTxResponse.GetReceipt(client)
		if err != nil {
			println(err.Error(), ": error retrieving topic message submit transaction receipt")
			return
		}

		time.Sleep(2 * time.Second)
	}
}

package main

import (
	"github.com/hashgraph/hedera-sdk-go"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")
	check, err := hedera.AccountIDFromString("0.0.0")
	if err != nil {
		panic(err)
	}

	if configOperatorID != "" && configOperatorKey != "" && client.GetOperatorPublicKey().Bytes() == nil && client.GetOperatorAccountID() == check {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	submitKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	println("acc", client.GetOperatorAccountID().String())

	transactionResponse, err := hedera.NewTopicCreateTransaction().
		SetTransactionMemo("HCS topic with submit key").
		SetSubmitKey(submitKey.PublicKey()).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	topicID := *transactionReceipt.TopicID

	println("Created new topic", topicID.String(), "with ED25519 submitKey of", submitKey.String())

	time.Sleep(5000)

	_, err = hedera.NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(client, func(message hedera.TopicMessage) {
			println(message.ConsensusTimestamp.String(), " received topic message:", string(message.Contents))
		})
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		message := "random message " + strconv.Itoa(rand.Int())

		println("Publishing message:", message)

		submitTx, err := hedera.NewTopicMessageSubmitTransaction().
			SetTopicID(topicID).
			SetMessage([]byte(message)).
			FreezeWith(client)
		if err != nil {
			panic(err)
		}

		submitTx.Sign(submitKey)
		submitTxResponse, err := submitTx.Execute(client)
		if err != nil {
			panic(err)
		}

		_, err = submitTxResponse.GetReceipt(client)
		if err != nil {
			panic(err)
		}

		time.Sleep(2500)
	}
}

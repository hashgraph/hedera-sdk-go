package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	initialAdminKeys := make([]hedera.PrivateKey, 3)

	// Generating the keys for the KeyList
	for i := range initialAdminKeys {
		key, err := hedera.GeneratePrivateKey()
		if err != nil {
			panic(fmt.Sprintf("%v : error generating PrivateKey", err))
		}
		initialAdminKeys[i] = key
	}

	// Creating KeyList with a threshold 2
	keyList := hedera.KeyListWithThreshold(2)
	for _, key := range initialAdminKeys {
		keyList.Add(key.PublicKey())
	}

	topicTx, err := hedera.NewTopicCreateTransaction().
		SetTopicMemo("demo topic").
		// Access control for UpdateTopicTransaction/DeleteTopicTransaction.
		// Anyone can increase the topic's expirationTime via UpdateTopicTransaction, regardless of the adminKey.
		// If no adminKey is specified, UpdateTopicTransaction may only be used to extend the topic's expirationTime,
		// and DeleteTopicTransaction is disallowed.
		SetAdminKey(keyList).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing topic create transaction", err))
	}

	// Signing ConsensusTopicCreateTransaction with initialAdminKeys
	for i := 0; i < 2; i++ {
		println("Signing ConsensusTopicCreateTransaction with key ", initialAdminKeys[i].String())
		topicTx.Sign(initialAdminKeys[i])
	}

	// Executing ConsensusTopicCreateTransaction
	response, err := topicTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating topic", err))
	}

	// Make sure it executed properly
	receipt, err := response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic creation receipt", err))
	}

	// Get the topic ID out of the receipt
	topicID := *receipt.TopicID

	println("Created new topic ", topicID.String(), " with 2-of-3 threshold key as adminKey.")

	newAdminKeys := make([]hedera.PrivateKey, 4)

	// Generating the keys
	for i := range newAdminKeys {
		key, err := hedera.GeneratePrivateKey()
		if err != nil {
			panic(fmt.Sprintf("%v : error generating PrivateKey", err))
		}
		newAdminKeys[i] = key
	}

	// Creating KeyList with a threshold 3
	keyList = hedera.KeyListWithThreshold(3)
	for _, key := range newAdminKeys {
		keyList.Add(key.PublicKey())
	}

	topicUpdate, err := hedera.NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetTopicMemo("updated topic demo").
		// Updating with new KeyList here
		SetAdminKey(keyList).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing topic update transaction", err))
	}

	// Have to sign with the initial admin keys first
	for i := 0; i < 2; i++ {
		println("Signing ConsensusTopicCreateTransaction with initial admin key ", initialAdminKeys[i].String())
		topicUpdate.Sign(initialAdminKeys[i])
	}

	// Then the new ones we updated the topic with
	for i := 0; i < 3; i++ {
		println("Signing ConsensusTopicCreateTransaction with new admin key ", newAdminKeys[i].String())
		topicUpdate.Sign(newAdminKeys[i])
	}

	// Now to execute the topic update transaction
	response, err = topicUpdate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating topic", err))
	}

	// Make sure the transaction ran properly
	receipt, err = response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic update receipt", err))
	}

	println("Updated topic ", topicID.String(), " with 3-of-4 threshold key as adminKey")

	// Make sure everything worked by checking the topic memo
	topicInfo, err := hedera.NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing topic info query", err))
	}

	// Should be "updated topic demo"
	println(topicInfo.TopicMemo)
}

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

	accountID, _ := hedera.AccountIDFromString("0.0.1999")
	description := "Hedera™ cryptocurrency"
	newDescription := "Hedera™ cryptocurrency - updated"

	ipv4 := hedera.IPv4Address{}
	ipv4.SetNetwork(0, 0).SetHost(0, 0)

	gossipEndpoint := hedera.Endpoint{}
	gossipEndpoint.SetAddress(ipv4).SetPort(50211)

	serviceEndpoint := hedera.Endpoint{}
	serviceEndpoint.SetAddress(ipv4).SetPort(50211)

	adminKey, _ := hedera.PrivateKeyGenerateEd25519()

	nodeCreateTransaction := hedera.NewNodeCreateTransaction().
		SetAccountID(accountID).
		SetDescription(description).
		SetGossipCaCertificate([]byte("gossipCaCertificate")).
		SetServiceEndpoints([]hedera.Endpoint{serviceEndpoint}).
		SetGossipEndpoints([]hedera.Endpoint{gossipEndpoint}).
		SetAdminKey(adminKey.PublicKey())

	resp, err := nodeCreateTransaction.Execute(client)
	fmt.Println(err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	fmt.Println(err)

	nodeUpdateTransaction := hedera.NewNodeUpdateTransaction().
		SetNodeID(123).
		SetDescription(newDescription).
		SetGossipCaCertificate([]byte("gossipCaCertificate")).
		SetServiceEndpoints([]hedera.Endpoint{serviceEndpoint}).
		SetGossipEndpoints([]hedera.Endpoint{gossipEndpoint}).
		SetAdminKey(adminKey.PublicKey())
	resp, err = nodeUpdateTransaction.Execute(client)
	fmt.Println(err)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	fmt.Println(err)

	nodeDeleteTransaction := hedera.NewNodeDeleteTransaction().
		SetNodeID(123)
	resp, err = nodeDeleteTransaction.Execute(client)
	fmt.Println(err)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	fmt.Println(err)
}

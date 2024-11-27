package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

func main() {
	var client *hiero.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	accountID, _ := hiero.AccountIDFromString("0.0.1999")
	description := "Hedera™ cryptocurrency"
	newDescription := "Hedera™ cryptocurrency - updated"

	ipv4 := []byte{127, 0, 0, 1}

	gossipEndpoint := hiero.Endpoint{}
	gossipEndpoint.SetAddress(ipv4).SetPort(50211)

	serviceEndpoint := hiero.Endpoint{}
	serviceEndpoint.SetAddress(ipv4).SetPort(50211)

	adminKey, _ := hiero.PrivateKeyGenerateEd25519()

	nodeCreateTransaction := hiero.NewNodeCreateTransaction().
		SetAccountID(accountID).
		SetDescription(description).
		SetGossipCaCertificate([]byte("gossipCaCertificate")).
		SetServiceEndpoints([]hiero.Endpoint{serviceEndpoint}).
		SetGossipEndpoints([]hiero.Endpoint{gossipEndpoint}).
		SetAdminKey(adminKey.PublicKey())

	resp, err := nodeCreateTransaction.Execute(client)
	fmt.Println(err)
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	fmt.Println(err)

	nodeUpdateTransaction := hiero.NewNodeUpdateTransaction().
		SetNodeID(123).
		SetDescription(newDescription).
		SetGossipCaCertificate([]byte("gossipCaCertificate")).
		SetServiceEndpoints([]hiero.Endpoint{serviceEndpoint}).
		SetGossipEndpoints([]hiero.Endpoint{gossipEndpoint}).
		SetAdminKey(adminKey.PublicKey())
	resp, err = nodeUpdateTransaction.Execute(client)
	fmt.Println(err)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	fmt.Println(err)

	nodeDeleteTransaction := hiero.NewNodeDeleteTransaction().
		SetNodeID(123)
	resp, err = nodeDeleteTransaction.Execute(client)
	fmt.Println(err)

	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	fmt.Println(err)
}

package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

func main() {
	var client *hiero.Client
	var err error

	// Initialize the client with the testnet mirror node. This will also get the address book from the mirror node and
	// use it to populate the Client's consensus network.
	client, err = hiero.ClientForMirrorNetwork([]string{"testnet.mirrornode.hedera.com:443"})
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

	privateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey()

	txResponse, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.NewHbar(1)).
		SetKey(publicKey).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create transaction", err))
	}

	receipt, err := txResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}

	fmt.Printf("New account id, %s", receipt.AccountID.String())
}

package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// The client comes with default logger.
	// We can also set a custom logger for this client.
	// The logger must implement the hedera.Logger interface.
	logger := hedera.NewLogger("Hedera sdk", hedera.LoggerLevelDebug)
	client.SetLogger(logger)

	// Set the logging level fot this client, to be used as default.
	// Individual log levels can be set for the Query or Transaction object by
	// chaining the SetLogLevel() function on the given Transaction or Query object.
	client.SetLogLevel(hedera.LoggerLevelTrace)

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

	// Generate new key to use with new account
	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey}", err))
	}

	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	// Transaction used to show default logging functionality from client
	_, err = hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create transaction}", err))
	}

	// Disable default logging on client, to show logging functionality from transaction
	client.SetLogLevel(hedera.LoggerLevelDisabled)

	// Create account transaction used to show logging functionality
	_, err = hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		// Set logging level for the specific transaction
		SetLogLevel(hedera.LoggerLevelTrace).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create transaction}", err))
	}
}

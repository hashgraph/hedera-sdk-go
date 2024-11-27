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

	// The client comes with default logger.
	// We can also set a custom logger for this client.
	// The logger must implement the hiero.Logger interface.
	logger := hiero.NewLogger("Hiero sdk", hiero.LoggerLevelDebug)
	client.SetLogger(logger)

	// Set the logging level fot this client, to be used as default.
	// Individual log levels can be set for the Query or Transaction object by
	// chaining the SetLogLevel() function on the given Transaction or Query object.
	client.SetLogLevel(hiero.LoggerLevelTrace)

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

	// Generate new key to use with new account
	newKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey}", err))
	}

	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	// Transaction used to show default logging functionality from client
	_, err = hiero.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create transaction}", err))
	}

	// Disable default logging on client, to show logging functionality from transaction
	client.SetLogLevel(hiero.LoggerLevelDisabled)

	// Create account transaction used to show logging functionality
	_, err = hiero.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		// Set logging level for the specific transaction
		SetLogLevel(hiero.LoggerLevelTrace).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create transaction}", err))
	}
}

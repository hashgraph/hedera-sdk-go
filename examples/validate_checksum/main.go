package main

import (
	"bufio"
	"os"
	"strings"

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

	println("An example of manual checksum validation.")
	reader := bufio.NewReader(os.Stdin)
	print("Enter an account ID with checksum: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")

	// Making a new account ID from a string that we read in
	id, err := hedera.AccountIDFromString(text)
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Getting the checksum of the account ID we just created
	println("The ID with no checksum is", id.String())
	idWithChecksum, err := id.ToStringWithChecksum(client)
	if err != nil {
		println(err.Error(), ": error generating ID checksum")
		return
	}

	println("The ID with correct checksum is", idWithChecksum)
	if id.GetChecksum() == nil {
		println("You must enter a checksum.")
		return
	}

	println("The checksum entered was", *id.GetChecksum())

	// Validating the checksum, will error if not valid
	err = id.ValidateChecksum(client)
	if err != nil {
		println(err.Error(), ": error validating checksum.")
		return
	}

	// Executing with the created account ID with autoValidateChecksum being false
	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(id).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account balance query.")
		return
	}
	println(balance.Hbars.String())

	println("An example of automatic checksum validation: ")

	// Setting autoValidateChecksum to true, any operation with an ID will fail if the checksum is wrong
	client.SetAutoValidateChecksums(true)

	print("Enter an account ID with checksum: ")
	text, _ = reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")

	// Making a new account ID from a string that we read in
	id, err = hedera.AccountIDFromString(text)
	if err != nil {
		println(err.Error(), ": error converting string to AccountID.")
		return
	}

	// Checking if checksum exists
	if id.GetChecksum() == nil {
		println("You must enter a checksum.")
		return
	}

	// Executing with the created account ID with autoValidateChecksum being true
	balance, err = hedera.NewAccountBalanceQuery().
		SetAccountID(id).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account balance query.")
		return
	}

	// Print out the account balance
	println(balance.Hbars.String())

	println("Example complete!")
	_ = client.Close()
}

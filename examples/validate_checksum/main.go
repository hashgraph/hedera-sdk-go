package main

import (
	"bufio"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
	"strings"
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

	println("An example of manual checksum validation.")
	reader := bufio.NewReader(os.Stdin)
	print("Enter an account ID with checksum: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")

	id, err := hedera.AccountIDFromString(text)
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

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

	err = id.ValidateChecksum(client)
	if err != nil {
		println(err.Error(), ": error validating checksum.")
		return
	}

	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(id).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account balance query.")
		return
	}
	println(balance.Hbars.String())

	println("An example of automatic checksum validation: ")

	client.SetAutoValidateChecksums(true)

	print("Enter an account ID with checksum: ")
	text, _ = reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")

	id, err = hedera.AccountIDFromString(text)
	if err != nil {
		println(err.Error(), ": error converting string to AccountID.")
		return
	}

	if id.GetChecksum() == nil {
		println("You must enter a checksum.")
		return
	}
	balance, err = hedera.NewAccountBalanceQuery().
		SetAccountID(id).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account balance query.")
		return
	}
	println(balance.Hbars.String())

	println("Example complete!")
	_ = client.Close()
}

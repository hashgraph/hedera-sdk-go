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


	// ## Example
	// Create a ECDSA private key
	// Extract the ECDSA public key public key
	// Extract the Ethereum public address
	// Use the `AccountCreateTransaction` and populate `setAlias(evmAddress)` field with the Ethereum public address
	// Sign the `AccountCreateTransaction` transaction with the new private key
	// Get the `AccountInfo` on the new account and show that the account has contractAccountId

	// Create a ECDSA private key
	privateKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		println(err.Error())
	}
	// Extract the ECDSA public key public key
	publicKey := privateKey.PublicKey()
	// Extract the Ethereum public address
	evmAddress := publicKey.ToEvmAddress()

	// Use the `AccountCreateTransaction` and set the EVM address field to the Ethereum public address
	response, err := hedera.NewAccountCreateTransaction().SetInitialBalance(hedera.HbarFromTinybar(100)).
		SetKey(operatorKey).SetAlias(evmAddress).Sign(privateKey).Execute(client)
	if err != nil {
		println(err.Error())
	}

	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt}", err))
	}

	newAccountId := *transactionReceipt.AccountID

	// Get the `AccountInfo` on the new account and show that the account has contractAccountId
	info, err := hedera.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		println(err.Error())
	}
	// Verify account is created with the provided EVM address
	fmt.Println(info.ContractAccountID == evmAddress)
	// Verify the account Id is the same from the create account transaction
	fmt.Println(info.AccountID.String() == newAccountId.String())
}

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
	// Transfer tokens using the `TransferTransaction` to the Etherum Account Address
	// The From field should be a complete account that has a public address
	// The To field should be to a public address (to create a new account)
	// Get the child receipt or child record to return the Hedera Account ID for the new account that was created
	// Get the `AccountInfo` on the new account and show it is a hollow account by not having a public key
	// This is a hollow account in this state
	// Use the hollow account as a transaction fee payer in a HAPI transaction
	// Sign the transaction with ECDSA private key
	// Get the `AccountInfo` of the account and show the account is now a complete account by returning the public key on the account

	// Create a ECDSA private key
	privateKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(err)
	}
	// Extract the ECDSA public key public key
	publicKey := privateKey.PublicKey()
	// Extract the Ethereum public address
	evmAddress := publicKey.ToEvmAddress()

	// Create an AccountID struct with EVM address
	evmAddressAccount, err := hedera.AccountIDFromEvmPublicAddress(evmAddress)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account from EVM address", err))
	}
	// Transfer tokens using the `TransferTransaction` to the Etherum Account Address
	tx, err := hedera.NewTransferTransaction().AddHbarTransfer(evmAddressAccount, hedera.NewHbar(4)).
		AddHbarTransfer(operatorAccountID, hedera.NewHbar(-4)).Execute(client)
	if err != nil {
		panic(err)
	}

	// Get the child receipt or child record to return the Hedera Account ID for the new account that was created
	receipt, err := tx.GetReceiptQuery().SetIncludeChildren(true).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error with receipt: ", err))
	}
	newAccountId := *receipt.Children[0].AccountID

	// Get the `AccountInfo` on the new account and show it is a hollow account by not having a public key
	info, err := hedera.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		panic(err)
	}
	// Verify account is created with the public address provided
	fmt.Println(info.ContractAccountID == publicKey.ToEvmAddress())
	// Verify the account Id is the same from the create account transaction
	fmt.Println(info.AccountID.String() == newAccountId.String())
	// Verify the account does not have a Hedera public key /hollow account/
	fmt.Println(info.Key.String() == "{[]}")

	// Use the hollow account as a transaction fee payer in a HAPI transaction
	// Sign the transaction with ECDSA private key
	client.SetOperator(newAccountId, privateKey)
	tx, err = hedera.NewTransferTransaction().AddHbarTransfer(operatorAccountID, hedera.NewHbar(1)).
		AddHbarTransfer(newAccountId, hedera.NewHbar(-1)).Execute(client)
	if err != nil {
		panic(err)
	}
	receipt, err = tx.GetReceiptQuery().Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error with receipt: ", err))
	}

	// Get the `AccountInfo` of the account and show the account is now a complete account by returning the public key on the account
	info, err = hedera.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		panic(err)
	}
	// Verify account is created with the public address provided
	fmt.Println(info.ContractAccountID == publicKey.ToEvmAddress())
	// Verify the account Id is the same from the create account transaction
	fmt.Println(info.AccountID.String() == newAccountId.String())
	// Verify the account does have a Hedera public key /complete Hedera account/
	fmt.Println(info.Key.String() == publicKey.String())

}

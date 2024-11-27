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

	// Generate new key to use with new account
	newKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	transactionResponse, err := hiero.NewAccountCreateTransaction().
		// The key that must sign each transfer out of the account.
		SetKey(newKey.PublicKey()).
		// If true, this account's key must sign any transaction depositing into this account (in
		// addition to all withdrawals)
		SetReceiverSignatureRequired(false).
		// The maximum number of tokens that an Account can be implicitly associated with. Defaults to 0
		// and up to a maximum value of 1000.
		SetMaxAutomaticTokenAssociations(1).
		// The memo associated with the account
		SetTransactionMemo("go sdk example create_account/main.go").
		// The account is charged to extend its expiration date every this many seconds. If it doesn't
		// have enough balance, it extends as long as possible. If it is empty when it expires, then it
		// is deleted.
		SetStakedAccountID(hiero.AccountID{Account: 3}).
		SetInitialBalance(hiero.NewHbar(20)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create transaction", err))
	}

	// Get receipt to see if transaction succeeded, and has the account ID
	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}

	accountID := *transactionReceipt.AccountID

	println("AccountID of staking account:", accountID.String())

	info, err := hiero.NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account info", err))
	}

	println("Staked Node Account ID:", info.StakingInfo.StakedAccountID.String())
}

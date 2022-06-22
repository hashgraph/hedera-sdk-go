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

	// Generate new key to use with new account
	newKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	// Create account
	// The only required property here is `key`
	transactionResponse, err := hedera.NewAccountCreateTransaction().
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
		SetStakedNodeID(3).
		SetInitialBalance(hedera.NewHbar(20)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account create transaction")
		return
	}

	// Get receipt to see if transaction succeeded, and has the account ID
	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt")
		return
	}

	accountID := *transactionReceipt.AccountID

	info, err := hedera.NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error retrieving account info")
		return
	}

	println("Staked node id:", *info.StakingInfo.StakedNodeID)

	freezeUpdate, err := hedera.NewAccountUpdateTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetAccountID(accountID).
		// Should use ClearStakedNodeID(), but it doesn't work for now.
		SetStakedNodeID(0).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account update transaction")
		return
	}

	freezeUpdate.Sign(newKey)

	transactionResponse, err = freezeUpdate.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account update transaction")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt}")
		return
	}

	info, err = hedera.NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error retrieving account info")
		return
	}

	println("Staked node id after update:", *info.StakingInfo.StakedNodeID)
}

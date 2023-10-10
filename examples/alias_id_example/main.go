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

	// Defaults the operator account ID and key such that all generated transactions will be paid for
	// by this account and be signed by this key
	client.SetOperator(operatorAccountID, operatorKey)

	/*
	 * Hedera supports a form of auto account creation.
	 *
	 * You can "create" an account by generating a private key, and then deriving the public key,
	 * without any need to interact with the Hedera network.  The public key more or less acts as the user's
	 * account ID.  This public key is an account's aliasKey: a public key that aliases (or will eventually alias)
	 * to a Hedera account.
	 *
	 * An AccountId takes one of two forms: a normal AccountId with a null aliasKey member takes the form 0.0.123,
	 * while an account ID with a non-null aliasKey member takes the form
	 * 0.0.302a300506032b6570032100114e6abc371b82dab5c15ea149f02d34a012087b163516dd70f44acafabf7777
	 * Note the prefix of "0.0." indicating the shard and realm.  Also note that the aliasKey is stringified
	 * as a hex-encoded ASN1 DER representation of the key.
	 *
	 * An AccountId with an aliasKey can be used just like a normal AccountId for the purposes of queries and
	 * transactions, however most queries and transactions involving such an AccountId won't work until Hbar has
	 * been transferred to the aliasKey account.
	 *
	 * There is no record in the Hedera network of an account associated with a given aliasKey
	 * until an amount of Hbar is transferred to the account.  The moment that Hbar is transferred to that aliasKey
	 * AccountId is the moment that that account actually begins to exist in the Hedera ledger.
	 */

	println("Creating a new account")

	key, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating private key", err))
	}
	publicKey := key.PublicKey()

	// Assuming that the target shard and realm are known.
	// For now they are virtually always 0 and 0.
	aliasAccountID := publicKey.ToAccountID(0, 0)

	println("New account ID:", aliasAccountID.String())
	println("Just the key:", aliasAccountID.AliasKey.String())

	/*
	*
	* Note that no queries or transactions have taken place yet.
	* This account "creation" process is entirely local.
	*
	* AccountId.fromString() can construct an AccountId with an aliasKey.
	* It expects a string of the form 0.0.123 in the case of a normal AccountId, or of the form
	* 0.0.302a300506032b6570032100114e6abc371b82dab5c15ea149f02d34a012087b163516dd70f44acafabf7777
	* in the case of an AccountId with aliasKey.  Note the prefix of "0.0." to indicate the shard and realm.
	*
	* If the shard and realm are known, you may use PublicKeyFromString() then PublicKey.toAccountId() to construct the
	* aliasKey AccountID
	* fromStr, err := hedera.AccountIDFromString("0.0.302a300506032b6570032100114e6abc371b82dab5c15ea149f02d34a012087b163516dd70f44acafabf7777")
	* publicKey2, err := hedera.PublicKeyFromString("302a300506032b6570032100114e6abc371b82dab5c15ea149f02d34a012087b163516dd70f44acafabf7777")
	* fromKeyString := publicKey2.ToAccountID(0,0)
	 */

	println("Transferring some Hbar to the new account")
	resp, err := hedera.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hedera.NewHbar(1).Negated()).
		AddHbarTransfer(*aliasAccountID, hedera.NewHbar(1)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing transfer transaction", err))
	}

	receipt, err := resp.GetReceipt(client)
	println(receipt.Status.String())
	if receipt.AccountID != nil {
		println(receipt.AccountID.String())
	}
	if err != nil {
		panic(fmt.Sprintf("%v : error getting transfer transaction receipt", err))
	}

	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(*aliasAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving balance", err))
	}

	println("Balance of the new account:", balance.Hbars.String())

	info, err := hedera.NewAccountInfoQuery().
		SetAccountID(*aliasAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account info", err))
	}

	/*
	 * Note that once an account exists in the ledger, it is assigned a normal AccountId, which can be retrieved
	 * via an AccountInfoQuery.
	 *
	 * Users may continue to refer to the account by its aliasKey AccountId, but they may also
	 * now refer to it by its normal AccountId
	 */

	println("New account info:")
	println("The normal account ID:", info.AccountID.String())
	println("The alias key:", info.AliasKey.String())
	println("Example complete")
	err = client.Close()
	if err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}
}

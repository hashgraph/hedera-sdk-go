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

	key1, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}
	key2, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	fmt.Printf("privateKey = %v\n", key1.String())
	fmt.Printf("publicKey = %v\n", key1.PublicKey().String())
	fmt.Printf("privateKey = %v\n", key2.String())
	fmt.Printf("publicKey = %v\n", key2.PublicKey().String())

	// Creating 2 accounts for transferring tokens
	transactionResponse, err := hedera.NewAccountCreateTransaction().
		// The key that must sign each transfer out of the account. If receiverSigRequired is true, then
		// it must also sign any transfer into the account.
		SetKey(key1.PublicKey()).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating account")
		return
	}

	// First receipt with account ID 1, will error if transaction failed
	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account creation receipt")
		return
	}

	// Retrieving account ID out of the first receipt
	accountID1 := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", accountID1.String())

	// Creating a new account for the token
	transactionResponse, err = hedera.NewAccountCreateTransaction().
		// The key that must sign each transfer out of the account. If receiverSigRequired is true, then
		// it must also sign any transfer into the account.
		SetKey(key2.PublicKey()).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating second account")
		return
	}

	// Second receipt with account ID 2, will error if transaction failed
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving second account creation receipt")
		return
	}

	// Retrieving account ID out of the second receipt
	accountID2 := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", accountID2.String())

	// Creating a new token
	transactionResponse, err = hedera.NewTokenCreateTransaction().
		// The publicly visible name of the token
		SetTokenName("ffff").
		// The publicly visible token symbol
		SetTokenSymbol("F").
		SetMaxTransactionFee(hedera.NewHbar(1000)).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// For tokens of type FUNGIBLE_COMMON - the number of decimal places a
		// token is divisible by. For tokens of type NON_FUNGIBLE_UNIQUE - value
		// must be 0
		SetDecimals(3).
		// Specifies the initial supply of tokens to be put in circulation. The
		// initial supply is sent to the Treasury Account. The supply is in the
		// lowest denomination possible. In the case for NON_FUNGIBLE_UNIQUE Type
		// the value must be 0
		SetInitialSupply(1000000).
		// The account which will act as a treasury for the token. This account
		// will receive the specified initial supply or the newly minted NFTs in
		// the case for NON_FUNGIBLE_UNIQUE Type
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		// The key which can perform update/delete operations on the token. If empty, the token can be
		// perceived as immutable (not being able to be updated/deleted)
		SetAdminKey(client.GetOperatorPublicKey()).
		// The key which can sign to freeze or unfreeze an account for token transactions. If empty,
		// freezing is not possible
		SetFreezeKey(client.GetOperatorPublicKey()).
		// The key which can wipe the token balance of an account. If empty, wipe is not possible
		SetWipeKey(client.GetOperatorPublicKey()).
		// The key which can grant or revoke KYC of an account for the token's transactions. If empty,
		// KYC is not required, and KYC grant or revoke operations are not possible.
		SetKycKey(client.GetOperatorPublicKey()).
		// The key which can change the supply of a token. The key is used to sign Token Mint/Burn
		// operations
		SetSupplyKey(client.GetOperatorPublicKey()).
		// The default Freeze status (frozen or unfrozen) of Hedera accounts relative to this token. If
		// true, an account must be unfrozen before it can receive the token
		SetFreezeDefault(false).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating token")
		return
	}

	// Make sure the token create transaction ran
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token creation receipt")
		return
	}

	// Retrieve the token out of the receipt
	tokenID := *transactionReceipt.TokenID

	fmt.Printf("token = %v\n", tokenID.String())

	// Associating the token with the first account, so it can interact with the token
	transaction, err := hedera.NewTokenAssociateTransaction().
		// The account ID to be associated
		SetAccountID(accountID1).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// The token ID that the account will be associated to
		SetTokenIDs(tokenID).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token associate transaction")
		return
	}

	// Has to be signed by the account1's key
	transactionResponse, err = transaction.
		Sign(key1).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error associating token")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token associate transaction receipt")
		return
	}

	fmt.Printf("Associated account %v with token %v\n", accountID1.String(), tokenID.String())

	// Associating the token with the second account, so it can interact with the token
	transaction, err = hedera.NewTokenAssociateTransaction().
		// The account ID to be associated
		SetAccountID(accountID2).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// The token ID that the account will be associated to
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing second token associate transaction")
		return
	}

	// Has to be signed by the account2's key
	transactionResponse, err = transaction.
		Sign(key2).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing second token associate transaction")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving second token associate transaction receipt")
		return
	}

	fmt.Printf("Associated account %v with token %v\n", accountID2.String(), tokenID.String())

	// This transaction grants Kyc to the first account
	// Must be signed by the Token's kycKey.
	transactionResponse, err = hedera.NewTokenGrantKycTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// The account that KYC is being granted to
		SetAccountID(accountID1).
		// As the token kyc key is client.GetOperatorPublicKey(), we don't have to explicitly sign with anything
		// as it's done automatically by execute for the operator
		Execute(client)
	if err != nil {
		println(err.Error(), ": error granting kyc")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving grant kyc transaction receipt")
		return
	}

	fmt.Printf("Granted KYC for account %v on token %v\n", accountID1.String(), tokenID.String())

	// This transaction grants Kyc to the second account
	transactionResponse, err = hedera.NewTokenGrantKycTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// The account that KYC is being granted to
		SetAccountID(accountID2).
		// As the token kyc key is client.GetOperatorPublicKey(), we don't have to explicitly sign with anything
		// as it's done automatically by execute for the operator
		Execute(client)
	if err != nil {
		println(err.Error(), ": error granting kyc to second account")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving grant kyc transaction receipt")
		return
	}

	fmt.Printf("Granted KYC for account %v on token %v\n", accountID2.String(), tokenID.String())

	transactionResponse, err = hedera.NewTransferTransaction().
		// Same as for Hbar transfer, token value has to be negated to denote they are being taken out
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -10).
		// Same as for Hbar transfer, the 2 transfers here have to be equal, otherwise it will lead to an error
		AddTokenTransfer(tokenID, accountID1, 10).
		// We don't have to sign this one as we are transferring tokens from the operator
		Execute(client)
	if err != nil {
		println(err.Error(), ": error transferring from operator to account1")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving transfer from operator to account1 receipt")
		return
	}

	fmt.Printf(
		"Sent 10 tokens from account %v to account %v on token %v\n",
		client.GetOperatorAccountID().String(),
		accountID1.String(),
		tokenID.String(),
	)

	transferTransaction, err := hedera.NewTransferTransaction().
		// 10 tokens from account 1
		AddTokenTransfer(tokenID, accountID1, -10).
		// 10 token to account 2
		AddTokenTransfer(tokenID, accountID2, 10).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer from account1 to account2")
		return
	}

	// As we are now transferring tokens from accountID1 to accountID2, this has to be signed by accountID1's key
	transferTransaction = transferTransaction.Sign(key1)

	// Execute the transfer transaction
	transactionResponse, err = transferTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error transferring from account1 to account2")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving transfer from account1 to account2 receipt")
		return
	}

	fmt.Printf(
		"Sent 10 tokens from account %v to account %v on token %v\n",
		accountID1.String(),
		accountID2.String(),
		tokenID.String(),
	)

	transferTransaction, err = hedera.NewTransferTransaction().
		// 10 tokens from account 2
		AddTokenTransfer(tokenID, accountID2, -10).
		// 10 token to account 1
		AddTokenTransfer(tokenID, accountID1, 10).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer from account2 to account1")
		return
	}

	// As we are now transferring tokens from accountID2 back to accountID1, this has to be signed by accountID2's key
	transferTransaction = transferTransaction.Sign(key2)

	// Executing the transfer transaction
	transactionResponse, err = transferTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error transferring from account2 to account1")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving transfer from account2 to account1 receipt")
		return
	}

	fmt.Printf(
		"Sent 10 tokens from account %v to account %v on token %v\n",
		accountID2.String(),
		accountID1.String(),
		tokenID.String(),
	)

	// Clean up

	// Now we can wipe the 10 tokens that are in possession of accountID1
	// Has to be signed by wipe key of the token, in this case it was the operator key
	transactionResponse, err = hedera.NewTokenWipeTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// From which account
		SetAccountID(accountID1).
		// For which token
		SetTokenID(tokenID).
		// How many
		SetAmount(10).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error wiping from token")
		return
	}

	// Make sure the transaction succeeded
	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token wipe transaction receipt")
		return
	}

	fmt.Printf("Wiped account %v on token %v\n", accountID1.String(), tokenID.String())

	// Now to delete the token
	// Has to be signed by admin key of the token, in this case it was the operator key
	transactionResponse, err = hedera.NewTokenDeleteTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting token")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token delete transaction receipt")
		return
	}

	fmt.Printf("DeletedAt token %v\n", tokenID.String())

	// Now that the tokens have been wiped from accountID1, we can safely delete it
	accountDeleteTx, err := hedera.NewAccountDeleteTransaction().
		SetAccountID(accountID1).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// Tp which account to transfer the account 1 balance
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account delete transaction")
		return
	}

	// Account deletion has to always be signed by the key for the account
	transactionResponse, err = accountDeleteTx.
		Sign(key1).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting account 1")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving transfer transaction receipt")
		return
	}

	fmt.Printf("DeletedAt account %v\n", accountID1.String())

	accountDeleteTx, err = hedera.NewAccountDeleteTransaction().
		SetAccountID(accountID2).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// Tp which account to transfer the account 2 balance
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account delete transaction")
		return
	}

	// Account deletion has to always be signed by the key for the account
	transactionResponse, err = accountDeleteTx.
		Sign(key2).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting account2")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account delete transaction receipt")
		return
	}

	fmt.Printf("DeletedAt account %v\n", accountID2.String())
}

package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
)

func main() {
	client, err := hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	id, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	key, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(id, key)

	/**     Example 1
	 *
	 * Step 1
	 *
	 * Create accounts A, B, and C
	 */

	firstAccountPrivateKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(err)
	}
	firstPublicKey := firstAccountPrivateKey.PublicKey()
	firstAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(firstPublicKey).
		SetInitialBalance(hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar)).
		Execute(client)
	if err != nil {
		panic(err)
	}
	receiptFirstAccount, err := firstAccount.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	//Get the new account ID from the receipt
	firstAccountId := *receiptFirstAccount.AccountID
	fmt.Println("firstAccountId: ", firstAccountId)

	secondAccountPrivateKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(err)
	}
	secondAccountPublicKey := secondAccountPrivateKey.PublicKey()
	secondAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(secondAccountPublicKey).
		SetInitialBalance(hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar)).
		Execute(client)
	if err != nil {
		panic(err)
	}
	receiptSecondAccount, err := secondAccount.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	//Get the new account ID from the receipt
	secondAccountId := *receiptSecondAccount.AccountID
	fmt.Println("secondAccountId: ", secondAccountId)

	thirdAccountPrivateKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(err)
	}
	thirdAccountPublicKey := thirdAccountPrivateKey.PublicKey()
	thirdAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(thirdAccountPublicKey).
		SetInitialBalance(hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar)).
		Execute(client)
	if err != nil {
		panic(err)
	}
	receiptThirdAccount, err := thirdAccount.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	//Get the new account ID from the receipt
	thirdAccountId := *receiptThirdAccount.AccountID
	fmt.Println("thirdAccountId: ", thirdAccountId)

	/**
	 * Step 2
	 *
	 * 2. Create a fungible token that has three fractional fees
	 * Fee #1 sends 1/100 of the transferred value to collector 0.0.A.
	 * Fee #2 sends 2/100 of the transferred value to collector 0.0.B.
	 * Fee #3 sends 3/100 of the transferred value to collector 0.0.C.
	 */

	fee1 := hedera.NewCustomFractionalFee().SetFeeCollectorAccountID(firstAccountId).SetNumerator(1).SetDenominator(100).SetAllCollectorsAreExempt(true)
	fee2 := hedera.NewCustomFractionalFee().SetFeeCollectorAccountID(secondAccountId).SetNumerator(2).SetDenominator(100).SetAllCollectorsAreExempt(true)
	fee3 := hedera.NewCustomFractionalFee().SetFeeCollectorAccountID(thirdAccountId).SetNumerator(3).SetDenominator(100).SetAllCollectorsAreExempt(true)
	tokenCreateTransaction, err := hedera.NewTokenCreateTransaction().
		SetTokenName("HIP-573 Token").SetTokenSymbol("H573").
		SetTokenType(hedera.TokenTypeFungibleCommon).
		SetTreasuryAccountID(id).SetAutoRenewAccount(id).
		SetAdminKey(key.PublicKey()).SetFreezeKey(key.PublicKey()).
		SetWipeKey(key.PublicKey()).SetInitialSupply(100000000). // Total supply = 100000000 / 10 ^ 2
		SetDecimals(2).SetCustomFees([]hedera.Fee{fee1, fee2, fee3}).FreezeWith(client)
	if err != nil {
		panic(err)
	}

	transactionResponse, err := tokenCreateTransaction.Sign(key).
		Sign(firstAccountPrivateKey).
		Sign(secondAccountPrivateKey).
		Sign(thirdAccountPrivateKey).Execute(client)
	if err != nil {
		panic(err)
	}
	receipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}
	tokenId := *receipt.TokenID
	fmt.Println("Created token with token id: ", tokenId)

	/**
	 * Step 3
	 *
	 * Collector 0.0.B sends 10_000 units of the token to 0.0.A.
	 */

	const amount = 10_000
	// First we transfer the amount from treasury account to second account
	treasuryTokenTransferTransaction, err := hedera.NewTransferTransaction().
		AddTokenTransfer(tokenId, id, -amount).AddTokenTransfer(tokenId, secondAccountId, amount).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	treasuryTokenTransferSubmit, err := treasuryTokenTransferTransaction.Sign(key).Execute(client)
	if err != nil {
		panic(err)
	}
	treasuryTransferReceipt, err := treasuryTokenTransferSubmit.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending from treasury account to the second account - 'TransferTransaction' status: ", treasuryTransferReceipt.Status)

	tokenTransferTx, err := hedera.NewTransferTransaction().
		AddTokenTransfer(tokenId, secondAccountId, -amount).
		AddTokenTransfer(tokenId, firstAccountId, amount).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	submitTransaction, err := tokenTransferTx.Sign(key).Sign(secondAccountPrivateKey).Execute(client)
	if err != nil {
		panic(err)
	}

	record, err := submitTransaction.GetRecord(client)
	if err != nil {
		panic(err)
	}
	fmt.Println("Transaction fee: ", record.TransactionFee)

	/**
	 * Step 5
	 *
	 * Show that the fee collector accounts in the custom fee list
	 * of the token that was created was not charged a custom fee in the transfer
	 */

	firstAccountBalanceAfter, err := hedera.NewAccountBalanceQuery().
		SetAccountID(firstAccountId).
		Execute(client)
	if err != nil {
		panic(err)
	}
	fmt.Println("first's balance:", firstAccountBalanceAfter.Tokens.Get(tokenId))

	secondAccountBalanceAfter, err := hedera.NewAccountBalanceQuery().
		SetAccountID(secondAccountId).
		Execute(client)
	if err != nil {
		panic(err)
	}
	fmt.Println("second's balance:", secondAccountBalanceAfter.Tokens.Get(tokenId))

	thirdAccountBalanceAfter, err := hedera.NewAccountBalanceQuery().
		SetAccountID(thirdAccountId).
		Execute(client)
	if err != nil {
		panic(err)
	}
	fmt.Println("third's balance:", secondAccountBalanceAfter.Tokens.Get(tokenId))

	if firstAccountBalanceAfter.Tokens.Get(tokenId) == amount && secondAccountBalanceAfter.Tokens.Get(tokenId) == 0 && thirdAccountBalanceAfter.Tokens.Get(tokenId) == 0 {
		fmt.Println("Fee collector accounts were not charged after transfer transaction")
	}

	client.Close()

}

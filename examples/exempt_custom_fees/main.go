package main

import (
	"fmt"
	"os"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	client, err := hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	id, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	key, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
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

	firstAccountBalanceBefore, err := hedera.NewAccountBalanceQuery().
		SetAccountID(firstAccountId).
		Execute(client)
	if err != nil {
		fmt.Println(err)
	}
	println("first's balance:", firstAccountBalanceBefore.Hbars.String())

	secondAccountBalanceBefore, err := hedera.NewAccountBalanceQuery().
		SetAccountID(secondAccountId).
		Execute(client)
	if err != nil {
		fmt.Println(err)
	}
	println("second's balance:", secondAccountBalanceBefore.Hbars.String())

	thirdAccountBalanceBefore, err := hedera.NewAccountBalanceQuery().
		SetAccountID(thirdAccountId).
		Execute(client)
	if err != nil {
		fmt.Println(err)
	}
	println("third's balance:", thirdAccountBalanceBefore.Hbars.String())

	/**
	 * Step 2
	 *
	 * 2. Create a fungible token that has three fractional fees
	 * Fee #1 sends 1/100 of the transferred value to collector 0.0.A.
	 * Fee #2 sends 2/100 of the transferred value to collector 0.0.B.
	 * Fee #3 sends 3/100 of the transferred value to collector 0.0.C.
	 */

	fee1 := hedera.NewCustomFractionalFee().SetFeeCollectorAccountID(firstAccountId).SetNumerator(2).SetDenominator(100)
	fee2 := hedera.NewCustomFractionalFee().SetFeeCollectorAccountID(secondAccountId).SetNumerator(3).SetDenominator(100)
	fee3 := hedera.NewCustomFractionalFee().SetFeeCollectorAccountID(thirdAccountId).SetNumerator(1).SetDenominator(100)
	tokenCreateTransaction, err := hedera.NewTokenCreateTransaction().
		SetTokenName("HIP-573 Token").SetTokenSymbol("H573").
		SetTokenType(hedera.TokenTypeFungibleCommon).
		SetTreasuryAccountID(secondAccountId).SetAutoRenewAccount(id).
		SetAdminKey(key.PublicKey()).SetFreezeKey(key.PublicKey()).
		SetWipeKey(key.PublicKey()).SetInitialSupply(100000000). // Total supply = 100000000 / 10 ^ 2
		SetDecimals(2).SetCustomFees([]hedera.Fee{fee1, fee2, fee3}).FreezeWith(client)
	if err != nil {
		fmt.Println(err)
	}
	tokenCreateTransaction = tokenCreateTransaction.Sign(key).
		Sign(firstAccountPrivateKey).
		Sign(secondAccountPrivateKey).
		Sign(thirdAccountPrivateKey)

	transactionResponse, err := tokenCreateTransaction.Execute(client)
	if err != nil {
		fmt.Println(err)
	}
	receipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		fmt.Println(err)
	}
	tokenId := *receipt.TokenID
	fmt.Println("Created token with token id: ", tokenId)

	/**
	 * Step 3
	 *
	 * Collector 0.0.B sends 10_000 units of the token to 0.0.A.
	 */
	tokenTransferTx, err := hedera.NewTransferTransaction().
		AddTokenTransfer(tokenId, secondAccountId, -10000).
		AddTokenTransfer(tokenId, firstAccountId, 10000).
		FreezeWith(client)
	if err != nil {
		fmt.Println(err)
	}

	submitTransaction, err := tokenTransferTx.Sign(key).Sign(secondAccountPrivateKey).Execute(client)
	if err != nil {
		fmt.Println(err)
	}

	record, err := submitTransaction.GetRecord(client)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}
	println("first's balance:", firstAccountBalanceAfter.Hbars.String())

	secondAccountBalanceAfter, err := hedera.NewAccountBalanceQuery().
		SetAccountID(secondAccountId).
		Execute(client)
	if err != nil {
		fmt.Println(err)
	}
	println("second's balance:", secondAccountBalanceAfter.Hbars.String())

	thirdAccountBalanceAfter, err := hedera.NewAccountBalanceQuery().
		SetAccountID(thirdAccountId).
		Execute(client)
	if err != nil {
		fmt.Println(err)
	}
	println("third's balance:", thirdAccountBalanceAfter.Hbars.String())

	if firstAccountBalanceBefore.Hbars == firstAccountBalanceAfter.Hbars &&
		secondAccountBalanceBefore.Hbars == secondAccountBalanceAfter.Hbars &&
		thirdAccountBalanceBefore.Hbars == thirdAccountBalanceAfter.Hbars {

		fmt.Println(`Fee collector accounts were not charged after transfer transaction`)
	}

	client.Close()

}

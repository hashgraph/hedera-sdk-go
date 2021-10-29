package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
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

	aliceKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	aliceAccountCreate, err := hedera.NewAccountCreateTransaction().
		SetInitialBalance(hedera.NewHbar(10)).
		SetKey(aliceKey).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account create for alice")
		return
	}

	aliceAccountCreate.Sign(aliceKey)
	resp, err := aliceAccountCreate.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account create for alice")
		return
	}

	receipt, err := resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for alice account create")
		return
	}

	var aliceId hedera.AccountID
	if receipt.AccountID != nil {
		aliceId = *receipt.AccountID
	} else {
		println("Receipt didn't return alice's ID")
		return
	}

	bobKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	bobAccountCreate, err := hedera.NewAccountCreateTransaction().
		SetInitialBalance(hedera.NewHbar(10)).
		SetKey(bobKey).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account create for bob")
		return
	}

	bobAccountCreate.Sign(bobKey)
	resp, err = bobAccountCreate.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account create for bob")
		return
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for bob account create")
		return
	}

	var bobId hedera.AccountID
	if receipt.AccountID != nil {
		bobId = *receipt.AccountID
	} else {
		println("Receipt didn't return bob's ID")
		return
	}

	charlieKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	charlieAccountCreate, err := hedera.NewAccountCreateTransaction().
		SetInitialBalance(hedera.NewHbar(10)).
		SetKey(charlieKey).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account create for charlie")
		return
	}

	charlieAccountCreate.Sign(aliceKey)
	resp, err = charlieAccountCreate.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account create for charlie")
		return
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for charlie account create")
		return
	}

	var charlieId hedera.AccountID
	if receipt.AccountID != nil {
		charlieId = *receipt.AccountID
	} else {
		println("Receipt didn't return charlie's ID")
		return
	}

	println("Alice:", aliceId.String())
	println("Bob:", bobId.String())
	println("Charlie:", charlieId.String())

	customHbarFee := hedera.NewCustomFixedFee().
		SetHbarAmount(hedera.NewHbar(1)).
		SetFeeCollectorAccountID(aliceId)

	tokenCreate, err := hedera.NewTokenCreateTransaction().
		SetTokenName("Example Token").
		SetTokenSymbol("EX").
		SetAdminKey(aliceKey).
		SetSupplyKey(aliceKey).
		SetFeeScheduleKey(aliceKey).
		SetTreasuryAccountID(aliceId).
		SetCustomFees([]hedera.Fee{*customHbarFee}).
		SetInitialSupply(100).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token create transaction")
		return
	}

	tokenCreate.Sign(aliceKey)
	resp, err = tokenCreate.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing token create transaction")
		return
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for token create transaction")
		return
	}

	var tokenId hedera.TokenID
	if receipt.TokenID != nil {
		tokenId = *receipt.TokenID
	} else {
		println("Token ID missing in the receipt")
		return
	}

	println("TokenID:", tokenId.String())

	tokenInfo1, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenId).
		Execute(client)

	println("Custom Fees according to TokenInfoQuery:")
	for _, i := range tokenInfo1.CustomFees {
		switch t := i.(type) {
		case hedera.CustomFixedFee:
			println(t.String())
		}
	}

	tokenAssociate, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(bobId).
		SetTokenIDs(tokenId).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token associate transaction for bob")
		return
	}

	tokenAssociate.Sign(bobKey)
	resp, err = tokenAssociate.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing token associate transaction for bob")
		return
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for token associate transaction for bob")
		return
	}

	tokenAssociate, err = hedera.NewTokenAssociateTransaction().
		SetAccountID(charlieId).
		SetTokenIDs(tokenId).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token associate transaction for charlie")
		return
	}

	tokenAssociate.Sign(charlieKey)
	resp, err = tokenAssociate.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing token associate transaction for charlie")
		return
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for token associate transaction for charlie")
		return
	}

	transferTransaction, err := hedera.NewTransferTransaction().
		AddTokenTransfer(tokenId, bobId, 100).
		AddTokenTransfer(tokenId, aliceId, -100).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token transfer transaction for alice")
		return
	}

	transferTransaction.Sign(aliceKey)
	resp, err = transferTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing token transfer transaction for alice")
		return
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for token transfer transaction for alice")
		return
	}

	aliceBalance1, err := hedera.NewAccountBalanceQuery().
		SetAccountID(aliceId).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting account balance 1 for alice")
		return
	}

	println("Alice's Hbar balance before Bob transfers 20 tokens to Charlie:", aliceBalance1.Hbars.String())

	transferTransaction, err = hedera.NewTransferTransaction().
		AddTokenTransfer(tokenId, bobId, -20).
		AddTokenTransfer(tokenId, charlieId, 20).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token transfer transaction for bob")
		return
	}

	transferTransaction.Sign(bobKey)
	resp, err = transferTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing token transfer transaction for bob")
		return
	}

	record1, err := resp.GetRecord(client)
	if err != nil {
		println(err.Error(), ": error getting record for token transfer transaction for bob")
		return
	}

	aliceBalance2, err := hedera.NewAccountBalanceQuery().
		SetAccountID(aliceId).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting account balance 2 for alice")
		return
	}

	println("Alice's Hbar balance after Bob transfers 20 tokens to Charlie:", aliceBalance2.Hbars.String())
	println("Assessed fees according to transaction record:")
	for _, k := range record1.AssessedCustomFees {
		println(k.String())
	}

	customFractionalFee := hedera.NewCustomFractionalFee().
		SetNumerator(1).
		SetDenominator(10).
		SetMin(1).
		SetMax(10).
		SetFeeCollectorAccountID(aliceId)

	tokenFeeUpdate, err := hedera.NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenId).
		SetCustomFees([]hedera.Fee{*customFractionalFee}).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token fee update")
		return
	}

	tokenFeeUpdate.Sign(aliceKey)
	resp, err = tokenFeeUpdate.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing token fee update")
		return
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt for token fee update")
		return
	}

	tokenInfo2, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenId).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting token info 2")
		return
	}

	println("Custom Fees according to TokenInfoQuery:")
	for _, i := range tokenInfo2.CustomFees {
		switch t := i.(type) {
		case hedera.CustomFractionalFee:
			println(t.String())
		}
	}

	aliceBalance3, err := hedera.NewAccountBalanceQuery().
		SetAccountID(aliceId).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting account balance 3 for alice")
		return
	}

	println("Alice's token balance before Bob transfers 20 tokens to Charlie:", aliceBalance3.Tokens.Get(tokenId))

	transferTransaction, err = hedera.NewTransferTransaction().
		AddTokenTransfer(tokenId, bobId, -20).
		AddTokenTransfer(tokenId, charlieId, 20).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token transfer transaction for bob")
		return
	}

	transferTransaction.Sign(bobKey)
	resp, err = transferTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing token transfer transaction for bob")
		return
	}

	record2, err := resp.GetRecord(client)
	if err != nil {
		println(err.Error(), ": error getting record for token transfer transaction for bob")
		return
	}

	aliceBalance4, err := hedera.NewAccountBalanceQuery().
		SetAccountID(aliceId).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error getting account balance 2 for alice")
		return
	}

	println("Alice's token balance after Bob transfers 20 tokens to Charlie:", aliceBalance4.Tokens.Get(tokenId))
	println("Token transfers according to transaction record:")
	for token, transfer := range record2.TokenTransfers {
		tokenT := ""
		for _, t := range transfer {
			tokenT = tokenT + " " + t.String()
		}
		println("for token", token.String()+":", tokenT)
	}
	println("Assessed fees according to transaction record:")
	for _, k := range record2.AssessedCustomFees {
		println(k.String())
	}

	//Clean up

	tokenDelete, _ := hedera.NewTokenDeleteTransaction().
		SetTokenID(tokenId).
		FreezeWith(client)

	tokenDelete.Sign(aliceKey)
	resp, _ = tokenDelete.Execute(client)
	_, _ = resp.GetReceipt(client)

	accDelete, _ := hedera.NewAccountDeleteTransaction().
		SetAccountID(charlieId).
		FreezeWith(client)

	accDelete.Sign(charlieKey)
	resp, _ = accDelete.Execute(client)
	_, _ = resp.GetReceipt(client)

	accDelete, _ = hedera.NewAccountDeleteTransaction().
		SetAccountID(bobId).
		FreezeWith(client)

	accDelete.Sign(bobKey)
	resp, _ = accDelete.Execute(client)
	_, _ = resp.GetReceipt(client)

	accDelete, _ = hedera.NewAccountDeleteTransaction().
		SetAccountID(aliceId).
		FreezeWith(client)

	accDelete.Sign(aliceKey)
	resp, _ = accDelete.Execute(client)
	_, _ = resp.GetReceipt(client)

	_ = client.Close()
}

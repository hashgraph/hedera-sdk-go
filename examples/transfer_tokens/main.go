package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
)

func main() {
	operatorID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	operatorKey := operatorPrivateKey.PublicKey()

	client := hedera.ClientForPreviewnet()
	client.SetOperator(operatorID, operatorPrivateKey)

	key, err := hedera.GenerateEd25519PrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Printf("privateKey = %v\n", key.String())
	fmt.Printf("publicKey = %v\n", key.PublicKey().String())

	response, err := hedera.NewAccountCreateTransaction().
		SetKey(key.PublicKey()).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err := response.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	accountID := receipt.GetAccountID()

	fmt.Printf("account = %v\n", accountID.String())

	response, err = hedera.NewTokenCreateTransaction().
		SetName("ffff").
		SetSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasury(operatorID).
		SetAdminKey(operatorKey).
		SetFreezeKey(operatorKey).
		SetWipeKey(operatorKey).
		SetKycKey(operatorKey).
		SetSupplyKey(operatorKey).
		SetFreezeDefault(false).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = response.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	tokenID := receipt.GetTokenID()

	fmt.Printf("token = %v\n", tokenID.String())

	transaction, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(accountID).
		AddTokenID(tokenID).
		Build(client)
	if err != nil {
		panic(err)
	}

	response, err = transaction.
		Sign(key).
		Sign(operatorPrivateKey).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = response.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Associated account %v with token %v\n", accountID.String(), tokenID.String())

	response, err = hedera.NewTokenGrantKycTransaction().
		SetTokenID(tokenID).
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = response.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Granted KYC for account %v on token %v\n", accountID.String(), tokenID.String())

	response, err = hedera.NewTransferTransaction().
		AddTokenSender(tokenID, operatorID, 10).
		AddTokenRecipient(tokenID, accountID, 10).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = response.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sent 10 tokens from account %v to account %v on token %v\n", operatorID.String(), accountID.String(), tokenID.String())

	response, err = hedera.NewTokenWipeTransaction().
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = response.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wiped balance of account %v\n", accountID.String())

	response, err = hedera.NewTokenDeleteTransaction().
		SetTokenID(tokenID).
		Execute(client)

	receipt, err = response.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted token %v\n", tokenID.String())

	transaction, err = hedera.NewAccountDeleteTransaction().
		SetDeleteAccountID(accountID).
		SetTransferAccountID(operatorID).
		Build(client)
	if err != nil {
		panic(err)
	}

	response, err = transaction.
		Sign(key).
		Sign(operatorPrivateKey).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = response.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted account %v\n", accountID.String())
}

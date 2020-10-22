package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
	"time"
)

func main() {
	client, err := hedera.ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = hedera.ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	var operatorAccountID hedera.AccountID
	var operatorPrivateKey hedera.PrivateKey

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err = hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorPrivateKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorPrivateKey)
	}

	operatorKey := client.GetOperatorKey()

	key, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Printf("privateKey = %v\n", key.String())
	fmt.Printf("publicKey = %v\n", key.PublicKey().String())

	resp, err := hedera.NewAccountCreateTransaction().
		SetKey(key.PublicKey()).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err := resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	accountID := *receipt.AccountID

	fmt.Printf("account = %v\n", accountID.String())

	nodeIDs := make([]hedera.AccountID, 1)
	nodeIDs[0] = resp.NodeID

	resp, err = hedera.NewTokenCreateTransaction().
		SetName("ffff").
		SetSymbol("F").
		SetNodeAccountIDs(nodeIDs).
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasury(operatorAccountID).
		SetAdminKey(operatorKey).
		SetFreezeKey(operatorKey).
		SetWipeKey(operatorKey).
		SetKycKey(operatorKey).
		SetSupplyKey(operatorKey).
		SetFreezeDefault(false).
		SetExpirationTime(uint64(time.Now().Add(7890000 * time.Second).Unix())).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	tokenID := *receipt.TokenID

	fmt.Printf("token = %v\n", tokenID.String())

	transaction, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	resp, err = transaction.
		Sign(key).
		Sign(operatorPrivateKey).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Associated account %v with token %v\n", accountID.String(), tokenID.String())

	resp, err = hedera.NewTokenGrantKycTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs(nodeIDs).
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Granted KYC for account %v on token %v\n", accountID.String(), tokenID.String())

	//resp, err = hedera.NewTokenTransferTransaction().
	//	AddSender(tokenID, operatorID, 10).
	//	AddRecipient(tokenID, accountID, 10).
	//	Execute(client)
	//if err != nil {
	//	panic(err)
	//}
	//
	//receipt, err = response.GetReceipt(client)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("Sent 10 tokens from account %v to account %v on token %v\n", operatorID.String(), accountID.String(), tokenID.String())

	resp, err = hedera.NewTokenWipeTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs(nodeIDs).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wiped balance of account %v\n", accountID.String())

	resp, err = hedera.NewTokenDeleteTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted token %v\n", tokenID.String())

	accountDeleteTx, err := hedera.NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetTransferAccountID(operatorAccountID).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	resp, err = accountDeleteTx.
		Sign(key).
		Sign(operatorPrivateKey).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted account %v\n", accountID.String())
}

package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
)

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			println(err.Error(), ": error setting up client from config file")
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" && client.GetOperatorPublicKey().Bytes() == nil {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			println(err.Error(), ": error converting string to AccountID")
			return
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			println(err.Error(), ": error converting string to PrivateKey")
			return
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

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

	transactionResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(key1.PublicKey()).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating account")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account creation receipt")
		return
	}

	accountID1 := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", accountID1.String())

	transactionResponse, err = hedera.NewAccountCreateTransaction().
		SetKey(key2.PublicKey()).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating second account")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving second account creation receipt")
		return
	}

	accountID2 := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", accountID2.String())

	transactionResponse, err = hedera.NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetMaxTransactionFee(hedera.NewHbar(1000)).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFreezeKey(client.GetOperatorPublicKey()).
		SetWipeKey(client.GetOperatorPublicKey()).
		SetKycKey(client.GetOperatorPublicKey()).
		SetSupplyKey(client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		SetMaxTransactionFee(hedera.NewHbar(1000)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating token")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token creation receipt")
		return
	}

	tokenID := *transactionReceipt.TokenID

	fmt.Printf("token = %v\n", tokenID.String())

	transaction, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(accountID1).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing token associate transaction")
		return
	}

	transactionResponse, err = transaction.
		Sign(key1).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error associating token")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token associate transaction receipt")
		return
	}

	fmt.Printf("Associated account %v with token %v\n", accountID1.String(), tokenID.String())

	transaction, err = hedera.NewTokenAssociateTransaction().
		SetAccountID(accountID2).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing second token associate transaction")
		return
	}

	transactionResponse, err = transaction.
		Sign(key2).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing second token associate transaction")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving second token associate transaction receipt")
		return
	}

	fmt.Printf("Associated account %v with token %v\n", accountID2.String(), tokenID.String())

	transactionResponse, err = hedera.NewTokenGrantKycTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetAccountID(accountID1).
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

	transactionResponse, err = hedera.NewTokenGrantKycTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetAccountID(accountID2).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error granting kyc to second account")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving grant kyc transaction receipt")
		return
	}

	fmt.Printf("Granted KYC for account %v on token %v\n", accountID2.String(), tokenID.String())

	transferTransaction, err := hedera.NewTransferTransaction().
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID1, 10).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer from operator to account1")
		return
	}

	transferTransaction = transferTransaction.Sign(key1)

	transactionResponse, err = transferTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error transferring from operator to account1")
		return
	}

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

	transferTransaction, err = hedera.NewTransferTransaction().
		AddTokenTransfer(tokenID, accountID1, -10).
		AddTokenTransfer(tokenID, accountID2, 10).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer from account1 to account2")
		return
	}

	transferTransaction = transferTransaction.Sign(key1)

	transactionResponse, err = transferTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error transferring from account1 to account2")
		return
	}

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
		AddTokenTransfer(tokenID, accountID2, -10).
		AddTokenTransfer(tokenID, accountID1, 10).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer from account2 to account1")
		return
	}

	transferTransaction = transferTransaction.Sign(key2)

	transactionResponse, err = transferTransaction.Execute(client)
	if err != nil {
		println(err.Error(), ": error transferring from account2 to account1")
		return
	}

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

	transactionResponse, err = hedera.NewTokenWipeTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetAccountID(accountID1).
		SetTokenID(tokenID).
		SetAmount(10).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error wiping from token")
		return
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token wipe transaction receipt")
		return
	}

	fmt.Printf("Wiped account %v on token %v\n", accountID1.String(), tokenID.String())

	transactionResponse, err = hedera.NewTokenDeleteTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting token")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token delete transaction receipt")
		return
	}

	fmt.Printf("Deleted token %v\n", tokenID.String())

	accountDeleteTx, err := hedera.NewAccountDeleteTransaction().
		SetAccountID(accountID1).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account delete transaction")
		return
	}

	transactionResponse, err = accountDeleteTx.
		Sign(key1).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting account 1")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving transfer transaction receipt")
		return
	}

	fmt.Printf("Deleted account %v\n", accountID1.String())

	accountDeleteTx, err = hedera.NewAccountDeleteTransaction().
		SetAccountID(accountID2).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account delete transaction")
		return
	}

	transactionResponse, err = accountDeleteTx.
		Sign(key2).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error deleting account2")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account delete transaction receipt")
		return
	}

	fmt.Printf("Deleted account %v\n", accountID2.String())
}

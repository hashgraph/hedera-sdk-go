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
	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	metadataList := []string{"meta1",
		"meta2",
		"meta3",
		"meta4",
		"meta5"}
	nftCreateTransaction, err := hedera.NewTokenCreateTransaction().
		SetTokenName("HIP-542 Example Collection").SetTokenSymbol("HIP-542").
		SetTokenType(hedera.TokenTypeNonFungibleUnique).SetDecimals(0).
		SetInitialSupply(0).SetMaxSupply(10).
		SetTreasuryAccountID(client.GetOperatorAccountID()).SetSupplyType(hedera.TokenSupplyTypeFinite).
		SetAdminKey(operatorKey).SetFreezeKey(operatorKey).SetSupplyKey(operatorKey).SetMetadataKey(operatorKey).FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token transaction", err))
	}
	// Sign the transaction with the operator key
	nftSignTransaction := nftCreateTransaction.Sign(operatorKey)
	// Submit the transaction to the Hedera network
	nftCreateSubmit, err := nftSignTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting transaction", err))
	}
	// Get transaction receipt information
	nftCreateReceipt, err := nftCreateSubmit.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error receiving receipt", err))
	}
	// Get token id from the transaction
	nftTokenID := *nftCreateReceipt.TokenID
	fmt.Println("Created NFT with token id: ", nftTokenID)

	// mint nfts
	nftCollection := []hedera.TransactionReceipt{}

	for i, s := range metadataList {
		mintTransaction, err := hedera.NewTokenMintTransaction().SetTokenID(nftTokenID).SetMetadata([]byte(s)).FreezeWith(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error creating mint transaction", err))
		}
		mintTransactionSubmit, err := mintTransaction.Sign(operatorKey).Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error submitting transaction", err))
		}
		receipt, err := mintTransactionSubmit.GetReceipt(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error receiving receipt", err))
		}
		nftCollection = append(nftCollection, receipt)
		fmt.Println("Created NFT ", nftTokenID.String(), " with serial: ", nftCollection[i].SerialNumbers[0])
	}
	exampleNftId := nftTokenID.Nft(nftCollection[0].SerialNumbers[0])

	//Returns the info for the specified NFT id
	nftInfo, err := hedera.NewTokenNftInfoQuery().SetNftID(exampleNftId).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error info query transaction", err))
	}
	fmt.Println(string(nftInfo[0].Metadata))

	metadataUpdateSubmit, err := hedera.NewTokenUpdateNfts().
		SetTokenID(nftTokenID).
		SetSerialNumbers([]int64{1, 2, 3}).
		SetMetadata([]byte("updated")).
		Sign(operatorKey).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting transaction", err))
	}

	receipt, err := metadataUpdateSubmit.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error receiving receipt", err))
	}
	fmt.Println("metadata update: ", receipt)
}

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

	operatorId, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromStringEd25519(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorId, operatorKey)

	supplyKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating supply key", err))
	}
	freezeKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating freeze key", err))
	}
	wipeKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating wipe key", err))
	}
	/**
	 *     Example 1
	 *
	 * Step 1
	 *
	 * Create an NFT using the Hedera Token Service
	 */
	fmt.Println("Example 1")
	// IPFS content identifiers for the NFT metadata
	cid := []string{"QmNPCiNA3Dsu3K5FxDPMG5Q3fZRwVTg14EXA92uqEeSRXn",
		"QmZ4dgAgt8owvnULxnKxNe8YqpavtVCXmc1Lt2XajFpJs9",
		"QmPzY5GxevjyfMUF5vEAjtyRoigzWp47MiKAtLBduLMC1T",
		"Qmd3kGgSrAwwSrhesYcY7K54f3qD7MDo38r7Po2dChtQx5",
		"QmWgkKz3ozgqtnvbCLeh7EaR1H8u5Sshx3ZJzxkcrT3jbw"}
	// Creating the transaction for token creation
	nftCreateTransaction, err := hedera.NewTokenCreateTransaction().
		SetTokenName("HIP-542 Example Collection").SetTokenSymbol("HIP-542").
		SetTokenType(hedera.TokenTypeNonFungibleUnique).SetDecimals(0).
		SetInitialSupply(0).SetMaxSupply(int64(len(cid))).
		SetTreasuryAccountID(operatorId).SetSupplyType(hedera.TokenSupplyTypeFinite).
		SetAdminKey(operatorKey).SetFreezeKey(freezeKey).SetWipeKey(wipeKey).SetSupplyKey(supplyKey).FreezeWith(client)
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

	/**
	 * Step 2
	 *
	 * Mint the NFT
	 */

	nftCollection := []hedera.TransactionReceipt{}

	for i, s := range cid {
		mintTransaction, err := hedera.NewTokenMintTransaction().SetTokenID(nftTokenID).SetMetadata([]byte(s)).FreezeWith(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error creating mint transaction", err))
		}
		mintTransactionSubmit, err := mintTransaction.Sign(supplyKey).Execute(client)
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
	/**
	 * Step 3
	 *
	 * Create an ECDSA public key alias
	 */

	fmt.Println("Creating new account...")
	privateKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating private key", err))
	}
	publicKey := privateKey.PublicKey()
	// Assuming that the target shard and realm are known.
	// For now they are virtually always 0 and 0.
	aliasAccountId := publicKey.ToAccountID(0, 0)
	fmt.Println("New account ID: ", aliasAccountId)
	fmt.Println("Just the aliasKey: ", aliasAccountId.AliasKey)

	/**
	 * Step 4
	 *
	 * Tranfer the NFT to the public key alias using the transfer transaction
	 */
	nftTransferTransaction, err := hedera.NewTransferTransaction().AddNftTransfer(exampleNftId, operatorId, *aliasAccountId).FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating transaction", err))
	}
	// Sign the transaction with the operator key
	nftTransferTransactionSign := nftTransferTransaction.Sign(operatorKey)
	// Submit the transaction to the Hedera network
	nftTransferTransactionSubmit, err := nftTransferTransactionSign.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting transaction", err))
	}
	// Get transaction receipt information here
	fmt.Println(nftTransferTransactionSubmit.GetReceipt(client))

	/**
	 * Step 5
	 *
	 * Return the new account ID in the child record
	 */

	//Returns the info for the specified NFT id
	nftInfo, err := hedera.NewTokenNftInfoQuery().SetNftID(exampleNftId).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error info query transaction", err))
	}
	nftOwnerAccountId := nftInfo[0].AccountID
	fmt.Println("Current owner account id: ", nftOwnerAccountId)

	/**
	 * Step 6
	 *
	 * Show the new account ID owns the NFT
	 */
	accountInfo, err := hedera.NewAccountInfoQuery().SetAccountID(*aliasAccountId).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error account info query", err))
	}
	fmt.Println("The normal account ID of the given alias ", accountInfo.AccountID)

	if nftOwnerAccountId == accountInfo.AccountID {
		fmt.Println("The NFT owner accountId matches the accountId created with the HTS")
	} else {
		fmt.Println("The two account IDs does not match")
	}

	/**
	 *     Example 2
	 *
	 * Step 1
	 *
	 * Create a fungible HTS token using the Hedera Token Service
	 */
	fmt.Println("Example 2")

	tokenCreateTransaction, err := hedera.NewTokenCreateTransaction().SetTokenName("HIP-542 Token").
		SetTokenSymbol("H542").SetTokenType(hedera.TokenTypeFungibleCommon).SetTreasuryAccountID(operatorId).
		SetInitialSupply(10000).SetDecimals(2).SetAutoRenewAccount(operatorId).FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating transaction", err))
	}
	// Sign the transaction with the operator key
	tokenCreateTransactionSign := tokenCreateTransaction.Sign(operatorKey)
	// Submit the transaction to the Hedera network
	tokenCreateTransactionSubmit, err := tokenCreateTransactionSign.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting transaction", err))
	}

	// Get transaction receipt information
	tokenCreateTransactionReceipt, err := tokenCreateTransactionSubmit.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving receipt", err))
	}
	tokenId := *tokenCreateTransactionReceipt.TokenID
	fmt.Println("Created token with token id: ", tokenId)

	/**
	 * Step 2
	 *
	 * Create an ECDSA public key alias
	 */
	privateKey2, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating private key", err))
	}
	publicKey2 := privateKey2.PublicKey()
	// Assuming that the target shard and realm are known.
	// For now they are virtually always 0 and 0.
	aliasAccountId2 := *publicKey2.ToAccountID(0, 0)
	fmt.Println("New account ID: ", aliasAccountId2)
	fmt.Println("Just the aliasKey: ", aliasAccountId2.AliasKey)

	/**
	 * Step 3
	 *
	 * Transfer the fungible token to the public key alias
	 */

	tokenTransferTransaction, err := hedera.NewTransferTransaction().
		AddTokenTransfer(tokenId, operatorId, -10).AddTokenTransfer(tokenId, aliasAccountId2, 10).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating transaction", err))
	}
	// Sign the transaction with the operator key
	tokenTransferTransactionSign := tokenTransferTransaction.Sign(operatorKey)
	// Submit the transaction to the Hedera network
	tokenTransferTransactionSubmit, err := tokenTransferTransactionSign.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting transaction", err))
	}
	// Get transaction receipt information
	fmt.Println(tokenTransferTransactionSubmit.GetReceipt(client))

	/**
	 * Step 4
	 *
	 * Return the new account ID in the child record
	 */

	accountId2Info, err := hedera.NewAccountInfoQuery().SetAccountID(aliasAccountId2).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing acount info query", err))
	}
	accountId2 := accountId2Info.AccountID
	fmt.Println("The normal account ID of the given alias: ", accountId2)

	/**
	 * Step 5
	 *
	 * Show the new account ID owns the fungible token
	 */

	accountBalances, err := hedera.NewAccountBalanceQuery().SetAccountID(aliasAccountId2).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error receiving account balance", err))
	}

	tokenBalanceAccountId2 := accountBalances.Tokens.Get(tokenId)
	if tokenBalanceAccountId2 == 10 {
		fmt.Println(`Account is created succesfully using HTS "TransferTransaction"`)
	} else {
		fmt.Println("Creating account with HTS using public key alias failed")
	}

}

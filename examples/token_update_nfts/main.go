package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

/**
 * @summary E2E-HIP-657 https://hips.hedera.com/hip/hip-657
 * @description Update nfts metadata of non-fungible token with metadata key
 */
func main() {
	var client *hiero.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	metadataKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating metadata key", err))
	}
	fmt.Println("create metadata key: ", metadataKey.String())

	var initialMetadataList = [][]byte{{2, 1}, {1, 2}}
	var updatedMetadata = []byte{22, 22}

	// Create token with metadata key
	nftCreateTransaction, err := hiero.NewTokenCreateTransaction().
		SetTokenName("HIP-542 Example Collection").SetTokenSymbol("HIP-542").
		SetTokenType(hiero.TokenTypeNonFungibleUnique).SetDecimals(0).
		SetInitialSupply(0).SetMaxSupply(10).
		SetTreasuryAccountID(client.GetOperatorAccountID()).SetSupplyType(hiero.TokenSupplyTypeFinite).
		SetAdminKey(operatorKey).SetFreezeKey(operatorKey).SetSupplyKey(operatorKey).SetMetadataKey(metadataKey).FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token transaction", err))
	}

	// Sign the transaction with the operator key
	nftSignTransaction := nftCreateTransaction.Sign(operatorKey)

	// Submit the transaction to the Hiero network
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

	tokenInfo, err := hiero.NewTokenInfoQuery().SetTokenID(nftTokenID).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("Token metadata key: ", tokenInfo.MetadataKey.String())

	// Mint nfts
	mintTransaction, _ := hiero.NewTokenMintTransaction().SetTokenID(nftTokenID).SetMetadatas(initialMetadataList).FreezeWith(client)
	for _, v := range mintTransaction.GetMetadatas() {
		fmt.Println("Set metadata: ", v)
	}

	mintTransactionSubmit, err := mintTransaction.Sign(operatorKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error minting NFT", err))
	}
	receipt, err := mintTransactionSubmit.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error minting NFT", err))
	}

	// Check that metadata was set correctly
	serials := receipt.SerialNumbers
	fmt.Println(serials)
	var metadataAfterMint = make([][]byte, len(initialMetadataList))
	for i, v := range serials {
		nftID := hiero.NftID{TokenID: nftTokenID, SerialNumber: v}
		nftInfo, err := hiero.NewTokenNftInfoQuery().SetNftID(nftID).Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error getting token info", err))
		}
		fmt.Println(nftInfo)
		metadataAfterMint[i] = nftInfo[0].Metadata
	}
	fmt.Println("Metadata after mint: ", metadataAfterMint)

	// Create account owner of nft
	accountCreateTransaction, err := hiero.NewAccountCreateTransaction().
		SetKey(operatorKey).SetMaxAutomaticTokenAssociations(10). // If the account does not have any automatic token association slots open ONLY then associate the NFT to the account
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}
	receipt, err = accountCreateTransaction.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error receiving receipt", err))
	}
	newAccountId := receipt.AccountID

	// Transfer the NFT to the new account
	tokenTransferTransaction, err := hiero.NewTransferTransaction().AddNftTransfer(nftTokenID.Nft(serials[0]), operatorAccountID, *newAccountId).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error transfering nft", err))
	}
	_, err = tokenTransferTransaction.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}

	// update nfts metadata
	metadataUpdateTransaction, err := hiero.NewTokenUpdateNftsTransaction().
		SetTokenID(nftTokenID).
		SetSerialNumbers(serials).
		SetMetadata(updatedMetadata).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating transaction", err))
	}
	fmt.Println("Updatad metadata: ", metadataUpdateTransaction.GetMetadata())
	metadataUpdateSubmit, err := metadataUpdateTransaction.Sign(metadataKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting transaction", err))
	}

	receipt, err = metadataUpdateSubmit.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error receiving receipt", err))
	}
	fmt.Println("Metadata update status: ", receipt.Status)

	// Check that metadata for the NFT was updated correctly
	for _, v := range serials {
		nftID := hiero.NftID{TokenID: nftTokenID, SerialNumber: v}
		nftInfo, err := hiero.NewTokenNftInfoQuery().SetNftID(nftID).Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error getting token info", err))
		}
		fmt.Println("Metadata after update for serial number ", v, ": ", nftInfo[0].Metadata)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

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

	updateMutableTokenMetadata(client)

	updateImmutableTokenMetadata(client)
}

func updateMutableTokenMetadata(client *hiero.Client) {
	// Create admin key
	adminKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating admin key", err))
	}
	fmt.Println("create admin key: ", adminKey.String())

	var initialMetadata = []byte{1, 2, 3}
	var newMetadata = []byte{3, 4, 5, 6}

	// Create the token
	tx, err := hiero.NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(hiero.TokenTypeFungibleCommon). // The same flow can be executed with a TokenTypeNonFungibleUnique (i.e. HIP-765)
		SetTokenMetadata(initialMetadata).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(adminKey).
		SetFreezeDefault(false).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	resp, err := tx.Sign(adminKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}

	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	fmt.Println("created token: ", receipt.TokenID.String())

	// Query the token info to get the metadata after creation
	info, err := hiero.NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's metadata after creation: ", info.Metadata)

	// Update the token's metadata
	tx1, err := hiero.NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	resp, err = tx1.Sign(adminKey).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	// Query the token info to get the metadata after update
	info, err = hiero.NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's metadata after update: ", info.Metadata)
}

func updateImmutableTokenMetadata(client *hiero.Client) {
	// Create metadata key
	metadataKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating admin key", err))
	}
	fmt.Println("create metadata key: ", metadataKey.String())

	var initialMetadata = []byte{1, 2, 3}
	var newMetadata = []byte{3, 4, 5, 6}

	// Create the token
	resp, err := hiero.NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(hiero.TokenTypeFungibleCommon). // The same flow can be executed with a TokenTypeNonFungibleUnique (i.e. HIP-765)
		SetTokenMetadata(initialMetadata).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetMetadataKey(metadataKey).
		SetFreezeDefault(false).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}

	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	fmt.Println("created token: ", receipt.TokenID.String())

	// Query the token info to get the metadata after creation
	info, err := hiero.NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's metadata after creation: ", info.Metadata)

	// Update the token's metadata
	tx, err := hiero.NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetTokenMetadata(newMetadata).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	resp, err = tx.Sign(metadataKey).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	// Query the token info to get the metadata after update
	info, err = hiero.NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's metadata after update: ", info.Metadata)
}

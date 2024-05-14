package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

/**
 * @summary HIP-540 https://hips.hedera.com/hip/hip-540
 * @description Change or remove existing keys from a token
 */
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

	// Create supply key
	supplyKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating supply key", err))
	}
	fmt.Println("create supply key: ", supplyKey.String())

	// Create wipe key
	wipeKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating supply key", err))
	}
	fmt.Println("create wipe key: ", supplyKey.String())

	// Create the token
	tx, err := hedera.NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(hedera.TokenTypeFungibleCommon).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetFreezeDefault(false).
		SetSupplyKey(supplyKey).
		SetWipeKey(wipeKey).
		SetAdminKey(operatorKey).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}

	receipt, err := tx.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	tokenID := receipt.TokenID
	fmt.Println("created token: ", tokenID)

	// Query the token info to get the supply key after creation
	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's supply key after creation: ", info.SupplyKey)
	fmt.Println("token's wipe key after creation: ", info.WipeKey)
	fmt.Println("token's admin key after creation: ", info.AdminKey)

	removeWipeKeyFullValidation(client, *receipt.TokenID, operatorKey)
	updateSupplyKeyFullValidation(client, *receipt.TokenID, supplyKey)
	updateSupplyKeyNoValidation(client, *receipt.TokenID, supplyKey)
}

// With "full" verification mode, our required key
// structure is a 1/2 threshold with components:
//   - Admin key
//   - A 2/2 list including the role key and its replacement key
func removeWipeKeyFullValidation(client *hedera.Client, tokenID hedera.TokenID, adminKey hedera.PrivateKey) {
	tx, err := hedera.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetWipeKey(hedera.NewKeyList()).
		SetKeyVerificationMode(hedera.FULL_VALIDATION).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	resp, err := tx.Sign(adminKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's wipe key after removal: ", info.WipeKey)
}

func updateSupplyKeyFullValidation(client *hedera.Client, tokenID hedera.TokenID, oldSupplyKey hedera.PrivateKey) {
	newValidKey, _ := hedera.GeneratePrivateKey()
	tx, err := hedera.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(newValidKey.PublicKey()).
		SetKeyVerificationMode(hedera.FULL_VALIDATION).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	_, err = newValidKey.SignTransaction(&tx.Transaction)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing tx", err))
	}
	_, err = oldSupplyKey.SignTransaction(&tx.Transaction)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing tx", err))
	}

	resp, err := tx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's supply key after update: ", info.SupplyKey)
}

func updateSupplyKeyNoValidation(client *hedera.Client, tokenID hedera.TokenID, oldSupplyKey hedera.PrivateKey) {
	newValidKey, _ := hedera.GeneratePrivateKey()
	tx1, err := hedera.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(newValidKey.PublicKey()).
		SetKeyVerificationMode(hedera.NO_VALIDATION).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	resp, err := tx1.Sign(oldSupplyKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's supply key after update: ", info.SupplyKey)
}

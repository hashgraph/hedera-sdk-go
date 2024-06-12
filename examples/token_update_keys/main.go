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

	// Create admin key
	adminKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating admin key", err))
	}
	fmt.Println("create wipe key: ", adminKey.String())

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
		SetAdminKey(adminKey).
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
	tokenID := receipt.TokenID
	fmt.Println("created token: ", tokenID)

	// Query the token info after creation
	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(*receipt.TokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's supply key after creation: ", info.SupplyKey)
	fmt.Println("token's wipe key after creation: ", info.WipeKey)
	fmt.Println("token's admin key after creation: ", info.AdminKey)

	removeWipeKeyFullValidation(client, *receipt.TokenID, adminKey)
	newSupplyKey := updateSupplyKeyFullValidation(client, *receipt.TokenID, supplyKey)
	removeSupplyKeyNoValidation(client, *receipt.TokenID, newSupplyKey)
	removeAdminKeyNoValidation(client, *receipt.TokenID, adminKey)
}

// With "full" verification mode, our required key
// structure is a 1/2 threshold with components:
//   - Admin key
//   - A 2/2 list including the role key and its replacement key
func removeWipeKeyFullValidation(client *hedera.Client, tokenID hedera.TokenID, adminKey hedera.PrivateKey) {
	// Remove wipe key by setting it to a empty key list
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

	// Query the token info to get the wipe key after removal
	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's wipe key after removal: ", info.WipeKey)
}

func updateSupplyKeyFullValidation(client *hedera.Client, tokenID hedera.TokenID, oldSupplyKey hedera.PrivateKey) hedera.PrivateKey {
	newSupplyKey, _ := hedera.GeneratePrivateKey()

	// Update  supply key by setting it to a new key
	tx, err := hedera.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(newSupplyKey.PublicKey()).
		SetKeyVerificationMode(hedera.FULL_VALIDATION).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	// Sign with old and new supply keys
	_, err = newSupplyKey.SignTransaction(&tx.Transaction)
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

	// Query the token info to get the supply key after update
	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's supply key after update: ", info.SupplyKey)

	return newSupplyKey
}

func removeSupplyKeyNoValidation(client *hedera.Client, tokenID hedera.TokenID, oldSupplyKey hedera.PrivateKey) {
	zeroNewKey, _ := hedera.ZeroKey()

	// Remove supply key by setting it to a zero key
	tx, err := hedera.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(zeroNewKey).
		SetKeyVerificationMode(hedera.NO_VALIDATION).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	// Sign with old supply key
	resp, err := tx.Sign(oldSupplyKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	// Query the token info to get the supply key after removal
	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}

	newSupplyKey := info.SupplyKey.(hedera.PublicKey)
	fmt.Println("token's supply key after zero out: ", newSupplyKey.StringRaw())
}

func removeAdminKeyNoValidation(client *hedera.Client, tokenID hedera.TokenID, oldAdminKey hedera.PrivateKey) {
	// Remove admin key by setting it to a zero key
	tx, err := hedera.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetAdminKey(hedera.NewKeyList()).
		SetKeyVerificationMode(hedera.NO_VALIDATION).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	// Sign with old supply key
	resp, err := tx.Sign(oldAdminKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	// Query the token info to get the admin key after removal
	info, err := hedera.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}

	fmt.Println("token's admin key after removal: ", info.AdminKey)
}

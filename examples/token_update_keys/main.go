package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

/**
 * @summary HIP-540 https://hips.hedera.com/hip/hip-540
 * @description Change or remove existing keys from a token
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

	// Create admin key
	adminKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating admin key", err))
	}
	fmt.Println("create wipe key: ", adminKey.String())

	// Create supply key
	supplyKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating supply key", err))
	}
	fmt.Println("create supply key: ", supplyKey.String())

	// Create wipe key
	wipeKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating supply key", err))
	}
	fmt.Println("create wipe key: ", supplyKey.String())

	// Create the token
	tx, err := hiero.NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetTokenType(hiero.TokenTypeFungibleCommon).
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
	info, err := hiero.NewTokenInfoQuery().
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
func removeWipeKeyFullValidation(client *hiero.Client, tokenID hiero.TokenID, adminKey hiero.PrivateKey) {
	// Remove wipe key by setting it to a empty key list
	tx, err := hiero.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetWipeKey(hiero.NewKeyList()).
		SetKeyVerificationMode(hiero.FULL_VALIDATION).
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
	info, err := hiero.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's wipe key after removal: ", info.WipeKey)
}

func updateSupplyKeyFullValidation(client *hiero.Client, tokenID hiero.TokenID, oldSupplyKey hiero.PrivateKey) hiero.PrivateKey {
	newSupplyKey, _ := hiero.GeneratePrivateKey()

	// Update  supply key by setting it to a new key
	tx, err := hiero.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(newSupplyKey.PublicKey()).
		SetKeyVerificationMode(hiero.FULL_VALIDATION).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}

	// Sign with old and new supply keys
	_, err = newSupplyKey.SignTransaction(tx)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing tx", err))
	}
	_, err = oldSupplyKey.SignTransaction(tx)
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
	info, err := hiero.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}
	fmt.Println("token's supply key after update: ", info.SupplyKey)

	return newSupplyKey
}

func removeSupplyKeyNoValidation(client *hiero.Client, tokenID hiero.TokenID, oldSupplyKey hiero.PrivateKey) {
	zeroNewKey, _ := hiero.ZeroKey()

	// Remove supply key by setting it to a zero key
	tx, err := hiero.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(zeroNewKey).
		SetKeyVerificationMode(hiero.NO_VALIDATION).
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
	info, err := hiero.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}

	newSupplyKey := info.SupplyKey.(hiero.PublicKey)
	fmt.Println("token's supply key after zero out: ", newSupplyKey.StringRaw())
}

func removeAdminKeyNoValidation(client *hiero.Client, tokenID hiero.TokenID, oldAdminKey hiero.PrivateKey) {
	// Remove admin key by setting it to a zero key
	tx, err := hiero.NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetAdminKey(hiero.NewKeyList()).
		SetKeyVerificationMode(hiero.NO_VALIDATION).
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
	info, err := hiero.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info", err))
	}

	fmt.Println("token's admin key after removal: ", info.AdminKey)
}

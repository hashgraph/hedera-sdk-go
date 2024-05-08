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

	// Update the token's supply key to zero
	zeroKey, _ := hedera.PrivateKeyFromString("0000000000000000000000000000000000000000000000000000000000000000")
	tx1, err := hedera.NewTokenUpdateTransaction().
		SetTokenID(*receipt.TokenID).
		SetSupplyKey(zeroKey).
		SetKeyVerificationMode(hedera.NO_VALIDATION).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	resp, err := tx1.Sign(supplyKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
	_, err = resp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating token", err))
	}
}

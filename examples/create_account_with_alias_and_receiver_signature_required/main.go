package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	// client := hedera.ClientForTestnet()
	client := hedera.ClientForPreviewnet()
	myAccountId, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	myPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	client.SetOperator(myAccountId, myPrivateKey)

	// ## Example
	// Create an ED25519 admin private key and ECSDA private key
	// Extract the ECDSA public key public key
	// Extract the Ethereum public address
	// Use the `AccountCreateTransaction` and populate `setAlias(evmAddress)` field with the Ethereum public address and the `setReceiverSignatureRequired` to `true`
	// Sign the `AccountCreateTransaction` transaction with both the new private key and the admin key
	// Get the `AccountInfo` on the new account and show that the account has contractAccountId

	// Create an ED25519 admin private key and ECSDA private key
	adminKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		println(err.Error())
		return
	}

	privateKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		println(err.Error())
		return
	}
	// Extract the ECDSA public key public key
	publicKey := privateKey.PublicKey()
	// Extract the Ethereum public address
	evmAddress := publicKey.ToEvmAddress()

	evmAddressBytes, err := hex.DecodeString(evmAddress)
	if err != nil {
		println(err.Error())
		return
	}

	// Use the `AccountCreateTransaction` and set the EVM address field to the Ethereum public address
	response, err := hedera.NewAccountCreateTransaction().SetReceiverSignatureRequired(true).SetInitialBalance(hedera.HbarFromTinybar(100)).
		SetKey(adminKey).SetAlias(evmAddressBytes).Sign(adminKey).Sign(privateKey).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}

	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt}")
		return
	}

	newAccountId := *transactionReceipt.AccountID

	// Get the `AccountInfo` on the new account and show it is a hollow account by not having a public key
	info, err := hedera.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}
	// Verify account is created with the public address provided
	fmt.Println(info.ContractAccountID == publicKey.ToEvmAddress())
	// Verify the account Id is the same from the create account transaction
	fmt.Println(info.AccountID.String() == newAccountId.String())
}

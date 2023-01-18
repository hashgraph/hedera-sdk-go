package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	// ## Example 1:
	// Create an ECDSA private key
	// Get the ECDSA public key
	// Use the `AccountCreateTransaction` and populate the `setAlias(<ECDSA_public_key)` field
	// Sign the `AccountCreateTransaction` using an existing Hedera account and key to pay for the transaction fee
	// Execute the transaction
	// Return the Hedera account ID from the receipt of the transaction
	// Get the `AccountInfo` using the new account ID
	// Get the `AccountInfo` using the account public key in `0.0.aliasPublicKey` format
	// Show the public key and the public key alias are the same on the account
	// Show this account has a corresponding EVM address in the mirror node

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

	// Create an ECDSA private key

	privateKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		println(err.Error())
		return
	}
	// Get the ECDSA public key
	publicKey := privateKey.PublicKey()

	// Use the `AccountCreateTransaction` and populate the `setAlias(<ECDSA_public_key)` field
	// Sign the `AccountCreateTransaction` using an existing Hedera account and key to pay for the transaction fee
	// Execute the transaction
	tx, err := hedera.NewAccountCreateTransaction().SetAliasKey(publicKey).Execute(client)
	if err != nil {
		println("error creating account ", err.Error())
		return
	}
	receipt, err := tx.GetReceipt(client)
	if err != nil {
		println(err.Error())
		return
	}
	// Return the Hedera account ID from the receipt of the transaction
	newAccountId := *receipt.AccountID
	// Get the `AccountInfo` using the new account ID
	info, err := hedera.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}

	idWithPublicKey, err := hedera.AccountIDFromString("0.0." + publicKey.StringRaw())
	if err != nil {
		println(err.Error())
		return
	}
	// Get the `AccountInfo` using the account public key in `0.0.aliasPublicKey` format
	info2, err := hedera.NewAccountInfoQuery().SetAccountID(idWithPublicKey).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}
	// Show the public key and the public key alias are the same on the account
	fmt.Println(info.AccountID.Account == info2.AccountID.Account)
	fmt.Println(info.AliasKey.StringRaw() == info2.AliasKey.StringRaw())
	fmt.Println(info.Key.String() == info2.Key.String())
	fmt.Println(info.Key.String() == info.AliasKey.String())
	// Wait for the mirror node to have the information
	time.Sleep(time.Second * 30)
	// Show this account has a corresponding EVM address in the mirror node
	mirrorNodeInfo, err := getAccountInfoFromMirrorNode(newAccountId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Verify we are getting the info for the same account Id
	fmt.Println(mirrorNodeInfo.Accounts[0].Account == newAccountId.String())
	// Verify the account does have an EVM address
	fmt.Println(mirrorNodeInfo.Accounts[0].EvmAddress == "0x"+publicKey.ToEvmAddress())
	// Verify the account does have Hedera account public Key of type ECDSA
	fmt.Println(mirrorNodeInfo.Accounts[0].Key.Type == "ECDSA_SECP256K1")
	// Verify the account public key coresponding to the evm address
	fmt.Println(mirrorNodeInfo.Accounts[0].Key.Key == publicKey.StringRaw())

	// ## Example 2:
	// Create an ED2519 private key
	// Get the ED2519 public key
	// Use the `AccountCreateTransaction` and populate the `setAlias(<ED2519_public_key)` field
	// Sign the `AccountCreateTransaction` using an existing Hedera account and key to pay for the transaction fee
	// Execute the transaction
	// Return the Hedera account ID from the receipt of the transaction
	// Get the `AccountInfo` using the new account ID
	// Get the `AccountInfo` using the account public key in `0.0.aliasPublicKey` format
	// Show the public key and the public key alias are the same on the account

	// Create an ED2519 private key
	privateKey, err = hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		println(err.Error())
		return
	}
	// Get the ED2519 public key
	publicKey = privateKey.PublicKey()
	// Use the `AccountCreateTransaction` and populate the `setAlias(<ED2519_public_key)` field
	// Sign the `AccountCreateTransaction` using an existing Hedera account and key to pay for the transaction fee
	// Execute the transaction
	tx, err = hedera.NewAccountCreateTransaction().SetAliasKey(publicKey).SetInitialBalance(hedera.NewHbar(5)).Execute(client)
	if err != nil {
		println("error creating account ", err.Error())
		return
	}
	receipt, err = tx.GetReceipt(client)
	if err != nil {
		println(err.Error())
		return
	}
	// Return the Hedera account ID from the receipt of the transaction
	newAccountId = *receipt.AccountID
	// Get the `AccountInfo` using the new account ID
	info, err = hedera.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}

	idWithPublicKey, err = hedera.AccountIDFromString("0.0." + publicKey.StringRaw())
	if err != nil {
		println(err.Error())
		return
	}
	// Get the `AccountInfo` using the account public key in `0.0.aliasPublicKey` format
	info2, err = hedera.NewAccountInfoQuery().SetAccountID(idWithPublicKey).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}
	// Show the public key and the public key alias are the same on the account
	fmt.Println(info.AccountID.Account == info2.AccountID.Account)
	// Verify the account does have Alias Key of type ED
	fmt.Println(info.AliasKey.StringRaw() == info2.AliasKey.StringRaw())
	// Verify the account does have Hedera account public Key of type ED
	fmt.Println(info.Key.String() == info2.Key.String())

	// Wait for the mirror node to have the information
	time.Sleep(time.Second * 30)
	mirrorNodeInfo, err = getAccountInfoFromMirrorNode(newAccountId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Verify we are getting the info for the same account Id
	fmt.Println(mirrorNodeInfo.Accounts[0].Account == newAccountId.String())
	// Verify the account does have Hedera account public Key of type ED25519
	fmt.Println(mirrorNodeInfo.Accounts[0].Key.Type == "ED25519")
	// Verify the account public key coresponding to the Public Key
	fmt.Println(mirrorNodeInfo.Accounts[0].Key.Key == publicKey.StringRaw())
}

type Response struct {
	Accounts []struct {
		Account    string `json:"account"`
		EvmAddress string `json:"evm_address"`
		Key        struct {
			Type string `json:"_type"`
			Key  string `json:"key"`
		} `json:"key"`
	} `json:"accounts"`
}

func getAccountInfoFromMirrorNode(acountId hedera.AccountID) (*Response, error) {
	url := fmt.Sprintf("https://previewnet.mirrornode.hedera.com/api/v1/accounts?account.id=%d", acountId.Account)
	response, err := http.Get(url)
	if err != nil {
		return &Response{}, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &Response{}, err
	}

	var resp Response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return &Response{}, err
	}
	defer response.Body.Close()
	return &resp, nil
}

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

type Response struct {
	Accounts []struct {
		Account    string      `json:"account"`
		Alias      interface{} `json:"alias"`
		EvmAddress string      `json:"evm_address"`
		Key        interface{} `json:"key"`
	} `json:"accounts"`
}

func main() {
	// ## Example
	// Create a ECSDA private key
	// Extract the ECDSA public key
	// Extract the Ethereum public address
	// Use the `AccountCreateTransaction` and set the EVM address field to the Ethereum public address
	// Sign the transaction with the key that us paying for the transaction
	// Get the account ID from the receipt
	// Get the `AccountInfo` and return the account details
	// Verify the evm address provided for the account matches what is in the mirror node

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

	// Create a ECSDA private key
	privateKey, err := hedera.PrivateKeyGenerateEcdsa()
	if err != nil {
		println(err.Error())
		return
	}

	// Extract the ECDSA public key
	publicKey := privateKey.PublicKey()
	// Create an AccountID struct with EVM address
	evmAddressAccount, err := hedera.AccountIDFromEvmPublicAddress(publicKey.ToEvmAddress())
	if err != nil {
		println("error creating account from EVM address", err.Error())
		return
	}
	// Use the `AccountCreateTransaction` and set the EVM address
	tx, err := hedera.NewAccountCreateTransaction().SetEvmAddress(*evmAddressAccount.EvmAddress).SetInitialBalance(hedera.NewHbar(4)).Execute(client)
	if err != nil {
		println("error creating account ", err.Error())
		return
	}

	// Get Receipt
	receipt, err := tx.GetReceipt(client)
	if err != nil {
		println(err.Error())
		return
	}
	// Get the account ID from the receipt
	accountId := receipt.AccountID
	info, err := hedera.NewAccountInfoQuery().SetAccountID(*accountId).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}
	// Verify account is created with the public address provided
	fmt.Println(info.ContractAccountID == publicKey.ToEvmAddress())
	// Verify the account Id is the same from the create account transaction
	fmt.Println(info.AccountID.String() == accountId.String())
	// Verify the account does not have a Hedera public key /hollow account/
	fmt.Println(info.Key.String() == "{[]}")
	// Wait for the mirror node to have the information
	time.Sleep(time.Second * 30)
	mirrorNodeInfo, err := getAccountInfoFromMirrorNode(*accountId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Verify we are getting the info for the same account Id
	fmt.Println(mirrorNodeInfo.Accounts[0].Account == accountId.String())
	// Verify the account does not have an Alias Public Key
	fmt.Println(mirrorNodeInfo.Accounts[0].Alias == nil)
	// Verify the account does have an EVM address
	fmt.Println(mirrorNodeInfo.Accounts[0].EvmAddress == "0x"+publicKey.ToEvmAddress())
	// Verify the account does not have a Hedera public key /hollow account/
	fmt.Println(mirrorNodeInfo.Accounts[0].Key == nil)

	// To enhance the hollow account to have a public key the hollow account needs to be specified as a transaction fee payer in a HAPI transaction
	client.SetOperator(*accountId, privateKey)
	// Create a HAPI transaction and set the hollow account as the transaction fee payer
	tx, err = hedera.NewTransferTransaction().AddHbarTransfer(*accountId, hedera.NewHbar(-2)).AddHbarTransfer(myAccountId, hedera.NewHbar(2)).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}
	receipt, err = tx.GetReceipt(client)
	if err != nil {
		println(err.Error())
		return
	}

	info, err = hedera.NewAccountInfoQuery().SetAccountID(*accountId).Execute(client)
	if err != nil {
		println(err.Error())
		return
	}

	// Verify account is created with the public address provided
	fmt.Println(info.ContractAccountID == publicKey.ToEvmAddress())
	// Verify the account Id is the same from the create account transaction
	fmt.Println(info.AccountID.String() == accountId.String())
	// Verify the account does not have a Hedera public key /hollow account/
	fmt.Println(info.Key.String() == publicKey.String())
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

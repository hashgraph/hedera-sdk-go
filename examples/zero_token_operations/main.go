package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/hashgraph/hedera-sdk-go/v2/examples/contract_helper"
	"os"
)

type Contract struct {
	Bytecode string `json:"bytecode"`
}

// Steps 1-5 are executed through ContractHelper and calling HIP564Example Contract.
// Step 6 is executed through the SDK
func main() {
	var client *hedera.Client
	var err error
	var contract Contract
	// Retrieving network type from environment variable HEDERA_NETWORK, i.e. testnet
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	//Grab your testnet account ID and private key from the environment variable
	myAccountId, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	myPrivateKey, err := hedera.PrivateKeyFromStringEd25519(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	//Print your testnet account ID and private key to the console to make sure there was no error
	fmt.Printf("The account ID is = %v\n", myAccountId)
	fmt.Printf("The private key is = %v\n", myPrivateKey)

	client.SetOperator(myAccountId, myPrivateKey)
	//Generate new keys for the account you will create
	alicePrivateKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(err)
	}

	newAccountPublicKey := alicePrivateKey.PublicKey()

	//Create new account and assign the public key
	aliceAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(newAccountPublicKey).
		SetInitialBalance(hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar)).
		Execute(client)
	if err != nil {
		panic(err)
	}
	//Request the receipt of the transaction
	receipt, err := aliceAccount.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	//Get the new account ID from the receipt
	aliceAccountId := *receipt.AccountID
	fmt.Println("aliceAcountid is: ", aliceAccountId)
	//Transfer hbar from your testnet account to the new account
	transaction := hedera.NewTransferTransaction().
		AddHbarTransfer(myAccountId, hedera.HbarFrom(-1000, hedera.HbarUnits.Tinybar)).
		AddHbarTransfer(aliceAccountId, hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar))

	//Submit the transaction to a Hedera network
	transaction.Execute(client)

	rawContract, err := os.ReadFile("../precompile_example/ZeroTokenOperations.json")
	if err != nil {
		println(err.Error(), ": error reading json")
		return
	}

	err = json.Unmarshal([]byte(rawContract), &contract)
	if err != nil {
		println(err.Error(), ": error unmarshaling the json file")
		return
	}
	params, err := hedera.NewContractFunctionParameters().AddAddress(myAccountId.ToSolidityAddress())
	if err != nil {
		println(err.Error(), ": error adding first address to contract function parameters")
		return
	}
	params, err = params.AddAddress(aliceAccountId.ToSolidityAddress())
	if err != nil {
		println(err.Error(), ": error adding second address to contract function parameters")
		return
	}

	helper := contract_helper.NewContractHelper([]byte(contract.Bytecode), *params, client)
	helper.SetPayableAmountForStep(0, hedera.NewHbar(20)).AddSignerForStep(1, alicePrivateKey)
	_, err = helper.ExecuteSteps(0, 5, client)
	if err != nil {
		println(err.Error(), ": error in helper")
		return
	}
	transactionResponse, err := hedera.NewTokenCreateTransaction().
		SetTokenName("Black Sea LimeChain Token").
		SetTokenSymbol("BSL").
		SetTreasuryAccountID(myAccountId).
		SetInitialSupply(10000).
		SetDecimals(2).
		SetAutoRenewAccount(myAccountId).
		SetFreezeDefault(false).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating token")
		return
	}

	// Make sure the token create transaction ran
	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving token creation receipt")
		return
	}

	// Retrieve the token out of the receipt
	tokenID := *transactionReceipt.TokenID

	fmt.Printf("token = %v\n", tokenID.String())

	// Associating the token with the second account, so it can interact with the token
	associatingTransaction, err := hedera.NewTokenAssociateTransaction().
		// The account ID to be associated
		SetAccountID(aliceAccountId).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		// The token ID that the account will be associated to
		SetTokenIDs(tokenID).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing second token associate transaction")
		return
	}
	// Has to be signed by the account2's key
	transactionResponse, err = associatingTransaction.
		Sign(alicePrivateKey).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing second token associate transaction")
		return
	}

	// Make sure the transaction succeeded
	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving second token associate transaction receipt")
		return
	}

	fmt.Printf("Associated account %v with token %v\n", aliceAccountId.String(), tokenID.String())

	// Transfer 0 tokens
	transactionResponse, err = hedera.NewTransferTransaction().
		AddTokenTransfer(tokenID, myAccountId, 0).AddTokenTransfer(tokenID, aliceAccountId, 0).Execute(client)
	if err != nil {
		println(err.Error(), ": error transfering token")
		return
	}
	_, err = transactionResponse.GetRecord(client)
	if err != nil {
		println(err.Error(), ": error retrieving transaction")
		return
	}
}

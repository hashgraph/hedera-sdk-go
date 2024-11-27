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

	// create treasury account
	treasuryKey, _ := hiero.PrivateKeyGenerateEd25519()
	accountCreate, _ := hiero.NewAccountCreateTransaction().
		SetKey(treasuryKey).
		SetMaxAutomaticTokenAssociations(100).
		Execute(client)
	receipt, err := accountCreate.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error creating account", err))
	}

	treasury := *receipt.AccountID

	// create receiver account with unlimited auto associations
	receiverKey, _ := hiero.PrivateKeyGenerateEd25519()
	accountCreate, _ = hiero.NewAccountCreateTransaction().
		SetKey(receiverKey).
		SetMaxAutomaticTokenAssociations(-1).
		Execute(client)
	receipt, err = accountCreate.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error creating account", err))
	}

	receiver := *receipt.AccountID

	// create fungible token
	frozenTokenCreate, _ := hiero.NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("asdf").
		SetDecimals(18).
		SetInitialSupply(1_000_000).
		SetTreasuryAccountID(treasury).
		SetAdminKey(client.GetOperatorPublicKey()).FreezeWith(client)
	tokenCreate, _ := frozenTokenCreate.Sign(treasuryKey).Execute(client)
	receipt, err = tokenCreate.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error creating token", err))
	}

	tokenID := *receipt.TokenID

	// create NFT
	frozenTokenCreate, _ = hiero.NewTokenCreateTransaction().
		SetTokenName("Example Collection").SetTokenSymbol("ABC").
		SetTokenType(hiero.TokenTypeNonFungibleUnique).
		SetDecimals(0).
		SetInitialSupply(0).
		SetMaxSupply(10).
		SetTreasuryAccountID(treasury).
		SetSupplyType(hiero.TokenSupplyTypeFinite).
		SetSupplyKey(treasuryKey).
		SetAdminKey(treasuryKey).FreezeWith(client)
	tokenCreate, _ = frozenTokenCreate.Sign(treasuryKey).Execute(client)
	receipt, err = tokenCreate.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error creating token", err))
	}

	nftID := *receipt.TokenID

	// mint 3 NFTs
	initialMetadataList := [][]byte{{1}, {2}, {3}}
	frozenMint, _ := hiero.NewTokenMintTransaction().SetTokenID(nftID).SetMetadatas(initialMetadataList).FreezeWith(client)
	mint, _ := frozenMint.Sign(treasuryKey).Execute(client)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error minting nfts", err))
	}

	serials := receipt.SerialNumbers

	// transfer the tokens to the receiver account
	frozenTransfer, _ := hiero.NewTransferTransaction().
		AddNftTransfer(nftID.Nft(serials[0]), treasury, receiver).
		AddNftTransfer(nftID.Nft(serials[1]), treasury, receiver).
		AddNftTransfer(nftID.Nft(serials[2]), treasury, receiver).
		AddTokenTransfer(tokenID, treasury, -1000000).
		AddTokenTransfer(tokenID, receiver, 1000000).
		FreezeWith(client)
	tokenTransfer, _ := frozenTransfer.Sign(treasuryKey).Execute(client)
	receipt, err = tokenTransfer.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error transferring tokens", err))
	}

	tokenBalance, err := hiero.NewAccountBalanceQuery().SetAccountID(receiver).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error getting balance", err))
	}

	fmt.Println("Fungible token balance for receiver account before reject: ", tokenBalance.Tokens.Get(tokenID))
	fmt.Println("NFT balance for receiver account before reject: ", tokenBalance.Tokens.Get(nftID))

	tokenBalance, err = hiero.NewAccountBalanceQuery().SetAccountID(treasury).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error getting balance", err))
	}
	fmt.Println("Fungible token balance for treasury account before reject: ", tokenBalance.Tokens.Get(tokenID))
	fmt.Println("NFT balance for receiver treasury before reject: ", tokenBalance.Tokens.Get(nftID))
	fmt.Println("-----------------------------------")

	// reject the fungible tokens
	frozenReject, _ := hiero.NewTokenRejectTransaction().
		SetOwnerID(receiver).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	reject, _ := frozenReject.Sign(receiverKey).Execute(client)
	receipt, err = reject.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error rejecting tokens", err))
	}

	// reject the NFTs
	frozenRejectFlow, _ := hiero.NewTokenRejectFlow().
		SetOwnerID(receiver).
		SetNftIDs(nftID.Nft(serials[0]), nftID.Nft(serials[1]), nftID.Nft(serials[2])).
		FreezeWith(client)
	reject, _ = frozenRejectFlow.Sign(receiverKey).Execute(client)
	receipt, err = reject.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error rejecting tokens", err))
	}

	tokenBalance, err = hiero.NewAccountBalanceQuery().SetAccountID(receiver).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error getting balance", err))
	}

	fmt.Println("Fungible token balance for receiver account after reject: ", tokenBalance.Tokens.Get(tokenID))
	fmt.Println("NFT balance for receiver account after reject: ", tokenBalance.Tokens.Get(nftID))

	tokenBalance, err = hiero.NewAccountBalanceQuery().SetAccountID(treasury).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : Error getting balance", err))
	}
	fmt.Println("Fungible token balance for treasury account after reject: ", tokenBalance.Tokens.Get(tokenID))
	fmt.Println("NFT balance for receiver treasury after reject: ", tokenBalance.Tokens.Get(nftID))
}

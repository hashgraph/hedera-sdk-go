package main

import (
	"fmt"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	// Generate 24 word mnemonic
	mnemonic24, err := hedera.GenerateMnemonic24()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating 24 word mnemonic", err))
	}

	// Generate 12 word mnemonic
	mnemonic12, err := hedera.GenerateMnemonic12()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating 12 word mnemonic", err))
	}

	// Given legacy string
	legacyString := "jolly,kidnap,tom,lawn,drunk,chick,optic,lust,mutter,mole,bride,galley,dense,member,sage,neural,widow,decide,curb,aboard,margin,manure"

	// Initializing a legacy mnemonic from legacy string
	mnemonicLegacy, err := hedera.NewMnemonic(strings.Split(legacyString, ","))
	if err != nil {
		panic(fmt.Sprintf("%v : error generating mnemonic from legacy string", err))
	}

	fmt.Printf("mnemonic 24 word = %v\n", mnemonic24)
	fmt.Printf("mnemonic 12 word = %v\n", mnemonic12)
	fmt.Printf("mnemonic legacy = %v\n", mnemonicLegacy)

	// Creating a Private Key from 24 word mnemonic with an optional passphrase
	privateKey24, err := mnemonic24.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		panic(fmt.Sprintf("%v : error converting 24 word mnemonic to PrivateKey", err))
	}

	// Creating a Private Key from 12 word mnemonic with an optional passphrase
	privateKey12, err := mnemonic12.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		panic(fmt.Sprintf("%v : error converting 12 word mnemonic to PrivateKey", err))
	}

	// ToLegacyPrivateKey() doesn't support a passphrase
	privateLegacy, err := mnemonicLegacy.ToLegacyPrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error converting legacy mnemonic to PrivateKey", err))
	}

	// Retrieving the Public Key
	publicKey24 := privateKey24.PublicKey()
	publicKey12 := privateKey12.PublicKey()
	publicLegacy := privateLegacy.PublicKey()

	fmt.Printf("private 24 word = %v\n", privateKey24)
	fmt.Printf("public 24 word = %v\n", publicKey24)

	fmt.Printf("private 12 word = %v\n", privateKey12)
	fmt.Printf("public 12 word = %v\n", publicKey12)

	fmt.Printf("private legacy = %v\n", privateLegacy)
	fmt.Printf("public legacy = %v\n", publicLegacy)
}

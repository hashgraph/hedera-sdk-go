package main

import (
	"fmt"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	mnemonic24, err := hedera.GenerateMnemonic24()
	if err != nil {
		println(err.Error(), ": error generating 24 word mnemonic")
		return
	}

	mnemonic12, err := hedera.GenerateMnemonic12()
	if err != nil {
		println(err.Error(), ": error generating 12 word mnemonic")
		return
	}

	legacyString := "jolly,kidnap,tom,lawn,drunk,chick,optic,lust,mutter,mole,bride,galley,dense,member,sage,neural,widow,decide,curb,aboard,margin,manure"

	mnemonicLegacy, err := hedera.NewMnemonic(strings.Split(legacyString, ","))
	if err != nil {
		println(err.Error(), ": error generating mnemonic from legacy string")
		return
	}

	fmt.Printf("mnemonic 24 word = %v\n", mnemonic24)
	fmt.Printf("mnemonic 12 word = %v\n", mnemonic12)
	fmt.Printf("mnemonic legacy = %v\n", mnemonicLegacy)

	privateKey24, err := mnemonic24.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		println(err.Error(), ": error converting 24 word mnemonic to PrivateKey")
		return
	}
	privateKey12, err := mnemonic12.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		println(err.Error(), ": error converting 12 word mnemonic to PrivateKey")
		return
	}
	privateLegacy, err := mnemonicLegacy.ToLegacyPrivateKey()
	if err != nil {
		println(err.Error(), ": error converting legacy mnemonic to PrivateKey")
		return
	}

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

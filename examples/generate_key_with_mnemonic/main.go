package main

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	mnemonic24, err := hedera.GenerateMnemonic24()
	if err != nil {
		panic(err)
	}
	mnemonic12, err := hedera.GenerateMnemonic12()
	if err != nil {
		panic(err)
	}

	fmt.Printf("mnemonic 24 word = %v\n", mnemonic24)
	fmt.Printf("mnemonic 12 word = %v\n", mnemonic12)

	privateKey24, err := mnemonic24.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		panic(err)
	}
	privateKey12, err := mnemonic12.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		panic(err)
	}

	publicKey24 := privateKey24.PublicKey()
	publicKey12 := privateKey12.PublicKey()

	fmt.Printf("private 24 word = %v\n", privateKey24)
	fmt.Printf("public 24 word = %v\n", publicKey24)

	fmt.Printf("private 12 word = %v\n", privateKey12)
	fmt.Printf("public 12 word = %v\n", publicKey12)
}

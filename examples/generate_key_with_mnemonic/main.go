package main

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	mnemonic, err := hedera.GenerateMnemonic24()
	if err != nil {
		panic(err)
	}

	fmt.Printf("mnemonic = %v\n", mnemonic)

	privateKey, err := mnemonic.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey()

	fmt.Printf("private = %v\n", privateKey)
	fmt.Printf("public = %v\n", publicKey)
}

package main

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	// Generating key
	privateKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	// Retrieve the public key
	publicKey := privateKey.PublicKey()

	fmt.Printf("private = %v\n", privateKey)
	fmt.Printf("public = %v\n", publicKey)
}

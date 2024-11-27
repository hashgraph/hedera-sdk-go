package main

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

func main() {
	// Generating key
	privateKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// Retrieve the public key
	publicKey := privateKey.PublicKey()

	fmt.Printf("private = %v\n", privateKey)
	fmt.Printf("public = %v\n", publicKey)
}

package main

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	privateKey := hedera.NewEd25519PrivateKey()
	fmt.Println(privateKey.PublicKey.KeyData)
}

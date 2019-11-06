package main

import (
	"bytes"
	"fmt"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	publicKey, privateKey := hedera.GenerateKeys()
	publicKeyBytes := hedera.FromBytes(publicKey)
	privateKeyBytes := hedera.FromBytes(privateKey)
	pubK2 := hedera.FromString(publicKeyBytes)
	priK2 := hedera.FromString(privateKeyBytes)
	if bytes.Equal(pubK2, publicKey) {
		fmt.Println("Pub same")
	} else {
		fmt.Println("Pub diff")
	}
	if bytes.Equal(priK2, privateKey) {
		fmt.Println("Pri same")
	} else {
		fmt.Println("Pri diff")
	}
}

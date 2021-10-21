package main

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	client, err = hedera.ClientFromConfigFile("client_config_with_operator_testnet.json")
	if err != nil {
		println(err.Error(), ": error initializing client")
		return
	}

	println("client_config_with_operator_testnet: ")
	println("Nodes: ")
	for address, node := range client.GetNetwork() {
		fmt.Printf("%s : %s\n", address, node.String())
	}

	println("Mirror ledgerID: ")
	for _, mir := range client.GetMirrorNetwork() {
		fmt.Printf("%s\n", mir)
	}

	fmt.Printf("Operator AccountID: %s\n", client.GetOperatorAccountID().String())
	fmt.Printf("Operator Public Key: %s\n", client.GetOperatorPublicKey().String())

	client, err = hedera.ClientFromConfigFile("client_config_simple.json")
	if err != nil {
		println(err.Error(), ": error initializing client")
		return
	}

	println("\nclient_config_simple: ")
	println("Nodes: ")
	for address, node := range client.GetNetwork() {
		fmt.Printf("%s : %s\n", address, node.String())
	}

	println("Mirror ledgerID: ")
	for _, mir := range client.GetMirrorNetwork() {
		fmt.Printf("%s\n", mir)
	}

	key, _ := hedera.PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")

	client.SetOperator(hedera.AccountID{Account: 3}, key)

	fmt.Printf("Operator AccountID: %s\n", client.GetOperatorAccountID().String())
	fmt.Printf("Operator Public Key: %s\n", client.GetOperatorPublicKey().String())

}

package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"

	"github.com/golang/protobuf/proto"
	protobuf "github.com/hashgraph/hedera-sdk-go/proto"
)

func main() {
	client := hedera.ClientForTestnet()

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	// Constructors exist for convenient files
	fileID := hedera.FileIDForAddressBook()
	// fileID := hedera.FileIDForFeeSchedule()
	// fileID := hedera.FileIDForExchangeRate()

	client.SetOperator(operatorAccountID, operatorPrivateKey)

	contents, err := hedera.NewFileContentsQuery().
		SetFileID(fileID).
		Execute(client)

	if err != nil {
		panic(err)
	}

    var book protobuf.NodeAddressBook
    proto.Unmarshal(contents, &book)

	fmt.Printf("contents for file %v :\n", fileID)
    for _, node := range book.NodeAddress {
        fmt.Printf("IpAddress: %v\n", node.IpAddress)
        fmt.Printf("Portno: %v\n", node.Portno)
        fmt.Printf("Memo: %v\n", string(node.Memo))
        fmt.Printf("RSA_PubKey: %v\n", node.RSA_PubKey)
        fmt.Printf("NodeId: %v\n", node.NodeId)
        fmt.Printf("NodeAccountId: %v\n", node.NodeAccountId)
        fmt.Printf("NodeCertHash: %v\n", hex.EncodeToString(node.NodeCertHash))
        fmt.Println()
    }
}

package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

// Convert hex string to byte array
func hexToBytes(hexStr string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return bytes, nil
}

func main() {
	// RLP encoded data in hex format
	encodedHex := "f8cf82012a800085d1385c7bf0830249f094000000000000000000000000000000000000041080b864368b87720000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b6e6577206d657373616765000000000000000000000000000000000000000000c001a02a75313521c309b0cb12f1242370c4dbdaf7b262b171bf2c92f61829769a402da017cade29bbc5ca38ce94cb55e786b172dd3c699d50bb9ff60c9b5011818157db"

	// Convert the hex string to bytes
	encodedBytes, err := hexToBytes(encodedHex)
	if err != nil {
		fmt.Println("Error converting hex to bytes:", err)
		return
	}

	// Create a new RLP item to decode into
	decodedItem := hedera.NewRLPItem(hedera.LIST_TYPE) // Assuming this is a list type based on your data
	if err := decodedItem.Read(encodedBytes); err != nil {
		fmt.Println("Error decoding RLP:", err)
		return
	}

	// Output the decoded values
	fmt.Printf("Decoded RLP:\n")
	for i, item := range decodedItem.GetChildItems() {
		fmt.Printf("Item %d: %x\n", i, item.GetItemValue())
	}
}

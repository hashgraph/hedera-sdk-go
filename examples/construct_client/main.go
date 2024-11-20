package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

func main() {
	// Create client for previewnet
	previewnetClient := hiero.ClientForPreviewnet()
	// Create client for testnet
	testnetClient := hiero.ClientForTestnet()
	// Create client for mainnet
	mainnetClient := hiero.ClientForMainnet()

	println("Client Construction Example.")

	// Creating client from the set HEDERA_NETWORK environment variable
	namedNetworkClient, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client for name", err))
	}

	// Creating account ID of 0.0.3
	id, err := hiero.AccountIDFromString("0.0.3")
	if err != nil {
		panic(fmt.Sprintf("%v : error creating AccountID from string", err))
	}

	// Creating a PrivateKey from a random key string we have
	key, err := hiero.PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	if err != nil {
		panic(fmt.Sprintf("%v : error creating PrivateKey from string", err))
	}

	// Set the operators for each client
	testnetClient.SetOperator(id, key)
	mainnetClient.SetOperator(id, key)
	previewnetClient.SetOperator(id, key)
	namedNetworkClient.SetOperator(id, key)

	// Create the network map to use
	customNetwork := map[string]hiero.AccountID{
		"2.testnet.hedera.com:50211": {Account: 5},
		"3.testnet.hedera.com:50211": {Account: 6},
	}

	// Set network for customClient which uses the above custom network
	customClient := hiero.ClientForNetwork(customNetwork)
	// Setting NetworkName for the CustomClient, is only needed if you need to validate ID checksums
	customClient.SetNetworkName(hiero.NetworkNameTestnet)

	if os.Getenv("CONFIG_FILE") != "" {
		// Creating client from a file specified in environment variable CONFIG_FILE
		configClient, err := hiero.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))
		if err != nil {
			panic(fmt.Sprintf("%v : error creating Client from config file", err))
		}

		// Closing the client from file
		err = configClient.Close()
		if err != nil {
			panic(fmt.Sprintf("%v : error closing configClient", err))
		}
	}

	// Clean up, closing each client
	// Can also do this by using defer in after setting up the client
	err = previewnetClient.Close()
	if err != nil {
		panic(fmt.Sprintf("%v : error closing previewnetClient", err))
	}
	err = testnetClient.Close()
	if err != nil {
		panic(fmt.Sprintf("%v : error closing testnetClient", err))
	}
	err = mainnetClient.Close()
	if err != nil {
		panic(fmt.Sprintf("%v : error closing mainnetClient", err))
	}
	err = namedNetworkClient.Close()
	if err != nil {
		panic(fmt.Sprintf("%v : error closing namedNetworkClient", err))
	}
	err = customClient.Close()
	if err != nil {
		panic(fmt.Sprintf("%v : error closing customClient", err))
	}

	println("Success!")
}

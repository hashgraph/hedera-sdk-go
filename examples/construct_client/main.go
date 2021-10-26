package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
)

func main() {
	previewnetClient := hedera.ClientForPreviewnet()
	testnetClient := hedera.ClientForTestnet()
	mainnetClient := hedera.ClientForMainnet()

	println("Client Construction Example.")

	namedNetworkClient, err := hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client for name")
		return
	}

	id, err := hedera.AccountIDFromString("0.0.3")
	if err != nil {
		println(err.Error(), ": error creating AccountID from string")
		return
	}
	key, err := hedera.PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	if err != nil {
		println(err.Error(), ": error creating PrivateKey from string")
		return
	}
	testnetClient.SetOperator(id, key)
	mainnetClient.SetOperator(id, key)
	previewnetClient.SetOperator(id, key)
	namedNetworkClient.SetOperator(id, key)

	customNetwork := map[string]hedera.AccountID{
		"2.testnet.hedera.com:50211": {Account: 5},
		"3.testnet.hedera.com:50211": {Account: 6},
	}

	customClient := hedera.ClientForNetwork(customNetwork)
	//This requires addressbook/* to be present
	customClient.SetNetworkName(hedera.NetworkNameTestnet)

	if os.Getenv("CONFIG_FILE") != "" {
		configClient, err := hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))
		if err != nil {
			println(err.Error(), ": error creating Client from config file")
			return
		}
		err = configClient.Close()
		if err != nil {
			println(err.Error(), ": error closing configClient")
			return
		}
	}

	err = previewnetClient.Close()
	if err != nil {
		println(err.Error(), ": error closing previewnetClient")
		return
	}
	err = testnetClient.Close()
	if err != nil {
		println(err.Error(), ": error closing testnetClient")
		return
	}
	err = mainnetClient.Close()
	if err != nil {
		println(err.Error(), ": error closing mainnetClient")
		return
	}
	err = namedNetworkClient.Close()
	if err != nil {
		println(err.Error(), ": error closing namedNetworkClient")
		return
	}
	err = customClient.Close()
	if err != nil {
		println(err.Error(), ": error closing customClient")
		return
	}

	println("Success!")
}

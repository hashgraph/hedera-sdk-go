package utils

import (
	"encoding/hex"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func getKeyListFromString(keyStr string) (hedera.Key, error) {
	bytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return hedera.KeyList{}, err
	}

	return hedera.KeyFromBytes(bytes)
}

func GetKeyFromString(keyStr string) (hedera.Key, error) {
	key, err := hedera.PublicKeyFromString(keyStr)
	if err != nil {
		key, err := hedera.PrivateKeyFromStringDer(keyStr)
		if err != nil {
			return getKeyListFromString(keyStr)
		}
		return key, nil
	}
	return key, nil
}

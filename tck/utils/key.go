package utils

import (
	"encoding/hex"
	"errors"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

var ErrFromKeyShouldBeProvided = errors.New("invalid parameters: fromKey should only be provided for ed25519PublicKey, ecdsaSecp256k1PublicKey, or evmAddress types")
var ErrThresholdTypeShouldBeProvided = errors.New("invalid parameters: threshold should only be provided for thresholdKey types")
var ErrKeysShouldBeProvided = errors.New("invalid parameters: keys should only be provided for keyList or thresholdKey types")
var ErrKeylistRequired = errors.New("invalid request: keys list is required for generating a KeyList type")
var ErrThresholdRequired = errors.New("invalid request: threshold is required for generating a ThresholdKey type")

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

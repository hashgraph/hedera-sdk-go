package utils

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"errors"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

var ErrFromKeyShouldBeProvided = errors.New("invalid parameters: fromKey should only be provided for ed25519PublicKey, ecdsaSecp256k1PublicKey, or evmAddress types")
var ErrThresholdTypeShouldBeProvided = errors.New("invalid parameters: threshold should only be provided for thresholdKey types")
var ErrKeysShouldBeProvided = errors.New("invalid parameters: keys should only be provided for keyList or thresholdKey types")
var ErrKeylistRequired = errors.New("invalid request: keys list is required for generating a KeyList type")
var ErrThresholdRequired = errors.New("invalid request: threshold is required for generating a ThresholdKey type")

func getKeyListFromString(keyStr string) (hiero.Key, error) {
	bytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return hiero.KeyList{}, err
	}

	return hiero.KeyFromBytes(bytes)
}

func GetKeyFromString(keyStr string) (hiero.Key, error) {
	key, err := hiero.PublicKeyFromString(keyStr)
	if err != nil {
		key, err := hiero.PrivateKeyFromStringDer(keyStr)
		if err != nil {
			return getKeyListFromString(keyStr)
		}
		return key, nil
	}
	return key, nil
}

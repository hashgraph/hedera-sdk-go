package param

// SPDX-License-Identifier: Apache-2.0

type KeyType string

const (
	ED25519_PRIVATE_KEY         KeyType = "ed25519PrivateKey"
	ED25519_PUBLIC_KEY          KeyType = "ed25519PublicKey"
	ECDSA_SECP256K1_PRIVATE_KEY KeyType = "ecdsaSecp256k1PrivateKey"
	ECDSA_SECP256K1_PUBLIC_KEY  KeyType = "ecdsaSecp256k1PublicKey"
	LIST_KEY                    KeyType = "keyList"
	THRESHOLD_KEY               KeyType = "thresholdKey"
	EVM_ADDRESS_KEY             KeyType = "evmAddress"
)

type KeyParams struct {
	Type      KeyType      `json:"type"`
	FromKey   *string      `json:"fromKey"`
	Threshold *int         `json:"threshold"`
	Keys      *[]KeyParams `json:"keys"`
}

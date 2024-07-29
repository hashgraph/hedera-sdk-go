package methods

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/tck/param"
	"github.com/hashgraph/hedera-sdk-go/tck/response"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

// GenerateKey generates key based on provided key params
func GenerateKey(_ context.Context, params param.KeyParams) (response.GenerateKeyResponse, error) {
	if params.FromKey != "" && params.Type != param.ED25519_PUBLIC_KEY && params.Type != param.ECDSA_SECP256K1_PUBLIC_KEY && params.Type != param.EVM_ADDRESS_KEY {
		return response.GenerateKeyResponse{}, errors.New("invalid parameters: fromKey should only be provided for ed25519PublicKey, ecdsaSecp256k1PublicKey, or evmAddress types")
	}

	if params.Threshold != 0 && params.Type != param.THRESHOLD_KEY {
		return response.GenerateKeyResponse{}, errors.New("invalid parameters: threshold should only be provided for thresholdKey types")
	}

	if params.Keys != nil && params.Type != param.LIST_KEY && params.Type != param.THRESHOLD_KEY {
		return response.GenerateKeyResponse{}, errors.New("invalid parameters: keys should only be provided for keyList or thresholdKey types")
	}

	if (params.Type == param.THRESHOLD_KEY || params.Type == param.LIST_KEY) && params.Keys == nil {
		return response.GenerateKeyResponse{}, errors.New("invalid request: keys list is required for generating a KeyList type")
	}

	if params.Type == param.THRESHOLD_KEY && params.Threshold == 0 {
		return response.GenerateKeyResponse{}, errors.New("invalid request: threshold is required for generating a ThresholdKey type")
	}

	resp := response.GenerateKeyResponse{}
	key, err := processKeyRecursively(params, &resp, false)
	if err != nil {
		return response.GenerateKeyResponse{}, err
	}
	resp.Key = key
	return resp, nil
}

func processKeyRecursively(params param.KeyParams, response *response.GenerateKeyResponse, isList bool) (string, error) {
	switch params.Type {
	case param.ED25519_PRIVATE_KEY, param.ECDSA_SECP256K1_PRIVATE_KEY:
		var privateKey string
		if params.Type == param.ED25519_PRIVATE_KEY {
			pk, _ := hedera.PrivateKeyGenerateEd25519()
			privateKey = pk.StringDer()
		} else {
			pk, _ := hedera.PrivateKeyGenerateEcdsa()
			privateKey = pk.StringDer()
		}
		if isList {
			response.PrivateKeys = append(response.PrivateKeys, privateKey)
		}
		return privateKey, nil

	case param.ED25519_PUBLIC_KEY, param.ECDSA_SECP256K1_PUBLIC_KEY:
		var publicKey string
		if params.FromKey != "" {
			if params.Type == param.ED25519_PUBLIC_KEY {
				pk, _ := hedera.PrivateKeyFromStringEd25519(params.FromKey)
				publicKey = pk.PublicKey().StringDer()
			} else {
				pk, _ := hedera.PrivateKeyFromStringECDSA(params.FromKey)
				publicKey = pk.PublicKey().StringDer()
			}

			return publicKey, nil
		}
		if params.Type == param.ED25519_PUBLIC_KEY {
			pk, _ := hedera.PrivateKeyGenerateEd25519()
			publicKey = pk.PublicKey().StringDer()
		} else {
			pk, _ := hedera.PrivateKeyGenerateEcdsa()
			publicKey = pk.PublicKey().StringDer()
		}
		if isList {
			response.PrivateKeys = append(response.PrivateKeys, publicKey)
		}
		return publicKey, nil

	case param.LIST_KEY, param.THRESHOLD_KEY:
		keyList := hedera.NewKeyList()
		for _, keyParams := range params.Keys {
			keyStr, err := processKeyRecursively(keyParams, response, true)
			if err != nil {
				return "", err
			}
			if strings.Contains(keyStr, "326d") {
				key, err := getKeyListFromString(keyStr)
				if err != nil {
					return "", err
				}
				keyList.Add(key)
			} else {
				key, err := getKeyFromString(keyStr)
				if err != nil {
					return "", err
				}
				keyList.Add(key)
			}
		}
		if params.Type == param.THRESHOLD_KEY {
			keyList.SetThreshold(params.Threshold)
		}

		keyListBytes, err := hedera.KeyToBytes(keyList)
		if err != nil {
			return "", err
		}

		return hex.EncodeToString(keyListBytes), nil

	case param.EVM_ADDRESS_KEY:
		if params.FromKey != "" {
			key, err := getKeyFromString(params.FromKey)
			if err != nil {
				return "", err
			}
			publicKey, ok := key.(hedera.PublicKey)
			if ok {
				return publicKey.ToEthereumAddress(), nil
			}

			privateKey, ok := key.(hedera.PrivateKey)
			if ok {
				return privateKey.PublicKey().ToEthereumAddress(), nil
			}
			return "", errors.New("invalid parameters: fromKey for evmAddress is not ECDSAsecp256k1")
		}
		privateKey, err := hedera.PrivateKeyGenerateEcdsa()
		if err != nil {
			return "", err
		}
		return privateKey.PublicKey().ToEthereumAddress(), nil

	default:
		return "", errors.New("invalid request: key type not recognized")
	}
}

func getKeyListFromString(keyStr string) (hedera.Key, error) {
	bytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return hedera.KeyList{}, err
	}

	return hedera.KeyFromBytes(bytes)
}

func getKeyFromString(keyStr string) (hedera.Key, error) {
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

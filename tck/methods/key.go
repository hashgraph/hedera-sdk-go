package methods

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/hashgraph/hedera-sdk-go/tck/param"
	"github.com/hashgraph/hedera-sdk-go/tck/response"
	"github.com/hashgraph/hedera-sdk-go/tck/utils"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

// GenerateKey generates key based on provided key params
func GenerateKey(_ context.Context, params param.KeyParams) (response.GenerateKeyResponse, error) {
	if params.FromKey != nil && params.Type != param.ED25519_PUBLIC_KEY && params.Type != param.ECDSA_SECP256K1_PUBLIC_KEY && params.Type != param.EVM_ADDRESS_KEY {
		return response.GenerateKeyResponse{}, utils.ErrFromKeyShouldBeProvided
	}

	if params.Threshold != nil && params.Type != param.THRESHOLD_KEY {
		return response.GenerateKeyResponse{}, utils.ErrThresholdTypeShouldBeProvided
	}

	if params.Keys != nil && params.Type != param.LIST_KEY && params.Type != param.THRESHOLD_KEY {
		return response.GenerateKeyResponse{}, utils.ErrKeysShouldBeProvided
	}

	if (params.Type == param.THRESHOLD_KEY || params.Type == param.LIST_KEY) && params.Keys == nil {
		return response.GenerateKeyResponse{}, utils.ErrKeylistRequired
	}

	if params.Type == param.THRESHOLD_KEY && params.Threshold == nil {
		return response.GenerateKeyResponse{}, utils.ErrThresholdRequired
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
		var publicKey, privateKey string

		setKeysFromKey := func(fromKey string, isEd25519 bool) {
			var pk hedera.PrivateKey
			if isEd25519 {
				pk, _ = hedera.PrivateKeyFromStringEd25519(fromKey)
			} else {
				pk, _ = hedera.PrivateKeyFromStringECDSA(fromKey)
			}
			privateKey = pk.StringDer()
			publicKey = pk.PublicKey().StringDer()
		}

		generateKeys := func(isEd25519 bool) {
			var pk hedera.PrivateKey
			if isEd25519 {
				pk, _ = hedera.PrivateKeyGenerateEd25519()
			} else {
				pk, _ = hedera.PrivateKeyGenerateEcdsa()
			}
			privateKey = pk.StringDer()
			publicKey = pk.PublicKey().StringDer()
		}

		isEd25519 := params.Type == param.ED25519_PUBLIC_KEY

		if params.FromKey != nil {
			setKeysFromKey(*params.FromKey, isEd25519)
		} else {
			generateKeys(isEd25519)
		}

		if isList {
			response.PrivateKeys = append(response.PrivateKeys, privateKey)
		}

		return publicKey, nil

	case param.LIST_KEY, param.THRESHOLD_KEY:
		keyList := hedera.NewKeyList()
		for _, keyParams := range *params.Keys {
			keyStr, err := processKeyRecursively(keyParams, response, true)
			if err != nil {
				return "", err
			}
			key, err := utils.GetKeyFromString(keyStr)
			if err != nil {
				return "", err
			}
			keyList.Add(key)
		}
		if params.Type == param.THRESHOLD_KEY {
			keyList.SetThreshold(*params.Threshold)
		}

		keyListBytes, err := hedera.KeyToBytes(keyList)
		if err != nil {
			return "", err
		}

		return hex.EncodeToString(keyListBytes), nil

	case param.EVM_ADDRESS_KEY:
		if params.FromKey != nil {
			key, err := utils.GetKeyFromString(*params.FromKey)
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

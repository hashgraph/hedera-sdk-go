package methods

import (
	"context"

	"github.com/hashgraph/hedera-sdk-go/tck/response"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type KeyParams struct {
	PrivateKey string `json:"privateKey"`
}

// GeneratePublicKey generates public key based on provided private key
func GeneratePublicKey(_ context.Context, params KeyParams) (string, error) {
	key, err := hedera.PrivateKeyFromString(params.PrivateKey)
	if err != nil {
		return "", response.HederaError.WithData(err.Error())
	}
	return key.PublicKey().String(), nil
}

// GeneratePrivateKey generate new private key and return it.
func GeneratePrivateKey(_ context.Context, empty any) (string, error) {
	pk, err := hedera.PrivateKeyGenerateEd25519()

	if err != nil {
		return "", response.HederaError.WithData(err.Error())
	}
	return pk.String(), nil
}

package methods

import (
	"context"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

type KeyParams struct {
	PrivateKey string `json:"privateKey"`
}

// GeneratePublicKey generates public key based on provided private key
func GeneratePublicKey(_ context.Context, params KeyParams) (string, error) {
	trimmedPk := strings.TrimSpace(params.PrivateKey)
	key, err := hedera.PrivateKeyFromString(trimmedPk)
	if err != nil {
		return "", err
	}
	return key.PublicKey().String(), nil
}

// GeneratePrivateKey generate new private key and return it.
func GeneratePrivateKey(_ context.Context, empty any) (string, error) {
	pk, err := hedera.PrivateKeyGenerateEd25519()

	if err != nil {
		return "", err
	}
	return pk.String(), nil
}

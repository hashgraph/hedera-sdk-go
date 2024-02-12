package methods

import (
	"context"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type KeyService struct {
}

type KeyParams struct {
	PrivateKey string `json:"privateKey"`
}

// GeneratePublicKey generates public key based on provided private key
func (_ KeyService) GeneratePublicKey(_ context.Context, params KeyParams) string {
	key, err := hedera.PrivateKeyFromString(params.PrivateKey)
	if err != nil {
		return ""
	}
	return key.PublicKey().String()
}

// GeneratePrivateKey generate new private key and return it.
// NOTE: params object is used only to maintain the communication with jRPC client, which is providing empty object when calling this func
func (s KeyService) GeneratePrivateKey(_ context.Context, params KeyParams) string {
	pk, err := hedera.PrivateKeyGenerateEd25519()

	if err != nil {
		return ""
	}
	return pk.String()
}

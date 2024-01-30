package methods

import (
	"context"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type KeyService struct {
}

type KeyParams struct {
	PrivateKey string `json:"privateKey"`
}

func (_ KeyService) GeneratePublicKey(_ context.Context, params KeyParams) string {
	fmt.Println("Key: ", params.PrivateKey)
	key, err := hedera.PrivateKeyFromString(params.PrivateKey)
	if err != nil {
		return ""
	}
	fmt.Println(key)
	return key.PublicKey().String()
}

func (s KeyService) GeneratePrivateKey(_ context.Context) string {
	pk, err := hedera.PrivateKeyGenerateEd25519()

	if err != nil {
		return ""
	}
	return pk.String()
}

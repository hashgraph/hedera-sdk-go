//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/assert"
)

func TestContractNonceInfoFromProtobuf(t *testing.T) {
	contractID := &ContractID{Shard: 0, Realm: 0, Contract: 123}
	nonce := int64(456)
	protobuf := &services.ContractNonceInfo{
		ContractId: contractID._ToProtobuf(),
		Nonce:      nonce,
	}

	result := _ContractNonceInfoFromProtobuf(protobuf)

	assert.Equal(t, contractID, result.ContractID)
	assert.Equal(t, nonce, result.Nonce)
}

func TestContractNonceInfoFromProtobuf_NilInput(t *testing.T) {
	result := _ContractNonceInfoFromProtobuf(nil)

	assert.Nil(t, result)
}
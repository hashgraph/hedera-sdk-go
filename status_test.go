package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestStatusFromProtoToString tests pulling all codes from the proto generated code,
// converting it to the sdk enum, and calling String()
//
// Ideally this will catch any changes to response codes when the protobufs get updated
func TestStatusFromProtoToString(t *testing.T) {
	for _, code := range proto.ResponseCodeEnum_value {
		status := Status(code)
		assert.NotPanics(t, func() { _ = status.String() })
	}
}

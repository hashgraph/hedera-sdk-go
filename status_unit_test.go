//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
)

// TestStatusFromProtoToString tests pulling all codes from the proto generated code,
// converting it to the sdk enum, and calling String()
//
// Ideally this will catch any changes to _Response codes when the protobufs get updated
func TestUnitStatusFromProtoToString(t *testing.T) {
	t.Parallel()

	for _, code := range services.ResponseCodeEnum_value {
		status := Status(code)
		assert.NotPanics(t, func() { _ = status.String() })
	}
}

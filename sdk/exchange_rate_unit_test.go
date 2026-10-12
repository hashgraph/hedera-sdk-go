//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitExchangeRateGetters(t *testing.T) {
	t.Parallel()

	rate := _ExchangeRateFromProtobuf(&services.ExchangeRate{
		HbarEquiv:      30000,
		CentEquiv:      150000,
		ExpirationTime: &services.TimestampSeconds{Seconds: 1700000000},
	})

	require.Equal(t, int32(30000), rate.Hbars)
	require.Equal(t, int32(150000), rate.GetCents())
	require.NotNil(t, rate.GetExpirationTime())
	require.True(t, rate.GetExpirationTime().Equal(time.Unix(1700000000, 0)))
}

func TestUnitExchangeRateGetExpirationTimeNil(t *testing.T) {
	t.Parallel()

	rate := _ExchangeRateFromProtobuf(&services.ExchangeRate{HbarEquiv: 1, CentEquiv: 12})

	require.Nil(t, rate.GetExpirationTime())
}

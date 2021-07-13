package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type ExchangeRate struct {
	Hbars          int32
	cents          int32
	expirationTime *services.TimestampSeconds
}

func newExchangeRate(hbars int32, cents int32, expirationTime int64) *ExchangeRate {
	exchange := ExchangeRate{
		Hbars:          hbars,
		cents:          cents,
		expirationTime: &services.TimestampSeconds{Seconds: expirationTime},
	}

	return &exchange
}

func exchangeRateFromProtobuf(protoExchange *services.ExchangeRate) ExchangeRate {
	if protoExchange == nil {
		return ExchangeRate{}
	}
	var expirationTime *services.TimestampSeconds
	if protoExchange.ExpirationTime != nil {
		expirationTime = protoExchange.ExpirationTime
	}

	return ExchangeRate{
		protoExchange.HbarEquiv,
		protoExchange.CentEquiv,
		expirationTime,
	}
}

func (exchange *ExchangeRate) toProtobuf() *services.ExchangeRate {
	return &services.ExchangeRate{
		HbarEquiv:      exchange.Hbars,
		CentEquiv:      exchange.cents,
		ExpirationTime: exchange.expirationTime,
	}
}

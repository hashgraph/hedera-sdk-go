package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ExchangeRate struct {
	Hbars int32
	cents int32
	expirationTime *proto.TimestampSeconds
}

func NewExchangeRate(hbars int32, cents int32, expirationTime int64) *ExchangeRate {
	exchange := ExchangeRate{
		Hbars: hbars,
		cents: cents,
		expirationTime: &proto.TimestampSeconds{Seconds: expirationTime},
	}

	return &exchange
}

func exchangeRateFromProtobuf(protoExchange *proto.ExchangeRate) ExchangeRate {
	return	ExchangeRate{
		protoExchange.HbarEquiv,
		protoExchange.CentEquiv,
		protoExchange.ExpirationTime,
	}
}

func (exchange *ExchangeRate) toProtobuf() *proto.ExchangeRate {
	return &proto.ExchangeRate{
		HbarEquiv: exchange.Hbars,
		CentEquiv: exchange.cents,
		ExpirationTime: exchange.expirationTime,
	}
}

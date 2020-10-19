package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ExchangeRate struct {
	Hbars          int32
	cents          int32
	expirationTime *proto.TimestampSeconds
}

func newExchangeRate(hbars int32, cents int32, expirationTime int64) *ExchangeRate {
	exchange := ExchangeRate{
		Hbars:          hbars,
		cents:          cents,
		expirationTime: &proto.TimestampSeconds{Seconds: expirationTime},
	}

	return &exchange
}

func exchangeRateFromProtobuf(protoExchange *proto.ExchangeRate) ExchangeRate {
	var expirationTime *proto.TimestampSeconds
	if protoExchange.ExpirationTime != nil{
		expirationTime = protoExchange.ExpirationTime
	}

	return ExchangeRate{
		protoExchange.HbarEquiv,
		protoExchange.CentEquiv,
		expirationTime,
	}
}

func (exchange *ExchangeRate) toProtobuf() *proto.ExchangeRate {
	return &proto.ExchangeRate{
		HbarEquiv:      exchange.Hbars,
		CentEquiv:      exchange.cents,
		ExpirationTime: exchange.expirationTime,
	}
}

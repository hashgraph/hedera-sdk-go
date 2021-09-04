package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ExchangeRate struct {
	Hbars          int32
	cents          int32
	expirationTime *proto.TimestampSeconds
}

func _ExchangeRateFromProtobuf(protoExchange *proto.ExchangeRate) ExchangeRate {
	if protoExchange == nil {
		return ExchangeRate{}
	}
	var expirationTime *proto.TimestampSeconds
	if protoExchange.ExpirationTime != nil {
		expirationTime = protoExchange.ExpirationTime
	}

	return ExchangeRate{
		protoExchange.HbarEquiv,
		protoExchange.CentEquiv,
		expirationTime,
	}
}

func (exchange *ExchangeRate) _ToProtobuf() *proto.ExchangeRate {
	return &proto.ExchangeRate{
		HbarEquiv:      exchange.Hbars,
		CentEquiv:      exchange.cents,
		ExpirationTime: exchange.expirationTime,
	}
}

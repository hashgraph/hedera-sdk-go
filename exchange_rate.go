package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
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

func (exchange *ExchangeRate) ToBytes() []byte {
	data, err := protobuf.Marshal(exchange._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func ExchangeRateFromBytes(data []byte) (ExchangeRate, error) {
	if data == nil {
		return ExchangeRate{}, errByteArrayNull
	}
	pb := proto.ExchangeRate{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ExchangeRate{}, err
	}

	exchangeRate := _ExchangeRateFromProtobuf(&pb)
	if err != nil {
		return ExchangeRate{}, err
	}

	return exchangeRate, nil
}

func (exchange *ExchangeRate) String() string {
	return fmt.Sprintf("Hbars: %d to Cents: %d, expires: %s", exchange.Hbars, exchange.cents, exchange.expirationTime.String())
}

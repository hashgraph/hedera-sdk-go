package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type ExchangeRate struct {
	Hbars          int32
	cents          int32
	expirationTime *services.TimestampSeconds
}

func _ExchangeRateFromProtobuf(protoExchange *services.ExchangeRate) ExchangeRate {
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

func (exchange *ExchangeRate) _ToProtobuf() *services.ExchangeRate {
	return &services.ExchangeRate{
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
	pb := services.ExchangeRate{}
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

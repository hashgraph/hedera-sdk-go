package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type AccountBalance struct {
	Hbars Hbar

	// Deprecated: Use `AccountBalance.Tokens` instead
	Token map[TokenID]uint64

	// Deprecated
	Tokens TokenBalanceMap
	// Deprecated
	TokenDecimals TokenDecimalMap
}

func _AccountBalanceFromProtobuf(pb *services.CryptoGetAccountBalanceResponse) AccountBalance { //nolint
	if pb == nil {
		return AccountBalance{}
	}
	var tokens map[TokenID]uint64
	if pb.TokenBalances != nil { //nolint
		tokens = make(map[TokenID]uint64, len(pb.TokenBalances)) //nolint
		for _, token := range pb.TokenBalances {                 //nolint
			if t := _TokenIDFromProtobuf(token.TokenId); t != nil {
				tokens[*t] = token.Balance
			}
		}
	}
	return AccountBalance{
		Hbars:         HbarFromTinybar(int64(pb.Balance)),
		Token:         tokens,
		Tokens:        _TokenBalanceMapFromProtobuf(pb.TokenBalances), //nolint
		TokenDecimals: _TokenDecimalMapFromProtobuf(pb.TokenBalances), //nolint
	}
}

func (balance *AccountBalance) _ToProtobuf() *services.CryptoGetAccountBalanceResponse { //nolint
	return &services.CryptoGetAccountBalanceResponse{
		Balance: uint64(balance.Hbars.AsTinybar()),
	}
}

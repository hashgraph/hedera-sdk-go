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

func _AccountBalanceFromProtobuf(pb *services.CryptoGetAccountBalanceResponse) AccountBalance {
	if pb == nil {
		return AccountBalance{}
	}

	return AccountBalance{
		Hbars: HbarFromTinybar(int64(pb.Balance)),
	}
}

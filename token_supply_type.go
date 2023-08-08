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

import "fmt"

type TokenSupplyType int32

const (
	TokenSupplyTypeInfinite TokenSupplyType = 0
	TokenSupplyTypeFinite   TokenSupplyType = 1
)

// String returns a string representation of the TokenSupplyType
func (tokenSupplyType TokenSupplyType) String() string {
	switch tokenSupplyType {
	case TokenSupplyTypeInfinite:
		return "TOKEN_SUPPLY_TYPE_INFINITE"
	case TokenSupplyTypeFinite:
		return "TOKEN_SUPPLY_TYPE_FINITE"
	}

	panic(fmt.Sprintf("unreachable: TokenType.String() switch statement is non-exhaustive. Status: %v", uint32(tokenSupplyType)))
}

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

type TokenType uint32

const (
	TokenTypeFungibleCommon    TokenType = 0
	TokenTypeNonFungibleUnique TokenType = 1
)

// String returns a string representation of the TokenType
func (tokenType TokenType) String() string {
	switch tokenType {
	case TokenTypeFungibleCommon:
		return "TOKEN_TYPE_FUNGIBLE_COMMON"
	case TokenTypeNonFungibleUnique:
		return "TOKEN_TYPE_NON_FUNGIBLE_UNIQUE"
	}

	panic(fmt.Sprintf("unreachable: TokenType.String() switch statement is non-exhaustive. Status: %v", uint32(tokenType)))
}

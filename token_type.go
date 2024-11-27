package hiero

// SPDX-License-Identifier: Apache-2.0

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

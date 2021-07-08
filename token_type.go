package hedera

import "fmt"

type TokenType uint32

const (
	TokenTypeFungibleCommon    TokenType = 0
	TokenTypeNonFungibleUnique TokenType = 1
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case TokenTypeFungibleCommon:
		return "TOKEN_TYPE_FUNGIBLE_COMMON"
	case TokenTypeNonFungibleUnique:
		return "TOKEN_TYPE_NON_FUNGIBLE_UNIQUE"
	}

	panic(fmt.Sprintf("unreacahble: TokenType.String() switch statement is non-exhaustive. Status: %v", uint32(tokenType)))
}

package hedera

import "fmt"

type TokenSupplyType int32

const (
	TokenSupplyTypeInfinite TokenSupplyType = 0
	TokenSupplyTypeFinite   TokenSupplyType = 1
)

func (tokenSupplyType TokenSupplyType) String() string {
	switch tokenSupplyType {
	case TokenSupplyTypeInfinite:
		return "TOKEN_SUPPLY_TYPE_INFINITE"
	case TokenSupplyTypeFinite:
		return "TOKEN_SUPPLY_TYPE_FINITE"
	}

	panic(fmt.Sprintf("unreacahble: TokenType.String() switch statement is non-exhaustive. Status: %v", uint32(tokenSupplyType)))
}

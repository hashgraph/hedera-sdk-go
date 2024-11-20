package hiero

// SPDX-License-Identifier: Apache-2.0

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

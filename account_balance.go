package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type AccountBalance struct {
	Hbars Hbar
	// Deprecated: Use `AccountBalance.Tokens` instead
	Token         map[TokenID]uint64
	Tokens        TokenBalanceMap
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

package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type AccountBalance struct {
	Hbars Hbar

	// Deprecated: Use `AccountBalance.Tokens` instead
	Token map[TokenID]uint64

	Tokens        TokenBalanceMap
	TokenDecimals TokenDecimalMap
}

func _AccountBalanceFromProtobuf(pb *services.CryptoGetAccountBalanceResponse) AccountBalance {
	if pb == nil {
		return AccountBalance{}
	}
	var tokens map[TokenID]uint64
	if pb.TokenBalances != nil {
		tokens = make(map[TokenID]uint64, len(pb.TokenBalances))
		for _, token := range pb.TokenBalances {
			if t := _TokenIDFromProtobuf(token.TokenId); t != nil {
				tokens[*t] = token.Balance
			}
		}
	}

	return AccountBalance{
		Hbars:         HbarFromTinybar(int64(pb.Balance)),
		Token:         tokens,
		Tokens:        _TokenBalanceMapFromProtobuf(pb.TokenBalances),
		TokenDecimals: _TokenDecimalMapFromProtobuf(pb.TokenBalances),
	}
}

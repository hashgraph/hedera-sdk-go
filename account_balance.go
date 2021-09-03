package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type AccountBalance struct {
	Hbars Hbar

	// Deprecated: Use `AccountBalance.Tokens` instead
	Token map[TokenID]uint64

	Tokens        TokenBalanceMap
	TokenDecimals TokenDecimalMap
}

func accountBalanceFromProtobuf(pb *proto.CryptoGetAccountBalanceResponse) AccountBalance {
	if pb == nil {
		return AccountBalance{}
	}
	var tokens map[TokenID]uint64
	if pb.TokenBalances != nil {
		tokens = make(map[TokenID]uint64, len(pb.TokenBalances))
		for _, token := range pb.TokenBalances {
			if t := tokenIDFromProtobuf(token.TokenId); t != nil {
				tokens[*t] = token.Balance
			}
		}
	}

	return AccountBalance{
		Hbars:         HbarFromTinybar(int64(pb.Balance)),
		Token:         tokens,
		Tokens:        tokenBalanceMapFromProtobuf(pb.TokenBalances),
		TokenDecimals: tokenDecimalMapFromProtobuf(pb.TokenBalances),
	}
}

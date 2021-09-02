package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ProxyStaker struct {
	AccountID AccountID
	Amount    Hbar
}

func (staker *ProxyStaker) toProtobuf() *proto.ProxyStaker { // nolint
	return &proto.ProxyStaker{
		AccountID: staker.AccountID.toProtobuf(),
		Amount:    staker.Amount.tinybar,
	}
}

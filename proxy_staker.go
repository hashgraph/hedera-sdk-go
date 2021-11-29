package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type ProxyStaker struct {
	AccountID AccountID
	Amount    Hbar
}

func (staker *ProxyStaker) _ToProtobuf() *services.ProxyStaker { // nolint
	return &services.ProxyStaker{
		AccountID: staker.AccountID._ToProtobuf(),
		Amount:    staker.Amount.tinybar,
	}
}

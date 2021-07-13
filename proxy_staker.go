package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type ProxyStaker struct {
	AccountID AccountID
	amount    Hbar
}

func newProxyStaker(accountId AccountID, amount int64) ProxyStaker {
	return ProxyStaker{
		AccountID: accountId,
		amount:    HbarFromTinybar(amount),
	}
}

func fromProtobuf(staker *services.ProxyStaker, networkName *NetworkName) ProxyStaker {
	return ProxyStaker{
		AccountID: accountIDFromProtobuf(staker.AccountID, networkName),
		amount:    HbarFromTinybar(staker.Amount),
	}
}

func (staker *ProxyStaker) toProtobuf() *services.ProxyStaker {
	return &services.ProxyStaker{
		AccountID: staker.AccountID.toProtobuf(),
		Amount:    staker.amount.tinybar,
	}
}

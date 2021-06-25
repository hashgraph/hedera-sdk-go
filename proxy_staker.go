package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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

func fromProtobuf(staker *proto.ProxyStaker, networkName *NetworkName) ProxyStaker {
	return ProxyStaker{
		AccountID: accountIDFromProtobuf(staker.AccountID, networkName),
		amount:    HbarFromTinybar(staker.Amount),
	}
}

func (staker *ProxyStaker) toProtobuf() *proto.ProxyStaker {
	return &proto.ProxyStaker{
		AccountID: staker.AccountID.toProtobuf(),
		Amount:    staker.amount.tinybar,
	}
}

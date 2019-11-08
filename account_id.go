package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/hedera_proto"
)

type AccountID struct {
	Shard   uint64
	Realm   uint64
	Account uint64
}

func (accountID AccountID) String() string {
	return fmt.Sprintf("%d.%d.%d", accountID.Shard, accountID.Realm, accountID.Account)
}

func (accountID AccountID) proto() *hedera_proto.AccountID {
	return &hedera_proto.AccountID{
		ShardNum:   int64(accountID.Shard),
		RealmNum:   int64(accountID.Realm),
		AccountNum: int64(accountID.Account),
	}
}

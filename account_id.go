package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/hedera_proto"
)

type AccountId struct {
	inner *hedera_proto.AccountID
}

func NewAccountID(shard, realm, account int64) AccountId {
	return AccountId{
		inner: &hedera_proto.AccountID{
			ShardNum:             shard,
			RealmNum:             realm,
			AccountNum:           account,
		},
	}
}

func (accountId *AccountId) Shard() int64 {
	return accountId.inner.ShardNum
}

func (accountId *AccountId) Realm() int64 {
	return accountId.inner.RealmNum
}

func (accountId *AccountId) Account() int64 {
	return accountId.inner.AccountNum
}

func (accountId *AccountId) String() string {
	return fmt.Sprintf("%d.%d.%d", accountId.inner.ShardNum, accountId.inner.RealmNum, accountId.inner.AccountNum)
}

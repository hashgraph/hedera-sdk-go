package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/hedera_proto"
)

type AccountId struct {
	inner *hedera_proto.AccountID
}

func NewAccountID(shard, realm, account uint64) AccountId {
	return AccountId{
		inner: &hedera_proto.AccountID{
			ShardNum:             int64(shard),
			RealmNum:             int64(realm),
			AccountNum:           int64(account),
		},
	}
}

func (accountId *AccountId) Shard() uint64 {
	return uint64(accountId.inner.ShardNum)
}

func (accountId *AccountId) Realm() uint64 {
	return uint64(accountId.inner.RealmNum)
}

func (accountId *AccountId) Account() uint64 {
	return uint64(accountId.inner.AccountNum)
}

func (accountId *AccountId) String() string {
	return fmt.Sprintf("%d.%d.%d", accountId.inner.ShardNum, accountId.inner.RealmNum, accountId.inner.AccountNum)
}

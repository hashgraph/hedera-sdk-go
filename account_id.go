package hedera

import (
	"fmt"
)

type AccountID struct {
	Shard   uint64
	Realm   uint64
	Account uint64
}

func (accountID *AccountID) String() string {
	return fmt.Sprintf("%d.%d.%d", accountID.Shard, accountID.Realm, accountID.Account)
}

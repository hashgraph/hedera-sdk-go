package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type TransactionID struct {
	AccountID  AccountID
	ValidStart time.Time
}

func generateTransactionID(accountID AccountID) TransactionID {
	// TODO: Less 10s
	now := time.Now()

	return TransactionID{accountID, now}
}

func (id TransactionID) Receipt(client *Client) (TransactionReceipt, error) {
	return NewTransactionReceiptQuery().
		SetTransactionID(id).
		Execute(client)
}

// TODO: #Record

func (id TransactionID) String() string {
	seconds := id.ValidStart.Unix()
	nanos := int32(id.ValidStart.UnixNano() - (id.ValidStart.Unix() * 1e+9))

	return fmt.Sprintf("%v@%v.%v", id.AccountID, seconds, nanos)
}

func (id TransactionID) toProto() *proto.TransactionID {
	return &proto.TransactionID{
		TransactionValidStart: &proto.Timestamp{
			Seconds: id.ValidStart.Unix(),
			Nanos:   int32(id.ValidStart.UnixNano() - (id.ValidStart.Unix() * 1e+9)),
		},
		AccountID: &proto.AccountID{
			ShardNum:   int64(id.AccountID.Shard),
			RealmNum:   int64(id.AccountID.Realm),
			AccountNum: int64(id.AccountID.Account),
		},
	}
}

func transactionIDFromProto(pb *proto.TransactionID) TransactionID {
	validStart := time.Unix(pb.TransactionValidStart.Seconds, int64(pb.TransactionValidStart.Nanos))
	accountID := accountIDFromProto(pb.AccountID)

	return TransactionID{accountID, validStart}
}

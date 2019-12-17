package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type FileCreateTransaction struct {
	TransactionBuilder
	pb *proto.FileCreateTransactionBody
}

func NewFileCreateTransaction() FileCreateTransaction {
	pb := &proto.FileCreateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileCreate{pb}

	builder := FileCreateTransaction{inner, pb}
	builder.SetExpirationTime(time.Now().Add(7890000 * time.Second))

	return builder
}

func (builder FileCreateTransaction) AddKey(publicKey Ed25519PublicKey) FileCreateTransaction {
	var keylist *proto.KeyList
	if builder.pb.Keys != nil {
		keylist = builder.pb.Keys
	} else {
		keylist = &proto.KeyList{}
	}

	var keyarray []*proto.Key
	if keylist.Keys != nil {
		keyarray = keylist.GetKeys()
	} else {
		keyarray = []*proto.Key{}
	}

	keylist.Keys = append(keyarray, publicKey.toProto())
	builder.pb.Keys = keylist

	return builder
}

func (builder FileCreateTransaction) SetExpirationTime(expiration time.Time) FileCreateTransaction {
	builder.pb.ExpirationTime = timeToProto(expiration)
	return builder
}

func (builder FileCreateTransaction) SetContents(contents []byte) FileCreateTransaction {
	builder.pb.Contents = contents
	return builder
}

func (builder FileCreateTransaction) Build(client *Client) Transaction {
	// If a shard/realm is not set, it is inferred from the Operator on the Client

	if builder.pb.ShardID == nil {
		builder.pb.ShardID = &proto.ShardID{
			ShardNum: int64(client.operator.accountID.Shard),
		}
	}

	if builder.pb.RealmID == nil {
		builder.pb.RealmID = &proto.RealmID{
			ShardNum: int64(client.operator.accountID.Shard),
			RealmNum: int64(client.operator.accountID.Realm),
		}
	}

	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder FileCreateTransaction) SetMaxTransactionFee(maxTransactionFee uint64) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder FileCreateTransaction) SetMemo(memo string) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetMemo(memo), builder.pb}
}

func (builder FileCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder FileCreateTransaction) SetTransactionID(transactionID TransactionID) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder FileCreateTransaction) SetNodeAccountID(nodeAccountID AccountID) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

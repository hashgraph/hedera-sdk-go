package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type FileUpdateTransaction struct {
	TransactionBuilder
	pb *proto.FileUpdateTransactionBody
}

func NewFileUpdateTransaction() FileUpdateTransaction {
	pb := &proto.FileUpdateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileUpdate{FileUpdate: pb}

	builder := FileUpdateTransaction{inner, pb}
	builder.SetExpirationTime(time.Now().Add(7890000 * time.Second))

	return builder
}

func (builder FileUpdateTransaction) SetFileID(id FileID) FileUpdateTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder FileUpdateTransaction) AddKey(publicKey PublicKey) FileUpdateTransaction {
	var keyList *proto.KeyList
	if builder.pb.Keys != nil {
		keyList = builder.pb.Keys
	} else {
		keyList = &proto.KeyList{}
	}

	var keyarray []*proto.Key
	if keyList.Keys != nil {
		keyarray = keyList.GetKeys()
	} else {
		keyarray = []*proto.Key{}
	}

	keyList.Keys = append(keyarray, publicKey.toProto())
	builder.pb.Keys = keyList

	return builder
}

func (builder FileUpdateTransaction) SetExpirationTime(expiration time.Time) FileUpdateTransaction {
	builder.pb.ExpirationTime = timeToProto(expiration)
	return builder
}

func (builder FileUpdateTransaction) SetContents(contents []byte) FileUpdateTransaction {
	builder.pb.Contents = contents
	return builder
}

func (builder FileUpdateTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder FileUpdateTransaction) SetMaxTransactionFee(maxTransactionFee uint64) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder FileUpdateTransaction) SetTransactionMemo(memo string) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder FileUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder FileUpdateTransaction) SetTransactionID(transactionID TransactionID) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder FileUpdateTransaction) SetNodeAccountID(nodeAccountID AccountID) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

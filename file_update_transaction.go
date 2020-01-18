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
	pb := &proto.FileUpdateTransactionBody{Keys: &proto.KeyList{Keys: []*proto.Key{}}}

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
	builder.pb.Keys.Keys = append(builder.pb.Keys.Keys, publicKey.toProto())
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

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder FileUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileUpdateTransaction {
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

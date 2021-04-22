package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
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

	return builder
}

func fileUpdateTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) FileUpdateTransaction {
	return FileUpdateTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetFileUpdate(),
	}
}

func (builder FileUpdateTransaction) SetFileID(id FileID) FileUpdateTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder FileUpdateTransaction) AddKey(publicKey PublicKey) FileUpdateTransaction {
	if builder.pb.Keys == nil {
		builder.pb.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}

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

func (builder FileUpdateTransaction) Build(client *Client) (Transaction, error) {
	return builder.TransactionBuilder.Build(client)
}

func (builder FileUpdateTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *FileUpdateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_FileUpdate{
			FileUpdate: &proto.FileUpdateTransactionBody{
				FileID:         builder.pb.GetFileID(),
				ExpirationTime: builder.pb.GetExpirationTime(),
				Keys:           builder.pb.GetKeys(),
				Contents:       builder.pb.GetContents(),
				Memo:           builder.pb.GetMemo(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder FileUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder FileUpdateTransaction) SetTransactionMemo(memo string) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder FileUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder FileUpdateTransaction) SetTransactionID(transactionID TransactionID) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder FileUpdateTransaction) SetNodeAccountID(nodeAccountID AccountID) FileUpdateTransaction {
	return FileUpdateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

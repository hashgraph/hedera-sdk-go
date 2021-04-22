package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type SystemDeleteTransaction struct {
	TransactionBuilder
	pb *proto.SystemDeleteTransactionBody
}

func NewSystemDeleteTransaction() SystemDeleteTransaction {
	pb := &proto.SystemDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_SystemDelete{SystemDelete: pb}

	builder := SystemDeleteTransaction{inner, pb}

	return builder
}

func systemDeleteTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) SystemDeleteTransaction {
	return SystemDeleteTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetSystemDelete(),
	}
}

func (builder SystemDeleteTransaction) SetExpirationTime(expiration time.Time) SystemDeleteTransaction {
	builder.pb.ExpirationTime = &proto.TimestampSeconds{
		Seconds: expiration.Unix(),
	}
	return builder
}

func (builder SystemDeleteTransaction) SetContractID(ID ContractID) SystemDeleteTransaction {
	builder.pb.Id = &proto.SystemDeleteTransactionBody_ContractID{ContractID: ID.toProto()}
	return builder
}

func (builder SystemDeleteTransaction) SetFileID(ID FileID) SystemDeleteTransaction {
	builder.pb.Id = &proto.SystemDeleteTransactionBody_FileID{FileID: ID.toProto()}
	return builder
}

func (builder SystemDeleteTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *SystemDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_SystemDelete{
			SystemDelete: &proto.SystemDeleteTransactionBody{
				Id:             nil,
				ExpirationTime: builder.pb.GetExpirationTime(),
			},
		},
	}

	switch builder.pb.GetId().(type) {
	case *proto.SystemDeleteTransactionBody_ContractID:
		body.GetSystemUndelete().Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: builder.pb.GetContractID()}
	case *proto.SystemDeleteTransactionBody_FileID:
		body.GetSystemUndelete().Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: builder.pb.GetFileID()}
	}

	return body, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder SystemDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder SystemDeleteTransaction) SetTransactionMemo(memo string) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder SystemDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder SystemDeleteTransaction) SetTransactionID(transactionID TransactionID) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder SystemDeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

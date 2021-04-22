package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ScheduleDeleteTransaction struct {
	TransactionBuilder
	pb *proto.ScheduleDeleteTransactionBody
}

func NewScheduleDeleteTransaction() ScheduleDeleteTransaction {
	pb := &proto.ScheduleDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ScheduleDelete{ScheduleDelete: pb}

	builder := ScheduleDeleteTransaction{inner, pb}

	return builder
}

func scheduleDeleteTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) ScheduleDeleteTransaction {
	return ScheduleDeleteTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetScheduleDelete(),
	}
}

func (builder ScheduleDeleteTransaction) SetScheduleID(scheduleID ScheduleID) ScheduleDeleteTransaction {
	builder.pb.ScheduleID = scheduleID.toProto()
	return builder
}

func (builder ScheduleDeleteTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder ScheduleDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_ScheduleDelete{
			ScheduleDelete: &proto.ScheduleDeleteTransactionBody{
				ScheduleID: builder.pb.GetScheduleID(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder ScheduleDeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ScheduleDeleteTransaction {
	return ScheduleDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder ScheduleDeleteTransaction) SetTransactionMemo(memo string) ScheduleDeleteTransaction {
	return ScheduleDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder ScheduleDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) ScheduleDeleteTransaction {
	return ScheduleDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder ScheduleDeleteTransaction) SetTransactionID(transactionID TransactionID) ScheduleDeleteTransaction {
	return ScheduleDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder ScheduleDeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) ScheduleDeleteTransaction {
	return ScheduleDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

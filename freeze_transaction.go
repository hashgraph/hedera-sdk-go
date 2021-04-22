package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FreezeTransaction struct {
	TransactionBuilder
	pb *proto.FreezeTransactionBody
}

func NewFreezeTransaction() FreezeTransaction {
	pb := &proto.FreezeTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_Freeze{Freeze: pb}

	builder := FreezeTransaction{inner, pb}

	return builder
}

func freezeTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) FreezeTransaction {
	return FreezeTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetFreeze(),
	}
}

func (builder FreezeTransaction) SetStartTime(hour uint8, minute uint8) FreezeTransaction {
	builder.pb.StartHour = int32(hour)
	builder.pb.StartMin = int32(minute)
	return builder
}

func (builder FreezeTransaction) SetEndTime(hour uint8, minute uint8) FreezeTransaction {
	builder.pb.EndHour = int32(hour)
	builder.pb.EndMin = int32(minute)
	return builder
}

func (builder FreezeTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *FreezeTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_Freeze{
			Freeze: &proto.FreezeTransactionBody{
				StartHour:  builder.pb.GetStartHour(),
				StartMin:   builder.pb.GetStartMin(),
				EndHour:    builder.pb.GetEndHour(),
				EndMin:     builder.pb.GetEndMin(),
				UpdateFile: builder.pb.GetUpdateFile(),
				FileHash:   builder.pb.GetFileHash(),
			},
		},
	}, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder FreezeTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder FreezeTransaction) SetTransactionMemo(memo string) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder FreezeTransaction) SetTransactionValidDuration(validDuration time.Duration) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder FreezeTransaction) SetTransactionID(transactionID TransactionID) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder FreezeTransaction) SetNodeAccountID(nodeAccountID AccountID) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

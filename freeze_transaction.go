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

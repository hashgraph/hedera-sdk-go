package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
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

func (builder FreezeTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder FreezeTransaction) SetTransactionMemo(memo string) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder FreezeTransaction) SetTransactionValidDuration(validDuration time.Duration) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder FreezeTransaction) SetTransactionID(transactionID TransactionID) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder FreezeTransaction) SetNodeAccountID(nodeAccountID AccountID) FreezeTransaction {
	return FreezeTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

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
	inner.pb.Data = &proto.TransactionBody_Freeze{pb}

	builder := FreezeTransaction{inner, pb}

	return builder
}

func (builder FreezeTransaction) SetStartTime(start time.Time) FreezeTransaction {
	builder.pb.StartHour = int32(start.Hour())
	builder.pb.StartMin = int32(start.Minute())
	return builder
}

func (builder FreezeTransaction) SetEndTime(end time.Time) FreezeTransaction {
	builder.pb.EndHour = int32(end.Hour())
	builder.pb.EndMin = int32(end.Minute())
	return builder
}

func (builder FreezeTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder FreezeTransaction) SetMaxTransactionFee(maxTransactionFee uint64) FreezeTransaction {
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

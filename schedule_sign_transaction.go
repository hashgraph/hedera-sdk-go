package hedera

import (
	"errors"
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ScheduleSignTransaction struct {
	TransactionBuilder
	pb *proto.ScheduleSignTransactionBody
}

func NewScheduleSignTransaction() ScheduleSignTransaction {
	pb := &proto.ScheduleSignTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ScheduleSign{ScheduleSign: pb}

	builder := ScheduleSignTransaction{inner, pb}

	return builder
}

func (builder ScheduleSignTransaction) SetScheduleID(id ScheduleID) ScheduleSignTransaction {
	builder.pb.ScheduleID = id.toProto()

	return builder
}

func (transaction *ScheduleSignTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleSignTransaction")
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder ScheduleSignTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ScheduleSignTransaction {
	return ScheduleSignTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder ScheduleSignTransaction) SetTransactionMemo(memo string) ScheduleSignTransaction {
	return ScheduleSignTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder ScheduleSignTransaction) SetTransactionValidDuration(validDuration time.Duration) ScheduleSignTransaction {
	return ScheduleSignTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder ScheduleSignTransaction) SetTransactionID(transactionID TransactionID) ScheduleSignTransaction {
	return ScheduleSignTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder ScheduleSignTransaction) SetNodeAccountID(nodeAccountID AccountID) ScheduleSignTransaction {
	return ScheduleSignTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

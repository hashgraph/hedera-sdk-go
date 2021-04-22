package hedera

import (
	"errors"
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ScheduleCreateTransaction struct {
	TransactionBuilder
	pb *proto.ScheduleCreateTransactionBody
}

func NewScheduleCreateTransaction() ScheduleCreateTransaction {
	pb := &proto.ScheduleCreateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ScheduleCreate{ScheduleCreate: pb}

	builder := ScheduleCreateTransaction{inner, pb}

	return builder
}

func (builder ScheduleCreateTransaction) SetScheduledTransaction(tx ITransaction) (ScheduleCreateTransaction, error) {
	scheduled, err := tx.constructScheduleProtobuf()
	if err != nil {
		return builder, err
	}

	builder.pb.ScheduledTransactionBody = scheduled
	return builder, nil
}

func (builder ScheduleCreateTransaction) SetPayerAccountID(id AccountID) ScheduleCreateTransaction {
	builder.pb.PayerAccountID = id.toProto()

	return builder
}

func (builder ScheduleCreateTransaction) SetAdminKey(key PublicKey) ScheduleCreateTransaction {
	builder.pb.AdminKey = key.toProto()

	return builder
}

func (builder ScheduleCreateTransaction) SetScheduleMemo(memo string) ScheduleCreateTransaction {
	builder.pb.Memo = memo

	return builder
}

func (builder ScheduleCreateTransaction) setSchedulableTransactionBody(txBody *proto.SchedulableTransactionBody) ScheduleCreateTransaction {
	builder.pb.ScheduledTransactionBody = txBody

	return builder
}

func (builder *ScheduleCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleCreateTransaction`")
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder ScheduleCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ScheduleCreateTransaction {
	return ScheduleCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder ScheduleCreateTransaction) SetTransactionMemo(memo string) ScheduleCreateTransaction {
	return ScheduleCreateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder ScheduleCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) ScheduleCreateTransaction {
	return ScheduleCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder ScheduleCreateTransaction) SetTransactionID(transactionID TransactionID) ScheduleCreateTransaction {
	return ScheduleCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder ScheduleCreateTransaction) SetNodeAccountID(nodeAccountID AccountID) ScheduleCreateTransaction {
	return ScheduleCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

package hedera

import (
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

func (builder ScheduleCreateTransaction) SetTransaction(transaction Transaction) ScheduleCreateTransaction {
	other := transaction.Schedule()
	builder.SetTransactionBody(other.TransactionBuilder.pb.GetScheduleCreate().TransactionBody)
	builder.pb.SigMap = other.TransactionBuilder.pb.GetScheduleCreate().SigMap
	return builder
}

func (builder ScheduleCreateTransaction) SetPayerAccountID(id AccountID) ScheduleCreateTransaction {
	builder.pb.PayerAccountID = id.toProto()

	return builder
}

func (builder ScheduleCreateTransaction) GetPayerAccountID() AccountID {
	return accountIDFromProto(builder.pb.PayerAccountID)
}

func (builder ScheduleCreateTransaction) SetAdminKey(key PublicKey) ScheduleCreateTransaction {
	builder.pb.AdminKey = key.toProto()

	return builder
}

func (builder ScheduleCreateTransaction) GetAdminKey() *PublicKey {
	key, err := publicKeyFromProto(builder.pb.GetAdminKey())
	if err != nil {
		return nil
	}
	return &key
}

func (builder ScheduleCreateTransaction) SetTransactionBody(dat []byte) ScheduleCreateTransaction {
	builder.pb.TransactionBody = dat

	return builder
}

func (builder ScheduleCreateTransaction) GetTransactionBody() []byte {
	return builder.pb.TransactionBody
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

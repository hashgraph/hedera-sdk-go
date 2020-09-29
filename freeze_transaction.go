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

	transaction := FreezeTransaction{inner, pb}

	return transaction
}

func (transaction FreezeTransaction) SetStartTime(hour uint8, minute uint8) FreezeTransaction {
	transaction.pb.StartHour = int32(hour)
	transaction.pb.StartMin = int32(minute)
	return transaction
}

func (transaction FreezeTransaction) SetEndTime(hour uint8, minute uint8) FreezeTransaction {
	transaction.pb.EndHour = int32(hour)
	transaction.pb.EndMin = int32(minute)
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction FreezeTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FreezeTransaction {
	return FreezeTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction FreezeTransaction) SetTransactionMemo(memo string) FreezeTransaction {
	return FreezeTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction FreezeTransaction) SetTransactionValidDuration(validDuration time.Duration) FreezeTransaction {
	return FreezeTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction FreezeTransaction) SetTransactionID(transactionID TransactionID) FreezeTransaction {
	return FreezeTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction FreezeTransaction) SetNodeID(nodeAccountID AccountID) FreezeTransaction {
	return FreezeTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}

package hedera

import (
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

func (builder ScheduleSignTransaction) AddScheduleSignature(key Ed25519PublicKey, signature []byte) ScheduleSignTransaction {
	var sigPair proto.SignaturePair

	sigPair.PubKeyPrefix = key.Bytes()

	switch sig := sigPair.Signature.(type) {
	case *proto.SignaturePair_Ed25519:
		sig.Ed25519 = signature
	}

	if builder.pb.SigMap != nil {
		if builder.pb.SigMap.SigPair != nil {
			builder.pb.SigMap.SigPair = append(builder.pb.SigMap.SigPair, &sigPair)
		} else {
			builder.pb.SigMap.SigPair = make([]*proto.SignaturePair, 0)
			builder.pb.SigMap.SigPair = append(builder.pb.SigMap.SigPair, &sigPair)
		}
	} else {
		builder.pb.SigMap = &proto.SignatureMap{
			SigPair: make([]*proto.SignaturePair, 0),
		}
		builder.pb.SigMap.SigPair = append(builder.pb.SigMap.SigPair, &sigPair)
	}

	return builder
}

func (builder ScheduleSignTransaction) SetScheduleID(id ScheduleID) ScheduleSignTransaction {
	builder.pb.ScheduleID = id.toProto()

	return builder
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

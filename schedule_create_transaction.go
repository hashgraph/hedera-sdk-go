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
	builder.pb.TransactionBody = other.TransactionBuilder.pb.GetScheduleCreate().TransactionBody
	builder.pb.SigMap = other.TransactionBuilder.pb.GetScheduleCreate().SigMap
	return builder
}

func (builder ScheduleCreateTransaction) SetPayerAccountID(id AccountID) ScheduleCreateTransaction {
	builder.pb.PayerAccountID = id.toProto()

	return builder
}

func (builder ScheduleCreateTransaction) SetAdminKey(key PublicKey) ScheduleCreateTransaction {
	builder.pb.AdminKey = key.toProto()

	return builder
}

func (builder *ScheduleCreateTransaction) GetScheduleSignatures() (map[*Ed25519PublicKey][]byte, error) {
	signMap := make(map[*Ed25519PublicKey][]byte, len(builder.pb.GetSigMap().GetSigPair()))

	for _, sigPair := range builder.pb.GetSigMap().GetSigPair() {
		key, err := Ed25519PublicKeyFromBytes(sigPair.PubKeyPrefix)
		if err != nil {
			return make(map[*Ed25519PublicKey][]byte, 0), err
		}
		switch sigPair.Signature.(type) {
		case *proto.SignaturePair_Contract:
			signMap[&key] = sigPair.GetContract()
		case *proto.SignaturePair_Ed25519:
			signMap[&key] = sigPair.GetEd25519()
		case *proto.SignaturePair_RSA_3072:
			signMap[&key] = sigPair.GetRSA_3072()
		case *proto.SignaturePair_ECDSA_384:
			signMap[&key] = sigPair.GetECDSA_384()
		}
	}

	return signMap, nil
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

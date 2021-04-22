package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type SystemUndeleteTransaction struct {
	TransactionBuilder
	pb *proto.SystemUndeleteTransactionBody
}

func NewSystemUndeleteTransaction() SystemUndeleteTransaction {
	pb := &proto.SystemUndeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_SystemUndelete{SystemUndelete: pb}

	builder := SystemUndeleteTransaction{inner, pb}

	return builder
}

func systemUndeleteTransactionFromProtobuf(transactionBuilder TransactionBuilder, pb *proto.TransactionBody) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{
		TransactionBuilder: transactionBuilder,
		pb:                 pb.GetSystemUndelete(),
	}
}

func (builder SystemUndeleteTransaction) SetContractID(ID ContractID) SystemUndeleteTransaction {
	builder.pb.Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: ID.toProto()}
	return builder
}

func (builder SystemUndeleteTransaction) SetFileID(ID FileID) SystemUndeleteTransaction {
	builder.pb.Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: ID.toProto()}
	return builder
}

func (builder SystemUndeleteTransaction) Schedule() (ScheduleCreateTransaction, error) {
	scheduled, err := builder.constructScheduleProtobuf()
	if err != nil {
		return ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (builder *SystemUndeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.SchedulableTransactionBody{
		TransactionFee: builder.TransactionBuilder.pb.GetTransactionFee(),
		Memo:           builder.TransactionBuilder.pb.GetMemo(),
		Data: &proto.SchedulableTransactionBody_SystemUndelete{
			SystemUndelete: &proto.SystemUndeleteTransactionBody{},
		},
	}

	switch builder.pb.GetId().(type) {
	case *proto.SystemUndeleteTransactionBody_ContractID:
		body.GetSystemUndelete().Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: builder.pb.GetContractID()}
	case *proto.SystemUndeleteTransactionBody_FileID:
		body.GetSystemUndelete().Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: builder.pb.GetFileID()}
	}

	return body, nil
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder SystemUndeleteTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder SystemUndeleteTransaction) SetTransactionMemo(memo string) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder SystemUndeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder SystemUndeleteTransaction) SetTransactionID(transactionID TransactionID) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder SystemUndeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

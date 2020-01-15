package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
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

func (builder SystemUndeleteTransaction) SetContractID(ID ContractID) SystemUndeleteTransaction {
	builder.pb.Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: ID.toProto()}
	return builder
}

func (builder SystemUndeleteTransaction) SetFileID(ID FileID) SystemUndeleteTransaction {
	builder.pb.Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: ID.toProto()}
	return builder
}

func (builder SystemUndeleteTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder SystemUndeleteTransaction) SetMaxTransactionFee(maxTransactionFee uint64) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder SystemUndeleteTransaction) SetTransactionMemo(memo string) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder SystemUndeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder SystemUndeleteTransaction) SetTransactionID(transactionID TransactionID) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder SystemUndeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

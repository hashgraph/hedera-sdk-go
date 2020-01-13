package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type SystemDeleteTransaction struct {
	TransactionBuilder
	pb *proto.SystemDeleteTransactionBody
}

func NewSystemDeleteTransaction() SystemDeleteTransaction {
	pb := &proto.SystemDeleteTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_SystemDelete{pb}

	builder := SystemDeleteTransaction{inner, pb}

	return builder
}

func (builder SystemDeleteTransaction) SetExpirationTime(expiration time.Time) SystemDeleteTransaction {
	builder.pb.ExpirationTime = &proto.TimestampSeconds{
		Seconds: expiration.Unix(),
	}
	return builder
}

func (builder SystemDeleteTransaction) SetID(id ContractIdOrFileID) SystemDeleteTransaction {
	file, contract, ty := id.toProtoContractIDOrFile()
	if ty == 0 {
		builder.pb.Id = &proto.SystemDeleteTransactionBody_FileID{FileID: file}
	} else {
		builder.pb.Id = &proto.SystemDeleteTransactionBody_ContractID{ContractID: contract}
	}

	return builder
}

func (builder SystemDeleteTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder SystemDeleteTransaction) SetMaxTransactionFee(maxTransactionFee uint64) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder SystemDeleteTransaction) SetTransactionMemo(memo string) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder SystemDeleteTransaction) SetTransactionValidDuration(validDuration time.Duration) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder SystemDeleteTransaction) SetTransactionID(transactionID TransactionID) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder SystemDeleteTransaction) SetNodeAccountID(nodeAccountID AccountID) SystemDeleteTransaction {
	return SystemDeleteTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ContractUpdateTransaction struct {
	TransactionBuilder
	pb *proto.ContractUpdateTransactionBody
}

func NewContractUpdateTransaction() ContractUpdateTransaction {
	pb := &proto.ContractUpdateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractUpdateInstance{pb}

	builder := ContractUpdateTransaction{inner, pb}

	return builder
}

func (builder ContractUpdateTransaction) SetContractID(id ContractID) ContractUpdateTransaction {
	builder.pb.ContractID = id.toProto()
	return builder
}

func (builder ContractUpdateTransaction) SetBytecodeFile(id FileID) ContractUpdateTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder ContractUpdateTransaction) SetAdminKey(publicKey Ed25519PublicKey) ContractUpdateTransaction {
	builder.pb.AdminKey = publicKey.toProto()
	return builder
}

func (builder ContractUpdateTransaction) SetProxyAccountID(id AccountID) ContractUpdateTransaction {
	builder.pb.ProxyAccountID = id.toProto()
	return builder
}

func (builder ContractUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) ContractUpdateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return builder
}

func (builder ContractUpdateTransaction) SetExpirationTime(expiration time.Time) ContractUpdateTransaction {
	builder.pb.ExpirationTime = timeToProto(expiration)
	return builder
}

func (builder ContractUpdateTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ContractUpdateTransaction) SetMaxTransactionFee(maxTransactionFee uint64) ContractUpdateTransaction {
	return ContractUpdateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder ContractUpdateTransaction) SetMemo(memo string) ContractUpdateTransaction {
	return ContractUpdateTransaction{builder.TransactionBuilder.SetMemo(memo), builder.pb}
}

func (builder ContractUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractUpdateTransaction {
	return ContractUpdateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder ContractUpdateTransaction) SetTransactionID(transactionID TransactionID) ContractUpdateTransaction {
	return ContractUpdateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ContractUpdateTransaction) SetNodeAccountID(nodeAccountID AccountID) ContractUpdateTransaction {
	return ContractUpdateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

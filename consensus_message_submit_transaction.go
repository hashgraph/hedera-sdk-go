package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ConsensusMessageSubmitTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusSubmitMessageTransactionBody
}

func NewConsensusMessageSubmitTransaction() ConsensusMessageSubmitTransaction {
	pb := &proto.ConsensusSubmitMessageTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusSubmitMessage{pb}

	builder := ConsensusMessageSubmitTransaction{inner, pb}

	return builder
}

func (builder ConsensusMessageSubmitTransaction) SetTopicID(id ConsensusTopicID) ConsensusMessageSubmitTransaction {
	builder.pb.TopicID = id.toProto()
	return builder
}

func (builder ConsensusMessageSubmitTransaction) SetMessage(message []byte) ConsensusMessageSubmitTransaction {
	builder.pb.Message = message
	return builder
}

func (builder ConsensusMessageSubmitTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ConsensusMessageSubmitTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder ConsensusMessageSubmitTransaction) SetMemo(memo string) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder ConsensusMessageSubmitTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder ConsensusMessageSubmitTransaction) SetTransactionID(transactionID TransactionID) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ConsensusMessageSubmitTransaction) SetNodeAccountID(nodeAccountID AccountID) ConsensusMessageSubmitTransaction {
	return ConsensusMessageSubmitTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

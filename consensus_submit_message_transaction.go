package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ConsensusSubmitMessageTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusSubmitMessageTransactionBody
}

func NewConsensusSubmitMessageTransaction() ConsensusSubmitMessageTransaction {
	pb := &proto.ConsensusSubmitMessageTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusSubmitMessage{pb}

	builder := ConsensusSubmitMessageTransaction{inner, pb}

	return builder
}

func (builder ConsensusSubmitMessageTransaction) SetTopicID(id ConsensusTopicID) ConsensusSubmitMessageTransaction {
	builder.pb.TopicID = id.toProto()
	return builder
}

func (builder ConsensusSubmitMessageTransaction) SetMessage(message []byte) ConsensusSubmitMessageTransaction {
	builder.pb.Message = message
	return builder
}

func (builder ConsensusSubmitMessageTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ConsensusSubmitMessageTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusSubmitMessageTransaction {
	return ConsensusSubmitMessageTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder ConsensusSubmitMessageTransaction) SetMemo(memo string) ConsensusSubmitMessageTransaction {
	return ConsensusSubmitMessageTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder ConsensusSubmitMessageTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusSubmitMessageTransaction {
	return ConsensusSubmitMessageTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder ConsensusSubmitMessageTransaction) SetTransactionID(transactionID TransactionID) ConsensusSubmitMessageTransaction {
	return ConsensusSubmitMessageTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ConsensusSubmitMessageTransaction) SetNodeAccountID(nodeAccountID AccountID) ConsensusSubmitMessageTransaction {
	return ConsensusSubmitMessageTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

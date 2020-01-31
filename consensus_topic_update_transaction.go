package hedera

import (
	"time"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ConsensusTopicUpdateTransaction struct {
	TransactionBuilder
	pb *proto.ConsensusUpdateTopicTransactionBody
}

func NewConsensusTopicUpdateTransaction() ConsensusTopicUpdateTransaction {
	pb := &proto.ConsensusUpdateTopicTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ConsensusUpdateTopic{pb}

	builder := ConsensusTopicUpdateTransaction{inner, pb}

	return builder
}

func (builder ConsensusTopicUpdateTransaction) SetTopicID(id ConsensusTopicID) ConsensusTopicUpdateTransaction {
	builder.pb.TopicID = id.toProto()
	return builder
}

func (builder ConsensusTopicUpdateTransaction) SetAdminKey(publicKey PublicKey) ConsensusTopicUpdateTransaction {
	builder.pb.AdminKey = publicKey.toProto()
	return builder
}

func (builder ConsensusTopicUpdateTransaction) SetSubmitKey(publicKey PublicKey) ConsensusTopicUpdateTransaction {
	builder.pb.SubmitKey = publicKey.toProto()
	return builder
}

func (builder ConsensusTopicUpdateTransaction) SetTopicMemo(memo string) ConsensusTopicUpdateTransaction {
	builder.pb.Memo = &wrappers.StringValue{Value: memo}
	return builder
}

func (builder ConsensusTopicUpdateTransaction) SetExpirationTime(expiration time.Time) ConsensusTopicUpdateTransaction {
	builder.pb.ExpirationTime = timeToProto(expiration)
	return builder
}

func (builder ConsensusTopicUpdateTransaction) SetAutoRenewPeriod(period time.Duration) ConsensusTopicUpdateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(period)
	return builder
}

func (builder ConsensusTopicUpdateTransaction) SetAutoRenewAccountID(id AccountID) ConsensusTopicUpdateTransaction {
	builder.pb.AutoRenewAccount = id.toProto()
	return builder
}

// ClearTopicMemo explicitly clears any memo on the topic by sending an empty string as the memo
func (builder ConsensusTopicUpdateTransaction) ClearTopicMemo() ConsensusTopicUpdateTransaction {
	return builder.SetTopicMemo("")
}

// ClearAdminKey explicitly clears any admin key on the topic by sending an empty key list as the key
func (builder ConsensusTopicUpdateTransaction) ClearAdminKey() ConsensusTopicUpdateTransaction {
	return builder.SetAdminKey(NewKeyList())
}

// ClearSubmitKey explicitly clears any submit key on the topic by sending an empty key list as the key
func (builder ConsensusTopicUpdateTransaction) ClearSubmitKey() ConsensusTopicUpdateTransaction {
	return builder.SetSubmitKey(NewKeyList())
}

// ClearAutoRenewAccountID explicitly clears any auto renew account ID on the topic by sending an empty accountID
func (builder ConsensusTopicUpdateTransaction) ClearAutoRenewAccountID() ConsensusTopicUpdateTransaction {
	builder.pb.AutoRenewAccount = &proto.AccountID{}

	return builder
}

func (builder ConsensusTopicUpdateTransaction) Build(client *Client) Transaction {
	return builder.TransactionBuilder.Build(client)
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder ConsensusTopicUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

func (builder ConsensusTopicUpdateTransaction) SetMemo(memo string) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

func (builder ConsensusTopicUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

func (builder ConsensusTopicUpdateTransaction) SetTransactionID(transactionID TransactionID) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

func (builder ConsensusTopicUpdateTransaction) SetNodeAccountID(nodeAccountID AccountID) ConsensusTopicUpdateTransaction {
	return ConsensusTopicUpdateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}

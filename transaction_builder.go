package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type TransactionBuilder struct {
	pb *proto.TransactionBody
}

func newTransactionBuilder() TransactionBuilder {
	builder := TransactionBuilder{&proto.TransactionBody{}}
	builder.SetTransactionValidDuration(120 * time.Second)

	return builder
}

func (builder TransactionBuilder) Build(client *Client) Transaction {
	if builder.pb.TransactionFee == 0 {
		builder.SetMaxTransactionFee(client.maxTransactionFee)
	}

	if builder.pb.NodeAccountID == nil {
		builder.SetNodeAccountID(client.randomNode().id)
	}

	if builder.pb.TransactionID == nil {
		builder.SetTransactionID(NewTransactionID(client.operator.accountID))
	}

	bodyBytes, err := protobuf.Marshal(builder.pb)
	if err != nil {
		// This should be unreachable
		// From the documentation this appears to only be possible if there are missing proto types
		panic(err)
	}

	pb := &proto.Transaction{
		BodyData: &proto.Transaction_BodyBytes{BodyBytes: bodyBytes},
		SigMap:   &proto.SignatureMap{SigPair: []*proto.SignaturePair{}},
	}

	return Transaction{pb, transactionIDFromProto(builder.pb.TransactionID)}
}

func (builder TransactionBuilder) Execute(client *Client) (TransactionID, error) {
	return builder.Build(client).Execute(client)
}

//
// Shared
//

func (builder TransactionBuilder) SetMaxTransactionFee(maxTransactionFee Hbar) TransactionBuilder {
	builder.pb.TransactionFee = uint64(maxTransactionFee.AsTinybar())
	return builder
}

func (builder TransactionBuilder) SetTransactionMemo(memo string) TransactionBuilder {
	builder.pb.Memo = memo
	return builder
}

func (builder TransactionBuilder) SetTransactionValidDuration(validDuration time.Duration) TransactionBuilder {
	builder.pb.TransactionValidDuration = durationToProto(validDuration)
	return builder
}

func (builder TransactionBuilder) SetTransactionID(transactionID TransactionID) TransactionBuilder {
	builder.pb.TransactionID = transactionID.toProto()
	return builder
}

func (builder TransactionBuilder) SetNodeAccountID(nodeAccountID AccountID) TransactionBuilder {
	builder.pb.NodeAccountID = nodeAccountID.toProto()
	return builder
}

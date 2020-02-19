package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
	"math"
	"time"
)

// TransactionBuilder is used to construct Transactions. The state is mutable through the various setter functions.
type TransactionBuilder struct {
	pb *proto.TransactionBody

	// unfortunately; this is required to prevent setting the max TXFee if it is purposely set to 0
	// (for example, when .Cost() is called)
	noTXFee bool
}

func newTransactionBuilder() TransactionBuilder {
	builder := TransactionBuilder{pb: &proto.TransactionBody{}}
	builder.SetTransactionValidDuration(120 * time.Second)

	return builder
}

// Build validates and finalizes the transaction's state and prepares it for execution, returning a Transaction.
// The inner state becomes immutable, however it can still be signed after building.
func (builder TransactionBuilder) Build(client *Client) (Transaction, error) {
	if client != nil && !builder.noTXFee {
		builder.SetMaxTransactionFee(client.maxTransactionFee)
	}

	if builder.pb.NodeAccountID == nil {
		if client != nil {
			builder.SetNodeAccountID(client.randomNode().id)
		}
	}

	if builder.pb.TransactionID == nil && client != nil && client.operator != nil {
		builder.SetTransactionID(NewTransactionID(client.operator.accountID))
	}

	// todo: add a validate function per transaction type
	if builder.pb.TransactionID == nil {
		return Transaction{}, newErrLocalValidationf(".setTransactionID() required")
	}

	if builder.pb.NodeAccountID == nil {
		return Transaction{}, newErrLocalValidationf(".setNodeAccountID() required")
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

	return Transaction{pb, transactionIDFromProto(builder.pb.TransactionID)}, nil
}

// Execute is a short hand function to build and execute a transaction. It first calls build on the TransactionBuilder
// and as long as validation passes it will then execute the resulting Transaction.
func (builder TransactionBuilder) Execute(client *Client) (TransactionID, error) {
	tx, err := builder.Build(client)

	if err != nil {
		return TransactionID{}, err
	}

	return tx.Execute(client)
}

// Cost returns the estimated cost of the transaction.
//
// NOTE: The actual cost returned by Hedera is within 99.8% to 99.9%  of the actual fee that will be assessed. We're
// unsure if this is because the fee fluctuates that much or if the calculations are simply incorrect on the server. To
// compensate for this we just bump by a 1% the value returned. As this would only ever be a maximum this will not cause
// you to be charged more.
func (builder TransactionBuilder) Cost(client *Client) (Hbar, error) {
	// An operator must be set on the client
	if client == nil || client.operator == nil {
		return ZeroHbar, newErrLocalValidationf("calling .Cost() requires client.SetOperator")
	}

	oldFee := builder.pb.TransactionFee
	oldTxID := builder.pb.TransactionID
	oldValidDuration := builder.pb.TransactionValidDuration
	oldTxFeeStatus := builder.noTXFee

	defer func() {
		// always reset the state of the builder before exiting this function
		builder.pb.TransactionFee = oldFee
		builder.pb.TransactionID = oldTxID
		builder.pb.TransactionValidDuration = oldValidDuration
		builder.noTXFee = oldTxFeeStatus
	}()

	builder.noTXFee = true

	costTx, err := builder.
		SetMaxTransactionFee(ZeroHbar).
		SetTransactionID(NewTransactionID(client.operator.accountID)).
		Build(client)

	if err != nil {
		return ZeroHbar, err
	}

	_, resp, err := costTx.
		executeForResponse(client)

	if err != nil {
		return ZeroHbar, err
	}

	status := Status(resp.NodeTransactionPrecheckCode)

	if status != StatusInsufficientTxFee {
		//  any status that is not insufficienttxfee should be considered an error in this case
		return ZeroHbar, newErrHederaPreCheckStatus(transactionIDFromProto(builder.pb.TransactionID), status)
	}

	return HbarFromTinybar(int64(math.Ceil(float64(resp.GetCost()) * 1.1))), nil
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

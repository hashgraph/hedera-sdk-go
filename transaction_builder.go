package hedera

import (
	"math"
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// TransactionBuilder is used to construct Transactions. The state is mutable through the various setter functions.
type TransactionBuilder struct {
	pb *proto.TransactionBody

	// unfortunately; this is required to prevent setting the max TXFee if it is purposely set to 0
	// (for example, when .GetCost() is called)
	noTXFee bool
}

func newTransactionBuilder() TransactionBuilder {
	transaction := TransactionBuilder{pb: &proto.TransactionBody{}}
	transaction.SetTransactionValidDuration(120 * time.Second)

	return transaction
}

// Build validates and finalizes the transaction's state and prepares it for execution, returning a Transaction.
// The inner state becomes immutable, however it can still be signed after building.
func (transaction TransactionBuilder) Build(client *Client) (Transaction, error) {
	if client != nil && !transaction.noTXFee && transaction.pb.TransactionFee == 0 {
		transaction.SetMaxTransactionFee(client.maxTransactionFee)
	}

	if transaction.pb.NodeAccountID == nil {
		if client != nil {
			transaction.SetNodeID(client.randomNode().id)
		}
	}

	if transaction.pb.TransactionID == nil && client != nil && client.operator != nil {
		transaction.SetTransactionID(NewTransactionID(client.operator.accountID))
	}

	// todo: add a validate function per transaction type
	if transaction.pb.TransactionID == nil {
		return Transaction{}, newErrLocalValidationf(".setTransactionID() required")
	}

	if transaction.pb.NodeAccountID == nil {
		return Transaction{}, newErrLocalValidationf(".setNodeAccountID() required")
	}

	bodyBytes, err := protobuf.Marshal(transaction.pb)
	if err != nil {
		// This should be unreachable
		// From the documentation this appears to only be possible if there are missing proto types
		panic(err)
	}

	pb := &proto.Transaction{
		BodyBytes: bodyBytes,
		SigMap:    &proto.SignatureMap{SigPair: []*proto.SignaturePair{}},
	}

	return Transaction{pb, transactionIDFromProto(transaction.pb.TransactionID)}, nil
}

// Execute is a short hand function to build and execute a transaction. It first calls build on the TransactionBuilder
// and as long as validation passes it will then execute the resulting Transaction.
func (transaction TransactionBuilder) Execute(client *Client) (TransactionID, error) {
	tx, err := transaction.Build(client)

	if err != nil {
		return TransactionID{}, err
	}

	return tx.Execute(client)
}

// GetCost returns the estimated cost of the transaction.
//
// NOTE: The actual cost returned by Hedera is within 99.8% to 99.9%  of the actual fee that will be assessed. We're
// unsure if this is because the fee fluctuates that much or if the calculations are simply incorrect on the server. To
// compensate for this we just bump by a 1% the value returned. As this would only ever be a maximum this will not cause
// you to be charged more.
func (transaction TransactionBuilder) GetCost(client *Client) (Hbar, error) {
	// An operator must be set on the client
	if client == nil || client.operator == nil {
		return ZeroHbar, newErrLocalValidationf("calling .GetCost() requires client.SetOperator")
	}

	oldFee := transaction.pb.TransactionFee
	oldTxID := transaction.pb.TransactionID
	oldValidDuration := transaction.pb.TransactionValidDuration
	oldTxFeeStatus := transaction.noTXFee

	defer func() {
		// always reset the state of the transaction before exiting this function
		transaction.pb.TransactionFee = oldFee
		transaction.pb.TransactionID = oldTxID
		transaction.pb.TransactionValidDuration = oldValidDuration
		transaction.noTXFee = oldTxFeeStatus
	}()

	transaction.noTXFee = true

	costTx, err := transaction.
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
		return ZeroHbar, newErrHederaPreCheckStatus(transactionIDFromProto(transaction.pb.TransactionID), status)
	}

	return HbarFromTinybar(int64(math.Ceil(float64(resp.GetCost()) * 1.1))), nil
}

//
// Shared
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction TransactionBuilder) SetMaxTransactionFee(maxTransactionFee Hbar) TransactionBuilder {
	transaction.pb.TransactionFee = uint64(maxTransactionFee.AsTinybar())
	return transaction
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction TransactionBuilder) SetTransactionMemo(memo string) TransactionBuilder {
	transaction.pb.Memo = memo
	return transaction
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction TransactionBuilder) SetTransactionValidDuration(validDuration time.Duration) TransactionBuilder {
	transaction.pb.TransactionValidDuration = durationToProto(validDuration)
	return transaction
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction TransactionBuilder) SetTransactionID(transactionID TransactionID) TransactionBuilder {
	transaction.pb.TransactionID = transactionID.toProto()
	return transaction
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction TransactionBuilder) SetNodeID(nodeAccountID AccountID) TransactionBuilder {
	transaction.pb.NodeAccountID = nodeAccountID.toProto()
	return transaction
}

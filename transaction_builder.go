package hedera

import (
	"strings"

	"github.com/hashgraph/hedera-sdk-go/hedera_proto"
)

type ErrorTransactionValidation struct {
	Messages []string
	Err      error
}

func (e *ErrorTransactionValidation) Error() string {
	return "The following requirements were not met: \n" + strings.Join(e.Messages, "\n")
}

type TransactionBuilderInterface interface {
	Validate() error
	Build() (*Transaction, error)
	Execute() (*TransactionID, error)
	ExecuteForReceipt() (*TransactionReceipt, error)
}

type TransactionBuilder struct {
	TransactionBuilderInterface
	client            *Client
	kind              TransactionKind
	MaxTransactionFee uint64
	body              hedera_proto.TransactionBody
}

func (tb TransactionBuilder) SetMemo(memo string) TransactionBuilder {
	tb.body.Memo = memo

	return tb
}

func (tb TransactionBuilder) SetMaxTransactionFee(fee uint64) TransactionBuilder {
	tb.MaxTransactionFee = fee

	return tb
}

func (tb TransactionBuilder) SetTransactionID(txID TransactionID) TransactionBuilder {
	tb.body.TransactionID = txID.proto()

	return tb
}

func (tb TransactionBuilder) SetTransactionValidDuration(seconds uint64) TransactionBuilder {
	tb.body.TransactionValidDuration = &hedera_proto.Duration{Seconds: int64(seconds)}

	return tb
}

func (tb TransactionBuilder) build(kind TransactionKind) (*Transaction, error) {
	if tb.client != nil {
		if tb.body.TransactionFee == 0 {
			tb.body.TransactionFee = tb.client.MaxTransactionFee()
		}

		if tb.body.TransactionValidDuration == nil {
			tb.body.TransactionValidDuration = &hedera_proto.Duration{Seconds: maxValidDuration}
		}

		if tb.body.NodeAccountID == nil {
			// let the client pick an actual node
			tb.body.NodeAccountID = tb.client.nodeID.proto()
		}
	}

	if tb.MaxTransactionFee == 0 {
		if tb.client != nil {
			tb.body.TransactionFee = tb.MaxTransactionFee
		}
	} else {
		tb.body.TransactionFee = tb.client.MaxTransactionFee()
	}

	protoBody := hedera_proto.Transaction_Body{
		Body: &tb.body,
	}

	tx := Transaction{
		Kind:   kind,
		client: tb.client,
		inner: hedera_proto.Transaction{
			BodyData: &protoBody,
		},
	}

	return &tx, nil
}

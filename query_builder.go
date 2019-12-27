package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
	"math"
	"math/rand"
	"time"
)

type QueryBuilder struct {
	pb       *proto.Query
	pbHeader *proto.QueryHeader

	maxPayment uint64
	payment    *uint64
}

func newQueryBuilder(pbHeader *proto.QueryHeader) QueryBuilder {
	builder := QueryBuilder{pb: &proto.Query{}, pbHeader: pbHeader}
	return builder
}

func (builder *QueryBuilder) SetMaxQueryPayment(maxPayment uint64) *QueryBuilder {
	builder.maxPayment = maxPayment
	return builder
}

func (builder *QueryBuilder) SetQueryPayment(paymentAmount uint64) *QueryBuilder {
	builder.payment = &paymentAmount
	return builder
}

func (builder *QueryBuilder) SetQueryPaymentTransaction(tx Transaction) *QueryBuilder {
	builder.pbHeader.Payment = tx.pb
	return builder
}

func (builder *QueryBuilder) Cost(client *Client) (uint64, error) {
	// Store the current response type and payment from the
	// query header
	currentResponseType := builder.pbHeader.ResponseType
	currentPayment := builder.pbHeader.Payment

	defer func() {
		// Reset the response type and payment transaction
		// on the query header
		builder.pbHeader.ResponseType = currentResponseType
		builder.pbHeader.Payment = currentPayment
	}()

	// Pick a random node for us to use
	node := client.randomNode()

	// COST_ANSWER tells Hedera to only return the cost for the given query body
	builder.pbHeader.ResponseType = proto.ResponseType_COST_ANSWER

	// COST_ANSWER requires a "null" payment (it checks for it but does not process it)
	tx := NewCryptoTransferTransaction().
		SetNodeAccountID(node.id).
		AddRecipient(node.id, 0).
		AddSender(client.operator.accountID, 0).
		Build(client)

	if client.operator.privateKey != nil {
		tx = tx.Sign(*client.operator.privateKey)
	} else {
		tx = tx.SignWith(client.operator.publicKey, client.operator.signer)
	}

	builder.pbHeader.Payment = tx.pb

	resp, err := execute(node, builder.pb, time.Now().Add(10*time.Second))
	if err != nil {
		return 0, err
	}

	return mapResponseHeader(resp).Cost, nil
}

func (builder *QueryBuilder) execute(client *Client) (*proto.Response, error) {
	var node *node

	if builder.isPaymentRequired() {
		if builder.pbHeader.Payment != nil {
			paymentBodyBytes := builder.pbHeader.Payment.GetBodyBytes()
			paymentBody := new(proto.TransactionBody)
			err := protobuf.Unmarshal(paymentBodyBytes, paymentBody)
			if err != nil {
				// The bodyBytes inside of the transaction at this point have been verified and this should be impossible
				panic(err)
			}

			nodeID := accountIDFromProto(paymentBody.NodeAccountID)
			node = client.node(nodeID)
		} else if builder.payment != nil {
			node = client.randomNode()
			builder.generatePaymentTransaction(client, node, *builder.payment)
		} else if builder.maxPayment > 0 || client.maxQueryPayment > 0 {
			node = client.randomNode()

			var maxPayment = builder.maxPayment
			if maxPayment == 0 {
				maxPayment = client.maxQueryPayment
			}

			actualCost, err := builder.Cost(client)
			if err != nil {
				return nil, err
			}

			if actualCost > maxPayment {
				return nil, fmt.Errorf("query cost of %v exceeds configured limit of %v", actualCost, maxPayment)
			}

			builder.generatePaymentTransaction(client, node, 0)
		}
	} else {
		node = client.randomNode()
	}

	var deadline time.Time

	switch builder.pb.Query.(type) {
	case *proto.Query_TransactionGetReceipt:
		// Receipt queries want to wait at most 2 minutes
		deadline = time.Now().Add(2 * time.Minute)

	default:
		// Most queries want to wait at most 10s
		deadline = time.Now().Add(10 * time.Second)
	}

	return execute(node, builder.pb, deadline)
}

func (builder *QueryBuilder) generatePaymentTransaction(client *Client, node *node, amount uint64) {
	tx := NewCryptoTransferTransaction().
		SetNodeAccountID(node.id).
		AddRecipient(node.id, amount).
		AddSender(client.operator.accountID, amount).
		Build(client)

	if client.operator.privateKey != nil {
		tx = tx.Sign(*client.operator.privateKey)
	} else {
		tx = tx.SignWith(client.operator.publicKey, client.operator.signer)
	}

	builder.pbHeader.Payment = tx.pb
}

func (builder *QueryBuilder) isPaymentRequired() bool {
	switch builder.pb.Query.(type) {
	case *proto.Query_TransactionGetReceipt:
		return false

	default:
		// All other queries cost
		return true
	}
}

func methodName(pb *proto.Query) string {
	switch pb.Query.(type) {
	case *proto.Query_TransactionGetReceipt:
		return "/proto.CryptoService/getTransactionReceipts"

	case *proto.Query_CryptogetAccountBalance:
		return "/proto.CryptoService/cryptoGetBalance"

	default:
		panic(fmt.Sprintf("[methodName] not implemented: %T", pb.Query))
	}
}

func mapResponseHeader(resp *proto.Response) *proto.ResponseHeader {
	switch resp.Response.(type) {
	case *proto.Response_TransactionGetReceipt:
		return resp.GetTransactionGetReceipt().Header

	case *proto.Response_CryptogetAccountBalance:
		return resp.GetCryptogetAccountBalance().Header

	default:
		panic(fmt.Sprintf("[mapResponseHeader] not implemented: %T", resp.Response))
	}
}

func isStatusUnknown(status proto.ResponseCodeEnum) bool {
	switch status {
	case proto.ResponseCodeEnum_UNKNOWN,
		proto.ResponseCodeEnum_RECEIPT_NOT_FOUND,
		proto.ResponseCodeEnum_RECORD_NOT_FOUND:
		return true

	default:
	}

	return false
}

func isResponseUnknown(resp *proto.Response) bool {
	switch resp.Response.(type) {
	case *proto.Response_TransactionGetReceipt:
		body := resp.GetTransactionGetReceipt()
		return isStatusUnknown(body.Header.NodeTransactionPrecheckCode) ||
			isStatusUnknown(body.Receipt.Status)

	case *proto.Response_TransactionGetRecord:
		body := resp.GetTransactionGetRecord()
		return isStatusUnknown(body.Header.NodeTransactionPrecheckCode) ||
			isStatusUnknown(body.TransactionRecord.Receipt.Status)

	default:
	}

	return false
}

func execute(node *node, pb *proto.Query, deadline time.Time) (*proto.Response, error) {
	methodName := methodName(pb)
	resp := new(proto.Response)

	for attempt := 0; true; attempt += 1 {
		if attempt > 0 && time.Now().After(deadline) {
			// Timed out
			break
		}

		if attempt > 0 {
			// After the first attempt, start an exponentially increasing delay
			delay := 500.0 * rand.Float64() * ((math.Pow(2, float64(attempt))) - 1)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		err := node.invoke(methodName, pb, resp)
		if err != nil {
			return nil, err
		}

		respHeader := mapResponseHeader(resp)

		if respHeader.NodeTransactionPrecheckCode == proto.ResponseCodeEnum_BUSY {
			// Try again (in a flash) on BUSY
			continue
		}

		if isResponseUnknown(resp) {
			// Receipts and Records can be flagged as unknown
			continue
		}

		if isStatusExceptional(respHeader.NodeTransactionPrecheckCode, true) {
			return nil, fmt.Errorf("%v", respHeader.NodeTransactionPrecheckCode)
		}

		return resp, nil
	}

	// Timed out
	// TODO: Better error here?
	respHeader := mapResponseHeader(resp)
	return nil, fmt.Errorf("%v", respHeader.NodeTransactionPrecheckCode)
}

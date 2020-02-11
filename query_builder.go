package hedera

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type QueryBuilder struct {
	pb       *proto.Query
	pbHeader *proto.QueryHeader

	maxPayment Hbar
	payment    *Hbar
}

func newQueryBuilder(pbHeader *proto.QueryHeader) QueryBuilder {
	builder := QueryBuilder{pb: &proto.Query{}, pbHeader: pbHeader}
	return builder
}

func (builder *QueryBuilder) SetMaxQueryPayment(maxPayment Hbar) *QueryBuilder {
	builder.maxPayment = maxPayment
	return builder
}

func (builder *QueryBuilder) SetQueryPayment(paymentAmount Hbar) *QueryBuilder {
	builder.payment = &paymentAmount
	return builder
}

func (builder *QueryBuilder) SetQueryPaymentTransaction(tx Transaction) *QueryBuilder {
	builder.pbHeader.Payment = tx.pb
	return builder
}

func (builder *QueryBuilder) Cost(client *Client) (Hbar, error) {
	// An operator must be set on the client
	if client == nil || client.operator == nil {
		return ZeroHbar, newErrLocalValidationf("calling .Cost() requires client.SetOperator")
	}

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
	tx, err := NewCryptoTransferTransaction().
		SetNodeAccountID(node.id).
		AddRecipient(node.id, ZeroHbar).
		AddSender(client.operator.accountID, ZeroHbar).
		Build(client)

	if err != nil {
		return ZeroHbar, err
	}

	tx = tx.signWithOperator(*client.operator)

	builder.pbHeader.Payment = tx.pb

	resp, err := execute(node, &tx.ID, builder.pb, time.Now().Add(10*time.Second))
	if err != nil {
		return ZeroHbar, err
	}

	tbCost := int64(mapResponseHeader(resp).Cost)

	// Some queries require more than the server requests, so 10% is added to the cost as an estimated max range
	switch builder.pb.Query.(type) {
	case *proto.Query_ContractCallLocal:
		tbCost = int64(float64(tbCost) * 1.1)
		break
	default:
		// do nothing -- add more if they are found
		break
	}

	return HbarFromTinybar(tbCost), nil
}

func (builder *QueryBuilder) execute(client *Client) (*proto.Response, error) {
	var node *node
	var payment *TransactionID

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
			txID := transactionIDFromProto(paymentBody.TransactionID)
			payment = &txID
			node = client.node(nodeID)
		} else if builder.payment != nil {
			node = client.randomNode()
			txID, err := builder.generatePaymentTransaction(client, node, *builder.payment)

			if err != nil {
				return nil, err
			}

			payment = &txID
		} else if builder.maxPayment.AsTinybar() > 0 || client.maxQueryPayment.AsTinybar() > 0 {
			node = client.randomNode()

			var maxPayment = builder.maxPayment
			if maxPayment.AsTinybar() == 0 {
				maxPayment = client.maxQueryPayment
			}

			actualCost, err := builder.Cost(client)

			if err != nil {
				return nil, err
			}

			if actualCost.AsTinybar() > maxPayment.AsTinybar() {
				return nil, newErrorMaxQueryPaymentExceeded(builder, actualCost, maxPayment)
			}

			txID, err := builder.generatePaymentTransaction(client, node, actualCost)
			if err != nil {
				return nil, err
			}

			payment = &txID

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

	return execute(node, payment, builder.pb, deadline)
}

func (builder *QueryBuilder) generatePaymentTransaction(client *Client, node *node, amount Hbar) (TransactionID, error) {
	tx, err := NewCryptoTransferTransaction().
		SetNodeAccountID(node.id).
		AddRecipient(node.id, amount).
		AddSender(client.operator.accountID, amount).
		SetMaxTransactionFee(HbarFrom(1, HbarUnits.Hbar)).
		Build(client)

	if err != nil {
		return TransactionID{}, err
	}

	if client.operator != nil {
		tx = tx.signWithOperator(*client.operator)
	}

	builder.pbHeader.Payment = tx.pb

	return tx.ID, nil
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

	// Crypto
	case *proto.Query_TransactionGetReceipt:
		return "/proto.CryptoService/getTransactionReceipts"

	case *proto.Query_CryptogetAccountBalance:
		return "/proto.CryptoService/cryptoGetBalance"

	case *proto.Query_CryptoGetAccountRecords:
		return "/proto.CryptoService/getAccountRecords"

	case *proto.Query_CryptoGetInfo:
		return "/proto.CryptoService/getAccountInfo"

	case *proto.Query_TransactionGetRecord:
		return "/proto.CryptoService/getTxRecordByTxID"

	case *proto.Query_CryptoGetProxyStakers:
		return "/proto.CryptoService/getStakersByAccountID"

	// Smart Contracts
	case *proto.Query_ContractCallLocal:
		return "/proto.SmartContractService/contractCallLocalMethod"

	case *proto.Query_ContractGetBytecode:
		return "/proto.SmartContractService/ContractGetBytecode"

	case *proto.Query_ContractGetInfo:
		return "/proto.SmartContractService/getContractInfo"

	case *proto.Query_ContractGetRecords:
		return "/proto.SmartContractService/getTxRecordByContractID"

	// File
	case *proto.Query_FileGetContents:
		return "/proto.FileService/getFileContent"

	case *proto.Query_FileGetInfo:
		return "/proto.FileService/getFileInfo"

	default:
		panic(fmt.Sprintf("[methodName] not implemented: %T", pb.Query))
	}
}

func mapResponseHeader(resp *proto.Response) *proto.ResponseHeader {
	switch resp.Response.(type) {

	// Crypto
	case *proto.Response_TransactionGetReceipt:
		return resp.GetTransactionGetReceipt().Header

	case *proto.Response_CryptogetAccountBalance:
		return resp.GetCryptogetAccountBalance().Header

	case *proto.Response_CryptoGetAccountRecords:
		return resp.GetCryptoGetAccountRecords().Header

	case *proto.Response_CryptoGetInfo:
		return resp.GetCryptoGetInfo().Header

	case *proto.Response_TransactionGetRecord:
		return resp.GetTransactionGetRecord().Header

	case *proto.Response_CryptoGetProxyStakers:
		return resp.GetCryptoGetProxyStakers().Header

	// Smart Contracts
	case *proto.Response_ContractCallLocal:
		return resp.GetContractCallLocal().Header

	case *proto.Response_ContractGetBytecodeResponse:
		return resp.GetContractGetBytecodeResponse().Header

	case *proto.Response_ContractGetInfo:
		return resp.GetContractGetInfo().Header

	case *proto.Response_ContractGetRecordsResponse:
		return resp.GetContractGetRecordsResponse().Header

	// File
	case *proto.Response_FileGetContents:
		return resp.GetFileGetContents().Header

	case *proto.Response_FileGetInfo:
		return resp.GetFileGetInfo().Header

	// HCS
	case *proto.Response_ConsensusGetTopicInfo:
		return resp.GetConsensusGetTopicInfo().Header

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

		if body.TransactionRecord == nil {
			return false
		}

		return isStatusUnknown(body.TransactionRecord.Receipt.Status)

	default:
	}

	return false
}

func execute(node *node, paymentID *TransactionID, pb *proto.Query, deadline time.Time) (*proto.Response, error) {
	methodName := methodName(pb)
	resp := new(proto.Response)

	for attempt := 0; true; attempt++ {
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

		status := Status(respHeader.NodeTransactionPrecheckCode)

		if status.isExceptional(true) {
			// precheck failed, paymentID should never be nil in this case
			return resp, newErrHederaPreCheckStatus(*paymentID, status)
		}

		// success
		return resp, nil
	}

	// Timed out
	respHeader := mapResponseHeader(resp)
	if paymentID != nil {
		return nil, newErrHederaPreCheckStatus(*paymentID, Status(respHeader.NodeTransactionPrecheckCode))
	}

	return nil, newErrHederaNetwork(fmt.Errorf("timed out with status %v", Status(respHeader.NodeTransactionPrecheckCode)))
}

package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransactionRecordQuery struct {
	Query
	pb *proto.TransactionGetRecordQuery
}

func NewTransactionRecordQuery() *TransactionRecordQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := &proto.TransactionGetRecordQuery{Header: &header}
	query.pb.Query = &proto.Query_TransactionGetRecord{
		TransactionGetRecord: pb,
	}

	return &TransactionRecordQuery{
		Query: query,
		pb:    pb,
	}
}

func (query *TransactionRecordQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return Hbar{}, err
	}

	query.pbHeader.Payment = paymentTransaction
	query.pbHeader.ResponseType = proto.ResponseType_COST_ANSWER
	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		transactionRecordQuery_getMethod,
		transactionRecordQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTransactionGetRecord().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func transactionRecordQuery_shouldRetry(status Status, response response) bool {
	switch status {
	case StatusBusy, StatusUnknown, StatusReceiptNotFound:
		return true
	}
	status = Status(response.query.GetTransactionGetRecord().TransactionRecord.Receipt.Status)

	switch status {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound, StatusRecordNotFound:
		return true
	default:
		return false
	}
}

func transactionRecordQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetTransactionGetRecord().Header.NodeTransactionPrecheckCode)
}

func transactionRecordQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetTxRecordByTxID,
	}
}

func (query *TransactionRecordQuery) SetTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	query.pb.TransactionID = transactionID.toProtobuf()
	return query
}

func (query *TransactionRecordQuery) GetTransactionID() TransactionID {
	return transactionIDFromProtobuf(query.pb.TransactionID)
}

func (query *TransactionRecordQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionRecordQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TransactionRecordQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

func (query *TransactionRecordQuery) SetQueryPayment(queryPayment Hbar) *TransactionRecordQuery {
	query.queryPayment = queryPayment
	return query
}

func (query *TransactionRecordQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionRecordQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *TransactionRecordQuery) Execute(client *Client) (TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return TransactionRecord{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		cost = client.maxQueryPayment

		actualCost, err := query.GetCost(client)
		if err != nil {
			return TransactionRecord{}, err
		}

		if cost.tinybar > actualCost.tinybar {
			return TransactionRecord{}, ErrMaxQueryPaymentExceeded{}
		}

		cost = actualCost
	}

	for _, nodeID := range query.nodeIDs {
		transaction, err := query_makePaymentTransaction(
			query.paymentTransactionID,
			nodeID,
			client.operator,
			cost,
		)
		if err != nil {
			return TransactionRecord{}, err
		}

		query.paymentTransactions = append(query.paymentTransactions, transaction)
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		transactionRecordQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		transactionRecordQuery_getMethod,
		transactionRecordQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return TransactionRecord{}, err
	}

	return TransactionRecordFromProtobuf(resp.query.GetTransactionGetRecord().TransactionRecord), nil
}

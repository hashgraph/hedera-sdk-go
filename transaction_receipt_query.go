package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransactionReceiptQuery struct {
	Query
	pb *proto.TransactionGetReceiptQuery
}

func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	header := proto.QueryHeader{}
	query := newQuery(false, &header)
	pb := &proto.TransactionGetReceiptQuery{Header: &header}
	query.pb.Query = &proto.Query_TransactionGetReceipt{
		TransactionGetReceipt: pb,
	}

	return &TransactionReceiptQuery{
		Query: query,
		pb:    pb,
	}
}

func (query *TransactionReceiptQuery) GetCost(client *Client) (Hbar, error) {
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
		transactionReceiptQuery_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		transactionReceiptQuery_getMethod,
		transactionReceiptQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTransactionGetReceipt().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func transactionReceiptQuery_shouldRetry(status Status, response response) bool {
	if status == StatusPlatformTransactionNotCreated {
		return true
	}

	switch status {
	case StatusBusy, StatusUnknown, StatusReceiptNotFound:
		return true
	case StatusOk:
		break
	default:
		return false
	}

	status = Status(response.query.GetTransactionGetReceipt().GetReceipt().GetStatus())

	switch status {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound:
		return true
	default:
		return false
	}
}

func transactionReceiptQuery_mapResponseStatus(_ request, response response) Status {
	status := Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode())

	switch status {
	case StatusBusy, StatusUnknown, StatusReceiptNotFound:
		return status
	case StatusOk:
		break
	default:
		return status
	}

	return Status(response.query.GetTransactionGetReceipt().GetReceipt().GetStatus())
}

func transactionReceiptQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetTransactionReceipts,
	}
}

func (query *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	query.pb.TransactionID = transactionID.toProtobuf()
	return query
}

func (query *TransactionReceiptQuery) GetTransactionID() TransactionID {
	return transactionIDFromProtobuf(query.pb.TransactionID)
}

func (query *TransactionReceiptQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionReceiptQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TransactionReceiptQuery) SetQueryPayment(queryPayment Hbar) *TransactionReceiptQuery {
	query.queryPayment = queryPayment
	return query
}

func (query *TransactionReceiptQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionReceiptQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *TransactionReceiptQuery) SetMaxRetry(count int) *TransactionReceiptQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	if client == nil || client.operator == nil {
		return TransactionReceipt{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		transactionReceiptQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		transactionReceiptQuery_getMethod,
		transactionReceiptQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		switch precheckErr := err.(type) {
		case ErrHederaPreCheckStatus:
			return TransactionReceipt{}, newErrHederaReceiptStatus(precheckErr.TxID, precheckErr.Status)
		}
		return TransactionReceipt{}, err
	}

	return transactionReceiptFromProtobuf(resp.query.GetTransactionGetReceipt().GetReceipt()), nil
}

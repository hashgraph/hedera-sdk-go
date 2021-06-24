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
		transactionReceiptQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTransactionGetReceipt().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func transactionReceiptQuery_shouldRetry(request request, response response) executionState {
	switch Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusOk:
		break
	default:
		return executionStateError
	}

	switch Status(response.query.GetTransactionGetReceipt().GetReceipt().GetStatus()) {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusSuccess:
		return executionStateFinished
	default:
		return executionStateError
	}
}

func transactionReceiptQuery_mapStatusError(request request, response response) error {
	switch Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound, StatusOk:
		break
	default:
		return ErrHederaPreCheckStatus{
			Status: Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()),
		}
	}

	return ErrHederaReceiptStatus{
		Status:  Status(response.query.GetTransactionGetReceipt().GetReceipt().GetStatus()),
		TxID:    transactionIDFromProtobuf(request.query.pb.GetTransactionGetReceipt().TransactionID),
		Receipt: transactionReceiptFromProtobuf(response.query.GetTransactionGetReceipt().GetReceipt()),
	}
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
	if client == nil {
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
		transactionReceiptQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		switch precheckErr := err.(type) {
		case ErrHederaPreCheckStatus:
			return TransactionReceipt{}, newErrHederaReceiptStatus(precheckErr.TxID, precheckErr.Status)
		}
		return TransactionReceipt{}, err
	}

	receipt := transactionReceiptFromProtobuf(resp.query.GetTransactionGetReceipt().GetReceipt())
	if receipt.TokenID != nil {
		receipt.TokenID.SetNetworkName(*client.networkName)
	}
	if receipt.TopicID != nil {
		receipt.TopicID.SetNetworkName(*client.networkName)
	}
	if receipt.FileID != nil {
		receipt.FileID.SetNetworkName(*client.networkName)
	}
	if receipt.ContractID != nil {
		receipt.ContractID.SetNetworkName(*client.networkName)
	}
	if receipt.ScheduleID != nil {
		receipt.ScheduleID.SetNetworkName(*client.networkName)
	}
	if receipt.AccountID != nil {
		receipt.AccountID.SetNetworkName(*client.networkName)
	}
	if receipt.ScheduledTransactionID != nil {
		receipt.ScheduledTransactionID.AccountID.SetNetworkName(*client.networkName)
	}

	return receipt, nil
}

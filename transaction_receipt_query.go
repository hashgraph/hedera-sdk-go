package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransactionReceiptQuery struct {
	Query
	transactionID TransactionID
	duplicates    *bool
}

func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	return &TransactionReceiptQuery{
		Query: newQuery(false),
	}
}

func (query *TransactionReceiptQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.transactionID.AccountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *TransactionReceiptQuery) setIncludeDuplicates(duplicates bool) *TransactionReceiptQuery {
	query.duplicates = &duplicates
	return query
}

func (query *TransactionReceiptQuery) getIncludeDuplicates() bool {
	if query.duplicates != nil {
		return *query.duplicates
	}

	return false
}

func (query *TransactionReceiptQuery) build() *proto.Query_TransactionGetReceipt {
	body := &proto.TransactionGetReceiptQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.transactionID.AccountID.isZero() {
		body.TransactionID = query.transactionID.toProtobuf()
	}
	if query.duplicates != nil {
		body.IncludeDuplicates = *query.duplicates
	}

	return &proto.Query_TransactionGetReceipt{
		TransactionGetReceipt: body,
	}
}

func (query *TransactionReceiptQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.TransactionGetReceipt.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.TransactionGetReceipt.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TransactionReceiptQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.TransactionGetReceipt.Header.Payment = paymentTransaction
	pb.TransactionGetReceipt.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TransactionReceiptQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query.costQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		transactionReceiptQuery_shouldRetry,
		protoReq,
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
		Receipt: transactionReceiptFromProtobuf(response.query.GetTransactionGetReceipt().GetReceipt()),
	}
}

func transactionReceiptQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetTransactionReceipts,
	}
}

func (query *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	query.transactionID = transactionID
	return query
}

func (query *TransactionReceiptQuery) GetTransactionID() TransactionID {
	return query.transactionID
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

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return TransactionReceipt{}, err
	}

	query.build()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		transactionReceiptQuery_shouldRetry,
		query.queryMakeRequest(),
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

	return transactionReceiptFromProtobuf(resp.query.GetTransactionGetReceipt().GetReceipt()), nil
}

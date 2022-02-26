package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TransactionReceiptQuery struct {
	Query
	transactionID *TransactionID
	childReceipts *bool
	duplicates    *bool
	timestamp     time.Time
}

func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	header := services.QueryHeader{}
	return &TransactionReceiptQuery{
		Query: _NewQuery(false, &header),
	}
}

func (query *TransactionReceiptQuery) SetIncludeChildren(includeChildReceipts bool) *TransactionReceiptQuery {
	query.childReceipts = &includeChildReceipts
	return query
}

func (query *TransactionReceiptQuery) GetIncludeChildren() bool {
	if query.childReceipts != nil {
		return *query.childReceipts
	}

	return false
}

func (query *TransactionReceiptQuery) SetIncludeDuplicates(includeDuplicates bool) *TransactionReceiptQuery {
	query.duplicates = &includeDuplicates
	return query
}

func (query *TransactionReceiptQuery) GetIncludeDuplicates() bool {
	if query.duplicates != nil {
		return *query.duplicates
	}

	return false
}

func (query *TransactionReceiptQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := query.transactionID.AccountID.ValidateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (query *TransactionReceiptQuery) _Build() *services.Query_TransactionGetReceipt {
	body := &services.TransactionGetReceiptQuery{
		Header: &services.QueryHeader{},
	}

	if query.transactionID.AccountID != nil {
		body.TransactionID = query.transactionID._ToProtobuf()
	}

	if query.duplicates != nil {
		body.IncludeDuplicates = *query.duplicates
	}

	if query.childReceipts != nil {
		body.IncludeChildReceipts = query.GetIncludeChildren()
	}

	return &services.Query_TransactionGetReceipt{
		TransactionGetReceipt: body,
	}
}

func (query *TransactionReceiptQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return Hbar{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	query.timestamp = time.Now()

	for range query.nodeAccountIDs {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.TransactionGetReceipt.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TransactionReceiptQueryShouldRetry,
		_CostQueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TransactionReceiptQueryGetMethod,
		_TransactionReceiptQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTransactionGetReceipt().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _TransactionReceiptQueryShouldRetry(logID string, request _Request, response _Response) _ExecutionState {
	status := Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode())
	logCtx.Trace().Str("requestId", logID).Str("status", status.String()).Msg("receipt precheck status received")

	switch status {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusOk:
		break
	default:
		return executionStateError
	}

	status = Status(response.query.GetTransactionGetReceipt().GetReceipt().GetStatus())
	logCtx.Trace().Str("requestId", logID).Str("status", status.String()).Msg("receipt status received")

	switch status {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	default:
		return executionStateFinished
	}
}

func _TransactionReceiptQueryMapStatusError(request _Request, response _Response) error {
	switch Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound, StatusOk:
		break
	default:
		return ErrHederaPreCheckStatus{
			Status: Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()),
		}
	}

	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetTransactionGetReceipt().GetReceipt().GetStatus()),
	}
}

func _TransactionReceiptQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTransactionReceipts,
	}
}

func (query *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	query.transactionID = &transactionID
	return query
}

func (query *TransactionReceiptQuery) GetTransactionID() TransactionID {
	if query.transactionID == nil {
		return TransactionID{}
	}

	return *query.transactionID
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

func (query *TransactionReceiptQuery) SetMaxBackoff(max time.Duration) *TransactionReceiptQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *TransactionReceiptQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *TransactionReceiptQuery) SetMinBackoff(min time.Duration) *TransactionReceiptQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *TransactionReceiptQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	if client == nil {
		return TransactionReceipt{}, errNoClientProvided
	}

	var err error
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return TransactionReceipt{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return TransactionReceipt{}, err
	}

	query.timestamp = time.Now()
	query.nextPaymentTransactionIndex = 0
	query.paymentTransactions = make([]*services.Transaction, 0)

	pb := query._Build()
	pb.TransactionGetReceipt.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TransactionReceiptQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TransactionReceiptQueryGetMethod,
		_TransactionReceiptQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
	)

	if err, ok := err.(ErrHederaPreCheckStatus); ok {
		if resp.query.GetTransactionGetReceipt() != nil {
			return _TransactionReceiptFromProtobuf(resp.query.GetTransactionGetReceipt()), err
		}

		return TransactionReceipt{}, err
	}

	return _TransactionReceiptFromProtobuf(resp.query.GetTransactionGetReceipt()), nil
}

func (query *TransactionReceiptQuery) _GetLogID() string {
	return fmt.Sprintf("TransactionReceiptQuery:%d", query.timestamp.UnixNano())
}

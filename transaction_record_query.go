package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransactionRecordQuery struct {
	Query
	transactionID *TransactionID
}

func NewTransactionRecordQuery() *TransactionRecordQuery {
	return &TransactionRecordQuery{
		Query: _NewQuery(true),
	}
}

func (query *TransactionRecordQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := query.transactionID.AccountID.Validate(client); err != nil {
		return err
	}

	return nil
}

func (query *TransactionRecordQuery) _Build() *proto.Query_TransactionGetRecord {
	body := &proto.TransactionGetRecordQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.transactionID.AccountID._IsZero() {
		body.TransactionID = query.transactionID._ToProtobuf()
	}

	return &proto.Query_TransactionGetRecord{
		TransactionGetRecord: body,
	}
}

func (query *TransactionRecordQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.TransactionGetRecord.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.TransactionGetRecord.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TransactionRecordQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.TransactionGetRecord.Header.Payment = paymentTransaction
	pb.TransactionGetRecord.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TransactionRecordQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network._GetNodeAccountIDsForExecute()

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query._CostQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TransactionRecordQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_TransactionRecordQueryGetMethod,
		_TransactionRecordQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTransactionGetRecord().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _TransactionRecordQueryShouldRetry(request _Request, response _Response) _ExecutionState {
	switch Status(response.query.GetTransactionGetRecord().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusOk:
		if response.query.GetTransactionGetRecord().GetHeader().ResponseType == proto.ResponseType_COST_ANSWER {
			return executionStateFinished
		}
	default:
		return executionStateError
	}

	switch Status(response.query.GetTransactionGetRecord().GetTransactionRecord().GetReceipt().GetStatus()) {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusSuccess:
		return executionStateFinished
	default:
		return executionStateError
	}
}

func _TransactionRecordQueryMapStatusError(request _Request, response _Response) error {
	switch Status(response.query.GetTransactionGetRecord().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound, StatusOk:
		break
	default:
		return ErrHederaPreCheckStatus{
			Status: Status(response.query.GetTransactionGetRecord().GetHeader().GetNodeTransactionPrecheckCode()),
		}
	}

	return ErrHederaReceiptStatus{
		Status: Status(response.query.GetTransactionGetRecord().GetTransactionRecord().GetReceipt().GetStatus()),
		// TxID:    _TransactionIDFromProtobuf(_Request.query.pb.GetTransactionGetRecord().TransactionID, networkName),
		Receipt: _TransactionReceiptFromProtobuf(response.query.GetTransactionGetReceipt().GetReceipt()),
	}
}

func _TransactionRecordQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTxRecordByTxID,
	}
}

func (query *TransactionRecordQuery) SetTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	query.transactionID = &transactionID
	return query
}

func (query *TransactionRecordQuery) GetTransactionID() TransactionID {
	if query.transactionID == nil {
		return TransactionID{}
	}

	return *query.transactionID
}

func (query *TransactionRecordQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionRecordQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TransactionRecordQuery) SetQueryPayment(queryPayment Hbar) *TransactionRecordQuery {
	query.queryPayment = queryPayment
	return query
}

func (query *TransactionRecordQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionRecordQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *TransactionRecordQuery) SetMaxRetry(count int) *TransactionRecordQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *TransactionRecordQuery) SetMaxBackoff(max time.Duration) *TransactionRecordQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *TransactionRecordQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *TransactionRecordQuery) SetMinBackoff(min time.Duration) *TransactionRecordQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *TransactionRecordQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *TransactionRecordQuery) Execute(client *Client) (TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return TransactionRecord{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return TransactionRecord{}, err
	}

	query._Build()

	query.paymentTransactionID = TransactionIDGenerate(client.GetOperatorAccountID())

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.maxQueryPayment
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return TransactionRecord{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TransactionRecord{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TransactionRecordQuery",
			}
		}

		cost = actualCost
	}

	for _, nodeID := range query.nodeIDs {
		transaction, err := _QueryMakePaymentTransaction(
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

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TransactionRecordQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TransactionRecordQueryGetMethod,
		_TransactionRecordQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		if precheckErr, ok := err.(ErrHederaPreCheckStatus); ok {
			return TransactionRecord{}, _NewErrHederaReceiptStatus(precheckErr.TxID, precheckErr.Status)
		}
		return TransactionRecord{}, err
	}

	record := _TransactionRecordFromProtobuf(resp.query.GetTransactionGetRecord().TransactionRecord)
	record.TransactionID.AccountID._SetNetworkWithClient(client)
	if record.Receipt.TokenID != nil {
		record.Receipt.TokenID._SetNetworkWithClient(client)
	}
	if record.Receipt.TopicID != nil {
		record.Receipt.TopicID._SetNetworkWithClient(client)
	}
	if record.Receipt.FileID != nil {
		record.Receipt.FileID._SetNetworkWithClient(client)
	}
	if record.Receipt.ContractID != nil {
		record.Receipt.ContractID._SetNetworkWithClient(client)
	}
	if record.Receipt.ScheduleID != nil {
		record.Receipt.ScheduleID._SetNetworkWithClient(client)
	}
	if record.Receipt.AccountID != nil {
		record.Receipt.AccountID._SetNetworkWithClient(client)
	}
	if record.Receipt.ScheduledTransactionID != nil {
		record.Receipt.ScheduledTransactionID.AccountID._SetNetworkWithClient(client)
	}

	return record, nil
}

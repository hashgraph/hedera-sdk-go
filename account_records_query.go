package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountRecordsQuery gets all of the records for an account for any transfers into it and out of
// it, that were above the threshold, during the last 25 hours.
type AccountRecordsQuery struct {
	Query
	accountID *AccountID
}

// NewAccountRecordsQuery creates an AccountRecordsQuery query which can be used to construct and execute
// an AccountRecordsQuery.
//
// It is recommended that you use this for creating new instances of an AccountRecordQuery
// instead of manually creating an instance of the struct.
func NewAccountRecordsQuery() *AccountRecordsQuery {
	return &AccountRecordsQuery{
		Query: _NewQuery(true),
	}
}

// SetAccountID sets the account ID for which the records should be retrieved.
func (query *AccountRecordsQuery) SetAccountID(accountID AccountID) *AccountRecordsQuery {
	query.accountID = &accountID
	return query
}

func (query *AccountRecordsQuery) GetAccountID() AccountID {
	if query.accountID == nil {
		return AccountID{}
	}

	return *query.accountID
}

func (query *AccountRecordsQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.accountID != nil {
		if err := query.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *AccountRecordsQuery) _Build() *services.Query_CryptoGetAccountRecords {
	pb := services.Query_CryptoGetAccountRecords{
		CryptoGetAccountRecords: &services.CryptoGetAccountRecordsQuery{
			Header: &services.QueryHeader{},
		},
	}

	if query.accountID != nil {
		pb.CryptoGetAccountRecords.AccountID = query.accountID._ToProtobuf()
	}

	return &pb
}

func (query *AccountRecordsQuery) GetCost(client *Client) (Hbar, error) {
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
		_AccountRecordsQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_AccountRecordsQueryGetMethod,
		_AccountRecordsQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptoGetAccountRecords().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _AccountRecordsQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetCryptoGetAccountRecords().Header.NodeTransactionPrecheckCode))
}

func _AccountRecordsQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptoGetAccountRecords().Header.NodeTransactionPrecheckCode),
	}
}

func _AccountRecordsQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetAccountRecords,
	}
}

func (query *AccountRecordsQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.CryptoGetAccountRecords.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.CryptoGetAccountRecords.Header.ResponseType = services.ResponseType_ANSWER_ONLY
	return _ProtoRequest{
		query: &services.Query{
			Query: pb,
		},
	}
}

func (query *AccountRecordsQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.CryptoGetAccountRecords.Header.Payment = paymentTransaction
	pb.CryptoGetAccountRecords.Header.ResponseType = services.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &services.Query{
			Query: pb,
		},
	}, nil
}

func (query *AccountRecordsQuery) Execute(client *Client) ([]TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return []TransactionRecord{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return []TransactionRecord{}, err
	}

	records := make([]TransactionRecord, 0)

	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

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
			return []TransactionRecord{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []TransactionRecord{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "AccountRecordsQuery",
			}
		}

		cost = actualCost
	}

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return []TransactionRecord{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountRecordsQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountRecordsQueryGetMethod,
		_AccountRecordsQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return []TransactionRecord{}, err
	}

	for _, element := range resp.query.GetCryptoGetAccountRecords().Records {
		record := _TransactionRecordFromProtobuf(element)
		records = append(records, record)
	}

	return records, err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *AccountRecordsQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountRecordsQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *AccountRecordsQuery) SetQueryPayment(paymentAmount Hbar) *AccountRecordsQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountRecordsQuery.
func (query *AccountRecordsQuery) SetNodeAccountIDs(accountID []AccountID) *AccountRecordsQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *AccountRecordsQuery) SetMaxRetry(count int) *AccountRecordsQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *AccountRecordsQuery) SetMaxBackoff(max time.Duration) *AccountRecordsQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *AccountRecordsQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *AccountRecordsQuery) SetMinBackoff(min time.Duration) *AccountRecordsQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *AccountRecordsQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

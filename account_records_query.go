package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// AccountRecordsQuery gets all of the records for an account for any transfers into it and out of
// it, that were above the threshold, during the last 25 hours.
type AccountRecordsQuery struct {
	Query
	pb        *proto.CryptoGetAccountRecordsQuery
	accountID AccountID
}

// NewAccountRecordsQuery creates an AccountRecordsQuery query which can be used to construct and execute
// an AccountRecordsQuery.
//
// It is recommended that you use this for creating new instances of an AccountRecordQuery
// instead of manually creating an instance of the struct.
func NewAccountRecordsQuery() *AccountRecordsQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.CryptoGetAccountRecordsQuery{Header: &header}
	query.pb.Query = &proto.Query_CryptoGetAccountRecords{
		CryptoGetAccountRecords: &pb,
	}

	return &AccountRecordsQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetAccountID sets the account ID for which the records should be retrieved.
func (query *AccountRecordsQuery) SetAccountID(id AccountID) *AccountRecordsQuery {
	query.accountID = id
	return query
}

func (query *AccountRecordsQuery) GetAccountID() AccountID {
	return query.accountID
}

func (query *AccountRecordsQuery) validateNetworkOnIDs(client *Client) error {
	if !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.accountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *AccountRecordsQuery) build() *AccountRecordsQuery {
	if !query.accountID.isZero() {
		query.pb.AccountID = query.accountID.toProtobuf()
	}

	return query
}

func (query *AccountRecordsQuery) GetCost(client *Client) (Hbar, error) {
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

	err = query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	query.build()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		accountRecordsQuery_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		accountRecordsQuery_getMethod,
		accountRecordsQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptoGetAccountRecords().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func accountRecordsQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetCryptoGetAccountRecords().Header.NodeTransactionPrecheckCode))
}

func accountRecordsQuery_mapStatusError(_ request, response response, _ *NetworkName) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptoGetAccountRecords().Header.NodeTransactionPrecheckCode),
	}
}

func accountRecordsQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetAccountRecords,
	}
}

func (query *AccountRecordsQuery) Execute(client *Client) ([]TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return []TransactionRecord{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return []TransactionRecord{}, err
	}

	query.build()

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

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return []TransactionRecord{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		accountRecordsQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		accountRecordsQuery_getMethod,
		accountRecordsQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return []TransactionRecord{}, err
	}

	for _, element := range resp.query.GetCryptoGetAccountRecords().Records {
		record := transactionRecordFromProtobuf(element)
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

// SetNodeAccountIDs sets the node AccountID for this AccountRecordsQuery.
func (query *AccountRecordsQuery) SetNodeAccountIDs(accountID []AccountID) *AccountRecordsQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *AccountRecordsQuery) SetMaxRetry(count int) *AccountRecordsQuery {
	query.Query.SetMaxRetry(count)
	return query
}

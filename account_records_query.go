package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// AccountRecordsQuery gets all of the records for an account for any transfers into it and out of
// it, that were above the threshold, during the last 25 hours.
type AccountRecordsQuery struct {
	Query
	pb *proto.CryptoGetAccountRecordsQuery
}

// NewAccountRecordsQuery creates an AccountRecordsQuery query which can be used to construct and execute
// an AccountRecordsQuery.
//
// It is recommended that you use this for creating new instances of an AccountRecordQuery
// instead of manually creating an instance of the struct.
func NewAccountRecordsQuery() *AccountRecordsQuery {
	header := proto.QueryHeader{}
	query := newQuery(false, &header)
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
	query.pb.AccountID = id.toProtobuf()
	return query
}

func accountRecordsQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetCryptoGetAccountRecords().Header.NodeTransactionPrecheckCode)
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

	var records = []TransactionRecord{}

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		cost = client.maxQueryPayment

		// actualCost := CostQuery()
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return []TransactionRecord{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeId,
		accountRecordsQuery_getMethod,
		accountRecordsQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return []TransactionRecord{}, err
	}

	for _, element := range resp.query.GetCryptoGetAccountRecords().Records {
		records = append(records, TransactionRecordFromProtobuf(element))
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

func (query *AccountRecordsQuery) SetNodeAccountID(accountID AccountID) *AccountRecordsQuery {
	query.Query.SetNodeAccountID(accountID)
	return query
}

func (query *AccountRecordsQuery) GetNodeAccountId() AccountID {
	return query.Query.GetNodeAccountId()
}

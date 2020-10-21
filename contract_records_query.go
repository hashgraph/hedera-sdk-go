package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// ContractRecordsQuery retrieves all of the records for a smart contract instance, for any function call
// (or the constructor call) during the last 25 hours, for which a Record was requested.
type ContractRecordsQuery struct {
	Query
	pb *proto.ContractGetRecordsQuery
}

// NewContractRecordsQuery creates a ContractRecordsQuery query which can be used to construct and execute a
// Contract Get Records Query
func NewContractRecordsQuery() *ContractRecordsQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.ContractGetRecordsQuery{Header: &header}
	query.pb.Query = &proto.Query_ContractGetRecords{
		ContractGetRecords: &pb,
	}

	return &ContractRecordsQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetContractID sets the smart contract instance for which the records should be retrieved
func (query *ContractRecordsQuery) SetContractID(id ContractID) *ContractRecordsQuery {
	query.pb.ContractID = id.toProtobuf()
	return query
}

func (query *ContractRecordsQuery) GetContractID(id ContractID) ContractID {
	return contractIDFromProtobuf(query.pb.GetContractID())
}

func contractRecordsQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetContractGetRecordsResponse().Header.NodeTransactionPrecheckCode)
}

func contractRecordsQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getContract().GetTxRecordByContractID,
	}
}

func (query *ContractRecordsQuery) Execute(client *Client) ([]TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return []TransactionRecord{}, errNoClientProvided
	}

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
		query_getNodeAccountID,
		contractRecordsQuery_getMethod,
		contractRecordsQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return []TransactionRecord{}, err
	}

	rawRecords := resp.query.GetContractGetRecordsResponse().Records
	records := make([]TransactionRecord, len(rawRecords))

	for i, element := range resp.query.GetContractGetRecordsResponse().Records {
		records[i] = TransactionRecordFromProtobuf(element)
	}

	return records, nil

}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *ContractRecordsQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractRecordsQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *ContractRecordsQuery) SetQueryPayment(paymentAmount Hbar) *ContractRecordsQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *ContractRecordsQuery) SetNodeAccountIDs(accountID []AccountID) *ContractRecordsQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractRecordsQuery) GetNodeAccountIds() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

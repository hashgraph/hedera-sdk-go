package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ContractInfoQuery retrieves information about a smart contract instance. This includes the account that it uses, the
// file containing its bytecode, and the time when it will expire.
type ContractInfoQuery struct {
	Query
	pb *proto.ContractGetInfoQuery
}

// NewContractInfoQuery creates a ContractInfoQuery query which can be used to construct and execute a
// Contract Get Info Query.
func NewContractInfoQuery() *ContractInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.ContractGetInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_ContractGetInfo{
		ContractGetInfo: &pb,
	}

	return &ContractInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetContractID sets the contract for which information is requested
func (query *ContractInfoQuery) SetContractID(id ContractID) *ContractInfoQuery {
	query.pb.ContractID = id.toProtobuf()
	return query
}

func contractInfoQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetContractGetInfo().Header.NodeTransactionPrecheckCode)
}

func contractInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getContract().GetContractInfo,
	}
}

func (query *ContractInfoQuery) Execute(client *Client) (ContractInfo, error) {
	if client == nil || client.operator == nil {
		return ContractInfo{}, errNoClientProvided
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
		return ContractInfo{}, err
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
		contractInfoQuery_getMethod,
		contractInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return ContractInfo{}, err
	}

	return contractInfoFromProtobuf(resp.query.GetContractGetInfo().ContractInfo)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *ContractInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *ContractInfoQuery) SetQueryPayment(paymentAmount Hbar) *ContractInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *ContractInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ContractInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractInfoQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

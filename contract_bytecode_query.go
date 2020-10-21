package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ContractBytecodeQuery retrieves the bytecode for a smart contract instance
type ContractBytecodeQuery struct {
	Query
	pb *proto.ContractGetBytecodeQuery
}

// NewContractBytecodeQuery creates a ContractBytecodeQuery query which can be used to construct and execute a
// Contract Get Bytecode Query.
func NewContractBytecodeQuery() *ContractBytecodeQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.ContractGetBytecodeQuery{Header: &header}
	query.pb.Query = &proto.Query_ContractGetBytecode{
		ContractGetBytecode: &pb,
	}

	return &ContractBytecodeQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetContractID sets the contract for which the bytecode is requested
func (query *ContractBytecodeQuery) SetContractID(id ContractID) *ContractBytecodeQuery {
	query.pb.ContractID = id.toProtobuf()
	return query
}

func contractBytecodeQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetContractGetBytecodeResponse().Header.NodeTransactionPrecheckCode)
}

func contractBytecodeQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getContract().ContractGetBytecode,
	}
}

func (query *ContractBytecodeQuery) Execute(client *Client) ([]byte, error) {
	if client == nil || client.operator == nil {
		return []byte{}, errNoClientProvided
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
		return []byte{}, err
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
		contractBytecodeQuery_getMethod,
		contractBytecodeQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return []byte{}, err
	}

	return resp.query.GetContractGetBytecodeResponse().Bytecode, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *ContractBytecodeQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractBytecodeQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *ContractBytecodeQuery) SetQueryPayment(paymentAmount Hbar) *ContractBytecodeQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *ContractBytecodeQuery) SetNodeAccountIDs(accountID []AccountID) *ContractBytecodeQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractBytecodeQuery) GetNodeAccountId() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

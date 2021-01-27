package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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

func (query *ContractBytecodeQuery) GetContractID() ContractID {
	return contractIDFromProtobuf(query.pb.ContractID)
}

func (query *ContractBytecodeQuery) GetCost(client *Client) (Hbar, error) {
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
		query_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		contractBytecodeQuery_getMethod,
		contractBytecodeQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetContractGetBytecodeResponse().Header.Cost)
	return HbarFromTinybar(cost), nil
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
		return make([]byte, 0), errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

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
			return []byte{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []byte{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "ContractBytecodeQuery",
			}
		}

		cost = actualCost
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

// SetNodeAccountIDs sets the node AccountID for this ContractBytecodeQuery.
func (query *ContractBytecodeQuery) SetNodeAccountIDs(accountID []AccountID) *ContractBytecodeQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractBytecodeQuery) SetMaxRetry(count int) *ContractBytecodeQuery {
	query.Query.SetMaxRetry(count)
	return query
}

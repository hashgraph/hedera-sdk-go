package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// ContractInfoQuery retrieves information about a smart contract instance. This includes the account that it uses, the
// file containing its bytecode, and the time when it will expire.
type ContractInfoQuery struct {
	Query
	contractID ContractID
}

// NewContractInfoQuery creates a ContractInfoQuery query which can be used to construct and execute a
// Contract Get Info Query.
func NewContractInfoQuery() *ContractInfoQuery {
	query := newQuery(true)

	query.SetMaxQueryPayment(NewHbar(2))

	return &ContractInfoQuery{
		Query: query,
	}
}

// SetContractID sets the contract for which information is requested
func (query *ContractInfoQuery) SetContractID(id ContractID) *ContractInfoQuery {
	query.contractID = id
	return query
}

func (query *ContractInfoQuery) GetContractID() ContractID {
	return query.contractID
}

func (query *ContractInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.contractID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *ContractInfoQuery) build() *proto.Query_ContractGetInfo {
	body := &proto.ContractGetInfoQuery{Header: &proto.QueryHeader{}}
	if !query.contractID.isZero() {
		body.ContractID = query.contractID.toProtobuf()
	}

	return &proto.Query_ContractGetInfo{
		ContractGetInfo: body,
	}
}

func (query *ContractInfoQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ContractGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ContractGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *ContractInfoQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.ContractGetInfo.Header.Payment = paymentTransaction
	pb.ContractGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *ContractInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoResp, err := query.costQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		contractInfoQuery_shouldRetry,
		protoResp,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		contractInfoQuery_getMethod,
		contractInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetContractGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func contractInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetContractGetInfo().Header.NodeTransactionPrecheckCode))
}

func contractInfoQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetContractGetInfo().Header.NodeTransactionPrecheckCode),
	}
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

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return ContractInfo{}, err
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
			return ContractInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return ContractInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "ContractInfoQuery",
			}
		}

		cost = actualCost
	}

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return ContractInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		contractInfoQuery_shouldRetry,
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		contractInfoQuery_getMethod,
		contractInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return ContractInfo{}, err
	}

	info, err := contractInfoFromProtobuf(resp.query.GetContractGetInfo().ContractInfo)
	if err != nil {
		return ContractInfo{}, err
	}

	return info, nil
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

// SetNodeAccountIDs sets the node AccountID for this ContractInfoQuery.
func (query *ContractInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ContractInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractInfoQuery) SetMaxRetry(count int) *ContractInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// ContractCallQuery calls a function of the given smart contract instance, giving it ContractFunctionParameters as its
// inputs. It will consume the entire given amount of gas.
//
// This is performed locally on the particular node that the client is communicating with. It cannot change the state of
// the contract instance (and so, cannot spend anything from the instance's Hedera account). It will not have a
// consensus timestamp. It cannot generate a record or a receipt. This is useful for calling getter functions, which
// purely read the state and don't change it. It is faster and cheaper than a ContractExecuteTransaction, because it is
// purely local to a single  node.
type ContractCallQuery struct {
	Query
	contractID         ContractID
	gas                uint64
	maxResultSize      uint64
	functionParameters []byte
}

// NewContractCallQuery creates a ContractCallQuery query which can be used to construct and execute a
// Contract Call Local Query.
func NewContractCallQuery() *ContractCallQuery {
	query := newQuery(true)
	query.SetMaxQueryPayment(NewHbar(2))

	return &ContractCallQuery{
		Query: query,
	}
}

// SetContractID sets the contract instance to call
func (query *ContractCallQuery) SetContractID(id ContractID) *ContractCallQuery {
	query.contractID = id
	return query
}

func (query *ContractCallQuery) GetContractID() ContractID {
	return query.contractID
}

// SetGas sets the amount of gas to use for the call. All of the gas offered will be charged for.
func (query *ContractCallQuery) SetGas(gas uint64) *ContractCallQuery {
	query.gas = gas
	return query
}

func (query *ContractCallQuery) GetGas() uint64 {
	return query.gas
}

// SetMaxResultSize sets the max number of bytes that the result might include. The run will fail if it would have
// returned more than this number of bytes.
func (query *ContractCallQuery) SetMaxResultSize(size uint64) *ContractCallQuery {
	query.maxResultSize = size
	return query
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (query *ContractCallQuery) SetFunction(name string, params *ContractFunctionParameters) *ContractCallQuery {
	if params == nil {
		params = NewContractFunctionParameters()
	}

	query.functionParameters = params.build(&name)
	return query
}

func (query *ContractCallQuery) SetFunctionParameters(byteArray []byte) *ContractCallQuery {
	query.functionParameters = byteArray
	return query
}

func (query *ContractCallQuery) GetFunctionParameters() []byte {
	return query.functionParameters
}

func (query *ContractCallQuery) validateNetworkOnIDs(client *Client) error {
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

func (query *ContractCallQuery) build() *proto.Query_ContractCallLocal {
	body := &proto.ContractCallLocalQuery{
		Header:        &proto.QueryHeader{},
		Gas:           int64(query.gas),
		MaxResultSize: int64(query.maxResultSize),
	}
	if !query.contractID.isZero() {
		body.ContractID = query.contractID.toProtobuf()
	}
	if len(query.functionParameters) > 0 {
		body.FunctionParameters = query.functionParameters
	}

	return &proto.Query_ContractCallLocal{
		ContractCallLocal: body,
	}
}

func (query *ContractCallQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ContractCallLocal.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ContractCallLocal.Header.ResponseType = proto.ResponseType_ANSWER_ONLY
	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *ContractCallQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.ContractCallLocal.Header.Payment = paymentTransaction
	pb.ContractCallLocal.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *ContractCallQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query.costQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		contractCallQuery_shouldRetry,
		protoReq,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		contractCallQuery_getMethod,
		contractCallQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetContractCallLocal().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func contractCallQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetContractCallLocal().Header.NodeTransactionPrecheckCode))
}

func contractCallQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetContractCallLocal().Header.NodeTransactionPrecheckCode),
	}
}

func contractCallQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getContract().ContractCallLocalMethod,
	}
}

func (query *ContractCallQuery) Execute(client *Client) (ContractFunctionResult, error) {
	if client == nil || client.operator == nil {
		return ContractFunctionResult{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return ContractFunctionResult{}, err
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
			return ContractFunctionResult{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return ContractFunctionResult{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "ContractCallQuery",
			}
		}

		cost = actualCost
	}

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return ContractFunctionResult{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		contractCallQuery_shouldRetry,
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		contractCallQuery_getMethod,
		contractCallQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return ContractFunctionResult{}, err
	}

	result := contractFunctionResultFromProtobuf(resp.query.GetContractCallLocal().FunctionResult)
	if result.ContractID != nil {
		result.ContractID.setNetworkWithClient(client)
	}

	return result, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *ContractCallQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractCallQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *ContractCallQuery) SetQueryPayment(paymentAmount Hbar) *ContractCallQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

// SetNodeAccountIDs sets the node AccountID for this ContractCallQuery.
func (query *ContractCallQuery) SetNodeAccountIDs(accountID []AccountID) *ContractCallQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractCallQuery) SetMaxRetry(count int) *ContractCallQuery {
	query.Query.SetMaxRetry(count)
	return query
}

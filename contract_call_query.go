package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// ContractCallQuery calls a function of the given smart contract instance, giving it ContractFunctionParameters as its
// inputs. It will consume the entire given amount of gas.
//
// This is performed locally on the particular _Node that the client is communicating with. It cannot change the state of
// the contract instance (and so, cannot spend anything from the instance's Hedera account). It will not have a
// consensus timestamp. It cannot generate a record or a receipt. This is useful for calling getter functions, which
// purely read the state and don't change it. It is faster and cheaper than a ContractExecuteTransaction, because it is
// purely local to a single  _Node.
type ContractCallQuery struct {
	Query
	contractID         *ContractID
	gas                uint64
	maxResultSize      uint64
	functionParameters []byte
}

// NewContractCallQuery creates a ContractCallQuery query which can be used to construct and execute a
// Contract Call Local Query.
func NewContractCallQuery() *ContractCallQuery {
	query := _NewQuery(true)
	query.SetMaxQueryPayment(NewHbar(2))

	return &ContractCallQuery{
		Query: query,
	}
}

// SetContractID sets the contract instance to call
func (query *ContractCallQuery) SetContractID(contractID ContractID) *ContractCallQuery {
	query.contractID = &contractID
	return query
}

func (query *ContractCallQuery) GetContractID() ContractID {
	if query.contractID == nil {
		return ContractID{}
	}

	return *query.contractID
}

// SetGas sets the amount of gas to use for the call. All of the gas offered will be charged for.
func (query *ContractCallQuery) SetGas(gas uint64) *ContractCallQuery {
	query.gas = gas
	return query
}

func (query *ContractCallQuery) GetGas() uint64 {
	return query.gas
}

// Deprecated
func (query *ContractCallQuery) SetMaxResultSize(size uint64) *ContractCallQuery {
	query.maxResultSize = size
	return query
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (query *ContractCallQuery) SetFunction(name string, params *ContractFunctionParameters) *ContractCallQuery {
	if params == nil {
		params = NewContractFunctionParameters()
	}

	query.functionParameters = params._Build(&name)
	return query
}

func (query *ContractCallQuery) SetFunctionParameters(byteArray []byte) *ContractCallQuery {
	query.functionParameters = byteArray
	return query
}

func (query *ContractCallQuery) GetFunctionParameters() []byte {
	return query.functionParameters
}

func (query *ContractCallQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.contractID != nil {
		if err := query.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *ContractCallQuery) _Build() *proto.Query_ContractCallLocal {
	pb := proto.Query_ContractCallLocal{
		ContractCallLocal: &proto.ContractCallLocalQuery{
			Header:        &proto.QueryHeader{},
			Gas:           int64(query.gas),
		},
	}

	if query.contractID != nil {
		pb.ContractCallLocal.ContractID = query.contractID._ToProtobuf()
	}

	if len(query.functionParameters) > 0 {
		pb.ContractCallLocal.FunctionParameters = query.functionParameters
	}

	return &pb
}

func (query *ContractCallQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ContractCallLocal.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ContractCallLocal.Header.ResponseType = proto.ResponseType_ANSWER_ONLY
	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *ContractCallQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.ContractCallLocal.Header.Payment = paymentTransaction
	pb.ContractCallLocal.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *ContractCallQuery) GetCost(client *Client) (Hbar, error) {
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
		_ContractCallQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_ContractCallQueryGetMethod,
		_ContractCallQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetContractCallLocal().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _ContractCallQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetContractCallLocal().Header.NodeTransactionPrecheckCode))
}

func _ContractCallQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetContractCallLocal().Header.NodeTransactionPrecheckCode),
	}
}

func _ContractCallQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().ContractCallLocalMethod,
	}
}

func (query *ContractCallQuery) Execute(client *Client) (ContractFunctionResult, error) {
	if client == nil || client.operator == nil {
		return ContractFunctionResult{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	err := query._ValidateNetworkOnIDs(client)
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

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return ContractFunctionResult{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ContractCallQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ContractCallQueryGetMethod,
		_ContractCallQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return ContractFunctionResult{}, err
	}

	return _ContractFunctionResultFromProtobuf(resp.query.GetContractCallLocal().FunctionResult), nil
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

// SetNodeAccountIDs sets the _Node AccountID for this ContractCallQuery.
func (query *ContractCallQuery) SetNodeAccountIDs(accountID []AccountID) *ContractCallQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractCallQuery) SetMaxRetry(count int) *ContractCallQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *ContractCallQuery) SetMaxBackoff(max time.Duration) *ContractCallQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *ContractCallQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *ContractCallQuery) SetMinBackoff(min time.Duration) *ContractCallQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *ContractCallQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

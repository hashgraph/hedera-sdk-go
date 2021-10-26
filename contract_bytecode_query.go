package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// ContractBytecodeQuery retrieves the bytecode for a smart contract instance
type ContractBytecodeQuery struct {
	Query
	contractID *ContractID
}

// NewContractBytecodeQuery creates a ContractBytecodeQuery query which can be used to construct and execute a
// Contract Get Bytecode Query.
func NewContractBytecodeQuery() *ContractBytecodeQuery {
	return &ContractBytecodeQuery{
		Query: _NewQuery(true),
	}
}

// SetContractID sets the contract for which the bytecode is requested
func (query *ContractBytecodeQuery) SetContractID(contractID ContractID) *ContractBytecodeQuery {
	query.contractID = &contractID
	return query
}

func (query *ContractBytecodeQuery) GetContractID() ContractID {
	if query.contractID == nil {
		return ContractID{}
	}

	return *query.contractID
}

func (query *ContractBytecodeQuery) _ValidateNetworkOnIDs(client *Client) error {
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

func (query *ContractBytecodeQuery) _Build() *proto.Query_ContractGetBytecode {
	pb := proto.Query_ContractGetBytecode{
		ContractGetBytecode: &proto.ContractGetBytecodeQuery{
			Header: &proto.QueryHeader{},
		},
	}

	if query.contractID != nil {
		pb.ContractGetBytecode.ContractID = query.contractID._ToProtobuf()
	}

	return &pb
}

func (query *ContractBytecodeQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ContractGetBytecode.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ContractGetBytecode.Header.ResponseType = proto.ResponseType_ANSWER_ONLY
	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *ContractBytecodeQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.ContractGetBytecode.Header.Payment = paymentTransaction
	pb.ContractGetBytecode.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *ContractBytecodeQuery) GetCost(client *Client) (Hbar, error) {
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
		_ContractBytecodeQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_ContractBytecodeQueryGetMethod,
		_ContractBytecodeQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetContractGetBytecodeResponse().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _ContractBytecodeQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetContractGetBytecodeResponse().Header.NodeTransactionPrecheckCode))
}

func _ContractBytecodeQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetContractGetBytecodeResponse().Header.NodeTransactionPrecheckCode),
	}
}

func _ContractBytecodeQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().ContractGetBytecode,
	}
}

func (query *ContractBytecodeQuery) Execute(client *Client) ([]byte, error) {
	if client == nil || client.operator == nil {
		return make([]byte, 0), errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return []byte{}, err
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

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return []byte{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ContractBytecodeQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ContractBytecodeQueryGetMethod,
		_ContractBytecodeQueryMapStatusError,
		_QueryMapResponse,
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

// SetNodeAccountIDs sets the _Node AccountID for this ContractBytecodeQuery.
func (query *ContractBytecodeQuery) SetNodeAccountIDs(accountID []AccountID) *ContractBytecodeQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractBytecodeQuery) SetMaxRetry(count int) *ContractBytecodeQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *ContractBytecodeQuery) SetMaxBackoff(max time.Duration) *ContractBytecodeQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *ContractBytecodeQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *ContractBytecodeQuery) SetMinBackoff(min time.Duration) *ContractBytecodeQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *ContractBytecodeQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

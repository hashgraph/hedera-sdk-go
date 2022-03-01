package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// ContractBytecodeQuery retrieves the bytecode for a smart contract instance
type ContractBytecodeQuery struct {
	Query
	contractID *ContractID
}

// NewContractBytecodeQuery creates a ContractBytecodeQuery query which can be used to construct and execute a
// Contract Get Bytecode Query.
func NewContractBytecodeQuery() *ContractBytecodeQuery {
	header := services.QueryHeader{}
	return &ContractBytecodeQuery{
		Query: _NewQuery(true, &header),
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

func (query *ContractBytecodeQuery) _Build() *services.Query_ContractGetBytecode {
	pb := services.Query_ContractGetBytecode{
		ContractGetBytecode: &services.ContractGetBytecodeQuery{
			Header: &services.QueryHeader{},
		},
	}

	if query.contractID != nil {
		pb.ContractGetBytecode.ContractID = query.contractID._ToProtobuf()
	}

	return &pb
}

func (query *ContractBytecodeQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return Hbar{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range query.nodeAccountIDs {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.ContractGetBytecode.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ContractBytecodeQueryShouldRetry,
		_CostQueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ContractBytecodeQueryGetMethod,
		_ContractBytecodeQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetContractGetBytecodeResponse().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _ContractBytecodeQueryShouldRetry(logID string, _ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(logID, Status(response.query.GetContractGetBytecodeResponse().Header.NodeTransactionPrecheckCode))
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

	var err error
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return []byte{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}
	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return []byte{}, err
	}

	if !query.lockedTransactionID {
		query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)
	}

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

	query.nextPaymentTransactionIndex = 0
	query.paymentTransactions = make([]*services.Transaction, 0)

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return []byte{}, err
	}

	pb := query._Build()
	pb.ContractGetBytecode.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ContractBytecodeQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ContractBytecodeQueryGetMethod,
		_ContractBytecodeQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
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

func (query *ContractBytecodeQuery) _GetLogID() string {
	timestamp := query.paymentTransactionID.ValidStart
	return fmt.Sprintf("ContractBytecodeQuery:%d", timestamp.UnixNano())
}

func (query *ContractBytecodeQuery) SetPaymentTransactionID(transactionID TransactionID) *ContractBytecodeQuery {
	if query.lockedTransactionID {
		panic("payment TransactionID is locked")
	}
	query.lockedTransactionID = true
	query.paymentTransactionID = transactionID
	return query
}

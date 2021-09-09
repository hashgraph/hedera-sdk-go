package hedera

import (
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// ContractInfoQuery retrieves information about a smart contract instance. This includes the account that it uses, the
// file containing its bytecode, and the time when it will expire.
type ContractInfoQuery struct {
	Query
	contractID *ContractID
}

// NewContractInfoQuery creates a ContractInfoQuery query which can be used to construct and execute a
// Contract Get Info Query.
func NewContractInfoQuery() *ContractInfoQuery {
	query := _NewQuery(true)

	query.SetMaxQueryPayment(NewHbar(2))

	return &ContractInfoQuery{
		Query: query,
	}
}

// SetContractID sets the contract for which information is requested
func (query *ContractInfoQuery) SetContractID(contractID ContractID) *ContractInfoQuery {
	query.contractID = &contractID
	return query
}

func (query *ContractInfoQuery) GetContractID() ContractID {
	if query.contractID == nil {
		return ContractID{}
	}

	return *query.contractID
}

func (query *ContractInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.contractID != nil {
		if err := query.contractID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *ContractInfoQuery) _Build() *proto.Query_ContractGetInfo {
	pb := proto.Query_ContractGetInfo{
		ContractGetInfo: &proto.ContractGetInfoQuery{
			Header: &proto.QueryHeader{},
		},
	}

	if query.contractID != nil {
		pb.ContractGetInfo.ContractID = query.contractID._ToProtobuf()
	}

	return &pb
}

func (query *ContractInfoQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	_ = query._BuildAllPaymentTransactions()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ContractGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ContractGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *ContractInfoQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionIDGenerate(client.GetOperatorAccountID()), AccountID{}, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	paymentBytes, err := protobuf.Marshal(paymentTransaction)
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.ContractGetInfo.Header.Payment = &proto.Transaction{
		SignedTransactionBytes: paymentBytes,
	}
	pb.ContractGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *ContractInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.nodeIDs = client.network._GetNodeAccountIDsForExecute()
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoResp, err := query._CostQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ContractInfoQueryShouldRetry,
		protoResp,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_ContractInfoQueryGetMethod,
		_ContractInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetContractGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}

	return HbarFromTinybar(cost), nil
}

func _ContractInfoQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetContractGetInfo().Header.NodeTransactionPrecheckCode))
}

func _ContractInfoQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetContractGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _ContractInfoQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().GetContractInfo,
	}
}

func (query *ContractInfoQuery) Execute(client *Client) (ContractInfo, error) {
	if client == nil || client.operator == nil {
		return ContractInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.nodeIDs = client.network._GetNodeAccountIDsForExecute()
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return ContractInfo{}, err
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

	query.actualCost = cost

	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return ContractInfo{}, err
		}
	}

	transactionID := query.paymentTransactionID

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		query.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_ContractInfoQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ContractInfoQueryGetMethod,
		_ContractInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return ContractInfo{}, err
	}

	info, err := _ContractInfoFromProtobuf(resp.query.GetContractGetInfo().ContractInfo)
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

// SetNodeAccountIDs sets the _Node AccountID for this ContractInfoQuery.
func (query *ContractInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ContractInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *ContractInfoQuery) SetMaxRetry(count int) *ContractInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *ContractInfoQuery) SetMaxBackoff(max time.Duration) *ContractInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *ContractInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *ContractInfoQuery) SetMinBackoff(min time.Duration) *ContractInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *ContractInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *ContractInfoQuery) IsFrozen() bool {
	return query._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (query *ContractInfoQuery) Sign(
	privateKey PrivateKey,
) *ContractInfoQuery {
	return query.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (query *ContractInfoQuery) SignWithOperator(
	client *Client,
) (*ContractInfoQuery, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return query, err
		}
	}
	return query.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (query *ContractInfoQuery) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractInfoQuery {
	if !query._KeyAlreadySigned(publicKey) {
		query._SignWith(publicKey, signer)
	}

	return query
}

func (query *ContractInfoQuery) Freeze() (*ContractInfoQuery, error) {
	return query.FreezeWith(nil)
}

func (query *ContractInfoQuery) FreezeWith(client *Client) (*ContractInfoQuery, error) {
	if query.IsFrozen() {
		return query, nil
	}
	if query.actualCost.AsTinybar() == 0 {
		if query.queryPayment.tinybar != 0 {
			query.actualCost = query.queryPayment
		} else {
			if query.maxQueryPayment.tinybar == 0 {
				query.actualCost = client.maxQueryPayment
			} else {
				query.actualCost = query.maxQueryPayment
			}

			actualCost, err := query.GetCost(client)
			if err != nil {
				return &ContractInfoQuery{}, err
			}

			if query.actualCost.tinybar < actualCost.tinybar {
				return &ContractInfoQuery{}, ErrMaxQueryPaymentExceeded{
					QueryCost:       actualCost,
					MaxQueryPayment: query.actualCost,
					query:           "ContractInfoQuery",
				}
			}

			query.actualCost = actualCost
		}
	}
	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ContractInfoQuery{}, err
	}
	if err := query._InitPaymentTransactionID(client); err != nil {
		return query, err
	}

	return query, _QueryGeneratePayments(&query.Query, query.actualCost)
}

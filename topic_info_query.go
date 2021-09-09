package hedera

import (
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TopicInfoQuery struct {
	Query
	topicID *TopicID
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//  Get Topic Info Query.
func NewTopicInfoQuery() *TopicInfoQuery {
	return &TopicInfoQuery{
		Query: _NewQuery(true),
	}
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (query *TopicInfoQuery) SetTopicID(topicID TopicID) *TopicInfoQuery {
	query.topicID = &topicID
	return query
}

func (query *TopicInfoQuery) GetTopicID() TopicID {
	if query.topicID == nil {
		return TopicID{}
	}

	return *query.topicID
}

func (query *TopicInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.topicID != nil {
		if err := query.topicID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *TopicInfoQuery) _Build() *proto.Query_ConsensusGetTopicInfo {
	body := &proto.ConsensusGetTopicInfoQuery{
		Header: &proto.QueryHeader{},
	}

	if query.topicID != nil {
		body.TopicID = query.topicID._ToProtobuf()
	}

	return &proto.Query_ConsensusGetTopicInfo{
		ConsensusGetTopicInfo: body,
	}
}

func (query *TopicInfoQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	_ = query._BuildAllPaymentTransactions()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.ConsensusGetTopicInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.ConsensusGetTopicInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TopicInfoQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionIDGenerate(client.GetOperatorAccountID()), AccountID{}, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	paymentBytes, err := protobuf.Marshal(paymentTransaction)
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.ConsensusGetTopicInfo.Header.Payment = &proto.Transaction{
		SignedTransactionBytes: paymentBytes,
	}
	pb.ConsensusGetTopicInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TopicInfoQuery) GetCost(client *Client) (Hbar, error) {
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

	protoReq, err := query._CostQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TopicInfoQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_TopicInfoQueryGetMethod,
		_TopicInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetConsensusGetTopicInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _TopicInfoQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetConsensusGetTopicInfo().Header.NodeTransactionPrecheckCode))
}

func _TopicInfoQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetConsensusGetTopicInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _TopicInfoQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetTopic().GetTopicInfo,
	}
}

// Execute executes the TopicInfoQuery using the provided client
func (query *TopicInfoQuery) Execute(client *Client) (TopicInfo, error) {
	if client == nil || client.operator == nil {
		return TopicInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.nodeIDs = client.network._GetNodeAccountIDsForExecute()
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return TopicInfo{}, err
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
			return TopicInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TopicInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TopicInfoQuery",
			}
		}

		cost = actualCost
	}

	query.actualCost = cost

	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return TopicInfo{}, err
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
		_TopicInfoQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TopicInfoQueryGetMethod,
		_TopicInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return TopicInfo{}, err
	}

	return _TopicInfoFromProtobuf(resp.query.GetConsensusGetTopicInfo().TopicInfo)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TopicInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TopicInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *TopicInfoQuery) SetQueryPayment(paymentAmount Hbar) *TopicInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *TopicInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TopicInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TopicInfoQuery) SetMaxRetry(count int) *TopicInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *TopicInfoQuery) SetMaxBackoff(max time.Duration) *TopicInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *TopicInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *TopicInfoQuery) SetMinBackoff(min time.Duration) *TopicInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *TopicInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *TopicInfoQuery) IsFrozen() bool {
	return query._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (query *TopicInfoQuery) Sign(
	privateKey PrivateKey,
) *TopicInfoQuery {
	return query.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (query *TopicInfoQuery) SignWithOperator(
	client *Client,
) (*TopicInfoQuery, error) {
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
func (query *TopicInfoQuery) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TopicInfoQuery {
	if !query._KeyAlreadySigned(publicKey) {
		query._SignWith(publicKey, signer)
	}

	return query
}

func (query *TopicInfoQuery) Freeze() (*TopicInfoQuery, error) {
	return query.FreezeWith(nil)
}

func (query *TopicInfoQuery) FreezeWith(client *Client) (*TopicInfoQuery, error) {
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
				return &TopicInfoQuery{}, err
			}

			if query.actualCost.tinybar < actualCost.tinybar {
				return &TopicInfoQuery{}, ErrMaxQueryPaymentExceeded{
					QueryCost:       actualCost,
					MaxQueryPayment: query.actualCost,
					query:           "TopicInfoQuery",
				}
			}

			query.actualCost = actualCost
		}
	}
	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TopicInfoQuery{}, err
	}
	if err = query._InitPaymentTransactionID(client); err != nil {
		return query, err
	}

	return query, _QueryGeneratePayments(&query.Query, query.actualCost)
}

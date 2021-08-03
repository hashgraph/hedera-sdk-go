package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type LiveHashQuery struct {
	Query
	pb        *proto.CryptoGetLiveHashQuery
	accountID AccountID
}

func NewLiveHashQuery() *LiveHashQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.CryptoGetLiveHashQuery{Header: &header}
	query.pb.Query = &proto.Query_CryptoGetLiveHash{
		CryptoGetLiveHash: &pb,
	}

	return &LiveHashQuery{
		Query: query,
		pb:    &pb,
	}
}

func (query *LiveHashQuery) SetAccountID(id AccountID) *LiveHashQuery {
	query.pb.AccountID = id.toProtobuf()
	return query
}

func (query *LiveHashQuery) GetAccountID() AccountID {
	return query.accountID
}

func (query *LiveHashQuery) SetHash(hash []byte) *LiveHashQuery {
	query.pb.Hash = hash
	return query
}

func (query *LiveHashQuery) GetGetHash() []byte {
	return query.pb.Hash
}

func (query *LiveHashQuery) validateChecksums(client *Client) error {
	var err error
	err = query.accountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *LiveHashQuery) build() *LiveHashQuery {
	if !query.accountID.isZero() {
		query.pb.AccountID = query.accountID.toProtobuf()
	}

	return query
}

func (query *LiveHashQuery) GetCost(client *Client) (Hbar, error) {
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

	err = query.validateChecksums(client)
	if err != nil {
		return Hbar{}, err
	}

	query.build()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		liveHashQuery_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		liveHashQuery_getMethod,
		liveHashQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptoGetLiveHash().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func liveHashQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode))
}

func liveHashQuery_mapStatusError(_ request, response response, _ *NetworkName) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode),
	}
}

func liveHashQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetLiveHash,
	}
}

func (query *LiveHashQuery) Execute(client *Client) (LiveHash, error) {
	if client == nil || client.operator == nil {
		return LiveHash{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateChecksums(client)
	if err != nil {
		return LiveHash{}, err
	}

	query.build()

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
			return LiveHash{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return LiveHash{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "LiveHashQuery",
			}
		}

		cost = actualCost
	}

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return LiveHash{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		liveHashQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		liveHashQuery_getMethod,
		liveHashQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return LiveHash{}, err
	}

	liveHash, err := liveHashFromProtobuf(resp.query.GetCryptoGetLiveHash().LiveHash, client.networkName)
	if err != nil {
		return LiveHash{}, err
	}

	return liveHash, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *LiveHashQuery) SetMaxQueryPayment(maxPayment Hbar) *LiveHashQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *LiveHashQuery) SetQueryPayment(paymentAmount Hbar) *LiveHashQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *LiveHashQuery) SetNodeAccountIDs(accountID []AccountID) *LiveHashQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

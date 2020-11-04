package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type LiveHashQuery struct {
	Query
	pb *proto.CryptoGetLiveHashQuery
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
	return accountIDFromProtobuf(query.pb.GetAccountID())
}

func (query *LiveHashQuery) SetHash(hash []byte) *LiveHashQuery {
	query.pb.Hash = hash
	return query
}

func (query *LiveHashQuery) GetGetHash() []byte {
	return query.pb.Hash
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

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		accountInfoQuery_getMethod,
		accountInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptoGetLiveHash().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func liveHashQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode)
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

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		cost = client.maxQueryPayment

		actualCost, err := query.GetCost(client)
		if err != nil {
			return LiveHash{}, err
		}

		if cost.tinybar > actualCost.tinybar {
			return LiveHash{}, ErrMaxQueryPaymentExceeded{}
		}

		cost = actualCost
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return LiveHash{}, err
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
		liveHashQuery_getMethod,
		liveHashQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return LiveHash{}, err
	}

	return liveHashFromProtobuf(resp.query.GetCryptoGetLiveHash().LiveHash)
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

func (query *LiveHashQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

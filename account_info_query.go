package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type AccountInfoQuery struct {
	Query
	pb *proto.CryptoGetInfoQuery
}

func NewAccountInfoQuery() *AccountInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.CryptoGetInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_CryptoGetInfo{
		CryptoGetInfo: &pb,
	}

	return &AccountInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func accountInfoQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetCryptoGetInfo().Header.NodeTransactionPrecheckCode)
}

func accountInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetAccountInfo,
	}
}

func (query *AccountInfoQuery) SetAccountID(accountID AccountID) *AccountInfoQuery {
	query.pb.AccountID = accountID.toProtobuf()
	return query
}

func (query *AccountInfoQuery) SetNodeAccountIDs(accountID []AccountID) *AccountInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *AccountInfoQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

func (query *AccountInfoQuery) SetQueryPayment(queryPayment Hbar) *AccountInfoQuery {
	query.queryPayment = queryPayment
	return query
}

func (query *AccountInfoQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *AccountInfoQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *AccountInfoQuery) Execute(client *Client) (AccountInfo, error) {
	if client == nil || client.operator == nil {
		return AccountInfo{}, errNoClientProvided
	}

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		cost = client.maxQueryPayment

		// actualCost := CostQuery()
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return AccountInfo{}, err
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
		accountInfoQuery_getMethod,
		accountInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return AccountInfo{}, err
	}

	return accountInfoFromProtobuf(resp.query.GetCryptoGetInfo().AccountInfo)
}

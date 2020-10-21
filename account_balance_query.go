package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
type AccountBalanceQuery struct {
	Query
	pb *proto.CryptoGetAccountBalanceQuery
}

// NewAccountBalanceQuery creates an AccountBalanceQuery query which can be used to construct and execute
// an AccountBalanceQuery.
// It is recommended that you use this for creating new instances of an AccountBalanceQuery
// instead of manually creating an instance of the struct.
func NewAccountBalanceQuery() *AccountBalanceQuery {
	header := proto.QueryHeader{}
	query := newQuery(false, &header)
	pb := proto.CryptoGetAccountBalanceQuery{Header: &header}
	query.pb.Query = &proto.Query_CryptogetAccountBalance{
		CryptogetAccountBalance: &pb,
	}

	return &AccountBalanceQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetAccountID sets the AccountID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this method.
func (query *AccountBalanceQuery) SetAccountID(id AccountID) *AccountBalanceQuery {
	query.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_AccountID{
		AccountID: id.toProtobuf(),
	}

	return query
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this method.
func (query *AccountBalanceQuery) SetContractID(id ContractID) *AccountBalanceQuery {
	query.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_ContractID{
		ContractID: id.toProtobuf(),
	}

	return query
}

func accountBalanceQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode)
}

func accountBalanceQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().CryptoGetBalance,
	}
}

func (query *AccountBalanceQuery) Execute(client *Client) (AccountBalance, error) {
	if client == nil || client.operator == nil {
		return AccountBalance{}, errNoClientProvided
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeId,
		accountBalanceQuery_getMethod,
		accountBalanceQuery_mapResponseStatus,
		query_mapResponse,
	)

	//for _, id := range transaction.nodeIDs {
	//	transaction.pbBody.NodeAccountID = id.toProtobuf()
	//	bodyBytes, err := protobuf.Marshal(transaction.pbBody)
	//	if err != nil {
	//		// This should be unreachable
	//		// From the documentation this appears to only be possible if there are missing proto types
	//		panic(err)
	//	}
	//
	//	sigmap := proto.SignatureMap{
	//		SigPair: make([]*proto.SignaturePair, 0),
	//	}
	//	transaction.signatures = append(transaction.signatures, &sigmap)
	//	transaction.transactions = append(transaction.transactions, &proto.Transaction{
	//		BodyBytes: bodyBytes,
	//		SigMap:    &sigmap,
	//	})
	//}

	if err != nil {
		return AccountBalance{}, err
	}

	var tokens []TokenBalance
	for i, token := range resp.query.GetCryptogetAccountBalance().TokenBalances {
		tokens[i] = tokenBalancesFromProtobuf(token)
	}

	return AccountBalance{
		Hbar:  HbarFromTinybar(int64(resp.query.GetCryptogetAccountBalance().Balance)),
		Token: &tokens,
	}, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *AccountBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountBalanceQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *AccountBalanceQuery) SetQueryPayment(paymentAmount Hbar) *AccountBalanceQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *AccountBalanceQuery) SetNodeAccountID(accountID AccountID) *AccountBalanceQuery {
	query.Query.SetNodeAccountID(accountID)
	return query
}

func (query *AccountBalanceQuery) GetNodeAccountId() AccountID {
	return query.Query.GetNodeAccountId()
}

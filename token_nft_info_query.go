package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenNftInfoQuery
// Applicable only to tokens of type NON_FUNGIBLE_UNIQUE.
// Gets info on a NFT for a given TokenID (of type NON_FUNGIBLE_UNIQUE) and serial number
type TokenNftInfoQuery struct {
	Query
	nftID *NftID
}

// NewTokenNftInfoQuery creates TokenNftInfoQuery which
// gets info on a NFT for a given TokenID (of type NON_FUNGIBLE_UNIQUE) and serial number
// Applicable only to tokens of type NON_FUNGIBLE_UNIQUE.
func NewTokenNftInfoQuery() *TokenNftInfoQuery {
	header := services.QueryHeader{}
	return &TokenNftInfoQuery{
		Query: _NewQuery(true, &header),
		nftID: nil,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *TokenNftInfoQuery) SetGrpcDeadline(deadline *time.Duration) *TokenNftInfoQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetNftID Sets the ID of the NFT
func (q *TokenNftInfoQuery) SetNftID(nftID NftID) *TokenNftInfoQuery {
	q.nftID = &nftID
	return q
}

// GetNftID returns the ID of the NFT
func (q *TokenNftInfoQuery) GetNftID() NftID {
	if q.nftID == nil {
		return NftID{}
	}

	return *q.nftID
}

// Deprecated
func (q *TokenNftInfoQuery) SetTokenID(id TokenID) *TokenNftInfoQuery {
	return q
}

// Deprecated
func (q *TokenNftInfoQuery) GetTokenID() TokenID {
	return TokenID{}
}

// Deprecated
func (q *TokenNftInfoQuery) SetAccountID(id AccountID) *TokenNftInfoQuery {
	return q
}

// Deprecated
func (q *TokenNftInfoQuery) GetAccountID() AccountID {
	return AccountID{}
}

// Deprecated
func (q *TokenNftInfoQuery) SetStart(start int64) *TokenNftInfoQuery {
	return q
}

// Deprecated
func (q *TokenNftInfoQuery) GetStart() int64 {
	return 0
}

// Deprecated
func (q *TokenNftInfoQuery) SetEnd(end int64) *TokenNftInfoQuery {
	return q
}

// Deprecated
func (q *TokenNftInfoQuery) GetEnd() int64 {
	return 0
}

// Deprecated
func (q *TokenNftInfoQuery) ByNftID(id NftID) *TokenNftInfoQuery {
	q.nftID = &id
	return q
}

// Deprecated
func (q *TokenNftInfoQuery) ByTokenID(id TokenID) *TokenNftInfoQuery {
	return q
}

// Deprecated
func (q *TokenNftInfoQuery) ByAccountID(id AccountID) *TokenNftInfoQuery {
	return q
}

func (q *TokenNftInfoQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *TokenNftInfoQuery) Execute(client *Client) ([]TokenNftInfo, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return []TokenNftInfo{}, err
	}

	tokenInfos := make([]TokenNftInfo, 0)
	tokenInfos = append(tokenInfos, _TokenNftInfoFromProtobuf(resp.GetTokenGetNftInfo().GetNft()))
	return tokenInfos, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *TokenNftInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenNftInfoQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *TokenNftInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenNftInfoQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenNftInfoQuery.
func (q *TokenNftInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TokenNftInfoQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *TokenNftInfoQuery) SetMaxRetry(count int) *TokenNftInfoQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *TokenNftInfoQuery) SetMaxBackoff(max time.Duration) *TokenNftInfoQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *TokenNftInfoQuery) SetMinBackoff(min time.Duration) *TokenNftInfoQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *TokenNftInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *TokenNftInfoQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *TokenNftInfoQuery) SetLogLevel(level LogLevel) *TokenNftInfoQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *TokenNftInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetToken().GetTokenNftInfo,
	}
}

func (q *TokenNftInfoQuery) getName() string {
	return "TokenNftInfoQuery"
}

func (q *TokenNftInfoQuery) buildQuery() *services.Query {
	body := &services.TokenGetNftInfoQuery{
		Header: q.pbHeader,
	}

	if q.nftID != nil {
		body.NftID = q.nftID._ToProtobuf()
	}

	return &services.Query{
		Query: &services.Query_TokenGetNftInfo{
			TokenGetNftInfo: body,
		},
	}
}

func (q *TokenNftInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.nftID != nil {
		if err := q.nftID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *TokenNftInfoQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetTokenGetNftInfo()
}

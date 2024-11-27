package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// AccountStakersQuery gets all of the accounts that are proxy staking to this account. For each of  them, the amount
// currently staked will be given. This is not yet implemented, but will be in a future version of the API.
type AccountStakersQuery struct {
	Query
	accountID *AccountID
}

// NewAccountStakersQuery creates an AccountStakersQuery query which can be used to construct and execute
// an AccountStakersQuery.
//
// It is recommended that you use this for creating new instances of an AccountStakersQuery
// instead of manually creating an instance of the struct.
func NewAccountStakersQuery() *AccountStakersQuery {
	header := services.QueryHeader{}
	return &AccountStakersQuery{
		Query: _NewQuery(true, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *AccountStakersQuery) SetGrpcDeadline(deadline *time.Duration) *AccountStakersQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetAccountID sets the Account ID for which the stakers should be retrieved
func (q *AccountStakersQuery) SetAccountID(accountID AccountID) *AccountStakersQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the AccountID for this AccountStakersQuery.
func (q *AccountStakersQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}

	return *q.accountID
}

func (q *AccountStakersQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *AccountStakersQuery) Execute(client *Client) ([]Transfer, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return []Transfer{}, err
	}

	var stakers = make([]Transfer, len(resp.GetCryptoGetProxyStakers().Stakers.ProxyStaker))

	// TODO: This is wrong, q _Method shold return `[]ProxyStaker` not `[]Transfer`
	for i, element := range resp.GetCryptoGetProxyStakers().Stakers.ProxyStaker {
		id := _AccountIDFromProtobuf(element.AccountID)
		accountID := AccountID{}

		if id == nil {
			accountID = *id
		}

		stakers[i] = Transfer{
			AccountID: accountID,
			Amount:    HbarFromTinybar(element.Amount),
		}
	}

	return stakers, err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *AccountStakersQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountStakersQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *AccountStakersQuery) SetQueryPayment(paymentAmount Hbar) *AccountStakersQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountStakersQuery.
func (q *AccountStakersQuery) SetNodeAccountIDs(accountID []AccountID) *AccountStakersQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *AccountStakersQuery) SetMaxRetry(count int) *AccountStakersQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountStakersQuery) SetMaxBackoff(max time.Duration) *AccountStakersQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *AccountStakersQuery) SetMinBackoff(min time.Duration) *AccountStakersQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountStakersQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountStakersQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *AccountStakersQuery) SetLogLevel(level LogLevel) *AccountStakersQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *AccountStakersQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetStakersByAccountID,
	}
}

func (q *AccountStakersQuery) getName() string {
	return "AccountStakersQuery"
}

func (q *AccountStakersQuery) buildQuery() *services.Query {
	pb := services.Query_CryptoGetProxyStakers{
		CryptoGetProxyStakers: &services.CryptoGetStakersQuery{
			Header: q.pbHeader,
		},
	}

	if q.accountID != nil {
		pb.CryptoGetProxyStakers.AccountID = q.accountID._ToProtobuf()
	}

	return &services.Query{
		Query: &pb,
	}
}

func (q *AccountStakersQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.accountID != nil {
		if err := q.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *AccountStakersQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetCryptoGetProxyStakers()
}

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
type AccountBalanceQuery struct {
	Query
	accountID  *AccountID
	contractID *ContractID
}

// NewAccountBalanceQuery creates an AccountBalanceQuery query which can be used to construct and execute
// an AccountBalanceQuery.
// It is recommended that you use this for creating new instances of an AccountBalanceQuery
// instead of manually creating an instance of the struct.
func NewAccountBalanceQuery() *AccountBalanceQuery {
	header := services.QueryHeader{}
	return &AccountBalanceQuery{
		Query: _NewQuery(false, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *AccountBalanceQuery) SetGrpcDeadline(deadline *time.Duration) *AccountBalanceQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetAccountID sets the AccountID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (q *AccountBalanceQuery) SetAccountID(accountID AccountID) *AccountBalanceQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the AccountID for which you wish to query the balance.
func (q *AccountBalanceQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}

	return *q.accountID
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (q *AccountBalanceQuery) SetContractID(contractID ContractID) *AccountBalanceQuery {
	q.contractID = &contractID
	return q
}

// GetContractID returns the ContractID for which you wish to query the balance.
func (q *AccountBalanceQuery) GetContractID() ContractID {
	if q.contractID == nil {
		return ContractID{}
	}

	return *q.contractID
}

func (q *AccountBalanceQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the query with the provided client
func (q *AccountBalanceQuery) Execute(client *Client) (AccountBalance, error) {
	if client == nil {
		return AccountBalance{}, errNoClientProvided
	}
	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return AccountBalance{}, err
	}

	q.paymentTransactions = make([]*services.Transaction, 0)
	q.pb = q.buildQuery()

	resp, err := _Execute(client, q)

	if err != nil {
		return AccountBalance{}, err
	}

	result := _AccountBalanceFromProtobuf(resp.(*services.Response).GetCryptogetAccountBalance())
	// Query account_balance_query for given network/account;

	fmt.Println(client.GetMirrorNetwork()[0])
	const localNetwork = "127.0.0.1"
	if client.GetMirrorNetwork()[0] == localNetwork+":5600" || client.GetMirrorNetwork()[0] == localNetwork+":443" {
		if q.accountID != nil {
			_, err = queryBalanceFromMirrorNode(localNetwork+"5551", q.accountID.String(), &result)
		} else {
			_, err = queryBalanceFromMirrorNode(localNetwork+"5551", q.contractID.String(), &result)
		}
		if err != nil {
			return AccountBalance{}, err
		}
	} else {
		if q.accountID != nil {
			_, err = queryBalanceFromMirrorNode(client.GetMirrorNetwork()[0], q.accountID.String(), &result)
		} else {
			_, err = queryBalanceFromMirrorNode(client.GetMirrorNetwork()[0], q.contractID.String(), &result)
		}
		if err != nil {
			return AccountBalance{}, err
		}
	}

	return result, nil
}

// Helper function, which quey the mirror node and if the balance has tokens, it iterate over them to obtain
// decimals for each token
func queryBalanceFromMirrorNode(network string, id string, result *AccountBalance) (*AccountBalance, error) {
	response, err := accountBalanceQuery(network, id)
	if err != nil {
		return result, err
	}
	// If user has tokens
	result.Tokens.balances = make(map[string]uint64)
	if tokens, ok := response["tokens"].([]map[string]interface{}); ok {

		for _, token := range tokens {
			for key, value := range token {
				result.Tokens.balances[key] = value.(uint64)
			}
		}
	}

	return result, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this query.
func (q *AccountBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountBalanceQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this query.
func (q *AccountBalanceQuery) SetQueryPayment(paymentAmount Hbar) *AccountBalanceQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountBalanceQuery.
func (q *AccountBalanceQuery) SetNodeAccountIDs(accountID []AccountID) *AccountBalanceQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *AccountBalanceQuery) SetMaxRetry(count int) *AccountBalanceQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountBalanceQuery) SetMaxBackoff(max time.Duration) *AccountBalanceQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *AccountBalanceQuery) SetMinBackoff(min time.Duration) *AccountBalanceQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountBalanceQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountBalanceQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *AccountBalanceQuery) SetLogLevel(level LogLevel) *AccountBalanceQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *AccountBalanceQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().CryptoGetBalance,
	}
}

func (q *AccountBalanceQuery) getName() string {
	return "AccountBalanceQuery"
}

func (q *AccountBalanceQuery) buildProtoBody() *services.CryptoGetAccountBalanceQuery {
	pb := services.CryptoGetAccountBalanceQuery{Header: &services.QueryHeader{}}

	if q.accountID != nil {
		pb.BalanceSource = &services.CryptoGetAccountBalanceQuery_AccountID{
			AccountID: q.accountID._ToProtobuf(),
		}
	}

	if q.contractID != nil {
		pb.BalanceSource = &services.CryptoGetAccountBalanceQuery_ContractID{
			ContractID: q.contractID._ToProtobuf(),
		}
	}
	return &pb
}
func (q *AccountBalanceQuery) buildQuery() *services.Query {
	pb := q.buildProtoBody()

	return &services.Query{
		Query: &services.Query_CryptogetAccountBalance{
			CryptogetAccountBalance: pb,
		},
	}
}

func (q *AccountBalanceQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.accountID != nil {
		if err := q.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if q.contractID != nil {
		if err := q.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *AccountBalanceQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetCryptogetAccountBalance()
}

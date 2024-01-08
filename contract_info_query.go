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
	"errors"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
	header := services.QueryHeader{}
	query := _NewQuery(true, &header)

	return &ContractInfoQuery{
		Query: query,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *ContractInfoQuery) SetGrpcDeadline(deadline *time.Duration) *ContractInfoQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetContractID sets the contract for which information is requested
func (q *ContractInfoQuery) SetContractID(contractID ContractID) *ContractInfoQuery {
	q.contractID = &contractID
	return q
}

func (q *ContractInfoQuery) GetContractID() ContractID {
	if q.contractID == nil {
		return ContractID{}
	}

	return *q.contractID
}

func (q *ContractInfoQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *ContractInfoQuery) Execute(client *Client) (ContractInfo, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return ContractInfo{}, err
	}

	info, err := _ContractInfoFromProtobuf(resp.GetContractGetInfo().ContractInfo)
	if err != nil {
		return ContractInfo{}, err
	}
	network := obtainUrlForMirrorNode(client)
	_, err = contractInfoqueryTokensRelationshipFromMirrorNode(network, q.contractID.String(), &info)
	if err != nil {
		return ContractInfo{}, err
	}

	return info, nil
}

// Helper function, which query the mirror node about tokenRelationship of for all tokens that the account is
// being associated with
func contractInfoqueryTokensRelationshipFromMirrorNode(network string, id string, result *ContractInfo) (*ContractInfo, error) {
	response, err := tokenReleationshipQuery(network, id)
	if err != nil {
		return result, err
	}
	tokens, ok := response["tokens"].([]interface{})
	if !ok {
		return result, errors.New("Ivalid tokens format")
	}
	mappedTokens, err := mapTokenRelationship(tokens)
	if err != nil {
		return &ContractInfo{}, err
	}
	result.TokenRelationships = mappedTokens
	return result, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ContractInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractInfoQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ContractInfoQuery) SetQueryPayment(paymentAmount Hbar) *ContractInfoQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractInfoQuery.
func (q *ContractInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ContractInfoQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ContractInfoQuery) SetMaxRetry(count int) *ContractInfoQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ContractInfoQuery) SetMaxBackoff(max time.Duration) *ContractInfoQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ContractInfoQuery) SetMinBackoff(min time.Duration) *ContractInfoQuery {
	q.Query.SetMinBackoff(min)
	return q
}

func (q *ContractInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *ContractInfoQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *ContractInfoQuery) SetLogLevel(level LogLevel) *ContractInfoQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *ContractInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().GetContractInfo,
	}
}

func (q *ContractInfoQuery) getName() string {
	return "ContractInfoQuery"
}

func (q *ContractInfoQuery) buildQuery() *services.Query {
	pb := services.Query_ContractGetInfo{
		ContractGetInfo: &services.ContractGetInfoQuery{
			Header: q.pbHeader,
		},
	}

	if q.contractID != nil {
		pb.ContractGetInfo.ContractID = q.contractID._ToProtobuf()
	}

	return &services.Query{
		Query: &pb,
	}
}

func (q *ContractInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.contractID != nil {
		if err := q.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *ContractInfoQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetContractGetInfo()
}

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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TopicInfo is the Query for retrieving information about a topic stored on the Hedera network.
type TopicInfoQuery struct {
	Query
	topicID *TopicID
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//
//	Get Topic Info Query.
func NewTopicInfoQuery() *TopicInfoQuery {
	header := services.QueryHeader{}
	return &TopicInfoQuery{
		Query: _NewQuery(true, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *TopicInfoQuery) SetGrpcDeadline(deadline *time.Duration) *TopicInfoQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (q *TopicInfoQuery) SetTopicID(topicID TopicID) *TopicInfoQuery {
	q.topicID = &topicID
	return q
}

// GetTopicID returns the TopicID for this TopicInfoQuery
func (q *TopicInfoQuery) GetTopicID() TopicID {
	if q.topicID == nil {
		return TopicID{}
	}

	return *q.topicID
}

func (q *TopicInfoQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the TopicInfoQuery using the provided client
func (q *TopicInfoQuery) Execute(client *Client) (TopicInfo, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return TopicInfo{}, err
	}

	return _TopicInfoFromProtobuf(resp.GetConsensusGetTopicInfo().TopicInfo)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *TopicInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TopicInfoQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *TopicInfoQuery) SetQueryPayment(paymentAmount Hbar) *TopicInfoQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this TopicInfoQuery.
func (q *TopicInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TopicInfoQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *TopicInfoQuery) SetMaxRetry(count int) *TopicInfoQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *TopicInfoQuery) SetMaxBackoff(max time.Duration) *TopicInfoQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *TopicInfoQuery) SetMinBackoff(min time.Duration) *TopicInfoQuery {
	q.Query.SetMinBackoff(min)
	return q
}

func (q *TopicInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *TopicInfoQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *TopicInfoQuery) SetLogLevel(level LogLevel) *TopicInfoQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *TopicInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetTopic().GetTopicInfo,
	}
}

func (q *TopicInfoQuery) getName() string {
	return "TopicInfoQuery"
}

func (q *TopicInfoQuery) buildQuery() *services.Query {
	body := &services.ConsensusGetTopicInfoQuery{
		Header: q.pbHeader,
	}

	if q.topicID != nil {
		body.TopicID = q.topicID._ToProtobuf()
	}

	return &services.Query{
		Query: &services.Query_ConsensusGetTopicInfo{
			ConsensusGetTopicInfo: body,
		},
	}
}

func (q *TopicInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.topicID != nil {
		if err := q.topicID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *TopicInfoQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetConsensusGetTopicInfo()
}

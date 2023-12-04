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

// TopicInfo is the Query for retrieving information about a topic stored on the Hedera network.
type TopicInfoQuery struct {
	query
	topicID *TopicID
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//
//	Get Topic Info Query.
func NewTopicInfoQuery() *TopicInfoQuery {
	header := services.QueryHeader{}
	result := TopicInfoQuery{
		query: _NewQuery(true, &header),
	}

	//	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *TopicInfoQuery) SetGrpcDeadline(deadline *time.Duration) *TopicInfoQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (this *TopicInfoQuery) SetTopicID(topicID TopicID) *TopicInfoQuery {
	this.topicID = &topicID
	return this
}

// GetTopicID returns the TopicID for this TopicInfoQuery
func (this *TopicInfoQuery) GetTopicID() TopicID {
	if this.topicID == nil {
		return TopicID{}
	}

	return *this.topicID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *TopicInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range this.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.ConsensusGetTopicInfo.Header = this.pbHeader

	this.pb = &services.Query{
		Query: pb,
	}

	this.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	this.paymentTransactionIDs._Advance()

	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetConsensusGetTopicInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the TopicInfoQuery using the provided client
func (this *TopicInfoQuery) Execute(client *Client) (TopicInfo, error) {
	if client == nil || client.operator == nil {
		return TopicInfo{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return TopicInfo{}, err
	}

	if !this.paymentTransactionIDs.locked {
		this.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if this.queryPayment.tinybar != 0 {
		cost = this.queryPayment
	} else {
		if this.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = this.maxQueryPayment
		}

		actualCost, err := this.GetCost(client)
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

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return TopicInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return TopicInfo{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.ConsensusGetTopicInfo.Header = this.pbHeader
	this.pb = &services.Query{
		Query: pb,
	}

	if this.isPaymentRequired && len(this.paymentTransactions) > 0 {
		this.paymentTransactionIDs._Advance()
	}
	this.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return TopicInfo{}, err
	}

	return _TopicInfoFromProtobuf(resp.(*services.Response).GetConsensusGetTopicInfo().TopicInfo)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *TopicInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TopicInfoQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *TopicInfoQuery) SetQueryPayment(paymentAmount Hbar) *TopicInfoQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this TopicInfoQuery.
func (this *TopicInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TopicInfoQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *TopicInfoQuery) SetMaxRetry(count int) *TopicInfoQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *TopicInfoQuery) SetMaxBackoff(max time.Duration) *TopicInfoQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *TopicInfoQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *TopicInfoQuery) SetMinBackoff(min time.Duration) *TopicInfoQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *TopicInfoQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *TopicInfoQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("TopicInfoQuery:%d", timestamp)
}

func (this *TopicInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *TopicInfoQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *TopicInfoQuery) SetLogLevel(level LogLevel) *TopicInfoQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *TopicInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetTopic().GetTopicInfo,
	}
}

func (this *TopicInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetConsensusGetTopicInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (this *TopicInfoQuery) getName() string {
	return "TopicInfoQuery"
}

func (this *TopicInfoQuery) build() *services.Query_ConsensusGetTopicInfo {
	body := &services.ConsensusGetTopicInfoQuery{
		Header: &services.QueryHeader{},
	}

	if this.topicID != nil {
		body.TopicID = this.topicID._ToProtobuf()
	}

	return &services.Query_ConsensusGetTopicInfo{
		ConsensusGetTopicInfo: body,
	}
}

func (this *TopicInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.topicID != nil {
		if err := this.topicID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *TopicInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetConsensusGetTopicInfo().Header.NodeTransactionPrecheckCode)
}

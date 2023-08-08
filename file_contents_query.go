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

// FileContentsQuery retrieves the contents of a file.
type FileContentsQuery struct {
	Query
	fileID *FileID
}

// NewFileContentsQuery creates a FileContentsQuery which retrieves the contents of a file.
func NewFileContentsQuery() *FileContentsQuery {
	header := services.QueryHeader{}
	return &FileContentsQuery{
		Query: _NewQuery(true, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (query *FileContentsQuery) SetGrpcDeadline(deadline *time.Duration) *FileContentsQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

// SetFileID sets the FileID of the file whose contents are requested.
func (query *FileContentsQuery) SetFileID(fileID FileID) *FileContentsQuery {
	query.fileID = &fileID
	return query
}

// GetFileID returns the FileID of the file whose contents are requested.
func (query *FileContentsQuery) GetFileID() FileID {
	if query.fileID == nil {
		return FileID{}
	}

	return *query.fileID
}

func (query *FileContentsQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.fileID != nil {
		if err := query.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *FileContentsQuery) _Build() *services.Query_FileGetContents {
	body := &services.FileGetContentsQuery{
		Header: &services.QueryHeader{},
	}

	if query.fileID != nil {
		body.FileID = query.fileID._ToProtobuf()
	}

	return &services.Query_FileGetContents{
		FileGetContents: body,
	}
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (query *FileContentsQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range query.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.FileGetContents.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_FileContentsQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_FileContentsQueryGetMethod,
		_FileContentsQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetFileGetContents().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _FileContentsQueryShouldRetry(_ interface{}, response interface{}) _ExecutionState {
	return _QueryShouldRetry(Status(response.(*services.Response).GetFileGetContents().Header.NodeTransactionPrecheckCode))
}

func _FileContentsQueryMapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetFileGetContents().Header.NodeTransactionPrecheckCode),
	}
}

func _FileContentsQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetFile().GetFileContent,
	}
}

// Execute executes the Query with the provided client
func (query *FileContentsQuery) Execute(client *Client) ([]byte, error) {
	if client == nil || client.operator == nil {
		return make([]byte, 0), errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return []byte{}, err
	}

	if !query.paymentTransactionIDs.locked {
		query.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return []byte{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []byte{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "FileContentsQuery",
			}
		}

		cost = actualCost
	}

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return []byte{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(query.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return []byte{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.FileGetContents.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_FileContentsQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_FileContentsQueryGetMethod,
		_FileContentsQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return []byte{}, err
	}

	return resp.(*services.Response).GetFileGetContents().FileContents.Contents, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *FileContentsQuery) SetMaxQueryPayment(maxPayment Hbar) *FileContentsQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *FileContentsQuery) SetQueryPayment(paymentAmount Hbar) *FileContentsQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

// SetNodeAccountIDs sets the _Node AccountID for this FileContentsQuery.
func (query *FileContentsQuery) SetNodeAccountIDs(accountID []AccountID) *FileContentsQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (query *FileContentsQuery) SetMaxRetry(count int) *FileContentsQuery {
	query.Query.SetMaxRetry(count)
	return query
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (query *FileContentsQuery) SetMaxBackoff(max time.Duration) *FileContentsQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (query *FileContentsQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (query *FileContentsQuery) SetMinBackoff(min time.Duration) *FileContentsQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (query *FileContentsQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *FileContentsQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("FileContentsQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (query *FileContentsQuery) SetPaymentTransactionID(transactionID TransactionID) *FileContentsQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}

func (query *FileContentsQuery) SetLogLevel(level LogLevel) *FileContentsQuery {
	query.Query.SetLogLevel(level)
	return query
}

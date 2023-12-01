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

// FileInfoQuery is a query which can be used to get all of the information about a file, except for its contents.
// When a file expires, it no longer exists, and there will be no info about it, and the fileInfo field will be blank.
// If a transaction or smart contract deletes the file, but it has not yet expired, then the
// fileInfo field will be non-empty, the deleted field will be true, its size will be 0,
// and its contents will be empty. Note that each file has a FileID, but does not have a filename.
type FileInfoQuery struct {
	query
	fileID *FileID
}

// NewFileInfoQuery creates a FileInfoQuery which can be used to get all of the information about a file, except for its contents.
func NewFileInfoQuery() *FileInfoQuery {
	header := services.QueryHeader{}
	result := FileInfoQuery{
		query: _NewQuery(true, &header),
	}
	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *FileInfoQuery) SetGrpcDeadline(deadline *time.Duration) *FileInfoQuery {
	q.query.SetGrpcDeadline(deadline)
	return q
}

// SetFileID sets the FileID of the file whose info is requested.
func (q *FileInfoQuery) SetFileID(fileID FileID) *FileInfoQuery {
	q.fileID = &fileID
	return q
}

// GetFileID returns the FileID of the file whose info is requested.
func (q *FileInfoQuery) GetFileID() FileID {
	if q.fileID == nil {
		return FileID{}
	}

	return *q.fileID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (q *FileInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range q.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.FileGetInfo.Header = q.pbHeader

	q.pb = &services.Query{
		Query: pb,
	}

	q.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	q.paymentTransactionIDs._Advance()

	resp, err := _Execute(
		client,
		q.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetFileGetInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (q *FileInfoQuery) Execute(client *Client) (FileInfo, error) {
	if client == nil || client.operator == nil {
		return FileInfo{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return FileInfo{}, err
	}

	if !q.paymentTransactionIDs.locked {
		q.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if q.queryPayment.tinybar != 0 {
		cost = q.queryPayment
	} else {
		if q.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = q.maxQueryPayment
		}

		actualCost, err := q.GetCost(client)
		if err != nil {
			return FileInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return FileInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "FileInfoQuery",
			}
		}

		cost = actualCost
	}

	q.paymentTransactions = make([]*services.Transaction, 0)

	if q.nodeAccountIDs.locked {
		err = q._QueryGeneratePayments(client, cost)
		if err != nil {
			if err != nil {
				return FileInfo{}, err
			}
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(q.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			if err != nil {
				return FileInfo{}, err
			}
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.FileGetInfo.Header = q.pbHeader
	q.pb = &services.Query{
		Query: pb,
	}

	if q.isPaymentRequired && len(q.paymentTransactions) > 0 {
		q.paymentTransactionIDs._Advance()
	}
	q.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(
		client,
		q.e,
	)

	if err != nil {
		return FileInfo{}, err
	}

	info, err := _FileInfoFromProtobuf(resp.(*services.Response).GetFileGetInfo().FileInfo)
	if err != nil {
		return FileInfo{}, err
	}

	return info, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *FileInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *FileInfoQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *FileInfoQuery) SetQueryPayment(paymentAmount Hbar) *FileInfoQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this FileInfoQuery.
func (q *FileInfoQuery) SetNodeAccountIDs(accountID []AccountID) *FileInfoQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// GetNodeAccountIDs returns the _Node AccountID for this FileInfoQuery.
func (q *FileInfoQuery) GetNodeAccountIDs() []AccountID {
	return q.query.GetNodeAccountIDs()
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *FileInfoQuery) SetMaxRetry(count int) *FileInfoQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *FileInfoQuery) SetMaxBackoff(max time.Duration) *FileInfoQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *FileInfoQuery) SetMinBackoff(min time.Duration) *FileInfoQuery {
	q.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *FileInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *FileInfoQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *FileInfoQuery) SetLogLevel(level LogLevel) *FileInfoQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *FileInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetFile().GetFileInfo,
	}
}

func (q *FileInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetFileGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (q *FileInfoQuery) getName() string {
	return "FileInfoQuery"
}

func (q *FileInfoQuery) build() *services.Query_FileGetInfo {
	body := &services.FileGetInfoQuery{
		Header: &services.QueryHeader{},
	}

	if q.fileID != nil {
		body.FileID = q.fileID._ToProtobuf()
	}

	return &services.Query_FileGetInfo{
		FileGetInfo: body,
	}
}

func (q *FileInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.fileID != nil {
		if err := q.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *FileInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetFileGetInfo().Header.NodeTransactionPrecheckCode)
}

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
func (this *FileInfoQuery) SetGrpcDeadline(deadline *time.Duration) *FileInfoQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetFileID sets the FileID of the file whose info is requested.
func (this *FileInfoQuery) SetFileID(fileID FileID) *FileInfoQuery {
	this.fileID = &fileID
	return this
}

// GetFileID returns the FileID of the file whose info is requested.
func (this *FileInfoQuery) GetFileID() FileID {
	if this.fileID == nil {
		return FileID{}
	}

	return *this.fileID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *FileInfoQuery) GetCost(client *Client) (Hbar, error) {
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
	pb.FileGetInfo.Header = this.pbHeader

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

	cost := int64(resp.(*services.Response).GetFileGetInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (this *FileInfoQuery) Execute(client *Client) (FileInfo, error) {
	if client == nil || client.operator == nil {
		return FileInfo{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return FileInfo{}, err
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

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			if err != nil {
				return FileInfo{}, err
			}
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			if err != nil {
				return FileInfo{}, err
			}
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.FileGetInfo.Header = this.pbHeader
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
		return FileInfo{}, err
	}

	info, err := _FileInfoFromProtobuf(resp.(*services.Response).GetFileGetInfo().FileInfo)
	if err != nil {
		return FileInfo{}, err
	}

	return info, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *FileInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *FileInfoQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *FileInfoQuery) SetQueryPayment(paymentAmount Hbar) *FileInfoQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this FileInfoQuery.
func (this *FileInfoQuery) SetNodeAccountIDs(accountID []AccountID) *FileInfoQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// GetNodeAccountIDs returns the _Node AccountID for this FileInfoQuery.
func (this *FileInfoQuery) GetNodeAccountIDs() []AccountID {
	return this.query.GetNodeAccountIDs()
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *FileInfoQuery) SetMaxRetry(count int) *FileInfoQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *FileInfoQuery) SetMaxBackoff(max time.Duration) *FileInfoQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *FileInfoQuery) GetMaxBackoff() time.Duration {
	return this.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *FileInfoQuery) SetMinBackoff(min time.Duration) *FileInfoQuery {
	this.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *FileInfoQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *FileInfoQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("FileInfoQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *FileInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *FileInfoQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *FileInfoQuery) SetLogLevel(level LogLevel) *FileInfoQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *FileInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetFile().GetFileInfo,
	}
}

func (this *FileInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetFileGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (this *FileInfoQuery) getName() string {
	return "FileInfoQuery"
}

func (this *FileInfoQuery) build() *services.Query_FileGetInfo {
	body := &services.FileGetInfoQuery{
		Header: &services.QueryHeader{},
	}

	if this.fileID != nil {
		body.FileID = this.fileID._ToProtobuf()
	}

	return &services.Query_FileGetInfo{
		FileGetInfo: body,
	}
}

func (this *FileInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.fileID != nil {
		if err := this.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *FileInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetFileGetInfo().Header.NodeTransactionPrecheckCode)
}

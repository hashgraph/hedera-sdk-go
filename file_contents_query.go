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
	query
	fileID *FileID
}

// NewFileContentsQuery creates a FileContentsQuery which retrieves the contents of a file.
func NewFileContentsQuery() *FileContentsQuery {
	header := services.QueryHeader{}
	return &FileContentsQuery{
		query: _NewQuery(true, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *FileContentsQuery) SetGrpcDeadline(deadline *time.Duration) *FileContentsQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetFileID sets the FileID of the file whose contents are requested.
func (this *FileContentsQuery) SetFileID(fileID FileID) *FileContentsQuery {
	this.fileID = &fileID
	return this
}

// GetFileID returns the FileID of the file whose contents are requested.
func (this *FileContentsQuery) GetFileID() FileID {
	if this.fileID == nil {
		return FileID{}
	}

	return *this.fileID
}

func (this *FileContentsQuery) validateNetworkOnIDs(client *Client) error {
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

func (this *FileContentsQuery) build() *services.Query_FileGetContents {
	body := &services.FileGetContentsQuery{
		Header: &services.QueryHeader{},
	}

	if this.fileID != nil {
		body.FileID = this.fileID._ToProtobuf()
	}

	return &services.Query_FileGetContents{
		FileGetContents: body,
	}
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *FileContentsQuery) GetCost(client *Client) (Hbar, error) {
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
	pb.FileGetContents.Header = this.pbHeader

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

	cost := int64(resp.(*services.Response).GetFileGetContents().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func (this *FileContentsQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetFileGetContents().Header.NodeTransactionPrecheckCode),
	}
}

func (this *FileContentsQuery) getMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetFile().GetFileContent,
	}
}
// Get the name of the query
func (this *FileContentsQuery) getName() string {
	return "FileContentsQuery"
}

// Execute executes the Query with the provided client
func (this *FileContentsQuery) Execute(client *Client) ([]byte, error) {
	if client == nil || client.operator == nil {
		return make([]byte, 0), errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return []byte{}, err
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

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&this.query, client, cost)
		if err != nil {
			return []byte{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return []byte{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.FileGetContents.Header = this.pbHeader
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
		return []byte{}, err
	}

	return resp.(*services.Response).GetFileGetContents().FileContents.Contents, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *FileContentsQuery) SetMaxQueryPayment(maxPayment Hbar) *FileContentsQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *FileContentsQuery) SetQueryPayment(paymentAmount Hbar) *FileContentsQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this FileContentsQuery.
func (this *FileContentsQuery) SetNodeAccountIDs(accountID []AccountID) *FileContentsQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *FileContentsQuery) SetMaxRetry(count int) *FileContentsQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *FileContentsQuery) SetMaxBackoff(max time.Duration) *FileContentsQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *FileContentsQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *FileContentsQuery) SetMinBackoff(min time.Duration) *FileContentsQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *FileContentsQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *FileContentsQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("FileContentsQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *FileContentsQuery) SetPaymentTransactionID(transactionID TransactionID) *FileContentsQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *FileContentsQuery) SetLogLevel(level LogLevel) *FileContentsQuery {
	this.query.SetLogLevel(level)
	return this
}

func (this *FileContentsQuery) getQueryStatus(response interface{}) (Status) {
	return Status(response.(*services.Response).GetFileGetContents().Header.NodeTransactionPrecheckCode)
}
package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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

type FileInfoQuery struct {
	Query
	fileID *FileID
}

func NewFileInfoQuery() *FileInfoQuery {
	header := services.QueryHeader{}
	return &FileInfoQuery{
		Query: _NewQuery(true, &header),
	}
}

func (query *FileInfoQuery) SetGrpcDeadline(deadline *time.Duration) *FileInfoQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

func (query *FileInfoQuery) SetFileID(fileID FileID) *FileInfoQuery {
	query.fileID = &fileID
	return query
}

func (query *FileInfoQuery) GetFileID() FileID {
	if query.fileID == nil {
		return FileID{}
	}

	return *query.fileID
}

func (query *FileInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
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

func (query *FileInfoQuery) _Build() *services.Query_FileGetInfo {
	body := &services.FileGetInfoQuery{
		Header: &services.QueryHeader{},
	}

	if query.fileID != nil {
		body.FileID = query.fileID._ToProtobuf()
	}

	return &services.Query_FileGetInfo{
		FileGetInfo: body,
	}
}

func (query *FileInfoQuery) GetCost(client *Client) (Hbar, error) {
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
	pb.FileGetInfo.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_FileInfoQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_FileInfoQueryGetMethod,
		_FileInfoQueryMapStatusError,
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

	cost := int64(resp.(*services.Response).GetFileGetInfo().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _FileInfoQueryShouldRetry(logID string, _ interface{}, response interface{}) _ExecutionState {
	return _QueryShouldRetry(logID, Status(response.(*services.Response).GetFileGetInfo().Header.NodeTransactionPrecheckCode))
}

func _FileInfoQueryMapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetFileGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _FileInfoQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetFile().GetFileInfo,
	}
}

func (query *FileInfoQuery) Execute(client *Client) (FileInfo, error) {
	if client == nil || client.operator == nil {
		return FileInfo{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return FileInfo{}, err
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

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			if err != nil {
				return FileInfo{}, err
			}
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(query.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			if err != nil {
				return FileInfo{}, err
			}
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.FileGetInfo.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_FileInfoQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_FileInfoQueryGetMethod,
		_FileInfoQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
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
func (query *FileInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *FileInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *FileInfoQuery) SetQueryPayment(paymentAmount Hbar) *FileInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *FileInfoQuery) SetNodeAccountIDs(accountID []AccountID) *FileInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *FileInfoQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

func (query *FileInfoQuery) SetMaxRetry(count int) *FileInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *FileInfoQuery) SetMaxBackoff(max time.Duration) *FileInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *FileInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *FileInfoQuery) SetMinBackoff(min time.Duration) *FileInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *FileInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *FileInfoQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("FileInfoQuery:%d", timestamp)
}

func (query *FileInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *FileInfoQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}

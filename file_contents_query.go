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

// FileContentsQuery retrieves the contents of a file.
type FileContentsQuery struct {
	query
	fileID *FileID
}

// NewFileContentsQuery creates a FileContentsQuery which retrieves the contents of a file.
func NewFileContentsQuery() *FileContentsQuery {
	header := services.QueryHeader{}
	result := FileContentsQuery{
		query: _NewQuery(true, &header),
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *FileContentsQuery) SetGrpcDeadline(deadline *time.Duration) *FileContentsQuery {
	q.query.SetGrpcDeadline(deadline)
	return q
}

// SetFileID sets the FileID of the file whose contents are requested.
func (q *FileContentsQuery) SetFileID(fileID FileID) *FileContentsQuery {
	q.fileID = &fileID
	return q
}

// GetFileID returns the FileID of the file whose contents are requested.
func (q *FileContentsQuery) GetFileID() FileID {
	if q.fileID == nil {
		return FileID{}
	}

	return *q.fileID
}

// Execute executes the Query with the provided client
func (q *FileContentsQuery) Execute(client *Client) ([]byte, error) {
	resp, err := q.query.execute(client)

	if err != nil {
		return []byte{}, err
	}

	return resp.GetFileGetContents().FileContents.Contents, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *FileContentsQuery) SetMaxQueryPayment(maxPayment Hbar) *FileContentsQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *FileContentsQuery) SetQueryPayment(paymentAmount Hbar) *FileContentsQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this FileContentsQuery.
func (q *FileContentsQuery) SetNodeAccountIDs(accountID []AccountID) *FileContentsQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *FileContentsQuery) SetMaxRetry(count int) *FileContentsQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *FileContentsQuery) SetMaxBackoff(max time.Duration) *FileContentsQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *FileContentsQuery) SetMinBackoff(min time.Duration) *FileContentsQuery {
	q.query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *FileContentsQuery) SetPaymentTransactionID(transactionID TransactionID) *FileContentsQuery {
	q.query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *FileContentsQuery) SetLogLevel(level LogLevel) *FileContentsQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *FileContentsQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetFile().GetFileContent,
	}
}

// Get the name of the query
func (q *FileContentsQuery) getName() string {
	return "FileContentsQuery"
}

func (q *FileContentsQuery) buildQuery() *services.Query {
	body := &services.FileGetContentsQuery{
		Header: q.pbHeader,
	}

	if q.fileID != nil {
		body.FileID = q.fileID._ToProtobuf()
	}

	return &services.Query{
		Query: &services.Query_FileGetContents{
			FileGetContents: body,
		},
	}
}

func (q *FileContentsQuery) validateNetworkOnIDs(client *Client) error {
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

func (q *FileContentsQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetFileGetContents()
}

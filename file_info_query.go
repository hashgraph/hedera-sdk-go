package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// FileInfoQuery is a query which can be used to get all of the information about a file, except for its contents.
// When a file expires, it no longer exists, and there will be no info about it, and the fileInfo field will be blank.
// If a transaction or smart contract deletes the file, but it has not yet expired, then the
// fileInfo field will be non-empty, the deleted field will be true, its size will be 0,
// and its contents will be empty. Note that each file has a FileID, but does not have a filename.
type FileInfoQuery struct {
	Query
	fileID *FileID
}

// NewFileInfoQuery creates a FileInfoQuery which can be used to get all of the information about a file, except for its contents.
func NewFileInfoQuery() *FileInfoQuery {
	header := services.QueryHeader{}
	return &FileInfoQuery{
		Query: _NewQuery(true, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *FileInfoQuery) SetGrpcDeadline(deadline *time.Duration) *FileInfoQuery {
	q.Query.SetGrpcDeadline(deadline)
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

func (q *FileInfoQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *FileInfoQuery) Execute(client *Client) (FileInfo, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return FileInfo{}, err
	}

	info, err := _FileInfoFromProtobuf(resp.GetFileGetInfo().FileInfo)
	if err != nil {
		return FileInfo{}, err
	}

	return info, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *FileInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *FileInfoQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *FileInfoQuery) SetQueryPayment(paymentAmount Hbar) *FileInfoQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this FileInfoQuery.
func (q *FileInfoQuery) SetNodeAccountIDs(accountID []AccountID) *FileInfoQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *FileInfoQuery) SetMaxRetry(count int) *FileInfoQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *FileInfoQuery) SetMaxBackoff(max time.Duration) *FileInfoQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *FileInfoQuery) SetMinBackoff(min time.Duration) *FileInfoQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *FileInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *FileInfoQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *FileInfoQuery) SetLogLevel(level LogLevel) *FileInfoQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *FileInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetFile().GetFileInfo,
	}
}

func (q *FileInfoQuery) getName() string {
	return "FileInfoQuery"
}

func (q *FileInfoQuery) buildQuery() *services.Query {
	body := &services.FileGetInfoQuery{
		Header: q.pbHeader,
	}

	if q.fileID != nil {
		body.FileID = q.fileID._ToProtobuf()
	}

	return &services.Query{
		Query: &services.Query_FileGetInfo{
			FileGetInfo: body,
		},
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

func (q *FileInfoQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetFileGetInfo()
}

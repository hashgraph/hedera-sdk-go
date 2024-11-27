package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"io"
	"math"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
	"google.golang.org/grpc/status"
)

// AddressBookQuery query an address book for its list of nodes
type AddressBookQuery struct {
	attempt     uint64
	maxAttempts uint64
	fileID      *FileID
	limit       int32
}

// Query the mirror node for the address book.
func NewAddressBookQuery() *AddressBookQuery {
	return &AddressBookQuery{
		fileID: nil,
		limit:  0,
	}
}

// SetFileID set the ID of the address book file on the network. Can be either 0.0.101 or 0.0.102.
func (q *AddressBookQuery) SetFileID(id FileID) *AddressBookQuery {
	q.fileID = &id
	return q
}

func (q *AddressBookQuery) GetFileID() FileID {
	if q.fileID == nil {
		return FileID{}
	}

	return *q.fileID
}

// SetLimit
// Set the maximum number of node addresses to receive before stopping.
// If not set or set to zero it will return all node addresses in the database.
func (q *AddressBookQuery) SetLimit(limit int32) *AddressBookQuery {
	q.limit = limit
	return q
}

func (q *AddressBookQuery) GetLimit() int32 {
	return q.limit
}

func (q *AddressBookQuery) SetMaxAttempts(maxAttempts uint64) *AddressBookQuery {
	q.maxAttempts = maxAttempts
	return q
}

func (q *AddressBookQuery) GetMaxAttempts() uint64 {
	return q.maxAttempts
}

func (q *AddressBookQuery) validateNetworkOnIDs(client *Client) error {
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

func (q *AddressBookQuery) build() *mirror.AddressBookQuery {
	body := &mirror.AddressBookQuery{
		Limit: q.limit,
	}
	if q.fileID != nil {
		body.FileId = q.fileID._ToProtobuf()
	}

	return body
}

// Execute executes the Query with the provided client
func (q *AddressBookQuery) Execute(client *Client) (NodeAddressBook, error) {
	var cancel func()
	var ctx context.Context
	var subClientError error
	err := q.validateNetworkOnIDs(client)
	if err != nil {
		return NodeAddressBook{}, err
	}

	pb := q.build()

	messages := make([]*services.NodeAddress, 0)

	channel, err := client.mirrorNetwork._GetNextMirrorNode()._GetNetworkServiceClient()
	if err != nil {
		return NodeAddressBook{}, err
	}
	ch := make(chan byte, 1)

	go func() {
		var subClient mirror.NetworkService_GetNodesClient
		var err error

		for {
			if err != nil {
				cancel()

				if grpcErr, ok := status.FromError(err); ok { // nolint
					if q.attempt < q.maxAttempts {
						subClient = nil

						delay := math.Min(250.0*math.Pow(2.0, float64(q.attempt)), 8000)
						time.Sleep(time.Duration(delay) * time.Millisecond)
						q.attempt++
					} else {
						subClientError = grpcErr.Err()
						break
					}
				} else if err == io.EOF {
					break
				} else {
					subClientError = err
					break
				}
			}

			if subClient == nil {
				ctx, cancel = context.WithCancel(client.networkUpdateContext)

				subClient, err = (*channel).GetNodes(ctx, pb)
				if err != nil {
					continue
				}
			}

			var resp *services.NodeAddress
			resp, err = subClient.Recv()

			if err != nil {
				continue
			}

			if pb.Limit > 0 {
				pb.Limit--
			}

			messages = append(messages, resp)
		}
		ch <- 1
	}()
	<-ch

	result := make([]NodeAddress, 0)

	for _, k := range messages {
		result = append(result, _NodeAddressFromProtobuf(k))
	}

	return NodeAddressBook{
		NodeAddresses: result,
	}, subClientError
}

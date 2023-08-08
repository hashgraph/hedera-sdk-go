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
	"context"
	"io"
	"math"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/hashgraph/hedera-protobufs-go/mirror"
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
func (query *AddressBookQuery) SetFileID(id FileID) *AddressBookQuery {
	query.fileID = &id
	return query
}

func (query *AddressBookQuery) GetFileID() FileID {
	if query.fileID == nil {
		return FileID{}
	}

	return *query.fileID
}

// SetLimit
// Set the maximum number of node addresses to receive before stopping.
// If not set or set to zero it will return all node addresses in the database.
func (query *AddressBookQuery) SetLimit(limit int32) *AddressBookQuery {
	query.limit = limit
	return query
}

func (query *AddressBookQuery) GetLimit() int32 {
	return query.limit
}

func (query *AddressBookQuery) SetMaxAttempts(maxAttempts uint64) *AddressBookQuery {
	query.maxAttempts = maxAttempts
	return query
}

func (query *AddressBookQuery) GetMaxAttempts() uint64 {
	return query.maxAttempts
}

func (query *AddressBookQuery) _ValidateNetworkOnIDs(client *Client) error {
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

func (query *AddressBookQuery) _Build() *mirror.AddressBookQuery {
	body := &mirror.AddressBookQuery{
		Limit: query.limit,
	}
	if query.fileID != nil {
		body.FileId = query.fileID._ToProtobuf()
	}

	return body
}

// Execute executes the Query with the provided client
func (query *AddressBookQuery) Execute(client *Client) (NodeAddressBook, error) {
	var cancel func()
	var ctx context.Context
	var subClientError error
	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return NodeAddressBook{}, err
	}

	pb := query._Build()

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
					if query.attempt < query.maxAttempts {
						subClient = nil

						delay := math.Min(250.0*math.Pow(2.0, float64(query.attempt)), 8000)
						time.Sleep(time.Duration(delay) * time.Millisecond)
						query.attempt++
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

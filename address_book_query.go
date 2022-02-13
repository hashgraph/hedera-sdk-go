package hedera

import (
	"context"
	"io"
	"math"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/hashgraph/hedera-protobufs-go/mirror"
	"google.golang.org/grpc/status"
)

type AddressBookQuery struct {
	attempt     uint64
	maxAttempts uint64
	fileID      *FileID
	limit       int32
}

func NewAddressBookQuery() *AddressBookQuery {
	return &AddressBookQuery{
		fileID: nil,
		limit:  0,
	}
}

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

func (query *AddressBookQuery) Execute(client *Client) (NodeAddressBook, error) {
	var cancel func()
	var ctx context.Context

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
						panic(grpcErr.Err())
					}
				} else if err == io.EOF {
					break
				} else {
					panic(err)
				}
			}

			if subClient == nil {
				ctx, cancel = context.WithCancel(context.TODO())

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
	}, nil
}

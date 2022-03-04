package hedera

import (
	"context"
	"github.com/hashgraph/hedera-protobufs-go/services"
	"io"
	"math"
	"regexp"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/mirror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var rstStream = regexp.MustCompile("(?i)\\brst[^0-9a-zA-Z]stream\\b") //nolint

type TopicMessageQuery struct {
	errorHandler      func(stat status.Status)
	completionHandler func()
	retryHandler      func(err error) bool
	attempt           uint64
	maxAttempts       uint64
	topicID           *TopicID
	startTime         *time.Time
	endTime           *time.Time
	limit             uint64
}

func NewTopicMessageQuery() *TopicMessageQuery {
	return &TopicMessageQuery{
		maxAttempts:       maxAttempts,
		errorHandler:      _DefaultErrorHandler,
		retryHandler:      _DefaultRetryHandler,
		completionHandler: _DefaultCompletionHandler,
	}
}

func (query *TopicMessageQuery) SetTopicID(topicID TopicID) *TopicMessageQuery {
	query.topicID = &topicID
	return query
}

func (query *TopicMessageQuery) GetTopicID() TopicID {
	if query.topicID == nil {
		return TopicID{}
	}

	return *query.topicID
}

func (query *TopicMessageQuery) SetStartTime(startTime time.Time) *TopicMessageQuery {
	query.startTime = &startTime
	return query
}

func (query *TopicMessageQuery) GetStartTime() time.Time {
	if query.startTime != nil {
		return *query.startTime
	}

	return time.Time{}
}

func (query *TopicMessageQuery) SetEndTime(endTime time.Time) *TopicMessageQuery {
	query.endTime = &endTime
	return query
}

func (query *TopicMessageQuery) GetEndTime() time.Time {
	if query.endTime != nil {
		return *query.endTime
	}

	return time.Time{}
}

func (query *TopicMessageQuery) SetLimit(limit uint64) *TopicMessageQuery {
	query.limit = limit
	return query
}

func (query *TopicMessageQuery) GetLimit() uint64 {
	return query.limit
}

func (query *TopicMessageQuery) SetMaxAttempts(maxAttempts uint64) *TopicMessageQuery {
	query.maxAttempts = maxAttempts
	return query
}

func (query *TopicMessageQuery) GetMaxAttempts() uint64 {
	return query.maxAttempts
}

func (query *TopicMessageQuery) SetErrorHandler(errorHandler func(stat status.Status)) *TopicMessageQuery {
	query.errorHandler = errorHandler
	return query
}

func (query *TopicMessageQuery) SetCompletionHandler(completionHandler func()) *TopicMessageQuery {
	query.completionHandler = completionHandler
	return query
}

func (query *TopicMessageQuery) SetRetryHandler(retryHandler func(err error) bool) *TopicMessageQuery {
	query.retryHandler = retryHandler
	return query
}

func (query *TopicMessageQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.topicID != nil {
		if err := query.topicID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *TopicMessageQuery) _Build() *mirror.ConsensusTopicQuery {
	body := &mirror.ConsensusTopicQuery{
		Limit: query.limit,
	}
	if query.topicID != nil {
		body.TopicID = query.topicID._ToProtobuf()
	}

	if query.startTime != nil {
		body.ConsensusStartTime = _TimeToProtobuf(*query.startTime)
	} else {
		body.ConsensusStartTime = &services.Timestamp{}
	}

	if query.endTime != nil {
		body.ConsensusEndTime = _TimeToProtobuf(*query.endTime)
	}

	return body
}

func (query *TopicMessageQuery) Subscribe(client *Client, onNext func(TopicMessage)) (SubscriptionHandle, error) {
	handle := SubscriptionHandle{}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return SubscriptionHandle{}, err
	}

	pb := query._Build()

	messages := make(map[string][]*mirror.ConsensusTopicResponse)

	channel, err := client.mirrorNetwork._GetNextMirrorNode()._GetConsensusServiceClient()
	if err != nil {
		return handle, err
	}

	go func() {
		var subClient mirror.ConsensusService_SubscribeTopicClient
		var err error

		for {
			if err != nil {
				handle.Unsubscribe()

				if grpcErr, ok := status.FromError(err); ok { // nolint
					if query.attempt < query.maxAttempts && query.retryHandler(err) {
						subClient = nil

						delay := math.Min(250.0*math.Pow(2.0, float64(query.attempt)), 8000)
						time.Sleep(time.Duration(delay) * time.Millisecond)
						query.attempt++
					} else {
						query.errorHandler(*grpcErr)
						break
					}
				} else if err == io.EOF {
					query.completionHandler()
					break
				} else {
					panic(err)
				}
			}

			if subClient == nil {
				ctx, cancel := context.WithCancel(context.TODO())
				handle.onUnsubscribe = cancel

				subClient, err = (*channel).SubscribeTopic(ctx, pb)

				if err != nil {
					continue
				}
			}

			var resp *mirror.ConsensusTopicResponse
			resp, err = subClient.Recv()

			if err != nil {
				continue
			}

			if resp.ConsensusTimestamp != nil {
				pb.ConsensusStartTime = _TimeToProtobuf(_TimeFromProtobuf(resp.ConsensusTimestamp).Add(1 * time.Nanosecond))
			}

			if pb.Limit > 0 {
				pb.Limit--
			}

			if resp.ChunkInfo == nil || resp.ChunkInfo.Total == 1 {
				onNext(_TopicMessageOfSingle(resp))
			} else {
				txID := _TransactionIDFromProtobuf(resp.ChunkInfo.InitialTransactionID).String()
				message, ok := messages[txID]
				if !ok {
					message = make([]*mirror.ConsensusTopicResponse, 0, resp.ChunkInfo.Total)
				}

				message = append(message, resp)
				messages[txID] = message

				if int32(len(message)) == resp.ChunkInfo.Total {
					delete(messages, txID)

					onNext(_TopicMessageOfMany(message))
				}
			}
		}
	}()

	return handle, nil
}

func _DefaultErrorHandler(stat status.Status) {
	println("Failed to subscribe to topic with status", stat.Code().String())
}

func _DefaultCompletionHandler() {
	println("Subscription to topic finished")
}

func _DefaultRetryHandler(err error) bool {
	code := status.Code(err)

	switch code {
	case codes.NotFound, codes.ResourceExhausted, codes.Unavailable:
		return true
	case codes.Internal:
		grpcErr, ok := status.FromError(err)

		if !ok {
			return false
		}

		return rstStream.FindIndex([]byte(grpcErr.Message())) != nil
	default:
		return false
	}
}

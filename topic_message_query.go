package hedera

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto/mirror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TopicMessageQuery struct {
	pb                *mirror.ConsensusTopicQuery
	errorHandler      func(stat status.Status)
	completionHandler func()
	retryHandler      func(err error) bool
	counter           uint64
	attempt           uint64
	maxAttempts       uint64
	limit             *uint64
}

func NewTopicMessageQuery() *TopicMessageQuery {
	return &TopicMessageQuery{
		pb:                &mirror.ConsensusTopicQuery{},
		maxAttempts:       maxAttempts,
		errorHandler:      nil,
		retryHandler:      defaultRetryHandler,
		completionHandler: nil,
	}
}

func (query *TopicMessageQuery) SetTopicID(topicID TopicID) *TopicMessageQuery {
	query.pb.TopicID = topicID.toProtobuf()
	return query
}

func (query *TopicMessageQuery) GetTopicID() TopicID {
	if query.pb.TopicID != nil {
		return topicIDFromProtobuf(query.pb.TopicID)
	} else {
		return TopicID{}
	}
}

func (query *TopicMessageQuery) SetStartTime(startTime time.Time) *TopicMessageQuery {
	query.pb.ConsensusStartTime = timeToProtobuf(startTime)
	return query
}

func (query *TopicMessageQuery) GetStartTime() time.Time {
	if query.pb.ConsensusStartTime != nil {
		return timeFromProtobuf(query.pb.ConsensusStartTime)
	} else {
		return time.Time{}
	}
}

func (query *TopicMessageQuery) SetEndTime(endTime time.Time) *TopicMessageQuery {
	query.pb.ConsensusEndTime = timeToProtobuf(endTime)
	return query
}

func (query *TopicMessageQuery) GetEndTime() time.Time {
	if query.pb.ConsensusEndTime != nil {
		return timeFromProtobuf(query.pb.ConsensusEndTime)
	} else {
		return time.Time{}
	}
}

func (query *TopicMessageQuery) SetLimit(limit uint64) *TopicMessageQuery {
	query.limit = &limit
	return query
}

func (query *TopicMessageQuery) GetLimit() uint64 {
	if query.limit != nil {
		return *query.limit
	} else {
		return 0
	}
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

func (query *TopicMessageQuery) Subscribe(client *Client, onNext func(TopicMessage)) (SubscriptionHandle, error) {
	handle := SubscriptionHandle{}

	messages := sync.Map{}
	messagesMutex := sync.Mutex{}

	channel, err := client.mirrorNetwork.getNextMirrorNode().getChannel()
	if err != nil {
		return handle, err
	}

	go func() {
		var subClient mirror.ConsensusService_SubscribeTopicClient
		var err error

		for {
			if query.attempt <= query.maxAttempts && subClient == nil {
				if query.limit != nil {
					query.pb.Limit = *query.limit - query.counter
				}

				ctx, cancel := context.WithCancel(context.TODO())
				handle.onUnsubscribe = cancel

				subClient, err = (*channel).SubscribeTopic(ctx, query.pb)

				if err != nil {
					handle.Unsubscribe()
					callErrorHandlerWithGrpcStatus(err, query.errorHandler)
					subClient = nil
				}
			}

			resp, err := subClient.Recv()

			if err != nil {
				if query.attempt <= query.maxAttempts && query.retryHandler(err) {
					handle.Unsubscribe()
					subClient = nil

					delay := 250.0 * math.Pow(2.0, float64(query.attempt))
					time.Sleep(time.Duration(delay) * time.Millisecond)
					query.attempt += 1
					continue
				}

				handle.Unsubscribe()
				callErrorHandlerWithGrpcStatus(err, query.errorHandler)
				break
			}

			if resp.ChunkInfo == nil || (resp.ChunkInfo != nil && resp.ChunkInfo.Total == 1) {
				query.counter += 1

				onNext(topicMessageOfSingle(resp))
			} else {
				messagesMutex.Lock()
				txID := transactionIDFromProtobuf(resp.ChunkInfo.InitialTransactionID).String()
				messageI, _ := messages.LoadOrStore(txID, make([]*mirror.ConsensusTopicResponse, 0, resp.ChunkInfo.Total))

				message := messageI.([]*mirror.ConsensusTopicResponse)
				message = append(message, resp)

				messages.Store(txID, message)

				if int32(len(message)) == resp.ChunkInfo.Total {
					query.counter += 1

					messages.Delete(txID)
					messagesMutex.Unlock()
					onNext(topicMessageOfMany(message))
				} else {
					messagesMutex.Unlock()
				}

			}

			if query.limit != nil && query.counter == *query.limit {
				query.completionHandler()
				break
			}
		}
	}()

	return handle, nil
}

func defaultRetryHandler(err error) bool {
	code := status.Code(err)

	switch code {
	case codes.NotFound, codes.ResourceExhausted, codes.Internal, codes.Unavailable:
		return true
	default:
		return false
	}
}

func callErrorHandlerWithGrpcStatus(err error, errorHandler func(stat status.Status)) {
	if grpcErr, ok := status.FromError(err); errorHandler != nil && ok {
		errorHandler(*grpcErr)
	}
}

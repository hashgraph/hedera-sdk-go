package hedera

import (
	"context"
	"io"
	"math"
	"regexp"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto/mirror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var RST_STREAM *regexp.Regexp

type TopicMessageQuery struct {
	pb                *mirror.ConsensusTopicQuery
	errorHandler      func(stat status.Status)
	completionHandler func()
	retryHandler      func(err error) bool
	attempt           uint64
	maxAttempts       uint64
	topicID           TopicID
}

func NewTopicMessageQuery() *TopicMessageQuery {
	return &TopicMessageQuery{
		pb:                &mirror.ConsensusTopicQuery{},
		maxAttempts:       maxAttempts,
		errorHandler:      defaultErrorHandler,
		retryHandler:      defaultRetryHandler,
		completionHandler: defaultCompletionHandler,
	}
}

func (query *TopicMessageQuery) SetTopicID(id TopicID) *TopicMessageQuery {
	query.topicID = id
	return query
}

func (query *TopicMessageQuery) GetTopicID() TopicID {
	return query.topicID
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
	query.pb.Limit = limit
	return query
}

func (query *TopicMessageQuery) GetLimit() uint64 {
	return query.pb.Limit
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

func (query *TopicMessageQuery) validateNetworkOnIDs(client *Client) error {
	var err error
	err = query.topicID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *TopicMessageQuery) build() *TopicMessageQuery {
	if !query.topicID.isZero() {
		query.pb.TopicID = query.topicID.toProtobuf()
	}

	return query
}

func (query *TopicMessageQuery) Subscribe(client *Client, onNext func(TopicMessage)) (SubscriptionHandle, error) {
	handle := SubscriptionHandle{}

	query.topicID.setNetworkWithClient(client)
	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return SubscriptionHandle{}, err
	}

	query.build()

	messages := make(map[string][]*mirror.ConsensusTopicResponse, 0)

	channel, err := client.mirrorNetwork.getNextMirrorNode().getChannel()
	if err != nil {
		return handle, err
	}

	go func() {
		var subClient mirror.ConsensusService_SubscribeTopicClient
		var err error

		for {
			if err != nil {
				if grpcErr, ok := status.FromError(err); ok {
					if query.attempt < query.maxAttempts && query.retryHandler(err) {
						subClient = nil

						delay := math.Min(250.0*math.Pow(2.0, float64(query.attempt)), 8000)
						time.Sleep(time.Duration(delay) * time.Millisecond)
						query.attempt += 1
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

				subClient, err = (*channel).SubscribeTopic(ctx, query.pb)

				if err != nil {
					continue
				}
			}

			var resp *mirror.ConsensusTopicResponse
			resp, err = subClient.Recv()

			if err != nil {
				continue
			}

			query.pb.ConsensusStartTime = resp.ConsensusTimestamp
			if query.pb.Limit > 0 {
				query.pb.Limit -= 1
			}

			if resp.ChunkInfo == nil || (resp.ChunkInfo != nil && resp.ChunkInfo.Total == 1) {
				onNext(topicMessageOfSingle(resp))
			} else {
				txID := transactionIDFromProtobuf(resp.ChunkInfo.InitialTransactionID, nil).String()
				message, ok := messages[txID]
				if !ok {
					message = make([]*mirror.ConsensusTopicResponse, 0)
				}

				message = append(message, resp)
				messages[txID] = message

				if int32(len(message)) == resp.ChunkInfo.Total {
					delete(messages, txID)

					onNext(topicMessageOfMany(message, client.networkName))
				}
			}
		}
	}()

	return handle, nil
}

func defaultErrorHandler(stat status.Status) {
	println("Failed to subscribe to topic with status", stat.Code().String())
}

func defaultCompletionHandler() {
	println("Subscription to topic finished")
}

func defaultRetryHandler(err error) bool {
	code := status.Code(err)

	switch code {
	case codes.NotFound, codes.ResourceExhausted, codes.Unavailable:
		return true
	case codes.Internal:
		if RST_STREAM == nil {
			var err1 error
			RST_STREAM, err1 = regexp.Compile(".*(rst.stream.*internal.error|internal.error.*rst.stream).*")
			if err1 != nil {
				panic(err1)
			}
		}

		grpcErr, ok := status.FromError(err)

		if !ok {
			return false
		}

		return RST_STREAM.FindIndex([]byte(grpcErr.Message())) != nil
	default:
		return false
	}
}

func callErrorHandlerWithGrpcStatus(err error, errorHandler func(stat status.Status)) {
	if grpcErr, ok := status.FromError(err); errorHandler != nil && ok {
		errorHandler(*grpcErr)
	}
}

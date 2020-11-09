package hedera

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TopicMessageQuery struct {
	pb *mirror.ConsensusTopicQuery
}

func NewTopicMessageQuery() *TopicMessageQuery {
	pb := mirror.ConsensusTopicQuery{}
	return &TopicMessageQuery{
		pb: &pb,
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
	query.pb.Limit = limit
	return query
}

func (query *TopicMessageQuery) GetLimit() uint64 {
	return query.pb.Limit
}

func (query *TopicMessageQuery) Subscribe(client *Client, onNext func(TopicMessage)) (SubscriptionHandle, error) {
	ctx, cancel := context.WithCancel(context.TODO())

	handle := newSubscriptionHandle(cancel)

	messages := sync.Map{}
	messagesMutex := sync.Mutex{}

	go func() {
		var subClient mirror.ConsensusService_SubscribeTopicClient
		var err error
		attempt := 0
		resubscribe := true
		channel, err := client.mirrorNetwork.getNextMirrorNode().getChannel()
		if err != nil {
			panic(err)
		}

		for {
			if resubscribe {
				subClient, err = (*channel).SubscribeTopic(ctx, query.pb)
				if err != nil {
					panic(err)
				}
			}

			resp, err := subClient.Recv()
			code := status.Code(err)

			if err != nil {
				if code == codes.NotFound || code == codes.Unavailable {
					if attempt >= 10 {
						cancel()
					} else {
						delay := 250.0 * math.Pow(2.0, float64(attempt))
						time.Sleep(time.Duration(delay) * time.Millisecond)
						attempt += 1
						continue
					}
					break
				} else {
					cancel()
					break
				}
			}

			resubscribe = false

			if resp.ChunkInfo == nil {
				onNext(topicMessageOfSingle(resp))
			} else {
				messagesMutex.Lock()
				txID := transactionIDFromProtobuf(resp.ChunkInfo.InitialTransactionID)
				messageI, _ := messages.LoadOrStore(txID, make([]*mirror.ConsensusTopicResponse, 0, resp.ChunkInfo.Total))

				message := messageI.([]*mirror.ConsensusTopicResponse)
				message = append(message, resp)

				messages.Store(txID, message)

				if int32(len(message)) == resp.ChunkInfo.Total {
					messages.Delete(txID)
					messagesMutex.Unlock()
					onNext(topicMessageOfMany(message))
				} else {
					messagesMutex.Unlock()
				}

			}
		}
	}()

	return handle, nil
}

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var rstStream = regexp.MustCompile("(?i)\\brst[^0-9a-zA-Z]stream\\b") //nolint

// TopicMessageQuery
// Query that listens to messages sent to the specific TopicID
type TopicMessageQuery struct {
	errorHandler      func(stat status.Status)
	completionHandler func()
	retryHandler      func(err error) bool
	maxAttempts       uint64
	topicID           *TopicID
	startTime         *time.Time
	endTime           *time.Time
	limit             uint64
}

// NewTopicMessageQuery creates TopicMessageQuery which
// listens to messages sent to the specific TopicID
func NewTopicMessageQuery() *TopicMessageQuery {
	return &TopicMessageQuery{
		maxAttempts:       maxAttempts,
		errorHandler:      _DefaultErrorHandler,
		retryHandler:      _DefaultRetryHandler,
		completionHandler: _DefaultCompletionHandler,
	}
}

// SetTopicID Sets topic ID to retrieve messages for.
// Required
func (query *TopicMessageQuery) SetTopicID(topicID TopicID) *TopicMessageQuery {
	query.topicID = &topicID
	return query
}

// GetTopicID returns the TopicID for this TopicMessageQuery
func (query *TopicMessageQuery) GetTopicID() TopicID {
	if query.topicID == nil {
		return TopicID{}
	}

	return *query.topicID
}

// SetStartTime Sets time for when to start listening for messages. Defaults to current time if
// not set.
func (query *TopicMessageQuery) SetStartTime(startTime time.Time) *TopicMessageQuery {
	query.startTime = &startTime
	return query
}

// GetStartTime returns the start time for this TopicMessageQuery
func (query *TopicMessageQuery) GetStartTime() time.Time {
	if query.startTime != nil {
		return *query.startTime
	}

	return time.Time{}
}

// SetEndTime Sets time when to stop listening for messages. If not set it will receive
// indefinitely.
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

// SetLimit Sets the maximum number of messages to receive before stopping. If not set or set to zero it will
// return messages indefinitely.
func (query *TopicMessageQuery) SetLimit(limit uint64) *TopicMessageQuery {
	query.limit = limit
	return query
}

func (query *TopicMessageQuery) GetLimit() uint64 {
	return query.limit
}

// SetMaxAttempts Sets the amount of consecutive attempts to reconnect after the subscription stream fails
func (query *TopicMessageQuery) SetMaxAttempts(maxAttempts uint64) *TopicMessageQuery {
	query.maxAttempts = maxAttempts
	return query
}

// GetMaxAttempts returns the amount of attempts to try to retrieve message
func (query *TopicMessageQuery) GetMaxAttempts() uint64 {
	return query.maxAttempts
}

// SetErrorHandler Sets the error handler for this query
func (query *TopicMessageQuery) SetErrorHandler(errorHandler func(stat status.Status)) *TopicMessageQuery {
	query.errorHandler = errorHandler
	return query
}

// SetCompletionHandler Sets the completion handler for this query
func (query *TopicMessageQuery) SetCompletionHandler(completionHandler func()) *TopicMessageQuery {
	query.completionHandler = completionHandler
	return query
}

// SetRetryHandler Sets the retry handler for this query
func (query *TopicMessageQuery) SetRetryHandler(retryHandler func(err error) bool) *TopicMessageQuery {
	query.retryHandler = retryHandler
	return query
}

func (query *TopicMessageQuery) validateNetworkOnIDs(client *Client) error {
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

func (query *TopicMessageQuery) build() *mirror.ConsensusTopicQuery {
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

// Subscribe subscribes to messages sent to the specific TopicID
func (query *TopicMessageQuery) Subscribe(client *Client, onNext func(TopicMessage)) (SubscriptionHandle, error) {
	handle := SubscriptionHandle{}

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return SubscriptionHandle{}, err
	}

	pbBody := query.build()

	mirrorNode, err := client.mirrorNetwork._GetNextMirrorNode()
	if err != nil {
		return handle, err
	}

	channel, err := mirrorNode._GetConsensusServiceClient()
	if err != nil {
		return handle, err
	}

	ctx, cancel := context.WithCancel(context.Background()) //nolint:gosec // G118: cancel is stored in handle.onUnsubscribe and called on unsubscribe
	handle.onUnsubscribe = cancel

	stream, err := channel.SubscribeTopic(ctx, pbBody)
	if err != nil {
		cancel()
		return SubscriptionHandle{}, err
	}
	messages := make(map[string][]*mirror.ConsensusTopicResponse)
	retryHandler, maxAttempts, limited := query.retryHandler, query.maxAttempts, query.limit > 0

	// consumeMessages reopens the stream from the last received message until the subscription ends.
	consumeMessages := func() {
		var attempt uint64
		var err error

		for {
			if err != nil {
				switch {
				case ctx.Err() != nil, err == io.EOF, limited && pbBody.Limit == 0:
					query.completionHandler()
					return
				case !retryHandler(err):
					query.errorHandler(_TopicMessageQueryErrorStatus(err))
					return
				case attempt >= maxAttempts:
					query.errorHandler(_TopicMessageQueryErrorStatus(fmt.Errorf("max attempts reached: %w", err)))
					return
				}

				if !waitMirrorStreamBackoff(ctx, attempt) {
					query.completionHandler()
					return
				}
				attempt++

				if stream, err = channel.SubscribeTopic(ctx, pbBody); err != nil {
					continue
				}
			}

			var resp *mirror.ConsensusTopicResponse
			if resp, err = stream.Recv(); err != nil {
				continue
			}
			attempt = 0

			if resp.ConsensusTimestamp != nil {
				pbBody.ConsensusStartTime = _TimeToProtobuf(_TimeFromProtobuf(resp.ConsensusTimestamp).Add(time.Nanosecond))
			}
			if pbBody.Limit > 0 {
				pbBody.Limit--
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
	}

	go consumeMessages()
	return handle, nil
}

func _TopicMessageQueryErrorStatus(err error) status.Status {
	if grpcErr, ok := status.FromError(err); ok {
		return *grpcErr
	}

	return *status.New(codes.Unknown, "Unknown error occurred")
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

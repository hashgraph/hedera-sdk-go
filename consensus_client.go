package hedera

import (
	"context"
	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
	"google.golang.org/grpc"
	"time"
)

type ErrorHandler func(error)
type Listener func(ConsensusMessage)

type ConsensusClient struct {
	client       mirror.ConsensusServiceClient
	errorHandler ErrorHandler
}

type ConsensusClientSubscription struct {
	topicID            ConsensusTopicID
	consensusStartTime *time.Time
	client             mirror.ConsensusService_SubscribeTopicClient
	errorHandler       ErrorHandler
	listener           Listener
}

func NewConsensusClient(endpoint string) (*ConsensusClient, error) {
	client, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &ConsensusClient{
		client:       mirror.NewConsensusServiceClient(client),
		errorHandler: nil,
	}, nil
}

func (self *ConsensusClient) SetErrorHandler(errorHandler ErrorHandler) *ConsensusClient {
	self.errorHandler = errorHandler
	return self
}

func (self *ConsensusClient) Subscribe(
	topicID ConsensusTopicID,
	startTime *time.Time,
	listener Listener,
) (ConsensusClientSubscription, error) {
	topicQuery := mirror.ConsensusTopicQuery{TopicID: topicID.toProto()}

	if startTime != nil {
		topicQuery.ConsensusStartTime = timeToProto(*startTime)
	}

	subscriptionClient, err := self.client.SubscribeTopic(context.TODO(), &topicQuery)
	if err != nil {
		return ConsensusClientSubscription{}, err
	}

	subscription := NewConsensusClientSubscription(
		topicID,
		startTime,
		subscriptionClient,
		self.errorHandler,
		listener,
	)

	go subscriptionHandler(subscription)
	return subscription, nil
}

func NewConsensusClientSubscription(
	topicID ConsensusTopicID,
	consensusStartTime *time.Time,
	client mirror.ConsensusService_SubscribeTopicClient,
	errorHandler ErrorHandler,
	listener Listener,
) ConsensusClientSubscription {
	return ConsensusClientSubscription{
		topicID:            topicID,
		consensusStartTime: consensusStartTime,
		client:             client,
		errorHandler:       errorHandler,
		listener:           listener,
	}
}

func (sub ConsensusClientSubscription) recv() (ConsensusMessage, error) {
	resp, err := sub.client.Recv()
	if err != nil {
		return ConsensusMessage{}, err
	} else {
		return NewConsensusMessage(sub.topicID, resp), nil
	}
}

func (sub ConsensusClientSubscription) Unsubscribe() {
	sub.client.SendMsg("unsubscribe from topic")
	sub.client.CloseSend()
}

func subscriptionHandler(sub ConsensusClientSubscription) {
	for {
		message, err := sub.recv()
		if err != nil {
			sub.errorHandler(err)
			sub.client.CloseSend()
			return
		} else {
			sub.listener(message)
		}
	}
}

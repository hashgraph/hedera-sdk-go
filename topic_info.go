package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type TopicInfo struct {
	Memo               string
	RunningHash        []byte
	SequenceNumber     uint64
	ExpirationTime     time.Time
	AdminKey           Key
	SubmitKey          Key
	AutoRenewPeriod    time.Duration
	AutoRenewAccountID *AccountID
}

func topicInfoFromProtobuf(topicInfo *proto.ConsensusTopicInfo) (TopicInfo, error) {
	var err error
	tempTopicInfo := TopicInfo{
		Memo:           topicInfo.Memo,
		RunningHash:    topicInfo.RunningHash,
		SequenceNumber: topicInfo.SequenceNumber,
		ExpirationTime: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
			time.Now().Hour(), time.Now().Minute(), int(topicInfo.ExpirationTime.Seconds),
			int(topicInfo.ExpirationTime.Nanos), time.Now().Location()),
		AutoRenewPeriod: durationFromProtobuf(topicInfo.AutoRenewPeriod),
	}

	if adminKey := topicInfo.AdminKey; adminKey != nil {
		tempTopicInfo.AdminKey, err = keyFromProtobuf(adminKey)
	}

	if submitKey := topicInfo.SubmitKey; submitKey != nil {
		tempTopicInfo.SubmitKey, err = keyFromProtobuf(submitKey)
	}

	if ARAccountID := topicInfo.AutoRenewAccount; ARAccountID != nil {
		ID := accountIDFromProtobuf(ARAccountID)

		tempTopicInfo.AutoRenewAccountID = &ID
	}
	return tempTopicInfo, err
}

func (topicInfo *TopicInfo) toProtobuf() *proto.ConsensusTopicInfo {
	return &proto.ConsensusTopicInfo{
		Memo:           topicInfo.Memo,
		RunningHash:    topicInfo.RunningHash,
		SequenceNumber: topicInfo.SequenceNumber,
		ExpirationTime: &proto.Timestamp{
			Seconds: int64(topicInfo.ExpirationTime.Second()),
			Nanos:   int32(topicInfo.ExpirationTime.Nanosecond()),
		},
		AdminKey:         topicInfo.AdminKey.toProtoKey(),
		SubmitKey:        topicInfo.SubmitKey.toProtoKey(),
		AutoRenewPeriod:  durationToProtobuf(topicInfo.AutoRenewPeriod),
		AutoRenewAccount: topicInfo.AutoRenewAccountID.toProtobuf(),
	}
}

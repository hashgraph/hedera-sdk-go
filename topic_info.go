package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type TopicInfo struct {
	TopicMemo          string
	RunningHash        []byte
	SequenceNumber     uint64
	ExpirationTime     time.Time
	AdminKey           Key
	SubmitKey          Key
	AutoRenewPeriod    time.Duration
	AutoRenewAccountID *AccountID
}

func topicInfoFromProtobuf(topicInfo *proto.ConsensusTopicInfo) (TopicInfo, error) {
	if topicInfo == nil {
		return TopicInfo{}, errParameterNull
	}
	var err error
	tempTopicInfo := TopicInfo{
		TopicMemo:      topicInfo.Memo,
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
		Memo:           topicInfo.TopicMemo,
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

func (topicInfo TopicInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(topicInfo.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TopicInfoFromBytes(data []byte) (TopicInfo, error) {
	if data == nil {
		return TopicInfo{}, errByteArrayNull
	}
	pb := proto.ConsensusTopicInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TopicInfo{}, err
	}

	info, err := topicInfoFromProtobuf(&pb)
	if err != nil {
		return TopicInfo{}, err
	}

	return info, nil
}

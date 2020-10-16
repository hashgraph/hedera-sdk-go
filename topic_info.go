package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type TopicInfo struct {
	Memo               string
	RunningHash        []byte
	SequenceNumber     uint64
	ExpirationTime     time.Time
	AdminKey           *PublicKey
	SubmitKey          *PublicKey
	AutoRenewPeriod    time.Duration
	AutoRenewAccountID *AccountID
}

func topicInfoFromProtobuf(topicInfo *proto.ConsensusTopicInfo) TopicInfo {
	return TopicInfo{
		Memo:               topicInfo.Memo,
		RunningHash:        topicInfo.RunningHash,
		SequenceNumber:     topicInfo.SequenceNumber,
		ExpirationTime:     time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
										time.Now().Hour(), time.Now().Minute(), int(topicInfo.ExpirationTime.Seconds),
										int(topicInfo.ExpirationTime.Nanos), time.Now().Location()),
		AdminKey:           &PublicKey{
			keyData: topicInfo.AdminKey.GetEd25519(),
		},
		SubmitKey:          &PublicKey{
			keyData: topicInfo.SubmitKey.GetEd25519(),
		},
		AutoRenewPeriod:    durationFromProtobuf(topicInfo.AutoRenewPeriod),
		AutoRenewAccountID: &AccountID{
			Shard:   uint64(topicInfo.GetAutoRenewAccount().ShardNum),
			Realm:   uint64(topicInfo.GetAutoRenewAccount().RealmNum),
			Account: uint64(topicInfo.GetAutoRenewAccount().AccountNum),
		},
	}
}

func (topicInfo *TopicInfo) toProtobuf() *proto.ConsensusTopicInfo {
	return &proto.ConsensusTopicInfo{
		Memo:             topicInfo.Memo,
		RunningHash:      topicInfo.RunningHash,
		SequenceNumber:   topicInfo.SequenceNumber,
		ExpirationTime:   &proto.Timestamp{
			Seconds: int64(topicInfo.ExpirationTime.Second()),
			Nanos:   int32(topicInfo.ExpirationTime.Nanosecond()),
		},
		AdminKey:         topicInfo.AdminKey.toProtoKey(),
		SubmitKey:        topicInfo.SubmitKey.toProtoKey(),
		AutoRenewPeriod:  durationToProtobuf(topicInfo.AutoRenewPeriod),
		AutoRenewAccount: topicInfo.AutoRenewAccountID.toProtobuf(),
	}
}

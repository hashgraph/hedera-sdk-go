package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

func durationToProto(duration time.Duration) *proto.Duration {
	return &proto.Duration{
		Seconds: int64(duration.Seconds()),
	}
}

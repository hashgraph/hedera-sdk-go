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

func durationFromProto(pb *proto.Duration) time.Duration {
	return time.Until(time.Unix(pb.Seconds, 0))
}

func timeToProto(t time.Time) *proto.Timestamp {
	return &proto.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.UnixNano() - (t.Unix() * 1e+9)),
	}
}

func timeFromProto(pb *proto.Timestamp) time.Time {
	return time.Unix(pb.Seconds, int64(pb.Nanos))
}

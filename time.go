package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

func durationToProtobuf(duration time.Duration) *proto.Duration {
	return &proto.Duration{
		Seconds: int64(duration.Seconds()),
	}
}

func durationFromProtobuf(pb *proto.Duration) time.Duration {
	if pb == nil {
		return time.Duration(0)
	}
	return time.Until(time.Unix(pb.Seconds, 0))
}

func timeToProtobuf(t time.Time) *proto.Timestamp {
	return &proto.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.UnixNano() - (t.Unix() * 1e+9)),
	}
}

func timeFromProtobuf(pb *proto.Timestamp) time.Time {
	if pb == nil {
		return time.Time{}
	}
	return time.Unix(pb.Seconds, int64(pb.Nanos))
}

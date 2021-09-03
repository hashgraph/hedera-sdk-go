package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

func _DurationToProtobuf(duration time.Duration) *proto.Duration {
	return &proto.Duration{
		Seconds: int64(duration.Seconds()),
	}
}

func _DurationFromProtobuf(pb *proto.Duration) time.Duration {
	if pb == nil {
		return time.Duration(0)
	}
	return time.Duration(pb.Seconds * int64(time.Second))
}

func _TimeToProtobuf(t time.Time) *proto.Timestamp {
	return &proto.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.UnixNano() - (t.Unix() * 1e+9)),
	}
}

func _TimeFromProtobuf(pb *proto.Timestamp) time.Time {
	if pb == nil {
		return time.Time{}
	}
	return time.Unix(pb.Seconds, int64(pb.Nanos))
}

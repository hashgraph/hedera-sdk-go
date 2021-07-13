package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"time"
)

func durationToProtobuf(duration time.Duration) *services.Duration {
	return &services.Duration{
		Seconds: int64(duration.Seconds()),
	}
}

func durationFromProtobuf(pb *services.Duration) time.Duration {
	if pb == nil {
		return time.Duration(0)
	}
	return time.Duration(pb.Seconds * int64(time.Second))
}

func timeToProtobuf(t time.Time) *services.Timestamp {
	return &services.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.UnixNano() - (t.Unix() * 1e+9)),
	}
}

func timeFromProtobuf(pb *services.Timestamp) time.Time {
	if pb == nil {
		return time.Time{}
	}
	return time.Unix(pb.Seconds, int64(pb.Nanos))
}

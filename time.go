package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

func _DurationToProtobuf(duration time.Duration) *services.Duration {
	return &services.Duration{
		Seconds: int64(duration.Seconds()),
	}
}

func _DurationFromProtobuf(pb *services.Duration) time.Duration {
	if pb == nil {
		return time.Duration(0)
	}
	return time.Duration(pb.Seconds * int64(time.Second))
}

func _TimeToProtobuf(t time.Time) *services.Timestamp {
	return &services.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.UnixNano() - (t.Unix() * 1e+9)),
	}
}

func _TimeFromProtobuf(pb *services.Timestamp) time.Time {
	if pb == nil {
		return time.Time{}
	}
	return time.Unix(pb.Seconds, int64(pb.Nanos))
}

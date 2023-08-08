package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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

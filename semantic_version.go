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
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type SemanticVersion struct {
	Major uint32
	Minor uint32
	Patch uint32
	Pre   string
	Build string
}

func _SemanticVersionFromProtobuf(version *services.SemanticVersion) SemanticVersion {
	if version == nil {
		return SemanticVersion{}
	}
	return SemanticVersion{
		Major: uint32(version.GetMajor()),
		Minor: uint32(version.GetMinor()),
		Patch: uint32(version.GetPatch()),
		Pre:   version.GetPre(),
		Build: version.GetBuild(),
	}
}

func (version *SemanticVersion) _ToProtobuf() *services.SemanticVersion {
	return &services.SemanticVersion{
		Major: int32(version.Major),
		Minor: int32(version.Minor),
		Patch: int32(version.Patch),
		Pre:   version.Pre,
		Build: version.Build,
	}
}

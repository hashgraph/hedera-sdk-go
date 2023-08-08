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

import "math/rand"

type _MirrorNetwork struct {
	_ManagedNetwork
}

func _NewMirrorNetwork() *_MirrorNetwork {
	return &_MirrorNetwork{
		_ManagedNetwork: _NewManagedNetwork(),
	}
}

func (network *_MirrorNetwork) _SetNetwork(newNetwork []string) (err error) {
	newMirrorNetwork := make(map[string]_IManagedNode)
	for _, url := range newNetwork {
		if newMirrorNetwork[url], err = _NewMirrorNode(url); err != nil {
			return err
		}
	}

	return network._ManagedNetwork._SetNetwork(newMirrorNetwork)
}

func (network *_MirrorNetwork) _GetNetwork() []string {
	temp := make([]string, 0)
	for url := range network._ManagedNetwork.network { //nolint
		temp = append(temp, url)
	}

	return temp
}

// nolint:unused
// Deprecated: _SetTransportSecurity is no longer supported, as only secured connections are now allowed.
func (network *_MirrorNetwork) _SetTransportSecurity(transportSecurity bool) *_MirrorNetwork {
	return network
}

func (network *_MirrorNetwork) _GetNextMirrorNode() *_MirrorNode {
	node := network._ManagedNetwork.healthyNodes[rand.Intn(len(network.healthyNodes))] // nolint
	if node, ok := node.(*_MirrorNode); ok {
		return node
	}
	return &_MirrorNode{}
}

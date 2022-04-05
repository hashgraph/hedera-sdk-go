package hedera

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

func (network *_MirrorNetwork) _SetTransportSecurity(transportSecurity bool) *_MirrorNetwork {
	_ = network._ManagedNetwork._SetTransportSecurity(transportSecurity)
	return network
}

func (network *_MirrorNetwork) _GetNextMirrorNode() *_MirrorNode {
	node := network._ManagedNetwork.healthyNodes[rand.Intn(len(network.healthyNodes))] // nolint
	if node, ok := node.(*_MirrorNode); ok {
		return node
	}
	return &_MirrorNode{}
}

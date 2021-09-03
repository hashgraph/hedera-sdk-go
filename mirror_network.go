package hedera

import (
	"math/rand"
)

type _MirrorNetwork struct {
	networkNodes map[string]*_MirrorNode
	network      []string
	index        uint
}

func _NewMirrorNetwork() *_MirrorNetwork {
	return &_MirrorNetwork{
		networkNodes: make(map[string]*_MirrorNode),
		network:      make([]string, 0),
		index:        0,
	}
}

func _Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (network *_MirrorNetwork) _SetNetwork(newNetwork []string) {
	for _, n := range network.network {
		if !_Contains(newNetwork, n) {
			_ = network.networkNodes[n]._Close()
			delete(network.networkNodes, n)
		}
	}

	network.network = newNetwork

	for _, url := range newNetwork {
		if _, ok := network.networkNodes[url]; !ok {
			network.networkNodes[url] = _NewMirrorNode(url)
		}
	}

	network.index = 0

	if len(network.network) > 1 {
		rand.Shuffle(len(network.network), func(i, j int) {
			network.network[i], network.network[j] = network.network[j], network.network[i]
		})
	}
}

func (network *_MirrorNetwork) _GetNextMirrorNode() *_MirrorNode {
	node := network.networkNodes[network.network[network.index]]
	network.index = (network.index + 1) % uint(len(network.network))
	return node
}

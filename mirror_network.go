package hedera

import (
	"math/rand"
)

type mirrorNetwork struct {
	networkNodes map[string]*mirrorNode
	network      []string
	index        uint
}

func newMirrorNetwork() *mirrorNetwork {
	return &mirrorNetwork{
		networkNodes: make(map[string]*mirrorNode),
		network:      make([]string, 0),
		index:        0,
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (network *mirrorNetwork) setNetwork(newNetwork []string) {
	for _, n := range network.network {
		if !contains(newNetwork, n) {
			delete(network.networkNodes, n)
		}
	}

	for _, url := range newNetwork {
		if !contains(network.network, url) {
			network.network = append(network.network, url)
			network.networkNodes[url] = newMirrorNode(url)
		}
	}

	network.index = 0

	if len(network.network) > 0 {
		rand.Shuffle(len(network.network), func(i, j int) {
			network.network[i], network.network[j] = network.network[j], network.network[i]
		})
	}
}

func (network *mirrorNetwork) getNextMirrorNode() *mirrorNode {
	node := network.networkNodes[network.network[network.index]]
	network.index = (network.index + 1) % uint(len(network.network))
	return node
}

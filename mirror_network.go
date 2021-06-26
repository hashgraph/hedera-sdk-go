package hedera

import (
	"fmt"
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
			_ = network.networkNodes[n].close()
			delete(network.networkNodes, n)
		}
	}

	network.network = newNetwork

	for _, url := range newNetwork {
		if _, ok := network.networkNodes[url]; !ok {
			network.networkNodes[url] = newMirrorNode(url)
		}
	}

	network.index = 0

	if len(network.network) > 1 {
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

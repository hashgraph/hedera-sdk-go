package hedera

import (
	"math"
	"math/rand"
	"time"
)

type node struct {
	accountID AccountID
	address   string
	delay     int64
	lastUsed  *int64
}

type nodes struct {
	nodes []node
}

func newNode(accountID AccountID, address string) node {
	return node{
		accountID: accountID,
		address:   address,
		delay:     250,
		lastUsed:  nil,
	}
}

func (node *node) isHealthy() bool {
	if node.lastUsed != nil {
		lastUsed := *node.lastUsed
		return lastUsed+node.delay < time.Now().UTC().UnixNano()
	}

	return true
}

func (node *node) increaseDelay() {
	lastUsed := time.Now().UTC().UnixNano()
	node.lastUsed = &lastUsed
	node.delay = int64(math.Min(float64(node.delay)*2, 8000))
}

func (node *node) decreaseDelay() {
	node.delay = int64(math.Max(float64(node.delay)/2, 250))
}

func (node *node) wait() {
	var delay int64
	if node.lastUsed != nil {
		delay = *node.lastUsed + node.delay - time.Now().UTC().UnixNano()
	} else {
		delay = 0 - time.Now().UTC().UnixNano()
	}

	time.Sleep(time.Duration(delay) * time.Nanosecond)
}

func (s nodes) Len() int {
	return len(s.nodes)
}
func (s nodes) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}
func (s nodes) Less(i, j int) bool {
	if s.nodes[i].isHealthy() && s.nodes[j].isHealthy() {
		return rand.Int63n(1) == 0
	} else if s.nodes[i].isHealthy() && !s.nodes[j].isHealthy() {
		return true
	} else if !s.nodes[i].isHealthy() && s.nodes[j].isHealthy() {
		return false
	} else {
		aLastUsed := int64(0)
		bLastUsed := int64(0)

		if s.nodes[i].lastUsed != nil {
			aLastUsed = *s.nodes[i].lastUsed
		}

		if s.nodes[i].lastUsed == nil {
			bLastUsed = *s.nodes[j].lastUsed
		}

		return aLastUsed+s.nodes[i].delay < bLastUsed+s.nodes[j].delay
	}
}

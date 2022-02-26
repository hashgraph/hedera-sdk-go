package hedera

import (
	"time"
)

type _IManagedNode interface {
	_SetMinBackoff(waitTime time.Duration)
	_GetMinBackoff() time.Duration
	_SetMaxBackoff(waitTime time.Duration)
	_GetMaxBackoff() time.Duration
	_InUse()
	_IsHealthy() bool
	_IncreaseDelay()
	_DecreaseDelay()
	_Wait() time.Duration
	_GetUseCount() int64
	_GetLastUsed() int64
	_GetAttempts() int64
	_GetAddress() string
	_ToSecure() _IManagedNode
	_ToInsecure() _IManagedNode
	_GetManagedNode() *_ManagedNode
	_Close() error
}

type _ManagedNode struct {
	address            *_ManagedNodeAddress
	currentBackoff     time.Duration
	lastUsed           int64
	backoffUntil       time.Time
	useCount           int64
	minBackoff         time.Duration
	maxBackoff         time.Duration
	badGrpcStatusCount int64
}

func (node *_ManagedNode) _GetAttempts() int64 {
	return node.badGrpcStatusCount
}

func (node *_ManagedNode) _GetAddress() string {
	if node.address != nil {
		return node.address._String()
	}

	return ""
}

func _NewManagedNode(address string, minBackoff time.Duration) _ManagedNode {
	return _ManagedNode{
		address:            _ManagedNodeAddressFromString(address),
		currentBackoff:     minBackoff,
		useCount:           0,
		minBackoff:         minBackoff,
		maxBackoff:         1 * time.Hour,
		badGrpcStatusCount: 0,
	}
}

func (node *_ManagedNode) _SetMinBackoff(waitTime time.Duration) {
	if node.currentBackoff == node.minBackoff {
		node.currentBackoff = node.minBackoff
	}

	node.minBackoff = waitTime
}

func (node *_ManagedNode) _GetMinBackoff() time.Duration {
	return node.minBackoff
}

func (node *_ManagedNode) _SetMaxBackoff(waitTime time.Duration) {
	node.maxBackoff = waitTime
}

func (node *_ManagedNode) _GetMaxBackoff() time.Duration {
	return node.maxBackoff
}

func (node *_ManagedNode) _InUse() {
	node.useCount++
	node.lastUsed = time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _IsHealthy() bool {
	return node.backoffUntil.UnixNano() < time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _IncreaseDelay() {
	node.badGrpcStatusCount++
	node.backoffUntil = time.Now().Add(node.currentBackoff)
	node.currentBackoff = node.currentBackoff * 2
	if node.currentBackoff > node.maxBackoff {
		node.currentBackoff = node.maxBackoff
	}
}

func (node *_ManagedNode) _DecreaseDelay() {
	node.currentBackoff = node.currentBackoff / 2
	if node.currentBackoff < node.minBackoff {
		node.currentBackoff = node.minBackoff
	}
}

func (node *_ManagedNode) _Wait() time.Duration {
	return time.Now().Sub(node.backoffUntil)
}

func (node *_ManagedNode) _GetUseCount() int64 {
	return node.useCount
}

func (node *_ManagedNode) _GetLastUsed() int64 {
	return node.lastUsed
}

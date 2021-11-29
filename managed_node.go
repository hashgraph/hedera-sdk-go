package hedera

import (
	"math"
	"time"
)

type _IManagedNode interface {
	_SetMinBackoff(waitTime int64)
	_GetMinBackoff() int64
	_SetMaxBackoff(waitTime int64)
	_GetMaxBackoff() int64
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
	address        *_ManagedNodeAddress
	currentBackoff int64
	lastUsed       int64
	backoffUntil   int64
	useCount       int64
	minBackoff     int64
	maxBackoff     int64
	attempts       int64
}

func (node *_ManagedNode) _GetAttempts() int64 {
	return node.attempts
}

func (node *_ManagedNode) _GetAddress() string {
	if node.address != nil {
		return node.address._String()
	}

	return ""
}

func _NewManagedNode(address string, minBackoff int64) _ManagedNode {
	return _ManagedNode{
		address:        _ManagedNodeAddressFromString(address),
		currentBackoff: 250,
		useCount:       0,
		minBackoff:     minBackoff,
		maxBackoff:     8000,
		attempts:       0,
	}
}

func (node *_ManagedNode) _SetMinBackoff(waitTime int64) {
	if node.currentBackoff == node.minBackoff {
		node.currentBackoff = node.minBackoff
	}

	node.minBackoff = waitTime
}

func (node *_ManagedNode) _GetMinBackoff() int64 {
	return node.minBackoff
}

func (node *_ManagedNode) _SetMaxBackoff(waitTime int64) {
	node.maxBackoff = waitTime
}

func (node *_ManagedNode) _GetMaxBackoff() int64 {
	return node.maxBackoff
}

func (node *_ManagedNode) _InUse() {
	node.useCount++
	node.lastUsed = time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _IsHealthy() bool {
	return node.backoffUntil < time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _IncreaseDelay() {
	node.attempts++
	node.backoffUntil = (node.currentBackoff * 1000000) + time.Now().UTC().UnixNano()
	node.currentBackoff = int64(math.Min(float64(node.currentBackoff)*2, float64(node.maxBackoff)))
}

func (node *_ManagedNode) _DecreaseDelay() {
	node.currentBackoff = int64(math.Max(float64(node.currentBackoff)/2, float64(node.minBackoff)))
}

func (node *_ManagedNode) _Wait() time.Duration {
	delay := node.backoffUntil - node.lastUsed
	return time.Duration(delay)
}

func (node *_ManagedNode) _GetUseCount() int64 {
	return node.useCount
}

func (node *_ManagedNode) _GetLastUsed() int64 {
	return node.lastUsed
}

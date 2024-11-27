package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"sync"
	"time"
)

type _IManagedNode interface {
	_GetKey() string
	_SetVerifyCertificate(verify bool)
	_GetVerifyCertificate() bool
	_SetMinBackoff(waitTime time.Duration)
	_GetMinBackoff() time.Duration
	_SetMaxBackoff(waitTime time.Duration)
	_GetMaxBackoff() time.Duration
	_InUse()
	_IsHealthy() bool
	_IncreaseBackoff()
	_DecreaseBackoff()
	_Wait() time.Duration
	_GetUseCount() int64
	_GetLastUsed() time.Time
	_GetAttempts() int64
	_GetReadmitTime() *time.Time
	_GetAddress() string
	_ToSecure() _IManagedNode
	_ToInsecure() _IManagedNode
	_GetManagedNode() *_ManagedNode
	_Close() error
}

type _ManagedNode struct {
	address            *_ManagedNodeAddress
	currentBackoff     time.Duration
	lastUsed           time.Time
	useCount           int64
	minBackoff         time.Duration
	maxBackoff         time.Duration
	badGrpcStatusCount int64
	readmitTime        *time.Time
	mutex              sync.RWMutex
}

func (node *_ManagedNode) _GetAttempts() int64 {
	node.mutex.RLock()
	defer node.mutex.RUnlock()
	return node.badGrpcStatusCount
}

func (node *_ManagedNode) _GetAddress() string {
	if node.address != nil {
		return node.address._String()
	}

	return ""
}

func (node *_ManagedNode) _GetReadmitTime() *time.Time {
	node.mutex.RLock()
	defer node.mutex.RUnlock()
	return node.readmitTime
}

func _NewManagedNode(address string, minBackoff time.Duration) (node *_ManagedNode, err error) {
	node = &_ManagedNode{
		currentBackoff:     minBackoff,
		lastUsed:           time.Now(),
		useCount:           0,
		minBackoff:         minBackoff,
		maxBackoff:         1 * time.Hour,
		badGrpcStatusCount: 0,
	}
	node.address, err = _ManagedNodeAddressFromString(address)
	return node, err
}

func (node *_ManagedNode) _SetMinBackoff(minBackoff time.Duration) {
	if node.currentBackoff == node.minBackoff {
		node.currentBackoff = node.minBackoff
	}

	node.minBackoff = minBackoff
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
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.useCount++
	node.lastUsed = time.Now()
}

func (node *_ManagedNode) _IsHealthy() bool {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	if node.readmitTime == nil {
		return true
	}

	return node.readmitTime.Before(time.Now())
}

func (node *_ManagedNode) _IncreaseBackoff() {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.badGrpcStatusCount++
	node.currentBackoff *= 2
	if node.currentBackoff > node.maxBackoff {
		node.currentBackoff = node.maxBackoff
	}
	readmitTime := time.Now().Add(node.currentBackoff)
	node.readmitTime = &readmitTime
}

func (node *_ManagedNode) _DecreaseBackoff() {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.currentBackoff /= 2
	if node.currentBackoff < node.minBackoff {
		node.currentBackoff = node.minBackoff
	}
}

func (node *_ManagedNode) _Wait() time.Duration {
	node.mutex.RLock()
	defer node.mutex.RUnlock()
	return node.readmitTime.Sub(node.lastUsed)
}

func (node *_ManagedNode) _GetUseCount() int64 {
	return node.useCount
}

func (node *_ManagedNode) _GetLastUsed() time.Time {
	return node.lastUsed
}

package hedera

import (
	"time"
)

type _IManagedNode interface {
	_GetKey() string
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
	_GetReadmitTime() time.Time
	_GetAddress() string
	_ToSecure() _IManagedNode
	_ToInsecure() _IManagedNode
	_GetManagedNode() *_ManagedNode
	_Close() error
}

type _ManagedNode struct {
	address        *_ManagedNodeAddress
	currentBackoff time.Duration
	lastUsed       time.Time
	useCount       int64
	minBackoff     time.Duration
	maxBackoff     time.Duration
	attempts       int64
	readmitTime    time.Time
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

func (node *_ManagedNode) _GetReadmitTime() time.Time {
	return node.readmitTime
}

func _NewManagedNode(address string, minBackoff time.Duration) (node *_ManagedNode, err error) {
	node = &_ManagedNode{
		currentBackoff: minBackoff,
		lastUsed:       time.Now(),
		useCount:       0,
		minBackoff:     minBackoff,
		maxBackoff:     1 * time.Hour,
		attempts:       0,
		readmitTime:    time.Now(),
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
	node.useCount++
	node.lastUsed = time.Now()
	node.readmitTime = time.Now()
}

func (node *_ManagedNode) _IsHealthy() bool {
	return node.readmitTime.Before(time.Now())
}

func (node *_ManagedNode) _IncreaseBackoff() {
	node.attempts++
	node.currentBackoff *= 2
	if node.currentBackoff > node.maxBackoff {
		node.currentBackoff = node.maxBackoff
	}
	node.readmitTime = time.Now().Add(node.currentBackoff)
}

func (node *_ManagedNode) _DecreaseBackoff() {
	node.currentBackoff /= 2
	if node.currentBackoff < node.minBackoff {
		node.currentBackoff = node.minBackoff
	}
}

func (node *_ManagedNode) _Wait() time.Duration {
	return node.readmitTime.Sub(node.lastUsed)
}

func (node *_ManagedNode) _GetUseCount() int64 {
	return node.useCount
}

func (node *_ManagedNode) _GetLastUsed() time.Time {
	return node.lastUsed
}

func _ManagedNodeCompare(i *_ManagedNode, j *_ManagedNode) int64 {
	iRemainingTime := i._Wait()
	jRemainingTime := j._Wait()
	comparison := iRemainingTime.Milliseconds() - jRemainingTime.Milliseconds()
	if iRemainingTime > 0 && jRemainingTime > 0 && comparison != 0 {
		return comparison
	}

	comparison = i.currentBackoff.Milliseconds() - j.currentBackoff.Milliseconds()
	if comparison != 0 {
		return comparison
	}

	comparison = i.attempts - j.attempts
	if comparison != 0 {
		return comparison
	}

	comparison = i.useCount - j.useCount
	if comparison != 0 {
		return comparison
	}

	return i.lastUsed.UnixNano() - j.lastUsed.UnixNano()
}

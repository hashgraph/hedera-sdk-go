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
	_GetReadmitTime() int64
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
	useCount       int64
	minBackoff     int64
	maxBackoff     int64
	attempts       int64
	readmitTime    int64
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

func (node *_ManagedNode) _GetReadmitTime() int64 {
	return node.readmitTime
}

func _NewManagedNode(address string, minBackoff int64) _ManagedNode {
	return _ManagedNode{
		address:        _ManagedNodeAddressFromString(address),
		currentBackoff: minBackoff,
		lastUsed:       time.Now().UnixNano(),
		useCount:       0,
		minBackoff:     minBackoff,
		maxBackoff:     1000 * 60 * 60 * time.Millisecond.Nanoseconds(),
		attempts:       0,
		readmitTime:    time.Now().UnixNano(),
	}
	node.address, err = _ManagedNodeAddressFromString(address)
	return node, err
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
	node.readmitTime = time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _IsHealthy() bool {
	return node.readmitTime < time.Now().UTC().UnixNano()
}

func (node *_ManagedNode) _IncreaseDelay() {
	node.attempts++
	node.currentBackoff = int64(math.Min(float64(node.currentBackoff)*2, float64(node.maxBackoff)))
	node.readmitTime = time.Now().UTC().UnixNano() + node.currentBackoff
}

func (node *_ManagedNode) _DecreaseDelay() {
	node.currentBackoff /= 2
	if node.currentBackoff < node.minBackoff {
		node.currentBackoff = node.minBackoff
	}
}

func (node *_ManagedNode) _Wait() time.Duration {
	delay := node.readmitTime - node.lastUsed
	return time.Duration(delay)
}

func (node *_ManagedNode) _GetUseCount() int64 {
	return node.useCount
}

func (node *_ManagedNode) _GetLastUsed() int64 {
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

	comparison = i.badGrpcStatusCount - j.badGrpcStatusCount
	if comparison != 0 {
		return comparison
	}

	comparison = i.useCount - j.useCount
	if comparison != 0 {
		return comparison
	}

	return i.lastUsed - j.lastUsed
}

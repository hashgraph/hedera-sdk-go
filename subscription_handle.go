package hiero

// SPDX-License-Identifier: Apache-2.0

type SubscriptionHandle struct {
	onUnsubscribe func()
}

// Unsubscribe removes the subscription from the client
func (handle SubscriptionHandle) Unsubscribe() {
	if handle.onUnsubscribe != nil {
		handle.onUnsubscribe()
	}
}

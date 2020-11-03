package hedera

type SubscriptionHandle struct {
	onUnsubscribe func()
}

func newSubscriptionHandle(onUnsubscribe func()) SubscriptionHandle {
	return SubscriptionHandle{onUnsubscribe: onUnsubscribe}
}

func (handle SubscriptionHandle) Unsubscribe() {
	handle.onUnsubscribe()
}

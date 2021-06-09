package hedera

type SubscriptionHandle struct {
	onUnsubscribe func()
}

func (handle SubscriptionHandle) Unsubscribe() {
	if handle.onUnsubscribe != nil {
		handle.onUnsubscribe()
	}
}

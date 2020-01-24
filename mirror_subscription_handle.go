package hedera

type MirrorSubscriptionHandle struct {
	onUnsubscribe func()
}

func newMirrorSubscriptionHandle(onUnsubscribe func()) MirrorSubscriptionHandle {
	return MirrorSubscriptionHandle{onUnsubscribe: onUnsubscribe}
}

func (handle MirrorSubscriptionHandle) Unsubscribe() {
	handle.onUnsubscribe()
}

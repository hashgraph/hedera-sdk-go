package hedera

type MirrorSubscriptionHandle struct {
	onUnsubscribe func() error
}

func NewMirrorSubscriptionHandle(onUnsubscribe func() error) MirrorSubscriptionHandle {
	return MirrorSubscriptionHandle{onUnsubscribe:onUnsubscribe}
}

func (handle MirrorSubscriptionHandle) Unsubscribe() error {
	return handle.onUnsubscribe()
}

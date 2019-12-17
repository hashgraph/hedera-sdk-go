package hedera

type EntityID struct {
	ty string
	id EntityIdLike
}

type EntityIdLike interface {
	isEntityIdLike()
}

func (id FileID) isEntityIdLike()     {}
func (id AccountID) isEntityIdLike()  {}
func (id ContractID) isEntityIdLike() {}

package hedera

type EntityID interface {
	isEntityID()
}

func (id FileID) isEntityID()     {}
func (id AccountID) isEntityID()  {}
func (id ContractID) isEntityID() {}

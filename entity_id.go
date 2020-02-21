package hedera

// EntityID is an interface for various IDs of entities (Account, Contract, File, etc)
type EntityID interface {
	isEntityID()
}

func (id FileID) isEntityID()     {}
func (id AccountID) isEntityID()  {}
func (id ContractID) isEntityID() {}

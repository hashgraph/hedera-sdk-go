package hedera

type BaseTokenAirdropTransaction struct {
	Transaction
	pendingAirdropIds []*PendingAirdropId
}

// SetPendingAirdropIds sets the pending airdrop IDs for this TokenCancelAirdropTransaction.
func (tx *BaseTokenAirdropTransaction) SetPendingAirdropIds(ids []*PendingAirdropId) *BaseTokenAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = ids
	return tx
}

// AddPendingAirdropId adds a pending airdrop ID to this TokenCancelAirdropTransaction.
func (tx *BaseTokenAirdropTransaction) AddPendingAirdropId(id PendingAirdropId) *BaseTokenAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = append(tx.pendingAirdropIds, &id)
	return tx
}

// GetPendingAirdropIds returns the pending airdrop IDs for this TokenCancelAirdropTransaction.
func (tx *BaseTokenAirdropTransaction) GetPendingAirdropIds() []*PendingAirdropId {
	return tx.pendingAirdropIds
}

func (tx *BaseTokenAirdropTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, pendingAirdropId := range tx.pendingAirdropIds {
		if pendingAirdropId.sender != nil {
			if err := pendingAirdropId.sender.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.receiver != nil {
			if err := pendingAirdropId.receiver.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.nftID != nil {
			if err := pendingAirdropId.nftID.Validate(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.tokenID != nil {
			if err := pendingAirdropId.tokenID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}
	return nil
}

// func (tx *BaseTokenAirdropTransaction) buildProtoBody() []*services.PendingAirdropId {
// 	pendingAirdrops := make([]*services.PendingAirdropId, len(tx.pendingAirdropIds))
// 	for i, pendingAirdropId := range tx.pendingAirdropIds {
// 		pendingAirdrops[i] = pendingAirdropId._ToProtobuf()
// 	}
// 	return pendingAirdrops
// }

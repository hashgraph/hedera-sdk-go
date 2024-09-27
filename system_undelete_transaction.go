package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

// Undelete a file or smart contract that was deleted by AdminDelete.
// Can only be done with a Hedera admin.
type SystemUndeleteTransaction struct {
	*Transaction[*SystemUndeleteTransaction]
	contractID *ContractID
	fileID     *FileID
}

// NewSystemUndeleteTransaction creates a SystemUndeleteTransaction transaction which can be
// used to construct and execute a System Undelete Transaction.
func NewSystemUndeleteTransaction() *SystemUndeleteTransaction {
	tx := &SystemUndeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _SystemUndeleteTransactionFromProtobuf(pb *services.TransactionBody) *SystemUndeleteTransaction {
	return &SystemUndeleteTransaction{
		contractID: _ContractIDFromProtobuf(pb.GetSystemUndelete().GetContractID()),
		fileID:     _FileIDFromProtobuf(pb.GetSystemUndelete().GetFileID()),
	}
}

// SetContractID sets the ContractID of the contract whose deletion is being undone.
func (tx *SystemUndeleteTransaction) SetContractID(contractID ContractID) *SystemUndeleteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// GetContractID returns the ContractID of the contract whose deletion is being undone.
func (tx *SystemUndeleteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// SetFileID sets the FileID of the file whose deletion is being undone.
func (tx *SystemUndeleteTransaction) SetFileID(fileID FileID) *SystemUndeleteTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID of the file whose deletion is being undone.
func (tx *SystemUndeleteTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// ----------- Overridden functions ----------------

func (tx *SystemUndeleteTransaction) getName() string {
	return "SystemUndeleteTransaction"
}

func (tx *SystemUndeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.contractID != nil {
		if err := tx.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.fileID != nil {
		if err := tx.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *SystemUndeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_SystemUndelete{
			SystemUndelete: tx.buildProtoBody(),
		},
	}
}

func (tx *SystemUndeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_SystemUndelete{
			SystemUndelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *SystemUndeleteTransaction) buildProtoBody() *services.SystemUndeleteTransactionBody {
	body := &services.SystemUndeleteTransactionBody{}
	if tx.contractID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_ContractID{
			ContractID: tx.contractID._ToProtobuf(),
		}
	}

	if tx.fileID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_FileID{
			FileID: tx.fileID._ToProtobuf(),
		}
	}

	return body
}

func (tx *SystemUndeleteTransaction) getMethod(channel *_Channel) _Method {
	if channel._GetContract() == nil {
		return _Method{
			transaction: channel._GetFile().SystemUndelete,
		}
	}

	return _Method{
		transaction: channel._GetContract().SystemUndelete,
	}
}

func (tx *SystemUndeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *SystemUndeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction[*SystemUndeleteTransaction](tx.Transaction)
}

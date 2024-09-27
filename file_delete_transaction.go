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

// FileDeleteTransaction Deletes the given file. After deletion, it will be marked as deleted and will have no contents.
// But information about it will continue to exist until it expires. A list of keys was given when
// the file was created. All the top level keys on that list must sign transactions to create or
// modify the file, but any single one of the top level keys can be used to delete the file. This
// transaction must be signed by 1-of-M KeyList keys. If keys contains additional KeyList or
// ThresholdKey then 1-of-M secondary KeyList or ThresholdKey signing requirements must be meet.
type FileDeleteTransaction struct {
	*Transaction[*FileDeleteTransaction]
	fileID *FileID
}

// NewFileDeleteTransaction creates a FileDeleteTransaction which deletes the given file. After deletion,
// it will be marked as deleted and will have no contents.
// But information about it will continue to exist until it expires. A list of keys was given when
// the file was created. All the top level keys on that list must sign transactions to create or
// modify the file, but any single one of the top level keys can be used to delete the file. This
// transaction must be signed by 1-of-M KeyList keys. If keys contains additional KeyList or
// ThresholdKey then 1-of-M secondary KeyList or ThresholdKey signing requirements must be meet.
func NewFileDeleteTransaction() *FileDeleteTransaction {
	tx := &FileDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _FileDeleteTransactionFromProtobuf(tx Transaction[*FileDeleteTransaction], pb *services.TransactionBody) *FileDeleteTransaction {
	return &FileDeleteTransaction{
		Transaction: &tx,
		fileID:      _FileIDFromProtobuf(pb.GetFileDelete().GetFileID()),
	}
}

// SetFileID Sets the FileID of the file to be deleted
func (tx *FileDeleteTransaction) SetFileID(fileID FileID) *FileDeleteTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID of the file to be deleted
func (tx *FileDeleteTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// ----------- Overridden functions ----------------

func (tx *FileDeleteTransaction) getName() string {
	return "FileDeleteTransaction"
}
func (tx *FileDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.fileID != nil {
		if err := tx.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *FileDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileDelete{
			FileDelete: tx.buildProtoBody(),
		},
	}
}

func (tx *FileDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileDelete{
			FileDelete: tx.buildProtoBody(),
		},
	}, nil
}
func (tx *FileDeleteTransaction) buildProtoBody() *services.FileDeleteTransactionBody {
	body := &services.FileDeleteTransactionBody{}
	if tx.fileID != nil {
		body.FileID = tx.fileID._ToProtobuf()
	}
	return body
}

func (tx *FileDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().DeleteFile,
	}
}
func (tx *FileDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *FileDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction[*FileDeleteTransaction](tx.Transaction)
}

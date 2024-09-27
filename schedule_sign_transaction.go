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
	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// ScheduleSignTransaction Adds zero or more signing keys to a schedule.
// If Long Term Scheduled Transactions are enabled and wait for expiry was set to true on the
// ScheduleCreate then the transaction will always wait till it's `expiration_time` to execute.
// Otherwise, if the resulting set of signing keys satisfy the
// scheduled transaction's signing requirements, it will be executed immediately after the
// triggering ScheduleSign.
// Upon SUCCESS, the receipt includes the scheduledTransactionID to use to query
// for the record of the scheduled transaction's execution (if it occurs).
type ScheduleSignTransaction struct {
	*Transaction[*ScheduleSignTransaction]
	scheduleID *ScheduleID
}

// NewScheduleSignTransaction creates ScheduleSignTransaction which adds zero or more signing keys to a schedule.
// If Long Term Scheduled Transactions are enabled and wait for expiry was set to true on the
// ScheduleCreate then the transaction will always wait till it's `expiration_time` to execute.
// Otherwise, if the resulting set of signing keys satisfy the
// scheduled transaction's signing requirements, it will be executed immediately after the
// triggering ScheduleSign.
// Upon SUCCESS, the receipt includes the scheduledTransactionID to use to query
// for the record of the scheduled transaction's execution (if it occurs).
func NewScheduleSignTransaction() *ScheduleSignTransaction {
	tx := &ScheduleSignTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _ScheduleSignTransactionFromProtobuf(pb *services.TransactionBody) *ScheduleSignTransaction {
	return &ScheduleSignTransaction{
		scheduleID: _ScheduleIDFromProtobuf(pb.GetScheduleSign().GetScheduleID()),
	}
}

// SetScheduleID Sets the id of the schedule to add signing keys to
func (tx *ScheduleSignTransaction) SetScheduleID(scheduleID ScheduleID) *ScheduleSignTransaction {
	tx._RequireNotFrozen()
	tx.scheduleID = &scheduleID
	return tx
}

// GetScheduleID returns the id of the schedule to add signing keys to
func (tx *ScheduleSignTransaction) GetScheduleID() ScheduleID {
	if tx.scheduleID == nil {
		return ScheduleID{}
	}

	return *tx.scheduleID
}

// ----------- Overridden functions ----------------

func (tx *ScheduleSignTransaction) getName() string {
	return "ScheduleSignTransaction"
}

func (tx *ScheduleSignTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.scheduleID != nil {
		if err := tx.scheduleID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *ScheduleSignTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ScheduleSign{
			ScheduleSign: tx.buildProtoBody(),
		},
	}
}

func (tx *ScheduleSignTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleSignTransaction")
}

func (tx *ScheduleSignTransaction) buildProtoBody() *services.ScheduleSignTransactionBody {
	body := &services.ScheduleSignTransactionBody{}
	if tx.scheduleID != nil {
		body.ScheduleID = tx.scheduleID._ToProtobuf()
	}

	return body
}

func (tx *ScheduleSignTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().SignSchedule,
	}
}

func (tx *ScheduleSignTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *ScheduleSignTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction[*ScheduleSignTransaction](tx.Transaction)
}

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
)

type ScheduleInfo struct {
	ScheduleID       ScheduleID
	CreatorAccountID AccountID
	PayerAccountID   AccountID
	ExecutedAt       *time.Time
	DeletedAt        *time.Time
	ExpirationTime   time.Time
	Signatories      *KeyList
	// Deprecated: Use ScheduleInfo.Signatories instead
	Signers                  *KeyList
	AdminKey                 Key
	Memo                     string
	ScheduledTransactionID   *TransactionID
	scheduledTransactionBody *services.SchedulableTransactionBody
	LedgerID                 LedgerID
	WaitForExpiry            bool
}

func _ScheduleInfoFromProtobuf(pb *services.ScheduleInfo) ScheduleInfo {
	if pb == nil {
		return ScheduleInfo{}
	}
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = _KeyFromProtobuf(pb.AdminKey)
	}

	var signatories KeyList
	if pb.Signers != nil {
		signatories, _ = _KeyListFromProtobuf(pb.Signers)
	}

	var scheduledTransactionID TransactionID
	if pb.ScheduledTransactionID != nil {
		scheduledTransactionID = _TransactionIDFromProtobuf(pb.ScheduledTransactionID)
	}

	var executed *time.Time
	var deleted *time.Time
	switch t := pb.Data.(type) {
	case *services.ScheduleInfo_ExecutionTime:
		temp := _TimeFromProtobuf(t.ExecutionTime)
		executed = &temp
	case *services.ScheduleInfo_DeletionTime:
		temp := _TimeFromProtobuf(t.DeletionTime)
		deleted = &temp
	}

	creatorAccountID := AccountID{}
	if pb.CreatorAccountID != nil {
		creatorAccountID = *_AccountIDFromProtobuf(pb.CreatorAccountID)
	}

	payerAccountID := AccountID{}
	if pb.PayerAccountID != nil {
		payerAccountID = *_AccountIDFromProtobuf(pb.PayerAccountID)
	}

	scheduleID := ScheduleID{}
	if pb.ScheduleID != nil {
		scheduleID = *_ScheduleIDFromProtobuf(pb.ScheduleID)
	}

	return ScheduleInfo{
		ScheduleID:               scheduleID,
		CreatorAccountID:         creatorAccountID,
		PayerAccountID:           payerAccountID,
		ExecutedAt:               executed,
		DeletedAt:                deleted,
		ExpirationTime:           _TimeFromProtobuf(pb.ExpirationTime),
		Signatories:              &signatories,
		Signers:                  &signatories,
		AdminKey:                 adminKey,
		Memo:                     pb.Memo,
		ScheduledTransactionID:   &scheduledTransactionID,
		scheduledTransactionBody: pb.ScheduledTransactionBody,
		LedgerID:                 LedgerID{pb.LedgerId},
		WaitForExpiry:            pb.WaitForExpiry,
	}
}

func (scheduleInfo *ScheduleInfo) _ToProtobuf() *services.ScheduleInfo { // nolint
	var adminKey *services.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey._ToProtoKey()
	}

	var signatories *services.KeyList
	if scheduleInfo.Signatories != nil {
		signatories = scheduleInfo.Signatories._ToProtoKeyList()
	} else if scheduleInfo.Signers != nil {
		signatories = scheduleInfo.Signers._ToProtoKeyList()
	}

	info := &services.ScheduleInfo{
		ScheduleID:               scheduleInfo.ScheduleID._ToProtobuf(),
		ExpirationTime:           _TimeToProtobuf(scheduleInfo.ExpirationTime),
		ScheduledTransactionBody: scheduleInfo.scheduledTransactionBody,
		Memo:                     scheduleInfo.Memo,
		AdminKey:                 adminKey,
		Signers:                  signatories,
		CreatorAccountID:         scheduleInfo.CreatorAccountID._ToProtobuf(),
		PayerAccountID:           scheduleInfo.PayerAccountID._ToProtobuf(),
		ScheduledTransactionID:   scheduleInfo.ScheduledTransactionID._ToProtobuf(),
		LedgerId:                 scheduleInfo.LedgerID.ToBytes(),
		WaitForExpiry:            scheduleInfo.WaitForExpiry,
	}

	if scheduleInfo.ExecutedAt != nil {
		info.Data = &services.ScheduleInfo_DeletionTime{
			DeletionTime: _TimeToProtobuf(*scheduleInfo.DeletedAt),
		}
	} else if scheduleInfo.DeletedAt != nil {
		info.Data = &services.ScheduleInfo_ExecutionTime{
			ExecutionTime: _TimeToProtobuf(*scheduleInfo.ExecutedAt),
		}
	}

	return info
}

// GetScheduledTransaction returns the scheduled transaction associated with this schedule
func (scheduleInfo *ScheduleInfo) GetScheduledTransaction() (ITransaction, error) { // nolint
	pb := scheduleInfo.scheduledTransactionBody

	pbBody := &services.TransactionBody{
		TransactionFee: pb.TransactionFee,
		Memo:           pb.Memo,
	}

	tx := Transaction{
		transactionFee: pb.GetTransactionFee(),
		memo:           pb.GetMemo(),
	}

	switch pb.Data.(type) {
	case *services.SchedulableTransactionBody_ContractCall:
		pbBody.Data = &services.TransactionBody_ContractCall{
			ContractCall: pb.GetContractCall(),
		}

		tx2 := _ContractExecuteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_ContractCreateInstance:
		pbBody.Data = &services.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: pb.GetContractCreateInstance(),
		}

		tx2 := _ContractCreateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_ContractUpdateInstance:
		pbBody.Data = &services.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: pb.GetContractUpdateInstance(),
		}

		tx2 := _ContractUpdateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_ContractDeleteInstance:
		pbBody.Data = &services.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: pb.GetContractDeleteInstance(),
		}

		tx2 := _ContractDeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_CryptoCreateAccount:
		pbBody.Data = &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: pb.GetCryptoCreateAccount(),
		}

		tx2 := _AccountCreateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_CryptoDelete:
		pbBody.Data = &services.TransactionBody_CryptoDelete{
			CryptoDelete: pb.GetCryptoDelete(),
		}

		tx2 := _AccountDeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_CryptoTransfer:
		pbBody.Data = &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: pb.GetCryptoTransfer(),
		}

		tx2 := _TransferTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_CryptoUpdateAccount:
		pbBody.Data = &services.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: pb.GetCryptoUpdateAccount(),
		}

		tx2 := _AccountUpdateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_CryptoApproveAllowance:
		pbBody.Data = &services.TransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: pb.GetCryptoApproveAllowance(),
		}

		tx2 := _AccountAllowanceApproveTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_CryptoDeleteAllowance:
		pbBody.Data = &services.TransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: pb.GetCryptoDeleteAllowance(),
		}

		tx2 := _AccountAllowanceDeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_FileAppend:
		pbBody.Data = &services.TransactionBody_FileAppend{
			FileAppend: pb.GetFileAppend(),
		}

		tx2 := _FileAppendTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_FileCreate:
		pbBody.Data = &services.TransactionBody_FileCreate{
			FileCreate: pb.GetFileCreate(),
		}

		tx2 := _FileCreateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_FileDelete:
		pbBody.Data = &services.TransactionBody_FileDelete{
			FileDelete: pb.GetFileDelete(),
		}

		tx2 := _FileDeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_FileUpdate:
		pbBody.Data = &services.TransactionBody_FileUpdate{
			FileUpdate: pb.GetFileUpdate(),
		}

		tx2 := _FileUpdateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_SystemDelete:
		pbBody.Data = &services.TransactionBody_SystemDelete{
			SystemDelete: pb.GetSystemDelete(),
		}

		tx2 := _SystemDeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_SystemUndelete:
		pbBody.Data = &services.TransactionBody_SystemUndelete{
			SystemUndelete: pb.GetSystemUndelete(),
		}

		tx2 := _SystemUndeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_Freeze:
		pbBody.Data = &services.TransactionBody_Freeze{
			Freeze: pb.GetFreeze(),
		}

		tx2 := _FreezeTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_ConsensusCreateTopic:
		pbBody.Data = &services.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: pb.GetConsensusCreateTopic(),
		}

		tx2 := _TopicCreateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_ConsensusUpdateTopic:
		pbBody.Data = &services.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: pb.GetConsensusUpdateTopic(),
		}

		tx2 := _TopicUpdateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_ConsensusDeleteTopic:
		pbBody.Data = &services.TransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: pb.GetConsensusDeleteTopic(),
		}

		tx2 := _TopicDeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_ConsensusSubmitMessage:
		pbBody.Data = &services.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: pb.GetConsensusSubmitMessage(),
		}

		tx2 := _TopicMessageSubmitTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenCreation:
		pbBody.Data = &services.TransactionBody_TokenCreation{
			TokenCreation: pb.GetTokenCreation(),
		}

		tx2 := _TokenCreateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenFreeze:
		pbBody.Data = &services.TransactionBody_TokenFreeze{
			TokenFreeze: pb.GetTokenFreeze(),
		}

		tx2 := _TokenFreezeTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenUnfreeze:
		pbBody.Data = &services.TransactionBody_TokenUnfreeze{
			TokenUnfreeze: pb.GetTokenUnfreeze(),
		}

		tx2 := _TokenUnfreezeTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenFeeScheduleUpdate:
		pbBody.Data = &services.TransactionBody_TokenFeeScheduleUpdate{
			TokenFeeScheduleUpdate: pb.GetTokenFeeScheduleUpdate(),
		}

		tx2 := _TokenFeeScheduleUpdateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenGrantKyc:
		pbBody.Data = &services.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: pb.GetTokenGrantKyc(),
		}

		tx2 := _TokenGrantKycTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenRevokeKyc:
		pbBody.Data = &services.TransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: pb.GetTokenRevokeKyc(),
		}

		tx2 := _TokenRevokeKycTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenDeletion:
		pbBody.Data = &services.TransactionBody_TokenDeletion{
			TokenDeletion: pb.GetTokenDeletion(),
		}

		tx2 := _TokenDeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenUpdate:
		pbBody.Data = &services.TransactionBody_TokenUpdate{
			TokenUpdate: pb.GetTokenUpdate(),
		}

		tx2 := _TokenUpdateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenMint:
		pbBody.Data = &services.TransactionBody_TokenMint{
			TokenMint: pb.GetTokenMint(),
		}

		tx2 := _TokenMintTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenBurn:
		pbBody.Data = &services.TransactionBody_TokenBurn{
			TokenBurn: pb.GetTokenBurn(),
		}

		tx2 := _TokenBurnTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenWipe:
		pbBody.Data = &services.TransactionBody_TokenWipe{
			TokenWipe: pb.GetTokenWipe(),
		}

		tx2 := _TokenWipeTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenAssociate:
		pbBody.Data = &services.TransactionBody_TokenAssociate{
			TokenAssociate: pb.GetTokenAssociate(),
		}

		tx2 := _TokenAssociateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_TokenDissociate:
		pbBody.Data = &services.TransactionBody_TokenDissociate{
			TokenDissociate: pb.GetTokenDissociate(),
		}

		tx2 := _TokenDissociateTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_ScheduleDelete:
		pbBody.Data = &services.TransactionBody_ScheduleDelete{
			ScheduleDelete: pb.GetScheduleDelete(),
		}

		tx2 := _ScheduleDeleteTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	case *services.SchedulableTransactionBody_UtilPrng:
		pbBody.Data = &services.TransactionBody_UtilPrng{
			UtilPrng: pb.GetUtilPrng(),
		}

		tx2 := _PrngTransactionFromProtobuf(tx, pbBody)
		return tx2, nil
	default:
		return nil, errors.New("(BUG) non-exhaustive switch statement")
	}
}

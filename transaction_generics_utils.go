package hedera

import "fmt"

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

// TODO find a way to test this for future tranasctions
func getInterfaceGenericTransaction(tx any) (*Transaction[TransactionInterface], error) { // nolint
	if s, ok := tx.(*Transaction[*ContractExecuteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenFeeScheduleUpdateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*AccountAllowanceDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*AccountAllowanceApproveTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*ContractCreateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*ContractUpdateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*ContractDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*LiveHashAddTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*AccountCreateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*AccountDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*LiveHashDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TransferTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*AccountUpdateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*FileAppendTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*FileCreateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*FileDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*FileUpdateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*SystemDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*SystemUndeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*FreezeTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TopicCreateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TopicUpdateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TopicDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TopicMessageSubmitTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenCreateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenFreezeTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenUnfreezeTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenGrantKycTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenRevokeKycTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenUpdateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenMintTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenBurnTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenWipeTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenAssociateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenDissociateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*ScheduleCreateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*ScheduleDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*ScheduleSignTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenPauseTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenUnpauseTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*EthereumTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*PrngTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenUpdateNfts]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenRejectTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*NodeCreateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*NodeUpdateTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*NodeDeleteTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenAirdropTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenClaimAirdropTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	if s, ok := tx.(*Transaction[*TokenCancelAirdropTransaction]); ok {
		return castFromConcreteToBaseTransaction(s), nil
	}
	return nil, fmt.Errorf("unsupported transaction type")
}

// Helper function to cast any type of transaction to the base transaction
func getConcreteGenericTransaction(baseTx Transaction[TransactionInterface]) (any, error) { // nolint
	switch baseTx.childTransaction.(type) {
	case *ContractExecuteTransaction:
		return *castFromBaseToConcreteTransaction[*ContractExecuteTransaction](baseTx), nil
	case *TokenFeeScheduleUpdateTransaction:
		return *castFromBaseToConcreteTransaction[*TokenFeeScheduleUpdateTransaction](baseTx), nil
	case *ContractCreateTransaction:
		return *castFromBaseToConcreteTransaction[*ContractCreateTransaction](baseTx), nil
	case *ContractUpdateTransaction:
		return *castFromBaseToConcreteTransaction[*ContractUpdateTransaction](baseTx), nil
	case *ContractDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*ContractDeleteTransaction](baseTx), nil
	case *AccountAllowanceApproveTransaction:
		return *castFromBaseToConcreteTransaction[*AccountAllowanceApproveTransaction](baseTx), nil
	case *AccountAllowanceDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*AccountAllowanceDeleteTransaction](baseTx), nil
	case *LiveHashAddTransaction:
		return *castFromBaseToConcreteTransaction[*LiveHashAddTransaction](baseTx), nil
	case *AccountCreateTransaction:
		return *castFromBaseToConcreteTransaction[*AccountCreateTransaction](baseTx), nil
	case *AccountDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*AccountDeleteTransaction](baseTx), nil
	case *LiveHashDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*LiveHashDeleteTransaction](baseTx), nil
	case *TransferTransaction:
		return *castFromBaseToConcreteTransaction[*TransferTransaction](baseTx), nil
	case *AccountUpdateTransaction:
		return *castFromBaseToConcreteTransaction[*AccountUpdateTransaction](baseTx), nil
	case *FileAppendTransaction:
		return *castFromBaseToConcreteTransaction[*FileAppendTransaction](baseTx), nil
	case *FileCreateTransaction:
		return *castFromBaseToConcreteTransaction[*FileCreateTransaction](baseTx), nil
	case *FileDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*FileDeleteTransaction](baseTx), nil
	case *FileUpdateTransaction:
		return *castFromBaseToConcreteTransaction[*FileUpdateTransaction](baseTx), nil
	case *SystemDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*SystemDeleteTransaction](baseTx), nil
	case *SystemUndeleteTransaction:
		return *castFromBaseToConcreteTransaction[*SystemUndeleteTransaction](baseTx), nil
	case *FreezeTransaction:
		return *castFromBaseToConcreteTransaction[*FreezeTransaction](baseTx), nil
	case *TopicCreateTransaction:
		return *castFromBaseToConcreteTransaction[*TopicCreateTransaction](baseTx), nil
	case *TopicUpdateTransaction:
		return *castFromBaseToConcreteTransaction[*TopicUpdateTransaction](baseTx), nil
	case *TopicDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*TopicDeleteTransaction](baseTx), nil
	case *TopicMessageSubmitTransaction:
		return *castFromBaseToConcreteTransaction[*TopicMessageSubmitTransaction](baseTx), nil
	case *TokenCreateTransaction:
		return *castFromBaseToConcreteTransaction[*TokenCreateTransaction](baseTx), nil
	case *TokenFreezeTransaction:
		return *castFromBaseToConcreteTransaction[*TokenFreezeTransaction](baseTx), nil
	case *TokenUnfreezeTransaction:
		return *castFromBaseToConcreteTransaction[*TokenUnfreezeTransaction](baseTx), nil
	case *TokenGrantKycTransaction:
		return *castFromBaseToConcreteTransaction[*TokenGrantKycTransaction](baseTx), nil
	case *TokenRevokeKycTransaction:
		return *castFromBaseToConcreteTransaction[*TokenRevokeKycTransaction](baseTx), nil
	case *TokenDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*TokenDeleteTransaction](baseTx), nil
	case *TokenUpdateTransaction:
		return *castFromBaseToConcreteTransaction[*TokenUpdateTransaction](baseTx), nil
	case *TokenMintTransaction:
		return *castFromBaseToConcreteTransaction[*TokenMintTransaction](baseTx), nil
	case *TokenBurnTransaction:
		return *castFromBaseToConcreteTransaction[*TokenBurnTransaction](baseTx), nil
	case *TokenWipeTransaction:
		return *castFromBaseToConcreteTransaction[*TokenWipeTransaction](baseTx), nil
	case *TokenAssociateTransaction:
		return *castFromBaseToConcreteTransaction[*TokenAssociateTransaction](baseTx), nil
	case *TokenDissociateTransaction:
		return *castFromBaseToConcreteTransaction[*TokenDissociateTransaction](baseTx), nil
	case *ScheduleCreateTransaction:
		return *castFromBaseToConcreteTransaction[*ScheduleCreateTransaction](baseTx), nil
	case *ScheduleDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*ScheduleDeleteTransaction](baseTx), nil
	case *ScheduleSignTransaction:
		return *castFromBaseToConcreteTransaction[*ScheduleSignTransaction](baseTx), nil
	case *TokenPauseTransaction:
		return *castFromBaseToConcreteTransaction[*TokenPauseTransaction](baseTx), nil
	case *TokenUnpauseTransaction:
		return *castFromBaseToConcreteTransaction[*TokenUnpauseTransaction](baseTx), nil
	case *EthereumTransaction:
		return *castFromBaseToConcreteTransaction[*EthereumTransaction](baseTx), nil
	case *PrngTransaction:
		return *castFromBaseToConcreteTransaction[*PrngTransaction](baseTx), nil
	case *TokenUpdateNfts:
		return *castFromBaseToConcreteTransaction[*TokenUpdateNfts](baseTx), nil
	case *TokenRejectTransaction:
		return *castFromBaseToConcreteTransaction[*TokenRejectTransaction](baseTx), nil
	case *NodeCreateTransaction:
		return *castFromBaseToConcreteTransaction[*NodeCreateTransaction](baseTx), nil
	case *NodeUpdateTransaction:
		return *castFromBaseToConcreteTransaction[*NodeUpdateTransaction](baseTx), nil
	case *NodeDeleteTransaction:
		return *castFromBaseToConcreteTransaction[*NodeDeleteTransaction](baseTx), nil
	case *TokenAirdropTransaction:
		return *castFromBaseToConcreteTransaction[*TokenAirdropTransaction](baseTx), nil
	case *TokenClaimAirdropTransaction:
		return *castFromBaseToConcreteTransaction[*TokenClaimAirdropTransaction](baseTx), nil
	case *TokenCancelAirdropTransaction:
		return *castFromBaseToConcreteTransaction[*TokenCancelAirdropTransaction](baseTx), nil
	default:
		return nil, fmt.Errorf("unsupported transaction type")
	}
}

// Helper function to cast the concrete Transaction to the generic Transaction
func castFromConcreteToBaseTransaction[T TransactionInterface](baseTx *Transaction[T]) *Transaction[TransactionInterface] {
	return &Transaction[TransactionInterface]{
		executable:               baseTx.executable,
		childTransaction:         baseTx.childTransaction,
		transactionFee:           baseTx.transactionFee,
		defaultMaxTransactionFee: baseTx.defaultMaxTransactionFee,
		memo:                     baseTx.memo,
		transactionValidDuration: baseTx.transactionValidDuration,
		transactionID:            baseTx.transactionID,
		transactions:             baseTx.transactions,
		signedTransactions:       baseTx.signedTransactions,
		publicKeys:               baseTx.publicKeys,
		transactionSigners:       baseTx.transactionSigners,
		freezeError:              baseTx.freezeError,
		regenerateTransactionID:  baseTx.regenerateTransactionID,
	}
}

// Helper function to cast the generic Transaction to another type
func castFromBaseToConcreteTransaction[T TransactionInterface](baseTx Transaction[TransactionInterface]) *Transaction[T] {
	concreteTx := &Transaction[T]{
		executable:               baseTx.executable,
		transactionFee:           baseTx.transactionFee,
		defaultMaxTransactionFee: baseTx.defaultMaxTransactionFee,
		memo:                     baseTx.memo,
		transactionValidDuration: baseTx.transactionValidDuration,
		transactionID:            baseTx.transactionID,
		transactions:             baseTx.transactions,
		signedTransactions:       baseTx.signedTransactions,
		publicKeys:               baseTx.publicKeys,
		transactionSigners:       baseTx.transactionSigners,
		freezeError:              baseTx.freezeError,
		regenerateTransactionID:  baseTx.regenerateTransactionID,
	}
	if baseTx.childTransaction != nil {
		concreteTx.childTransaction = baseTx.childTransaction.(T)
	}
	return concreteTx
}

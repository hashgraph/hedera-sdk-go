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

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
	"bytes"
	"crypto/sha512"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"time"

	"github.com/hashgraph/hedera-protobufs-go/sdk"
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// transaction contains the protobuf of a prepared transaction which can be signed and executed.

type ITransaction interface {
	_ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error)
}

type TransactionInterface interface {
	Executable

	build() *services.TransactionBody
	buildScheduled() (*services.SchedulableTransactionBody, error)
	preFreezeWith(*Client)
	regenerateID(*Client) bool
	getBaseTransaction() *Transaction[TransactionInterface]
}

// Transaction is base struct for all transactions that may be built and submitted to Hedera.
type Transaction[T TransactionInterface] struct {
	executable
	childTransaction T

	transactionFee           uint64
	defaultMaxTransactionFee uint64
	memo                     string
	transactionValidDuration *time.Duration
	transactionID            TransactionID

	transactions       *_LockableSlice
	signedTransactions *_LockableSlice

	publicKeys         []PublicKey
	transactionSigners []TransactionSigner

	freezeError error

	regenerateTransactionID bool
}

func _NewTransaction[T TransactionInterface](concreteTransaction T) *Transaction[T] {
	duration := 120 * time.Second
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	return &Transaction[T]{
		childTransaction:         concreteTransaction,
		transactionValidDuration: &duration,
		transactions:             _NewLockableSlice(),
		signedTransactions:       _NewLockableSlice(),
		freezeError:              nil,
		regenerateTransactionID:  true,
		executable: executable{
			transactionIDs: _NewLockableSlice(),
			nodeAccountIDs: _NewLockableSlice(),
			minBackoff:     &minBackoff,
			maxBackoff:     &maxBackoff,
			maxRetry:       10,
		},
	}
}

func (tx *Transaction[T]) GetSignedTransactionBodyBytes(transactionIndex int) []byte {
	return tx.signedTransactions._Get(transactionIndex).(*services.SignedTransaction).GetBodyBytes()
}

// Helper function to cast any type of transaction to the base transaction
func castFromAnyToBaseTransaction(tx any) (*Transaction[TransactionInterface], error) { // nolint
	switch t := tx.(type) {
	case *Transaction[*ContractExecuteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*ContractCreateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*ContractUpdateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*ContractDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*LiveHashAddTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*AccountCreateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*AccountDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*LiveHashDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TransferTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*AccountUpdateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*FileAppendTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*FileCreateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*FileDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*FileUpdateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*SystemDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*SystemUndeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*FreezeTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TopicCreateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TopicUpdateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*AccountAllowanceApproveTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*AccountAllowanceDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TopicDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TopicMessageSubmitTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenCreateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenFreezeTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenUnfreezeTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenGrantKycTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenRevokeKycTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenUpdateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenMintTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenBurnTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenWipeTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenAssociateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenDissociateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*ScheduleCreateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenFeeScheduleUpdateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*ScheduleDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*ScheduleSignTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenPauseTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenUnpauseTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*EthereumTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*PrngTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenUpdateNfts]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenRejectTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*NodeCreateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*NodeUpdateTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*NodeDeleteTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenAirdropTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenClaimAirdropTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	case *Transaction[*TokenCancelAirdropTransaction]:
		return castFromConcreteToBaseTransaction(t), nil
	default:
		return nil, fmt.Errorf("unsupported transaction type")
	}
}

// Helper function to cast the concrete Transaction to the generic Transaction
func castFromConcreteToBaseTransaction[T TransactionInterface](tx *Transaction[T]) *Transaction[TransactionInterface] {
	return &Transaction[TransactionInterface]{
		executable:               tx.executable,
		childTransaction:         tx.childTransaction,
		transactionFee:           tx.transactionFee,
		defaultMaxTransactionFee: tx.defaultMaxTransactionFee,
		memo:                     tx.memo,
		transactionValidDuration: tx.transactionValidDuration,
		transactionID:            tx.transactionID,
		transactions:             tx.transactions,
		signedTransactions:       tx.signedTransactions,
		publicKeys:               tx.publicKeys,
		transactionSigners:       tx.transactionSigners,
		freezeError:              tx.freezeError,
		regenerateTransactionID:  tx.regenerateTransactionID,
	}
}

// Helper function to cast the generic Transaction to another type
func castFromBaseTransactionToConcreteTransaction[T TransactionInterface](tx Transaction[TransactionInterface]) Transaction[T] {
	return Transaction[T]{
		executable:               tx.executable,
		transactionFee:           tx.transactionFee,
		defaultMaxTransactionFee: tx.defaultMaxTransactionFee,
		memo:                     tx.memo,
		transactionValidDuration: tx.transactionValidDuration,
		transactionID:            tx.transactionID,
		transactions:             tx.transactions,
		signedTransactions:       tx.signedTransactions,
		publicKeys:               tx.publicKeys,
		transactionSigners:       tx.transactionSigners,
		freezeError:              tx.freezeError,
		regenerateTransactionID:  tx.regenerateTransactionID,
	}
}

// TransactionFromBytes converts transaction bytes to a related *transaction.
func TransactionFromBytes(data []byte) (interface{}, error) { // nolint
	list := sdk.TransactionList{}
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	err := protobuf.Unmarshal(data, &list)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing from bytes to transaction List")
	}

	transactions := _NewLockableSlice()

	for _, transaction := range list.TransactionList {
		transactions._Push(transaction)
	}

	tx := Transaction[TransactionInterface]{
		transactions:            transactions,
		signedTransactions:      _NewLockableSlice(),
		publicKeys:              make([]PublicKey, 0),
		transactionSigners:      make([]TransactionSigner, 0),
		freezeError:             nil,
		regenerateTransactionID: true,
		executable: executable{
			transactionIDs: _NewLockableSlice(),
			nodeAccountIDs: _NewLockableSlice(),
			minBackoff:     &minBackoff,
			maxBackoff:     &maxBackoff,
			maxRetry:       10,
		},
	}

	comp, err := _TransactionCompare(&list)
	if err != nil {
		return nil, err
	}

	if !comp {
		return nil, errors.New("failed to validate transaction bodies")
	}

	var first *services.TransactionBody = nil
	// We introduce a boolean value to distinguish flow for signed tx vs unsigned transactions
	txIsSigned := true

	for i, transactionFromList := range list.TransactionList {
		var signedTransaction services.SignedTransaction
		var body services.TransactionBody

		// If the transaction is not signed/locked:
		if len(transactionFromList.SignedTransactionBytes) == 0 {
			txIsSigned = false
			if err := protobuf.Unmarshal(transactionFromList.BodyBytes, &body); err != nil { // nolint
				return nil, errors.Wrap(err, "error deserializing BodyBytes in TransactionFromBytes")
			}
		} else { // If the transaction is signed/locked
			if err := protobuf.Unmarshal(transactionFromList.SignedTransactionBytes, &signedTransaction); err != nil {
				return nil, errors.Wrap(err, "error deserializing SignedTransactionBytes in TransactionFromBytes")
			}
		}

		if txIsSigned {
			tx.signedTransactions = tx.signedTransactions._Push(&signedTransaction)

			if i == 0 {
				for _, sigPair := range signedTransaction.GetSigMap().GetSigPair() {
					key, err := PublicKeyFromBytes(sigPair.GetPubKeyPrefix())
					if err != nil {
						return nil, err
					}

					tx.publicKeys = append(tx.publicKeys, key)
					tx.transactionSigners = append(tx.transactionSigners, nil)
				}
			}

			if err := protobuf.Unmarshal(signedTransaction.GetBodyBytes(), &body); err != nil {
				return nil, errors.Wrap(err, "error deserializing BodyBytes in TransactionFromBytes")
			}
		}

		if first == nil {
			first = &body
		}
		var transactionID TransactionID
		var nodeAccountID AccountID

		if body.GetTransactionValidDuration() != nil {
			duration := _DurationFromProtobuf(body.GetTransactionValidDuration())
			tx.transactionValidDuration = &duration
		}

		if body.GetTransactionID() != nil {
			transactionID = _TransactionIDFromProtobuf(body.GetTransactionID())
		}

		if body.GetNodeAccountID() != nil {
			nodeAccountID = *_AccountIDFromProtobuf(body.GetNodeAccountID())
		}

		// If the transaction was serialised, without setting "NodeId", or "TransactionID", we should leave them empty
		if transactionID.AccountID.Account != 0 {
			tx.transactionIDs = tx.transactionIDs._Push(transactionID)
		}
		if !nodeAccountID._IsZero() {
			tx.nodeAccountIDs = tx.nodeAccountIDs._Push(nodeAccountID)
		}

		if i == 0 {
			tx.memo = body.Memo
			if body.TransactionFee != 0 {
				tx.transactionFee = body.TransactionFee
			}
		}
	}

	if txIsSigned {
		if tx.transactionIDs._Length() > 0 {
			tx.transactionIDs.locked = true
		}

		if tx.nodeAccountIDs._Length() > 0 {
			tx.nodeAccountIDs.locked = true
		}
	}

	if first == nil {
		return nil, errNoTransactionInBytes
	}

	switch first.Data.(type) {
	case *services.TransactionBody_ContractCall:
		return _ContractExecuteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*ContractExecuteTransaction](tx), first), nil
	case *services.TransactionBody_ContractCreateInstance:
		return _ContractCreateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*ContractCreateTransaction](tx), first), nil
	case *services.TransactionBody_ContractUpdateInstance:
		return _ContractUpdateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*ContractUpdateTransaction](tx), first), nil
	case *services.TransactionBody_CryptoApproveAllowance:
		return _AccountAllowanceApproveTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*AccountAllowanceApproveTransaction](tx), first), nil
	case *services.TransactionBody_CryptoDeleteAllowance:
		return _AccountAllowanceDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*AccountAllowanceDeleteTransaction](tx), first), nil
	case *services.TransactionBody_ContractDeleteInstance:
		return _ContractDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*ContractDeleteTransaction](tx), first), nil
	case *services.TransactionBody_CryptoAddLiveHash:
		return _LiveHashAddTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*LiveHashAddTransaction](tx), first), nil
	case *services.TransactionBody_CryptoCreateAccount:
		return _AccountCreateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*AccountCreateTransaction](tx), first), nil
	case *services.TransactionBody_CryptoDelete:
		return _AccountDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*AccountDeleteTransaction](tx), first), nil
	case *services.TransactionBody_CryptoDeleteLiveHash:
		return _LiveHashDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*LiveHashDeleteTransaction](tx), first), nil
	case *services.TransactionBody_CryptoTransfer:
		return _TransferTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TransferTransaction](tx), first), nil
	case *services.TransactionBody_CryptoUpdateAccount:
		return _AccountUpdateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*AccountUpdateTransaction](tx), first), nil
	case *services.TransactionBody_FileAppend:
		return _FileAppendTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*FileAppendTransaction](tx), first), nil
	case *services.TransactionBody_FileCreate:
		return _FileCreateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*FileCreateTransaction](tx), first), nil
	case *services.TransactionBody_FileDelete:
		return _FileDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*FileDeleteTransaction](tx), first), nil
	case *services.TransactionBody_FileUpdate:
		return _FileUpdateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*FileUpdateTransaction](tx), first), nil
	case *services.TransactionBody_SystemDelete:
		return _SystemDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*SystemDeleteTransaction](tx), first), nil
	case *services.TransactionBody_SystemUndelete:
		return _SystemUndeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*SystemUndeleteTransaction](tx), first), nil
	case *services.TransactionBody_Freeze:
		return _FreezeTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*FreezeTransaction](tx), first), nil
	case *services.TransactionBody_ConsensusCreateTopic:
		return _TopicCreateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TopicCreateTransaction](tx), first), nil
	case *services.TransactionBody_ConsensusUpdateTopic:
		return _TopicUpdateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TopicUpdateTransaction](tx), first), nil
	case *services.TransactionBody_ConsensusDeleteTopic:
		return _TopicDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TopicDeleteTransaction](tx), first), nil
	case *services.TransactionBody_ConsensusSubmitMessage:
		return _TopicMessageSubmitTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TopicMessageSubmitTransaction](tx), first), nil
	case *services.TransactionBody_TokenCreation:
		return _TokenCreateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenCreateTransaction](tx), first), nil
	case *services.TransactionBody_TokenFreeze:
		return _TokenFreezeTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenFreezeTransaction](tx), first), nil
	case *services.TransactionBody_TokenUnfreeze:
		return _TokenUnfreezeTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenUnfreezeTransaction](tx), first), nil
	case *services.TransactionBody_TokenGrantKyc:
		return _TokenGrantKycTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenGrantKycTransaction](tx), first), nil
	case *services.TransactionBody_TokenRevokeKyc:
		return _TokenRevokeKycTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenRevokeKycTransaction](tx), first), nil
	case *services.TransactionBody_TokenDeletion:
		return _TokenDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenDeleteTransaction](tx), first), nil
	case *services.TransactionBody_TokenUpdate:
		return _TokenUpdateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenUpdateTransaction](tx), first), nil
	case *services.TransactionBody_TokenMint:
		return _TokenMintTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenMintTransaction](tx), first), nil
	case *services.TransactionBody_TokenBurn:
		return _TokenBurnTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenBurnTransaction](tx), first), nil
	case *services.TransactionBody_TokenWipe:
		return _TokenWipeTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenWipeTransaction](tx), first), nil
	case *services.TransactionBody_TokenAssociate:
		return _TokenAssociateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenAssociateTransaction](tx), first), nil
	case *services.TransactionBody_TokenDissociate:
		return _TokenDissociateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenDissociateTransaction](tx), first), nil
	case *services.TransactionBody_ScheduleCreate:
		return _ScheduleCreateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*ScheduleCreateTransaction](tx), first), nil
	case *services.TransactionBody_TokenFeeScheduleUpdate:
		return _TokenFeeScheduleUpdateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenFeeScheduleUpdateTransaction](tx), first), nil
	case *services.TransactionBody_ScheduleDelete:
		return _ScheduleDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*ScheduleDeleteTransaction](tx), first), nil
	case *services.TransactionBody_ScheduleSign:
		return _ScheduleSignTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*ScheduleSignTransaction](tx), first), nil
	case *services.TransactionBody_TokenPause:
		return _TokenPauseTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenPauseTransaction](tx), first), nil
	case *services.TransactionBody_TokenUnpause:
		return _TokenUnpauseTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenUnpauseTransaction](tx), first), nil
	case *services.TransactionBody_EthereumTransaction:
		return _EthereumTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*EthereumTransaction](tx), first), nil
	case *services.TransactionBody_UtilPrng:
		return _PrngTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*PrngTransaction](tx), first), nil
	case *services.TransactionBody_TokenUpdateNfts:
		return _NewTokenUpdateNftsTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenUpdateNfts](tx), first), nil
	case *services.TransactionBody_TokenReject:
		return _TokenRejectTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenRejectTransaction](tx), first), nil
	case *services.TransactionBody_NodeCreate:
		return _NodeCreateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*NodeCreateTransaction](tx), first), nil
	case *services.TransactionBody_NodeUpdate:
		return _NodeUpdateTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*NodeUpdateTransaction](tx), first), nil
	case *services.TransactionBody_NodeDelete:
		return _NodeDeleteTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*NodeDeleteTransaction](tx), first), nil
	case *services.TransactionBody_TokenAirdrop:
		return _TokenAirdropTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenAirdropTransaction](tx), first), nil
	case *services.TransactionBody_TokenClaimAirdrop:
		return _TokenClaimAirdropTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenClaimAirdropTransaction](tx), first), nil
	case *services.TransactionBody_TokenCancelAirdrop:
		return _TokenCancelAirdropTransactionFromProtobuf(castFromBaseTransactionToConcreteTransaction[*TokenCancelAirdropTransaction](tx), first), nil
	default:
		return nil, errFailedToDeserializeBytes
	}
}

func _TransactionCompare(list *sdk.TransactionList) (bool, error) {
	signed := make([]*services.SignedTransaction, 0)
	var err error
	for _, s := range list.TransactionList {
		temp := services.SignedTransaction{}
		err = protobuf.Unmarshal(s.SignedTransactionBytes, &temp)
		if err != nil {
			return false, err
		}
		signed = append(signed, &temp)
	}
	body := make([]*services.TransactionBody, 0)
	for _, s := range signed {
		temp := services.TransactionBody{}
		err = protobuf.Unmarshal(s.BodyBytes, &temp)
		if err != nil {
			return false, err
		}
		body = append(body, &temp)
	}

	for i := 1; i < len(body); i++ {
		// #nosec G602
		if reflect.TypeOf(body[0].Data) != reflect.TypeOf(body[i].Data) {
			return false, nil
		}
	}

	return true, nil
}

// GetSignatures Gets all of the signatures stored in the transaction
func (tx *Transaction[T]) GetSignatures() (map[AccountID]map[*PublicKey][]byte, error) {
	returnMap := make(map[AccountID]map[*PublicKey][]byte, tx.nodeAccountIDs._Length())

	if tx.signedTransactions._Length() == 0 {
		return returnMap, nil
	}

	for i, nodeID := range tx.nodeAccountIDs.slice {
		var sigMap *services.SignatureMap
		var tempID AccountID
		switch k := tx.signedTransactions._Get(i).(type) { //nolint
		case *services.SignedTransaction:
			sigMap = k.SigMap
		}

		switch k := nodeID.(type) { //nolint
		case AccountID:
			tempID = k
		}
		inner := make(map[*PublicKey][]byte, len(sigMap.SigPair))

		for _, sigPair := range sigMap.SigPair {
			key, err := PublicKeyFromBytes(sigPair.PubKeyPrefix)
			if err != nil {
				return make(map[AccountID]map[*PublicKey][]byte), err
			}
			switch sigPair.Signature.(type) {
			case *services.SignaturePair_Contract:
				inner[&key] = sigPair.GetContract()
			case *services.SignaturePair_Ed25519:
				inner[&key] = sigPair.GetEd25519()
			case *services.SignaturePair_RSA_3072:
				inner[&key] = sigPair.GetRSA_3072()
			case *services.SignaturePair_ECDSA_384:
				inner[&key] = sigPair.GetECDSA_384()
			}
		}

		returnMap[tempID] = inner
	}
	tx.transactionIDs.locked = true

	return returnMap, nil
}

func (tx *Transaction[T]) GetTransactionHash() ([]byte, error) {
	current, err := tx._BuildTransaction(0)
	if err != nil {
		return nil, err
	}
	hash := sha512.New384()
	_, err = hash.Write(current.GetSignedTransactionBytes())
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func (tx *Transaction[T]) GetTransactionHashPerNode() (map[AccountID][]byte, error) {
	transactionHash := make(map[AccountID][]byte)
	if !tx.IsFrozen() {
		return transactionHash, errTransactionIsNotFrozen
	}

	allTx, err := tx._BuildAllTransactions()
	if err != nil {
		return transactionHash, err
	}
	tx.transactionIDs.locked = true

	for i, node := range tx.nodeAccountIDs.slice {
		switch n := node.(type) { //nolint
		case AccountID:
			hash := sha512.New384()
			_, err := hash.Write(allTx[i].GetSignedTransactionBytes())
			if err != nil {
				return transactionHash, err
			}

			finalHash := hash.Sum(nil)

			transactionHash[n] = finalHash
		}
	}

	return transactionHash, nil
}

// Sets the maxTransaction fee based on priority:
// 1. Explicitly set for this Transaction
// 2. Client has a default value set for all transactions
// 3. The default for this type of Transaction, which is set during creation
func (tx *Transaction[T]) _InitFee(client *Client) {
	if tx.transactionFee == 0 {
		if client != nil && client.GetDefaultMaxTransactionFee().AsTinybar() != 0 {
			tx.SetMaxTransactionFee(client.GetDefaultMaxTransactionFee())
		} else {
			tx.SetMaxTransactionFee(tx.GetDefaultMaxTransactionFee())
		}
	}
}

func (tx *Transaction[T]) _InitTransactionID(client *Client) error {
	if tx.transactionIDs._Length() == 0 {
		if client != nil {
			if client.operator != nil {
				tx.transactionIDs = _NewLockableSlice()
				tx.transactionIDs = tx.transactionIDs._Push(TransactionIDGenerate(client.operator.accountID))
			} else {
				return errNoClientOrTransactionID
			}
		} else {
			return errNoClientOrTransactionID
		}
	}

	tx.transactionID = tx.transactionIDs._GetCurrent().(TransactionID)
	return nil
}

func (tx *Transaction[T]) IsFrozen() bool {
	return tx.signedTransactions._Length() > 0
}

func (tx *Transaction[T]) _RequireFrozen() {
	if !tx.IsFrozen() {
		tx.freezeError = errTransactionIsNotFrozen
	}
}

func (tx *Transaction[T]) _RequireNotFrozen() {
	if tx.IsFrozen() {
		tx.freezeError = errTransactionIsFrozen
	}
}

func (tx *Transaction[T]) _RequireOneNodeAccountID() {
	if tx.nodeAccountIDs._Length() != 1 {
		panic("transaction has more than one _Node ID set")
	}
}

func (tx *Transaction[T]) _SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) {
	tx.transactions = _NewLockableSlice()
	tx.publicKeys = append(tx.publicKeys, publicKey)
	tx.transactionSigners = append(tx.transactionSigners, signer)
}

func (tx *Transaction[T]) _KeyAlreadySigned(
	pk PublicKey,
) bool {
	for _, key := range tx.publicKeys {
		if key.String() == pk.String() {
			return true
		}
	}

	return false
}

// String returns a string representation of the transaction
func (tx *Transaction[T]) String() string {
	switch sig := tx.signedTransactions._Get(0).(type) { //nolint
	case *services.SignedTransaction:
		return fmt.Sprintf("%+v", sig)
	}

	return ""
}

// ToBytes Builds then converts the current transaction to []byte
// Requires transaction to be frozen
func (tx *Transaction[T]) ToBytes() ([]byte, error) {
	var pbTransactionList []byte
	var allTx []*services.Transaction
	var err error
	// If transaction is frozen, build all transactions and "signedTransactions"
	if tx.IsFrozen() {
		allTx, err = tx._BuildAllTransactions()
		tx.transactionIDs.locked = true
	} else { // Build only onlt "BodyBytes" for each transaction in the list
		allTx, err = tx.buildAllUnsignedTransactions()
	}
	// If error has occurred, when building transactions
	if err != nil {
		return make([]byte, 0), err
	}

	pbTransactionList, err = protobuf.Marshal(&sdk.TransactionList{
		TransactionList: allTx,
	})
	if err != nil {
		return make([]byte, 0), errors.Wrap(err, "error serializing tx list")
	}
	return pbTransactionList, nil
}

func (tx *Transaction[T]) buildAllUnsignedTransactions() ([]*services.Transaction, error) {
	// All unsigned transactions would always be exactly 1
	allTx := make([]*services.Transaction, 0)
	if tx.nodeAccountIDs._IsEmpty() {
		t, err := tx.buildUnsignedTransaction(0)
		if err != nil {
			return allTx, err
		}
		allTx = append(allTx, t)
	} else { // If we have set some node account ids, we have to make one transaction copy per node account
		for range tx.nodeAccountIDs.slice {
			t, err := tx.buildUnsignedTransaction(tx.nodeAccountIDs.index)
			tx.nodeAccountIDs._Advance()
			if err != nil {
				return allTx, err
			}
			allTx = append(allTx, t)
		}
	}
	return allTx, nil
}

func (tx *Transaction[T]) buildUnsignedTransaction(index int) (*services.Transaction, error) {
	body := tx.childTransaction.build()
	if body.NodeAccountID == nil && !tx.nodeAccountIDs._IsEmpty() {
		body.NodeAccountID = tx.nodeAccountIDs._Get(index).(AccountID)._ToProtobuf()
	}

	bodyBytes, err := protobuf.Marshal(body)
	if err != nil {
		return &services.Transaction{}, errors.Wrap(err, "failed to update tx ID")
	}

	return &services.Transaction{BodyBytes: bodyBytes}, nil
}

func (tx *Transaction[T]) _SignTransaction(index int) {
	initialTx := tx.signedTransactions._Get(index).(*services.SignedTransaction)
	bodyBytes := initialTx.GetBodyBytes()
	if len(initialTx.SigMap.SigPair) != 0 {
		for i, key := range tx.publicKeys {
			if tx.transactionSigners[i] != nil {
				if key.ed25519PublicKey != nil {
					if bytes.Equal(initialTx.SigMap.SigPair[0].PubKeyPrefix, key.ed25519PublicKey.keyData) {
						if !tx.regenerateTransactionID {
							return
						}
						switch t := initialTx.SigMap.SigPair[0].Signature.(type) { //nolint
						case *services.SignaturePair_Ed25519:
							if bytes.Equal(t.Ed25519, tx.transactionSigners[0](bodyBytes)) && len(t.Ed25519) > 0 {
								return
							}
						}
					}
				}
				if key.ecdsaPublicKey != nil {
					if bytes.Equal(initialTx.SigMap.SigPair[0].PubKeyPrefix, key.ecdsaPublicKey._BytesRaw()) {
						if !tx.regenerateTransactionID {
							return
						}
						switch t := initialTx.SigMap.SigPair[0].Signature.(type) { //nolint
						case *services.SignaturePair_ECDSASecp256K1:
							if bytes.Equal(t.ECDSASecp256K1, tx.transactionSigners[0](bodyBytes)) && len(t.ECDSASecp256K1) > 0 {
								return
							}
						}
					}
				}
			}
		}
	}

	if tx.regenerateTransactionID && !tx.transactionIDs.locked {
		modifiedTx := tx.signedTransactions._Get(index).(*services.SignedTransaction)
		modifiedTx.SigMap.SigPair = make([]*services.SignaturePair, 0)
		tx.signedTransactions._Set(index, modifiedTx)
	}

	for i := 0; i < len(tx.publicKeys); i++ {
		publicKey := tx.publicKeys[i]
		signer := tx.transactionSigners[i]

		if signer == nil {
			continue
		}

		modifiedTx := tx.signedTransactions._Get(index).(*services.SignedTransaction)
		modifiedTx.SigMap.SigPair = append(modifiedTx.SigMap.SigPair, publicKey._ToSignaturePairProtobuf(signer(bodyBytes)))
		tx.signedTransactions._Set(index, modifiedTx)
	}
}

func (tx *Transaction[T]) _BuildAllTransactions() ([]*services.Transaction, error) {
	allTx := make([]*services.Transaction, 0)
	for i := 0; i < tx.signedTransactions._Length(); i++ {
		curr, err := tx._BuildTransaction(i)
		tx.transactionIDs._Advance()
		if err != nil {
			return []*services.Transaction{}, err
		}
		allTx = append(allTx, curr)
	}

	return allTx, nil
}

func (tx *Transaction[T]) _BuildTransaction(index int) (*services.Transaction, error) {
	signedTx := tx.signedTransactions._Get(index).(*services.SignedTransaction)

	txID := tx.transactionIDs._GetCurrent().(TransactionID)
	originalBody := services.TransactionBody{}
	_ = protobuf.Unmarshal(signedTx.BodyBytes, &originalBody)

	if originalBody.NodeAccountID == nil {
		originalBody.NodeAccountID = tx.nodeAccountIDs._GetCurrent().(AccountID)._ToProtobuf()
	}

	if originalBody.TransactionID.String() != txID._ToProtobuf().String() {
		originalBody.TransactionID = txID._ToProtobuf()
	}

	originalBody.Memo = tx.memo
	if tx.transactionFee != 0 {
		originalBody.TransactionFee = tx.transactionFee
	} else {
		originalBody.TransactionFee = tx.defaultMaxTransactionFee
	}

	updatedBody, err := protobuf.Marshal(&originalBody)
	if err != nil {
		return &services.Transaction{}, errors.Wrap(err, "failed to update tx ID")
	}

	// Bellow are checks whether we need to sign the transaction or we already have the same signed
	if bytes.Equal(signedTx.BodyBytes, updatedBody) {
		sigPairLen := len(signedTx.SigMap.GetSigPair())
		// For cases where we need more than 1 signature
		if sigPairLen > 0 && sigPairLen == len(tx.publicKeys) {
			data, err := protobuf.Marshal(signedTx)
			if err != nil {
				return &services.Transaction{}, errors.Wrap(err, "failed to serialize transactions for building")
			}
			transaction := &services.Transaction{
				SignedTransactionBytes: data,
			}

			return transaction, nil
		}
	}

	signedTx.BodyBytes = updatedBody
	tx.signedTransactions._Set(index, signedTx)
	tx._SignTransaction(index)

	signed := tx.signedTransactions._Get(index).(*services.SignedTransaction)
	data, err := protobuf.Marshal(signed)
	if err != nil {
		return &services.Transaction{}, errors.Wrap(err, "failed to serialize transactions for building")
	}

	transaction := &services.Transaction{
		SignedTransactionBytes: data,
	}

	return transaction, nil
}

//
// Shared
//

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *Transaction[T]) GetMaxTransactionFee() Hbar {
	return HbarFromTinybar(int64(tx.transactionFee))
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *Transaction[T]) SetMaxTransactionFee(fee Hbar) T {
	tx.transactionFee = uint64(fee.AsTinybar())
	return tx.childTransaction
}
func (tx *Transaction[T]) GetDefaultMaxTransactionFee() Hbar {
	return HbarFromTinybar(int64(tx.defaultMaxTransactionFee))
}

// SetMaxTransactionFee sets the max Transaction fee for this Transaction.
func (tx *Transaction[T]) _SetDefaultMaxTransactionFee(fee Hbar) {
	tx.defaultMaxTransactionFee = uint64(fee.AsTinybar())
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled
func (tx *Transaction[T]) GetRegenerateTransactionID() bool {
	return tx.regenerateTransactionID
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when \`TRANSACTION_EXPIRED\` is received
func (tx *Transaction[T]) SetRegenerateTransactionID(regenerateTransactionID bool) T {
	tx.regenerateTransactionID = regenerateTransactionID
	return tx.childTransaction
}

// GetTransactionMemo returns the memo for this	transaction.
func (tx *Transaction[T]) GetTransactionMemo() string {
	return tx.memo
}

// SetTransactionMemo sets the memo for this transaction.
func (tx *Transaction[T]) SetTransactionMemo(memo string) T {
	tx.memo = memo
	return tx.childTransaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (tx *Transaction[T]) GetTransactionValidDuration() time.Duration {
	if tx.transactionValidDuration != nil {
		return *tx.transactionValidDuration
	}

	return 0
}

// SetTransactionValidDuration sets the valid duration for this transaction.
func (tx *Transaction[T]) SetTransactionValidDuration(duration time.Duration) T {
	tx.transactionValidDuration = &duration
	return tx.childTransaction
}

// GetTransactionID gets the TransactionID for this	transaction.
func (tx *Transaction[T]) GetTransactionID() TransactionID {
	if tx.transactionIDs._Length() > 0 {
		t := tx.transactionIDs._GetCurrent().(TransactionID)
		return t
	}

	return TransactionID{}
}

// SetTransactionID sets the TransactionID for this transaction.
func (tx *Transaction[T]) SetTransactionID(transactionID TransactionID) T {
	tx.transactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return tx.childTransaction
}

// SetNodeAccountIDs sets the node AccountID for this transaction.
func (tx *Transaction[T]) SetNodeAccountIDs(nodeAccountIDs []AccountID) T {
	for _, nodeAccountID := range nodeAccountIDs {
		tx.nodeAccountIDs._Push(nodeAccountID)
	}
	tx.nodeAccountIDs._SetLocked(true)
	return tx.childTransaction
}

// ------------ Transaction methdos ---------------
func (tx *Transaction[T]) Sign(privateKey PrivateKey) T {
	return tx.SignWith(privateKey.PublicKey(), privateKey.Sign)
}
func (tx *Transaction[T]) SignWithOperator(client *Client) (T, error) { // nolint
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	if client == nil {
		return tx.childTransaction, errNoClientProvided
	} else if client.operator == nil {
		return tx.childTransaction, errClientOperatorSigning
	}

	if !tx.IsFrozen() {
		_, err := tx.FreezeWith(client)
		if err != nil {
			return tx.childTransaction, err
		}
	}
	return tx.SignWith(client.operator.publicKey, client.operator.signer), nil
}
func (tx *Transaction[T]) SignWith(publicKey PublicKey, signer TransactionSigner) T {
	// We need to make sure the request is frozen
	tx._RequireFrozen()

	if !tx._KeyAlreadySigned(publicKey) {
		tx._SignWith(publicKey, signer)
	}

	return tx.childTransaction
}
func (tx *Transaction[T]) AddSignature(publicKey PublicKey, signature []byte) T {
	tx._RequireOneNodeAccountID()

	if tx._KeyAlreadySigned(publicKey) {
		return tx.childTransaction
	}

	if tx.signedTransactions._Length() == 0 {
		return tx.childTransaction
	}

	tx.transactions = _NewLockableSlice()
	tx.publicKeys = append(tx.publicKeys, publicKey)
	tx.transactionSigners = append(tx.transactionSigners, nil)
	tx.transactionIDs.locked = true

	for index := 0; index < tx.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := tx.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		tx.signedTransactions._Set(index, temp)
	}

	return tx.childTransaction
}

func (tx *Transaction[T]) preFreezeWith(*Client) {
	// NO-OP
	// TODO
}

func (tx *Transaction[T]) getLogID(transactionInterface Executable) string {
	timestamp := tx.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("%s:%d", transactionInterface.getName(), timestamp.UnixNano())
}

// ------------ Executable Functions ------------
func (tx *Transaction[T]) shouldRetry(_ Executable, response interface{}) _ExecutionState {
	status := Status(response.(*services.TransactionResponse).NodeTransactionPrecheckCode)

	retryableStatuses := map[Status]bool{
		StatusPlatformTransactionNotCreated: true,
		StatusPlatformNotActive:             true,
		StatusBusy:                          true,
	}

	if retryableStatuses[status] {
		return executionStateRetry
	}

	if status == StatusTransactionExpired {
		return executionStateExpired
	}

	if status == StatusOk {
		return executionStateFinished
	}

	return executionStateError
}

func (tx *Transaction[T]) makeRequest() interface{} {
	index := tx.nodeAccountIDs._Length()*tx.transactionIDs.index + tx.nodeAccountIDs.index
	built, _ := tx._BuildTransaction(index)

	return built
}

func (tx *Transaction[T]) advanceRequest() {
	tx.nodeAccountIDs._Advance()
	tx.signedTransactions._Advance()
}

func (tx *Transaction[T]) getNodeAccountID() AccountID {
	return tx.nodeAccountIDs._GetCurrent().(AccountID)
}

func (tx *Transaction[T]) mapStatusError(
	_ Executable,
	response interface{},
) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.TransactionResponse).NodeTransactionPrecheckCode),
		//NodeID: request.transaction.nodeAccountIDs,
		TxID: tx.GetTransactionID(),
	}
}

func (tx *Transaction[T]) mapResponse(_ interface{}, nodeID AccountID, protoRequest interface{}) (interface{}, error) {
	hash := sha512.New384()
	_, err := hash.Write(protoRequest.(*services.Transaction).SignedTransactionBytes)
	if err != nil {
		return nil, err
	}

	return TransactionResponse{
		NodeID:        nodeID,
		TransactionID: tx.transactionIDs._GetNext().(TransactionID),
		Hash:          hash.Sum(nil),
	}, nil
}

func (tx *Transaction[T]) isTransaction() bool {
	return true
}

func (tx *Transaction[T]) getTransactionIDAndMessage() (string, string) {
	return tx.GetTransactionID().String(), "transaction status received"
}

func (tx *Transaction[T]) regenerateID(client *Client) bool {
	if !client.GetOperatorAccountID()._IsZero() && tx.regenerateTransactionID && !tx.transactionIDs.locked {
		tx.transactionIDs._Set(tx.transactionIDs.index, TransactionIDGenerate(client.GetOperatorAccountID()))
		return true
	}
	return false
}

func (tx *Transaction[T]) Execute(client *Client) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if tx.freezeError != nil {
		return TransactionResponse{}, tx.freezeError
	}

	if !tx.IsFrozen() {
		_, err := tx.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := tx.transactionIDs._GetCurrent().(TransactionID)

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		tx.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	if tx.grpcDeadline == nil {
		tx.grpcDeadline = client.requestTimeout
	}

	resp, err := _Execute(client, tx.childTransaction)

	if err != nil {
		return TransactionResponse{
			TransactionID:  tx.GetTransactionID(),
			NodeID:         resp.(TransactionResponse).NodeID,
			ValidateStatus: true,
		}, err
	}
	originalTxID := tx.GetTransactionID()
	tx.childTransaction.regenerateID(client)
	return TransactionResponse{
		TransactionID:  originalTxID,
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
		// set the tx in the response, in case of throttle error in the receipt
		// we can use this to re-submit the transaction
		Transaction: *tx.childTransaction.getBaseTransaction(),
	}, nil
}

func (tx *Transaction[T]) Freeze() (T, error) {
	return tx.FreezeWith(nil)
}

func (tx *Transaction[T]) FreezeWith(client *Client) (T, error) {
	if tx.IsFrozen() {
		return tx.childTransaction, nil
	}

	tx.childTransaction.preFreezeWith(client)

	tx._InitFee(client)
	if err := tx._InitTransactionID(client); err != nil {
		return tx.childTransaction, err
	}

	err := tx.childTransaction.validateNetworkOnIDs(client)
	if err != nil {
		return tx.childTransaction, err
	}
	body := tx.childTransaction.build()

	if tx.nodeAccountIDs._IsEmpty() {
		if client != nil {
			for _, nodeAccountID := range client.network._GetNodeAccountIDsForExecute() {
				tx.nodeAccountIDs._Push(nodeAccountID)
			}
		} else {
			return tx.childTransaction, errNoClientOrTransactionIDOrNodeId
		}
	}

	if client != nil {
		if client.defaultRegenerateTransactionIDs != tx.regenerateTransactionID {
			tx.regenerateTransactionID = client.defaultRegenerateTransactionIDs
		}
	}

	for _, nodeAccountID := range tx.nodeAccountIDs.slice {
		body.NodeAccountID = nodeAccountID.(AccountID)._ToProtobuf()
		bodyBytes, err := protobuf.Marshal(body)

		if err != nil {
			// This should be unreachable
			// From the documentation this appears to only be possible if there are missing proto types
			panic(err)
		}
		tx.signedTransactions = tx.signedTransactions._Push(&services.SignedTransaction{
			BodyBytes: bodyBytes,
			SigMap: &services.SignatureMap{
				SigPair: make([]*services.SignaturePair, 0),
			},
		})
	}

	return tx.childTransaction, nil
}

func (tx *Transaction[T]) Schedule() (*ScheduleCreateTransaction, error) {
	tx._RequireNotFrozen()

	scheduled, err := tx.childTransaction.buildScheduled()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (tx *Transaction[T]) GetMaxBackoff() time.Duration {
	if tx.maxBackoff != nil {
		return *tx.maxBackoff
	}

	return 8 * time.Second
}

func (tx *Transaction[T]) GetMinBackoff() time.Duration {
	if tx.minBackoff != nil {
		return *tx.minBackoff
	}

	return 250 * time.Millisecond
}

func (tx *Transaction[T]) SetMaxBackoff(max time.Duration) T {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < tx.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	tx.maxBackoff = &max
	return tx.childTransaction
}

func (tx *Transaction[T]) SetMinBackoff(min time.Duration) T {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if tx.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	tx.minBackoff = &min
	return tx.childTransaction
}

// GetGrpcDeadline returns the grpc deadline
func (tx *Transaction[T]) GetGrpcDeadline() *time.Duration {
	return tx.grpcDeadline
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *Transaction[T]) SetGrpcDeadline(deadline *time.Duration) T {
	tx.grpcDeadline = deadline
	return tx.childTransaction
}

// GetMaxRetry returns the max number of errors before execution will fail.
func (tx *Transaction[T]) GetMaxRetry() int {
	return tx.maxRetry
}

func (tx *Transaction[T]) SetMaxRetry(max int) T {
	tx.maxRetry = max
	return tx.childTransaction
}

func (tx *Transaction[T]) GetLogLevel() *LogLevel {
	return tx.logLevel
}

func (tx *Transaction[T]) SetLogLevel(level LogLevel) T {
	tx.logLevel = &level
	return tx.childTransaction
}

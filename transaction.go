package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// Interface that all concrete transactions must implement, eg. TransferTransaction, ContractCreateTransaction, etc.
type TransactionInterface interface {
	// common methods for all executables
	Executable

	// methods implemented by the parent transaction
	regenerateID(*Client) bool // creates new transaction ID

	// methods implemented by every concrete transaction
	build() *services.TransactionBody                                         // build a protobuf payload for the transaction
	buildScheduled() (*services.SchedulableTransactionBody, error)            // builds the protobuf payload for the scheduled transaction
	preFreezeWith(*Client, TransactionInterface)                              // utility method to set the transaction fields before freezing
	constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) // TODO remove this method if possible
	// NOTE: Any changes to the baseTransaction retuned by getBaseTransaction()
	// will be reflected in the transaction object
	getBaseTransaction() *Transaction[TransactionInterface]
}

// BaseTransaction contains all the common fields for all transactions.
type BaseTransaction struct {
	transactionFee           uint64
	defaultMaxTransactionFee uint64
	memo                     string
	transactionValidDuration *time.Duration
	transactionID            TransactionID

	transactions       *_LockableSlice
	signedTransactions *_LockableSlice

	publicKeys         []PublicKey
	transactionSigners []TransactionSigner
}

// Transaction is base struct for all transactions that may be built and submitted to hiero.
// It's generic over the type of transaction it contains. Example: TransferTransaction, ContractCreateTransaction, etc.
type Transaction[T TransactionInterface] struct {
	*executable
	*BaseTransaction
	childTransaction T

	freezeError error

	regenerateTransactionID bool
}

// Creates new transaction, embedding the concrete transaction.
func _NewTransaction[T TransactionInterface](concreteTransaction T) *Transaction[T] {
	duration := 120 * time.Second
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	return &Transaction[T]{
		BaseTransaction: &BaseTransaction{
			transactionValidDuration: &duration,
			transactions:             _NewLockableSlice(),
			signedTransactions:       _NewLockableSlice(),
		},
		childTransaction:        concreteTransaction,
		freezeError:             nil,
		regenerateTransactionID: true,
		executable: &executable{
			transactionIDs: _NewLockableSlice(),
			nodeAccountIDs: _NewLockableSlice(),
			minBackoff:     &minBackoff,
			maxBackoff:     &maxBackoff,
			maxRetry:       10,
		},
	}
}

// TransactionFromBytes converts transaction bytes to a related *transaction.
func TransactionFromBytes(data []byte) (TransactionInterface, error) { // nolint
	list := sdk.TransactionList{}
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	publicKeys := make([]PublicKey, 0)
	transactionSigners := make([]TransactionSigner, 0)
	err := protobuf.Unmarshal(data, &list)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing from bytes to transaction List")
	}

	transactions := _NewLockableSlice()

	for _, transaction := range list.TransactionList {
		transactions._Push(transaction)
	}

	baseTx := Transaction[TransactionInterface]{
		BaseTransaction: &BaseTransaction{
			signedTransactions: _NewLockableSlice(),
			publicKeys:         publicKeys,
			transactionSigners: transactionSigners,
			transactions:       transactions,
		},
		freezeError:             nil,
		regenerateTransactionID: true,
		executable: &executable{
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
			baseTx.signedTransactions = baseTx.signedTransactions._Push(&signedTransaction)

			if i == 0 {
				for _, sigPair := range signedTransaction.GetSigMap().GetSigPair() {
					key, err := PublicKeyFromBytes(sigPair.GetPubKeyPrefix())
					if err != nil {
						return nil, err
					}

					baseTx.publicKeys = append(baseTx.publicKeys, key)
					baseTx.transactionSigners = append(baseTx.transactionSigners, nil)
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
			baseTx.transactionValidDuration = &duration
		}

		if body.GetTransactionID() != nil {
			transactionID = _TransactionIDFromProtobuf(body.GetTransactionID())
		}

		if body.GetNodeAccountID() != nil {
			nodeAccountID = *_AccountIDFromProtobuf(body.GetNodeAccountID())
		}

		// If the transaction was serialised, without setting "NodeId", or "TransactionID", we should leave them empty
		if transactionID.AccountID.Account != 0 {
			baseTx.transactionIDs = baseTx.transactionIDs._Push(transactionID)
		}
		if !nodeAccountID._IsZero() {
			baseTx.nodeAccountIDs = baseTx.nodeAccountIDs._Push(nodeAccountID)
		}

		if i == 0 {
			baseTx.memo = body.Memo
			if body.TransactionFee != 0 {
				baseTx.transactionFee = body.TransactionFee
			}
		}
	}

	if txIsSigned {
		if baseTx.transactionIDs._Length() > 0 {
			baseTx.transactionIDs.locked = true
		}

		if baseTx.nodeAccountIDs._Length() > 0 {
			baseTx.nodeAccountIDs.locked = true
		}
	}

	if first == nil {
		return nil, errNoTransactionInBytes
	}

	var childTx TransactionInterface

	switch first.Data.(type) {
	case *services.TransactionBody_ContractCall:
		childTx = _ContractExecuteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ContractExecuteTransaction](baseTx), first)
	case *services.TransactionBody_ContractCreateInstance:
		childTx = _ContractCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ContractCreateTransaction](baseTx), first)
	case *services.TransactionBody_ContractUpdateInstance:
		childTx = _ContractUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ContractUpdateTransaction](baseTx), first)
	case *services.TransactionBody_CryptoApproveAllowance:
		childTx = _AccountAllowanceApproveTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountAllowanceApproveTransaction](baseTx), first)
	case *services.TransactionBody_CryptoDeleteAllowance:
		childTx = _AccountAllowanceDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountAllowanceDeleteTransaction](baseTx), first)
	case *services.TransactionBody_ContractDeleteInstance:
		childTx = _ContractDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ContractDeleteTransaction](baseTx), first)
	case *services.TransactionBody_CryptoAddLiveHash:
		childTx = _LiveHashAddTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*LiveHashAddTransaction](baseTx), first)
	case *services.TransactionBody_CryptoCreateAccount:
		childTx = _AccountCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountCreateTransaction](baseTx), first)
	case *services.TransactionBody_CryptoDelete:
		childTx = _AccountDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountDeleteTransaction](baseTx), first)
	case *services.TransactionBody_CryptoDeleteLiveHash:
		childTx = _LiveHashDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*LiveHashDeleteTransaction](baseTx), first)
	case *services.TransactionBody_CryptoTransfer:
		childTx = _TransferTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TransferTransaction](baseTx), first)
	case *services.TransactionBody_CryptoUpdateAccount:
		childTx = _AccountUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountUpdateTransaction](baseTx), first)
	case *services.TransactionBody_FileAppend:
		childTx = _FileAppendTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FileAppendTransaction](baseTx), first)
	case *services.TransactionBody_FileCreate:
		childTx = _FileCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FileCreateTransaction](baseTx), first)
	case *services.TransactionBody_FileDelete:
		childTx = _FileDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FileDeleteTransaction](baseTx), first)
	case *services.TransactionBody_FileUpdate:
		childTx = _FileUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FileUpdateTransaction](baseTx), first)
	case *services.TransactionBody_SystemDelete:
		childTx = _SystemDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*SystemDeleteTransaction](baseTx), first)
	case *services.TransactionBody_SystemUndelete:
		childTx = _SystemUndeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*SystemUndeleteTransaction](baseTx), first)
	case *services.TransactionBody_Freeze:
		childTx = _FreezeTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FreezeTransaction](baseTx), first)
	case *services.TransactionBody_ConsensusCreateTopic:
		childTx = _TopicCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TopicCreateTransaction](baseTx), first)
	case *services.TransactionBody_ConsensusUpdateTopic:
		childTx = _TopicUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TopicUpdateTransaction](baseTx), first)
	case *services.TransactionBody_ConsensusDeleteTopic:
		childTx = _TopicDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TopicDeleteTransaction](baseTx), first)
	case *services.TransactionBody_ConsensusSubmitMessage:
		childTx = _TopicMessageSubmitTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TopicMessageSubmitTransaction](baseTx), first)
	case *services.TransactionBody_TokenCreation:
		childTx = _TokenCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenCreateTransaction](baseTx), first)
	case *services.TransactionBody_TokenFreeze:
		childTx = _TokenFreezeTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenFreezeTransaction](baseTx), first)
	case *services.TransactionBody_TokenUnfreeze:
		childTx = _TokenUnfreezeTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenUnfreezeTransaction](baseTx), first)
	case *services.TransactionBody_TokenGrantKyc:
		childTx = _TokenGrantKycTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenGrantKycTransaction](baseTx), first)
	case *services.TransactionBody_TokenRevokeKyc:
		childTx = _TokenRevokeKycTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenRevokeKycTransaction](baseTx), first)
	case *services.TransactionBody_TokenDeletion:
		childTx = _TokenDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenDeleteTransaction](baseTx), first)
	case *services.TransactionBody_TokenUpdate:
		childTx = _TokenUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenUpdateTransaction](baseTx), first)
	case *services.TransactionBody_TokenMint:
		childTx = _TokenMintTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenMintTransaction](baseTx), first)
	case *services.TransactionBody_TokenBurn:
		childTx = _TokenBurnTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenBurnTransaction](baseTx), first)
	case *services.TransactionBody_TokenWipe:
		childTx = _TokenWipeTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenWipeTransaction](baseTx), first)
	case *services.TransactionBody_TokenAssociate:
		childTx = _TokenAssociateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenAssociateTransaction](baseTx), first)
	case *services.TransactionBody_TokenDissociate:
		childTx = _TokenDissociateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenDissociateTransaction](baseTx), first)
	case *services.TransactionBody_ScheduleCreate:
		childTx = _ScheduleCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ScheduleCreateTransaction](baseTx), first)
	case *services.TransactionBody_ScheduleDelete:
		childTx = _ScheduleDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ScheduleDeleteTransaction](baseTx), first)
	case *services.TransactionBody_ScheduleSign:
		childTx = _ScheduleSignTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ScheduleSignTransaction](baseTx), first)
	case *services.TransactionBody_TokenPause:
		childTx = _TokenPauseTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenPauseTransaction](baseTx), first)
	case *services.TransactionBody_TokenUnpause:
		childTx = _TokenUnpauseTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenUnpauseTransaction](baseTx), first)
	case *services.TransactionBody_EthereumTransaction:
		childTx = _EthereumTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*EthereumTransaction](baseTx), first)
	case *services.TransactionBody_UtilPrng:
		childTx = _PrngTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*PrngTransaction](baseTx), first)
	case *services.TransactionBody_TokenReject:
		childTx = _TokenRejectTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenRejectTransaction](baseTx), first)
	case *services.TransactionBody_TokenFeeScheduleUpdate:
		childTx = _TokenFeeScheduleUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenFeeScheduleUpdateTransaction](baseTx), first)
	case *services.TransactionBody_TokenUpdateNfts:
		childTx = _TokenUpdateNftsTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenUpdateNfts](baseTx), first)
	case *services.TransactionBody_NodeCreate:
		childTx = _NodeCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*NodeCreateTransaction](baseTx), first)
	case *services.TransactionBody_NodeUpdate:
		childTx = _NodeUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*NodeUpdateTransaction](baseTx), first)
	case *services.TransactionBody_NodeDelete:
		childTx = _NodeDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*NodeDeleteTransaction](baseTx), first)
	case *services.TransactionBody_TokenAirdrop:
		childTx = _TokenAirdropTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenAirdropTransaction](baseTx), first)
	case *services.TransactionBody_TokenCancelAirdrop:
		childTx = _TokenCancelAirdropTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenCancelAirdropTransaction](baseTx), first)
	case *services.TransactionBody_TokenClaimAirdrop:
		childTx = _TokenClaimAirdropTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenClaimAirdropTransaction](baseTx), first)
	default:
		return nil, errFailedToDeserializeBytes
	}

	// --- //

	return childTx, nil
}

// Creates a new transaction from a scheduled transaction body
func transactionFromScheduledTransaction(scheduledBody *services.SchedulableTransactionBody) (TransactionInterface, error) { // nolint
	pbBody := &services.TransactionBody{}

	memo := scheduledBody.GetMemo()
	baseTx := Transaction[TransactionInterface]{
		BaseTransaction: &BaseTransaction{
			memo:           memo,
			transactionFee: scheduledBody.GetTransactionFee(),
		},
	}

	var tx TransactionInterface

	switch scheduledBody.Data.(type) {
	case *services.SchedulableTransactionBody_ContractCall:
		pbBody.Data = &services.TransactionBody_ContractCall{
			ContractCall: scheduledBody.GetContractCall(),
		}
		tx = _ContractExecuteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ContractExecuteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_ContractCreateInstance:
		pbBody.Data = &services.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: scheduledBody.GetContractCreateInstance(),
		}
		tx = _ContractCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ContractCreateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_ContractUpdateInstance:
		pbBody.Data = &services.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: scheduledBody.GetContractUpdateInstance(),
		}
		tx = _ContractUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ContractUpdateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_CryptoApproveAllowance:
		pbBody.Data = &services.TransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: scheduledBody.GetCryptoApproveAllowance(),
		}
		tx = _AccountAllowanceApproveTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountAllowanceApproveTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_CryptoDeleteAllowance:
		pbBody.Data = &services.TransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: scheduledBody.GetCryptoDeleteAllowance(),
		}
		tx = _AccountAllowanceDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountAllowanceDeleteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_ContractDeleteInstance:
		pbBody.Data = &services.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: scheduledBody.GetContractDeleteInstance(),
		}
		tx = _ContractDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ContractDeleteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_CryptoCreateAccount:
		pbBody.Data = &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: scheduledBody.GetCryptoCreateAccount(),
		}
		tx = _AccountCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountCreateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_CryptoDelete:
		pbBody.Data = &services.TransactionBody_CryptoDelete{
			CryptoDelete: scheduledBody.GetCryptoDelete(),
		}
		tx = _AccountDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountDeleteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_CryptoTransfer:
		pbBody.Data = &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: scheduledBody.GetCryptoTransfer(),
		}
		tx = _TransferTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TransferTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_CryptoUpdateAccount:
		pbBody.Data = &services.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: scheduledBody.GetCryptoUpdateAccount(),
		}
		tx = _AccountUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*AccountUpdateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_FileAppend:
		pbBody.Data = &services.TransactionBody_FileAppend{
			FileAppend: scheduledBody.GetFileAppend(),
		}
		tx = _FileAppendTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FileAppendTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_FileCreate:
		pbBody.Data = &services.TransactionBody_FileCreate{
			FileCreate: scheduledBody.GetFileCreate(),
		}
		tx = _FileCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FileCreateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_FileDelete:
		pbBody.Data = &services.TransactionBody_FileDelete{
			FileDelete: scheduledBody.GetFileDelete(),
		}
		tx = _FileDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FileDeleteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_FileUpdate:
		pbBody.Data = &services.TransactionBody_FileUpdate{
			FileUpdate: scheduledBody.GetFileUpdate(),
		}
		tx = _FileUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FileUpdateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_SystemDelete:
		pbBody.Data = &services.TransactionBody_SystemDelete{
			SystemDelete: scheduledBody.GetSystemDelete(),
		}
		tx = _SystemDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*SystemDeleteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_SystemUndelete:
		pbBody.Data = &services.TransactionBody_SystemUndelete{
			SystemUndelete: scheduledBody.GetSystemUndelete(),
		}
		tx = _SystemUndeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*SystemUndeleteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_Freeze:
		pbBody.Data = &services.TransactionBody_Freeze{
			Freeze: scheduledBody.GetFreeze(),
		}
		tx = _FreezeTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*FreezeTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_ConsensusCreateTopic:
		pbBody.Data = &services.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: scheduledBody.GetConsensusCreateTopic(),
		}
		tx = _TopicCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TopicCreateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_ConsensusUpdateTopic:
		pbBody.Data = &services.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: scheduledBody.GetConsensusUpdateTopic(),
		}
		tx = _TopicUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TopicUpdateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_ConsensusDeleteTopic:
		pbBody.Data = &services.TransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: scheduledBody.GetConsensusDeleteTopic(),
		}
		tx = _TopicDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TopicDeleteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_ConsensusSubmitMessage:
		pbBody.Data = &services.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: scheduledBody.GetConsensusSubmitMessage(),
		}
		tx = _TopicMessageSubmitTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TopicMessageSubmitTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenCreation:
		pbBody.Data = &services.TransactionBody_TokenCreation{
			TokenCreation: scheduledBody.GetTokenCreation(),
		}
		tx = _TokenCreateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenCreateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenFreeze:
		pbBody.Data = &services.TransactionBody_TokenFreeze{
			TokenFreeze: scheduledBody.GetTokenFreeze(),
		}
		tx = _TokenFreezeTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenFreezeTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenUnfreeze:
		pbBody.Data = &services.TransactionBody_TokenUnfreeze{
			TokenUnfreeze: scheduledBody.GetTokenUnfreeze(),
		}
		tx = _TokenUnfreezeTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenUnfreezeTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenGrantKyc:
		pbBody.Data = &services.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: scheduledBody.GetTokenGrantKyc(),
		}
		tx = _TokenGrantKycTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenGrantKycTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenRevokeKyc:
		pbBody.Data = &services.TransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: scheduledBody.GetTokenRevokeKyc(),
		}
		tx = _TokenRevokeKycTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenRevokeKycTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenDeletion:
		pbBody.Data = &services.TransactionBody_TokenDeletion{
			TokenDeletion: scheduledBody.GetTokenDeletion(),
		}
		tx = _TokenDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenDeleteTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenUpdate:
		pbBody.Data = &services.TransactionBody_TokenUpdate{
			TokenUpdate: scheduledBody.GetTokenUpdate(),
		}
		tx = _TokenUpdateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenUpdateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenMint:
		pbBody.Data = &services.TransactionBody_TokenMint{
			TokenMint: scheduledBody.GetTokenMint(),
		}
		tx = _TokenMintTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenMintTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenBurn:
		pbBody.Data = &services.TransactionBody_TokenBurn{
			TokenBurn: scheduledBody.GetTokenBurn(),
		}
		tx = _TokenBurnTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenBurnTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenWipe:
		pbBody.Data = &services.TransactionBody_TokenWipe{
			TokenWipe: scheduledBody.GetTokenWipe(),
		}
		tx = _TokenWipeTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenWipeTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenAssociate:
		pbBody.Data = &services.TransactionBody_TokenAssociate{
			TokenAssociate: scheduledBody.GetTokenAssociate(),
		}
		tx = _TokenAssociateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenAssociateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_TokenDissociate:
		pbBody.Data = &services.TransactionBody_TokenDissociate{
			TokenDissociate: scheduledBody.GetTokenDissociate(),
		}
		tx = _TokenDissociateTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*TokenDissociateTransaction](baseTx), pbBody)
	case *services.SchedulableTransactionBody_ScheduleDelete:
		pbBody.Data = &services.TransactionBody_ScheduleDelete{
			ScheduleDelete: scheduledBody.GetScheduleDelete(),
		}
		tx = _ScheduleDeleteTransactionFromProtobuf(*castFromBaseToConcreteTransaction[*ScheduleDeleteTransaction](baseTx), pbBody)
	default:
		return nil, errors.New("unrecognized transaction type")
	}

	return tx, nil
}

// Private methods //

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
		publicKey := (tx.publicKeys)[i]
		signer := tx.transactionSigners[i]

		if signer == nil {
			continue
		}

		signature := signer(bodyBytes)
		if len(signature) == 65 {
			signature = signature[1:]
		}
		modifiedTx := tx.signedTransactions._Get(index).(*services.SignedTransaction)
		modifiedTx.SigMap.SigPair = append(modifiedTx.SigMap.SigPair, publicKey._ToSignaturePairProtobuf(signature))
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

// GetSignedTransactionBodyBytes
func (tx *Transaction[T]) GetSignedTransactionBodyBytes(transactionIndex int) []byte {
	return tx.signedTransactions._Get(transactionIndex).(*services.SignedTransaction).GetBodyBytes()
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
		return *new(T), errNoClientProvided
	} else if client.operator == nil {
		return *new(T), errClientOperatorSigning
	}

	if !tx.IsFrozen() {
		_, err := tx.FreezeWith(client)
		if err != nil {
			return *new(T), err
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

func (tx *Transaction[T]) preFreezeWith(*Client, TransactionInterface) {
	// No-op for every transaction except TokenCreateTransaction
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
		TxID:   tx.GetTransactionID(),
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
	tx.regenerateID(client)
	return TransactionResponse{
		TransactionID:  originalTxID,
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
		// set the tx in the response, in case of throttle error in the receipt
		// we can use this to re-submit the transaction
		Transaction: tx.childTransaction,
	}, nil
}

func (tx *Transaction[T]) Freeze() (T, error) {
	return tx.FreezeWith(nil)
}

func (tx *Transaction[T]) FreezeWith(client *Client) (T, error) {
	if tx.IsFrozen() {
		return tx.childTransaction, nil
	}

	tx.childTransaction.preFreezeWith(client, tx.childTransaction)

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

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *Transaction[T]) SetMaxRetry(max int) T {
	tx.maxRetry = max
	return tx.childTransaction
}

// GetNodeAccountIDs returns the node AccountID for this transaction.
func (tx *Transaction[T]) GetLogLevel() *LogLevel {
	return tx.logLevel
}

// SetNodeAccountIDs sets the node AccountID for this transaction.
func (tx *Transaction[T]) SetLogLevel(level LogLevel) T {
	tx.logLevel = &level
	return tx.childTransaction
}

// Static Utility functions //

func TransactionExecute(tx TransactionInterface, client *Client) (TransactionResponse, error) {
	return tx.getBaseTransaction().Execute(client)
}

func TransactionSign(tx TransactionInterface, key PrivateKey) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.Sign(key)

	return tx, nil
}

func TransactionAddSignature(tx TransactionInterface, publicKey PublicKey, signature []byte) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.AddSignature(publicKey, signature)

	return tx, nil
}

func TransactionToBytes(tx TransactionInterface) ([]byte, error) {
	return tx.getBaseTransaction().ToBytes()
}

func TransactionString(tx TransactionInterface) (string, error) {
	return tx.getBaseTransaction().String(), nil
}

func TransactionGetMaxBackoff(tx TransactionInterface) (time.Duration, error) {
	return tx.getBaseTransaction().GetMaxBackoff(), nil
}

func TransactionSetMaxBackoff(tx TransactionInterface, maxBackoff time.Duration) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.SetMaxBackoff(maxBackoff)

	return tx, nil
}

func TransactionGetMinBackoff(tx TransactionInterface) (time.Duration, error) {
	return tx.getBaseTransaction().GetMinBackoff(), nil
}

func TransactionSetMinBackoff(tx TransactionInterface, minBackoff time.Duration) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.SetMinBackoff(minBackoff)

	return tx, nil
}

func TransactionGetTransactionHashPerNode(tx TransactionInterface) (map[AccountID][]byte, error) {
	return tx.getBaseTransaction().GetTransactionHashPerNode()
}

func TransactionGetTransactionHash(tx TransactionInterface) ([]byte, error) {
	return tx.getBaseTransaction().GetTransactionHash()
}

func TransactionGetNodeAccountIDs(tx TransactionInterface) ([]AccountID, error) {
	return tx.getBaseTransaction().GetNodeAccountIDs(), nil
}

func TransactionSetNodeAccountIDs(tx TransactionInterface, nodeAccountIDs []AccountID) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.SetNodeAccountIDs(nodeAccountIDs)

	return tx, nil
}

func TransactionGetTransactionValidDuration(tx TransactionInterface) (time.Duration, error) {
	return tx.getBaseTransaction().GetTransactionValidDuration(), nil
}

func TransactionSetTransactionValidDuration(tx TransactionInterface, transactionValidDuration time.Duration) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.SetTransactionValidDuration(transactionValidDuration)

	return tx, nil
}

func TransactionGetMaxTransactionFee(tx TransactionInterface) (Hbar, error) {
	return tx.getBaseTransaction().GetMaxTransactionFee(), nil
}

func TransactionSetMaxTransactionFee(tx TransactionInterface, maxTransactionFee Hbar) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.SetMaxTransactionFee(maxTransactionFee)

	return tx, nil
}

func TransactionGetTransactionMemo(tx TransactionInterface) (string, error) {
	return tx.getBaseTransaction().GetTransactionMemo(), nil
}

func TransactionSetTransactionMemo(tx TransactionInterface, transactionMemo string) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.SetTransactionMemo(transactionMemo)
	return tx, nil
}

func TransactionGetTransactionID(tx TransactionInterface) (TransactionID, error) {
	return tx.getBaseTransaction().GetTransactionID(), nil
}

func TransactionSetTransactionID(tx TransactionInterface, transactionID TransactionID) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.SetTransactionID(transactionID)
	return tx, nil
}

func TransactionGetSignatures(tx TransactionInterface) (map[AccountID]map[*PublicKey][]byte, error) {
	return tx.getBaseTransaction().GetSignatures()
}

func TransactionFreezeWith(tx TransactionInterface, client *Client) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	_, err := baseTx.FreezeWith(client)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func TransactionSignWithOperator(tx TransactionInterface, client *Client) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	_, err := baseTx.SignWithOperator(client)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func TransactionSignWth(tx TransactionInterface, publicKKey PublicKey, signer TransactionSigner) (TransactionInterface, error) {
	baseTx := tx.getBaseTransaction()
	baseTx.SignWith(publicKKey, signer)

	return tx, nil
}

// Helper function to cast the concrete Transaction to the generic Transaction
func castFromConcreteToBaseTransaction[T TransactionInterface](baseTx *Transaction[T], tx TransactionInterface) *Transaction[TransactionInterface] {
	return &Transaction[TransactionInterface]{
		executable:              baseTx.executable,
		BaseTransaction:         baseTx.BaseTransaction,
		childTransaction:        tx,
		freezeError:             baseTx.freezeError,
		regenerateTransactionID: baseTx.regenerateTransactionID,
	}
}

// Helper function to cast the generic Transaction to another type
func castFromBaseToConcreteTransaction[T TransactionInterface](baseTx Transaction[TransactionInterface]) *Transaction[T] {
	concreteTx := &Transaction[T]{
		executable:              baseTx.executable,
		BaseTransaction:         baseTx.BaseTransaction,
		freezeError:             baseTx.freezeError,
		regenerateTransactionID: baseTx.regenerateTransactionID,
	}
	if baseTx.childTransaction != nil {
		concreteTx.childTransaction = baseTx.childTransaction.(T)
	}
	return concreteTx
}

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
}

// Transaction is base struct for all transactions that may be built and submitted to Hedera.
type Transaction struct {
	executable

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

func _NewTransaction() Transaction {
	duration := 120 * time.Second
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	return Transaction{
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

func (tx *Transaction) GetSignedTransactionBodyBytes(transactionIndex int) []byte {
	return tx.signedTransactions._Get(transactionIndex).(*services.SignedTransaction).GetBodyBytes()
}

// TransactionFromBytes converts transaction bytes to a related *transaction.
func TransactionFromBytes(data []byte) (interface{}, error) { // nolint
	list := sdk.TransactionList{}
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	err := protobuf.Unmarshal(data, &list)
	if err != nil {
		return Transaction{}, errors.Wrap(err, "error deserializing from bytes to transaction List")
	}
	if err != nil {
		return Transaction{}, err
	}

	transactions := _NewLockableSlice()
	for _, transaction := range list.TransactionList {
		transactions._Push(transaction)
	}

	tx := Transaction{
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
		return Transaction{}, err
	}

	if !comp {
		return Transaction{}, errors.New("failed to validate transaction bodies")
	}

	var first *services.TransactionBody = nil

	for i, transactionFromList := range list.TransactionList {
		var signedTransaction services.SignedTransaction
		if err := protobuf.Unmarshal(transactionFromList.SignedTransactionBytes, &signedTransaction); err != nil {
			return Transaction{}, errors.Wrap(err, "error deserializing SignedTransactionBytes in TransactionFromBytes")
		}

		tx.signedTransactions = tx.signedTransactions._Push(&signedTransaction)
		if err != nil {
			return Transaction{}, err
		}

		if i == 0 {
			for _, sigPair := range signedTransaction.GetSigMap().GetSigPair() {
				key, err := PublicKeyFromBytes(sigPair.GetPubKeyPrefix())
				if err != nil {
					return Transaction{}, err
				}

				tx.publicKeys = append(tx.publicKeys, key)
				tx.transactionSigners = append(tx.transactionSigners, nil)
			}
		}

		var body services.TransactionBody
		if err := protobuf.Unmarshal(signedTransaction.GetBodyBytes(), &body); err != nil {
			return Transaction{}, errors.Wrap(err, "error deserializing BodyBytes in TransactionFromBytes")
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

		found := false

		for _, value := range tx.transactionIDs.slice {
			id := value.(TransactionID)
			if id.AccountID != nil && transactionID.AccountID != nil &&
				id.AccountID._Equals(*transactionID.AccountID) &&
				id.ValidStart != nil && transactionID.ValidStart != nil &&
				id.ValidStart.Equal(*transactionID.ValidStart) {
				found = true
				break
			}
		}

		if !found {
			tx.transactionIDs = tx.transactionIDs._Push(transactionID)
		}

		for _, id := range tx.GetNodeAccountIDs() {
			if id._Equals(nodeAccountID) {
				found = true
				break
			}
		}

		if !found {
			tx.nodeAccountIDs = tx.nodeAccountIDs._Push(nodeAccountID)
		}

		if i == 0 {
			tx.memo = body.Memo
			if body.TransactionFee != 0 {
				tx.transactionFee = body.TransactionFee
			}
		}
	}

	if tx.transactionIDs._Length() > 0 {
		tx.transactionIDs.locked = true
	}

	if tx.nodeAccountIDs._Length() > 0 {
		tx.nodeAccountIDs.locked = true
	}

	if first == nil {
		return nil, errNoTransactionInBytes
	}

	switch first.Data.(type) {
	case *services.TransactionBody_ContractCall:
		return _ContractExecuteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ContractCreateInstance:
		return _ContractCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ContractUpdateInstance:
		return _ContractUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ContractDeleteInstance:
		return _ContractDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoAddLiveHash:
		return _LiveHashAddTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoCreateAccount:
		return _AccountCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoDelete:
		return _AccountDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoDeleteLiveHash:
		return _LiveHashDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoTransfer:
		return _TransferTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoUpdateAccount:
		return _AccountUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoApproveAllowance:
		return _AccountAllowanceApproveTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoDeleteAllowance:
		return _AccountAllowanceDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_FileAppend:
		return _FileAppendTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_FileCreate:
		return _FileCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_FileDelete:
		return _FileDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_FileUpdate:
		return _FileUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_SystemDelete:
		return _SystemDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_SystemUndelete:
		return _SystemUndeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_Freeze:
		return _FreezeTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ConsensusCreateTopic:
		return _TopicCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ConsensusUpdateTopic:
		return _TopicUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ConsensusDeleteTopic:
		return _TopicDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ConsensusSubmitMessage:
		return _TopicMessageSubmitTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenCreation:
		return _TokenCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenFreeze:
		return _TokenFreezeTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenUnfreeze:
		return _TokenUnfreezeTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenGrantKyc:
		return _TokenGrantKycTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenRevokeKyc:
		return _TokenRevokeKycTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenDeletion:
		return _TokenDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenUpdate:
		return _TokenUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenMint:
		return _TokenMintTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenBurn:
		return _TokenBurnTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenWipe:
		return _TokenWipeTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenAssociate:
		return _TokenAssociateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenDissociate:
		return _TokenDissociateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ScheduleCreate:
		return _ScheduleCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ScheduleSign:
		return _ScheduleSignTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ScheduleDelete:
		return _ScheduleDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenFeeScheduleUpdate:
		return _TokenFeeScheduleUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenPause:
		return _TokenPauseTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenUnpause:
		return _TokenUnpauseTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_EthereumTransaction:
		return _EthereumTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_UtilPrng:
		return _PrngTransactionFromProtobuf(tx, first), nil
	default:
		return Transaction{}, errFailedToDeserializeBytes
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
func (tx *Transaction) GetSignatures() (map[AccountID]map[*PublicKey][]byte, error) {
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

func (tx *Transaction) GetTransactionHash() ([]byte, error) {
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

func (tx *Transaction) GetTransactionHashPerNode() (map[AccountID][]byte, error) {
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
func (tx *Transaction) _InitFee(client *Client) {
	if tx.transactionFee == 0 {
		if client != nil && client.GetDefaultMaxTransactionFee().AsTinybar() != 0 {
			tx.SetMaxTransactionFee(client.GetDefaultMaxTransactionFee())
		} else {
			tx.SetMaxTransactionFee(tx.GetDefaultMaxTransactionFee())
		}
	}
}

func (tx *Transaction) _InitTransactionID(client *Client) error {
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

func (tx *Transaction) IsFrozen() bool {
	return tx.signedTransactions._Length() > 0
}

func (tx *Transaction) _RequireNotFrozen() {
	if tx.IsFrozen() {
		tx.freezeError = errTransactionIsFrozen
	}
}

func (tx *Transaction) _RequireOneNodeAccountID() {
	if tx.nodeAccountIDs._Length() != 1 {
		panic("transaction has more than one _Node ID set")
	}
}

func _TransactionFreezeWith(
	transaction *Transaction,
	client *Client,
	body *services.TransactionBody,
) error {
	if transaction.nodeAccountIDs._IsEmpty() {
		if client != nil {
			for _, nodeAccountID := range client.network._GetNodeAccountIDsForExecute() {
				transaction.nodeAccountIDs._Push(nodeAccountID)
			}
		} else {
			return errNoClientOrTransactionIDOrNodeId
		}
	}

	if client != nil {
		if client.defaultRegenerateTransactionIDs != transaction.regenerateTransactionID {
			transaction.regenerateTransactionID = client.defaultRegenerateTransactionIDs
		}
	}

	for _, nodeAccountID := range transaction.nodeAccountIDs.slice {
		body.NodeAccountID = nodeAccountID.(AccountID)._ToProtobuf()
		bodyBytes, err := protobuf.Marshal(body)
		if err != nil {
			// This should be unreachable
			// From the documentation this appears to only be possible if there are missing proto types
			panic(err)
		}
		transaction.signedTransactions = transaction.signedTransactions._Push(&services.SignedTransaction{
			BodyBytes: bodyBytes,
			SigMap: &services.SignatureMap{
				SigPair: make([]*services.SignaturePair, 0),
			},
		})
	}

	return nil
}

func (tx *Transaction) _SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) {
	tx.transactions = _NewLockableSlice()
	tx.publicKeys = append(tx.publicKeys, publicKey)
	tx.transactionSigners = append(tx.transactionSigners, signer)
}

func (tx *Transaction) _KeyAlreadySigned(
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
func (tx *Transaction) String() string {
	switch sig := tx.signedTransactions._Get(0).(type) { //nolint
	case *services.SignedTransaction:
		return fmt.Sprintf("%+v", sig)
	}

	return ""
}

// ToBytes Builds then converts the current transaction to []byte
// Requires transaction to be frozen
func (tx *Transaction) ToBytes() ([]byte, error) {
	if !tx.IsFrozen() {
		return make([]byte, 0), errTransactionIsNotFrozen
	}

	allTx, err := tx._BuildAllTransactions()
	if err != nil {
		return make([]byte, 0), err
	}
	tx.transactionIDs.locked = true

	pbTransactionList, lastError := protobuf.Marshal(&sdk.TransactionList{
		TransactionList: allTx,
	})

	if lastError != nil {
		return make([]byte, 0), errors.Wrap(err, "error serializing tx list")
	}

	return pbTransactionList, nil
}

func (tx *Transaction) _SignTransaction(index int) {
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

func (tx *Transaction) _BuildAllTransactions() ([]*services.Transaction, error) {
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

func (tx *Transaction) _BuildTransaction(index int) (*services.Transaction, error) {
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
func (tx *Transaction) GetMaxTransactionFee() Hbar {
	return HbarFromTinybar(int64(tx.transactionFee))
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *Transaction) SetMaxTransactionFee(fee Hbar) *Transaction {
	tx.transactionFee = uint64(fee.AsTinybar())
	return tx
}
func (tx *Transaction) GetDefaultMaxTransactionFee() Hbar {
	return HbarFromTinybar(int64(tx.defaultMaxTransactionFee))
}

// SetMaxTransactionFee sets the max Transaction fee for this Transaction.
func (tx *Transaction) _SetDefaultMaxTransactionFee(fee Hbar) {
	tx.defaultMaxTransactionFee = uint64(fee.AsTinybar())
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled
func (tx *Transaction) GetRegenerateTransactionID() bool {
	return tx.regenerateTransactionID
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when \`TRANSACTION_EXPIRED\` is received
func (tx *Transaction) SetRegenerateTransactionID(regenerateTransactionID bool) *Transaction {
	tx.regenerateTransactionID = regenerateTransactionID
	return tx
}

// GetTransactionMemo returns the memo for this	transaction.
func (tx *Transaction) GetTransactionMemo() string {
	return tx.memo
}

// SetTransactionMemo sets the memo for this transaction.
func (tx *Transaction) SetTransactionMemo(memo string) *Transaction {
	tx.memo = memo
	return tx
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (tx *Transaction) GetTransactionValidDuration() time.Duration {
	if tx.transactionValidDuration != nil {
		return *tx.transactionValidDuration
	}

	return 0
}

// SetTransactionValidDuration sets the valid duration for this transaction.
func (tx *Transaction) SetTransactionValidDuration(duration time.Duration) *Transaction {
	tx.transactionValidDuration = &duration
	return tx
}

// GetTransactionID gets the TransactionID for this	transaction.
func (tx *Transaction) GetTransactionID() TransactionID {
	if tx.transactionIDs._Length() > 0 {
		t := tx.transactionIDs._GetCurrent().(TransactionID)
		return t
	}

	return TransactionID{}
}

// SetTransactionID sets the TransactionID for this transaction.
func (tx *Transaction) SetTransactionID(transactionID TransactionID) *Transaction {
	tx.transactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return tx
}

// SetNodeAccountIDs sets the node AccountID for this transaction.
func (tx *Transaction) SetNodeAccountIDs(nodeAccountIDs []AccountID) *Transaction {
	for _, nodeAccountID := range nodeAccountIDs {
		tx.nodeAccountIDs._Push(nodeAccountID)
	}
	tx.nodeAccountIDs._SetLocked(true)
	return tx
}

// ------------ Transaction methdos ---------------
func (tx *Transaction) Sign(privateKey PrivateKey) TransactionInterface {
	return tx.SignWith(privateKey.PublicKey(), privateKey.Sign)
}
func (tx *Transaction) signWithOperator(client *Client, e TransactionInterface) (TransactionInterface, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !tx.IsFrozen() {
		_, err := tx.freezeWith(client, e)
		if err != nil {
			return tx, err
		}
	}
	return tx.SignWith(client.operator.publicKey, client.operator.signer), nil
}
func (tx *Transaction) SignWith(publicKey PublicKey, signer TransactionSigner) TransactionInterface {
	if !tx._KeyAlreadySigned(publicKey) {
		tx._SignWith(publicKey, signer)
	}

	return tx
}
func (tx *Transaction) AddSignature(publicKey PublicKey, signature []byte) TransactionInterface {
	tx._RequireOneNodeAccountID()

	if tx._KeyAlreadySigned(publicKey) {
		return tx
	}

	if tx.signedTransactions._Length() == 0 {
		return tx
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

	return tx
}

// Building empty object as "default" implementation. All inhertents must implement their own implementation.
func (tx *Transaction) build() *services.TransactionBody {
	return &services.TransactionBody{}
}

// Building empty object as "default" implementation. All inhertents must implement their own implementation.
func (tx *Transaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{}, nil
}

// -------------------------------------

func TransactionSign(transaction interface{}, privateKey PrivateKey) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.Sign(privateKey), nil
	case AccountDeleteTransaction:
		return i.Sign(privateKey), nil
	case AccountUpdateTransaction:
		return i.Sign(privateKey), nil
	case AccountAllowanceApproveTransaction:
		return i.Sign(privateKey), nil
	case AccountAllowanceDeleteTransaction:
		return i.Sign(privateKey), nil
	case ContractCreateTransaction:
		return i.Sign(privateKey), nil
	case ContractDeleteTransaction:
		return i.Sign(privateKey), nil
	case ContractExecuteTransaction:
		return i.Sign(privateKey), nil
	case ContractUpdateTransaction:
		return i.Sign(privateKey), nil
	case FileAppendTransaction:
		return i.Sign(privateKey), nil
	case FileCreateTransaction:
		return i.Sign(privateKey), nil
	case FileDeleteTransaction:
		return i.Sign(privateKey), nil
	case FileUpdateTransaction:
		return i.Sign(privateKey), nil
	case LiveHashAddTransaction:
		return i.Sign(privateKey), nil
	case LiveHashDeleteTransaction:
		return i.Sign(privateKey), nil
	case ScheduleCreateTransaction:
		return i.Sign(privateKey), nil
	case ScheduleDeleteTransaction:
		return i.Sign(privateKey), nil
	case ScheduleSignTransaction:
		return i.Sign(privateKey), nil
	case SystemDeleteTransaction:
		return i.Sign(privateKey), nil
	case SystemUndeleteTransaction:
		return i.Sign(privateKey), nil
	case TokenAssociateTransaction:
		return i.Sign(privateKey), nil
	case TokenBurnTransaction:
		return i.Sign(privateKey), nil
	case TokenCreateTransaction:
		return i.Sign(privateKey), nil
	case TokenDeleteTransaction:
		return i.Sign(privateKey), nil
	case TokenDissociateTransaction:
		return i.Sign(privateKey), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.Sign(privateKey), nil
	case TokenFreezeTransaction:
		return i.Sign(privateKey), nil
	case TokenGrantKycTransaction:
		return i.Sign(privateKey), nil
	case TokenMintTransaction:
		return i.Sign(privateKey), nil
	case TokenRevokeKycTransaction:
		return i.Sign(privateKey), nil
	case TokenUnfreezeTransaction:
		return i.Sign(privateKey), nil
	case TokenUpdateTransaction:
		return i.Sign(privateKey), nil
	case TokenWipeTransaction:
		return i.Sign(privateKey), nil
	case TopicCreateTransaction:
		return i.Sign(privateKey), nil
	case TopicDeleteTransaction:
		return i.Sign(privateKey), nil
	case TopicMessageSubmitTransaction:
		return i.Sign(privateKey), nil
	case TopicUpdateTransaction:
		return i.Sign(privateKey), nil
	case TransferTransaction:
		return i.Sign(privateKey), nil
	case *AccountCreateTransaction:
		return i.Sign(privateKey), nil
	case *AccountDeleteTransaction:
		return i.Sign(privateKey), nil
	case *AccountUpdateTransaction:
		return i.Sign(privateKey), nil
	case *AccountAllowanceApproveTransaction:
		return i.Sign(privateKey), nil
	case *AccountAllowanceDeleteTransaction:
		return i.Sign(privateKey), nil
	case *ContractCreateTransaction:
		return i.Sign(privateKey), nil
	case *ContractDeleteTransaction:
		return i.Sign(privateKey), nil
	case *ContractExecuteTransaction:
		return i.Sign(privateKey), nil
	case *ContractUpdateTransaction:
		return i.Sign(privateKey), nil
	case *FileAppendTransaction:
		return i.Sign(privateKey), nil
	case *FileCreateTransaction:
		return i.Sign(privateKey), nil
	case *FileDeleteTransaction:
		return i.Sign(privateKey), nil
	case *FileUpdateTransaction:
		return i.Sign(privateKey), nil
	case *LiveHashAddTransaction:
		return i.Sign(privateKey), nil
	case *LiveHashDeleteTransaction:
		return i.Sign(privateKey), nil
	case *ScheduleCreateTransaction:
		return i.Sign(privateKey), nil
	case *ScheduleDeleteTransaction:
		return i.Sign(privateKey), nil
	case *ScheduleSignTransaction:
		return i.Sign(privateKey), nil
	case *SystemDeleteTransaction:
		return i.Sign(privateKey), nil
	case *SystemUndeleteTransaction:
		return i.Sign(privateKey), nil
	case *TokenAssociateTransaction:
		return i.Sign(privateKey), nil
	case *TokenBurnTransaction:
		return i.Sign(privateKey), nil
	case *TokenCreateTransaction:
		return i.Sign(privateKey), nil
	case *TokenDeleteTransaction:
		return i.Sign(privateKey), nil
	case *TokenDissociateTransaction:
		return i.Sign(privateKey), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.Sign(privateKey), nil
	case *TokenFreezeTransaction:
		return i.Sign(privateKey), nil
	case *TokenGrantKycTransaction:
		return i.Sign(privateKey), nil
	case *TokenMintTransaction:
		return i.Sign(privateKey), nil
	case *TokenRevokeKycTransaction:
		return i.Sign(privateKey), nil
	case *TokenUnfreezeTransaction:
		return i.Sign(privateKey), nil
	case *TokenUpdateTransaction:
		return i.Sign(privateKey), nil
	case *TokenWipeTransaction:
		return i.Sign(privateKey), nil
	case *TopicCreateTransaction:
		return i.Sign(privateKey), nil
	case *TopicDeleteTransaction:
		return i.Sign(privateKey), nil
	case *TopicMessageSubmitTransaction:
		return i.Sign(privateKey), nil
	case *TopicUpdateTransaction:
		return i.Sign(privateKey), nil
	case *TransferTransaction:
		return i.Sign(privateKey), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSignWth(transaction interface{}, publicKKey PublicKey, signer TransactionSigner) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case AccountDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case AccountUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case AccountAllowanceApproveTransaction:
		return i.SignWith(publicKKey, signer), nil
	case AccountAllowanceDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case ContractCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case ContractDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case ContractExecuteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case ContractUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case FileAppendTransaction:
		return i.SignWith(publicKKey, signer), nil
	case FileCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case FileDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case FileUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case LiveHashAddTransaction:
		return i.SignWith(publicKKey, signer), nil
	case LiveHashDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case ScheduleCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case ScheduleDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case ScheduleSignTransaction:
		return i.SignWith(publicKKey, signer), nil
	case SystemDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case SystemUndeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenAssociateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenBurnTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenDissociateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenFreezeTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenGrantKycTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenMintTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenRevokeKycTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenUnfreezeTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TokenWipeTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TopicCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TopicDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TopicMessageSubmitTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TopicUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case TransferTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *AccountCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *AccountDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *AccountUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *AccountAllowanceApproveTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *AccountAllowanceDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *ContractCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *ContractDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *ContractExecuteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *ContractUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *FileAppendTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *FileCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *FileDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *FileUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *LiveHashAddTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *LiveHashDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *ScheduleCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *ScheduleDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *ScheduleSignTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *SystemDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *SystemUndeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenAssociateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenBurnTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenDissociateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenFreezeTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenGrantKycTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenMintTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenRevokeKycTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenUnfreezeTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TokenWipeTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TopicCreateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TopicDeleteTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TopicMessageSubmitTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TopicUpdateTransaction:
		return i.SignWith(publicKKey, signer), nil
	case *TransferTransaction:
		return i.SignWith(publicKKey, signer), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSignWithOperator(transaction interface{}, client *Client) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SignWithOperator(client)
	case AccountDeleteTransaction:
		return i.SignWithOperator(client)
	case AccountUpdateTransaction:
		return i.SignWithOperator(client)
	case AccountAllowanceApproveTransaction:
		return i.SignWithOperator(client)
	case AccountAllowanceDeleteTransaction:
		return i.SignWithOperator(client)
	case ContractCreateTransaction:
		return i.SignWithOperator(client)
	case ContractDeleteTransaction:
		return i.SignWithOperator(client)
	case ContractExecuteTransaction:
		return i.SignWithOperator(client)
	case ContractUpdateTransaction:
		return i.SignWithOperator(client)
	case FileAppendTransaction:
		return i.SignWithOperator(client)
	case FileCreateTransaction:
		return i.SignWithOperator(client)
	case FileDeleteTransaction:
		return i.SignWithOperator(client)
	case FileUpdateTransaction:
		return i.SignWithOperator(client)
	case LiveHashAddTransaction:
		return i.SignWithOperator(client)
	case LiveHashDeleteTransaction:
		return i.SignWithOperator(client)
	case ScheduleCreateTransaction:
		return i.SignWithOperator(client)
	case ScheduleDeleteTransaction:
		return i.SignWithOperator(client)
	case ScheduleSignTransaction:
		return i.SignWithOperator(client)
	case SystemDeleteTransaction:
		return i.SignWithOperator(client)
	case SystemUndeleteTransaction:
		return i.SignWithOperator(client)
	case TokenAssociateTransaction:
		return i.SignWithOperator(client)
	case TokenBurnTransaction:
		return i.SignWithOperator(client)
	case TokenCreateTransaction:
		return i.SignWithOperator(client)
	case TokenDeleteTransaction:
		return i.SignWithOperator(client)
	case TokenDissociateTransaction:
		return i.SignWithOperator(client)
	case TokenFeeScheduleUpdateTransaction:
		return i.SignWithOperator(client)
	case TokenFreezeTransaction:
		return i.SignWithOperator(client)
	case TokenGrantKycTransaction:
		return i.SignWithOperator(client)
	case TokenMintTransaction:
		return i.SignWithOperator(client)
	case TokenRevokeKycTransaction:
		return i.SignWithOperator(client)
	case TokenUnfreezeTransaction:
		return i.SignWithOperator(client)
	case TokenUpdateTransaction:
		return i.SignWithOperator(client)
	case TokenWipeTransaction:
		return i.SignWithOperator(client)
	case TopicCreateTransaction:
		return i.SignWithOperator(client)
	case TopicDeleteTransaction:
		return i.SignWithOperator(client)
	case TopicMessageSubmitTransaction:
		return i.SignWithOperator(client)
	case TopicUpdateTransaction:
		return i.SignWithOperator(client)
	case TransferTransaction:
		return i.SignWithOperator(client)
	case *AccountCreateTransaction:
		return i.SignWithOperator(client)
	case *AccountDeleteTransaction:
		return i.SignWithOperator(client)
	case *AccountUpdateTransaction:
		return i.SignWithOperator(client)
	case *AccountAllowanceApproveTransaction:
		return i.SignWithOperator(client)
	case *AccountAllowanceDeleteTransaction:
		return i.SignWithOperator(client)
	case *ContractCreateTransaction:
		return i.SignWithOperator(client)
	case *ContractDeleteTransaction:
		return i.SignWithOperator(client)
	case *ContractExecuteTransaction:
		return i.SignWithOperator(client)
	case *ContractUpdateTransaction:
		return i.SignWithOperator(client)
	case *FileAppendTransaction:
		return i.SignWithOperator(client)
	case *FileCreateTransaction:
		return i.SignWithOperator(client)
	case *FileDeleteTransaction:
		return i.SignWithOperator(client)
	case *FileUpdateTransaction:
		return i.SignWithOperator(client)
	case *LiveHashAddTransaction:
		return i.SignWithOperator(client)
	case *LiveHashDeleteTransaction:
		return i.SignWithOperator(client)
	case *ScheduleCreateTransaction:
		return i.SignWithOperator(client)
	case *ScheduleDeleteTransaction:
		return i.SignWithOperator(client)
	case *ScheduleSignTransaction:
		return i.SignWithOperator(client)
	case *SystemDeleteTransaction:
		return i.SignWithOperator(client)
	case *SystemUndeleteTransaction:
		return i.SignWithOperator(client)
	case *TokenAssociateTransaction:
		return i.SignWithOperator(client)
	case *TokenBurnTransaction:
		return i.SignWithOperator(client)
	case *TokenCreateTransaction:
		return i.SignWithOperator(client)
	case *TokenDeleteTransaction:
		return i.SignWithOperator(client)
	case *TokenDissociateTransaction:
		return i.SignWithOperator(client)
	case *TokenFeeScheduleUpdateTransaction:
		return i.SignWithOperator(client)
	case *TokenFreezeTransaction:
		return i.SignWithOperator(client)
	case *TokenGrantKycTransaction:
		return i.SignWithOperator(client)
	case *TokenMintTransaction:
		return i.SignWithOperator(client)
	case *TokenRevokeKycTransaction:
		return i.SignWithOperator(client)
	case *TokenUnfreezeTransaction:
		return i.SignWithOperator(client)
	case *TokenUpdateTransaction:
		return i.SignWithOperator(client)
	case *TokenWipeTransaction:
		return i.SignWithOperator(client)
	case *TopicCreateTransaction:
		return i.SignWithOperator(client)
	case *TopicDeleteTransaction:
		return i.SignWithOperator(client)
	case *TopicMessageSubmitTransaction:
		return i.SignWithOperator(client)
	case *TopicUpdateTransaction:
		return i.SignWithOperator(client)
	case *TransferTransaction:
		return i.SignWithOperator(client)
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionAddSignature(transaction interface{}, publicKey PublicKey, signature []byte) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case AccountDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case AccountUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case AccountAllowanceApproveTransaction:
		return i.AddSignature(publicKey, signature), nil
	case AccountAllowanceDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case ContractCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case ContractDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case ContractExecuteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case ContractUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case FileAppendTransaction:
		return i.AddSignature(publicKey, signature), nil
	case FileCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case FileDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case FileUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case LiveHashAddTransaction:
		return i.AddSignature(publicKey, signature), nil
	case LiveHashDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case SystemDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case SystemUndeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenAssociateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenBurnTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenDissociateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenFreezeTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenGrantKycTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenMintTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenRevokeKycTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenUnfreezeTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TokenWipeTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TopicCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TopicDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TopicMessageSubmitTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TopicUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case TransferTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *AccountCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *AccountDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *AccountUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *AccountAllowanceApproveTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *AccountAllowanceDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *ContractCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *ContractDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *ContractExecuteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *ContractUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *FileAppendTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *FileCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *FileDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *FileUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *LiveHashAddTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *LiveHashDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *SystemDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *SystemUndeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenAssociateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenBurnTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenDissociateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenFreezeTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenGrantKycTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenMintTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenRevokeKycTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenUnfreezeTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TokenWipeTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TopicCreateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TopicDeleteTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TopicMessageSubmitTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TopicUpdateTransaction:
		return i.AddSignature(publicKey, signature), nil
	case *TransferTransaction:
		return i.AddSignature(publicKey, signature), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetSignatures(transaction interface{}) (map[AccountID]map[*PublicKey][]byte, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetSignatures()
	case AccountDeleteTransaction:
		return i.GetSignatures()
	case AccountUpdateTransaction:
		return i.GetSignatures()
	case AccountAllowanceApproveTransaction:
		return i.GetSignatures()
	case AccountAllowanceDeleteTransaction:
		return i.GetSignatures()
	case ContractCreateTransaction:
		return i.GetSignatures()
	case ContractDeleteTransaction:
		return i.GetSignatures()
	case ContractExecuteTransaction:
		return i.GetSignatures()
	case ContractUpdateTransaction:
		return i.GetSignatures()
	case FileAppendTransaction:
		return i.GetSignatures()
	case FileCreateTransaction:
		return i.GetSignatures()
	case FileDeleteTransaction:
		return i.GetSignatures()
	case FileUpdateTransaction:
		return i.GetSignatures()
	case LiveHashAddTransaction:
		return i.GetSignatures()
	case LiveHashDeleteTransaction:
		return i.GetSignatures()
	case ScheduleCreateTransaction:
		return i.GetSignatures()
	case ScheduleDeleteTransaction:
		return i.GetSignatures()
	case ScheduleSignTransaction:
		return i.GetSignatures()
	case SystemDeleteTransaction:
		return i.GetSignatures()
	case SystemUndeleteTransaction:
		return i.GetSignatures()
	case TokenAssociateTransaction:
		return i.GetSignatures()
	case TokenBurnTransaction:
		return i.GetSignatures()
	case TokenCreateTransaction:
		return i.GetSignatures()
	case TokenDeleteTransaction:
		return i.GetSignatures()
	case TokenDissociateTransaction:
		return i.GetSignatures()
	case TokenFeeScheduleUpdateTransaction:
		return i.GetSignatures()
	case TokenFreezeTransaction:
		return i.GetSignatures()
	case TokenGrantKycTransaction:
		return i.GetSignatures()
	case TokenMintTransaction:
		return i.GetSignatures()
	case TokenRevokeKycTransaction:
		return i.GetSignatures()
	case TokenUnfreezeTransaction:
		return i.GetSignatures()
	case TokenUpdateTransaction:
		return i.GetSignatures()
	case TokenWipeTransaction:
		return i.GetSignatures()
	case TopicCreateTransaction:
		return i.GetSignatures()
	case TopicDeleteTransaction:
		return i.GetSignatures()
	case TopicMessageSubmitTransaction:
		return i.GetSignatures()
	case TopicUpdateTransaction:
		return i.GetSignatures()
	case TransferTransaction:
		return i.GetSignatures()
	case *AccountCreateTransaction:
		return i.GetSignatures()
	case *AccountDeleteTransaction:
		return i.GetSignatures()
	case *AccountUpdateTransaction:
		return i.GetSignatures()
	case *AccountAllowanceApproveTransaction:
		return i.GetSignatures()
	case *AccountAllowanceDeleteTransaction:
		return i.GetSignatures()
	case *ContractCreateTransaction:
		return i.GetSignatures()
	case *ContractDeleteTransaction:
		return i.GetSignatures()
	case *ContractExecuteTransaction:
		return i.GetSignatures()
	case *ContractUpdateTransaction:
		return i.GetSignatures()
	case *FileAppendTransaction:
		return i.GetSignatures()
	case *FileCreateTransaction:
		return i.GetSignatures()
	case *FileDeleteTransaction:
		return i.GetSignatures()
	case *FileUpdateTransaction:
		return i.GetSignatures()
	case *LiveHashAddTransaction:
		return i.GetSignatures()
	case *LiveHashDeleteTransaction:
		return i.GetSignatures()
	case *ScheduleCreateTransaction:
		return i.GetSignatures()
	case *ScheduleDeleteTransaction:
		return i.GetSignatures()
	case *ScheduleSignTransaction:
		return i.GetSignatures()
	case *SystemDeleteTransaction:
		return i.GetSignatures()
	case *SystemUndeleteTransaction:
		return i.GetSignatures()
	case *TokenAssociateTransaction:
		return i.GetSignatures()
	case *TokenBurnTransaction:
		return i.GetSignatures()
	case *TokenCreateTransaction:
		return i.GetSignatures()
	case *TokenDeleteTransaction:
		return i.GetSignatures()
	case *TokenDissociateTransaction:
		return i.GetSignatures()
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetSignatures()
	case *TokenFreezeTransaction:
		return i.GetSignatures()
	case *TokenGrantKycTransaction:
		return i.GetSignatures()
	case *TokenMintTransaction:
		return i.GetSignatures()
	case *TokenRevokeKycTransaction:
		return i.GetSignatures()
	case *TokenUnfreezeTransaction:
		return i.GetSignatures()
	case *TokenUpdateTransaction:
		return i.GetSignatures()
	case *TokenWipeTransaction:
		return i.GetSignatures()
	case *TopicCreateTransaction:
		return i.GetSignatures()
	case *TopicDeleteTransaction:
		return i.GetSignatures()
	case *TopicMessageSubmitTransaction:
		return i.GetSignatures()
	case *TopicUpdateTransaction:
		return i.GetSignatures()
	case *TransferTransaction:
		return i.GetSignatures()
	default:
		return nil, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSetTransactionID(transaction interface{}, transactionID TransactionID) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case AccountDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case AccountUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case AccountAllowanceApproveTransaction:
		return i.SetTransactionID(transactionID), nil
	case AccountAllowanceDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case ContractCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case ContractDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case ContractExecuteTransaction:
		return i.SetTransactionID(transactionID), nil
	case ContractUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case FileAppendTransaction:
		return i.SetTransactionID(transactionID), nil
	case FileCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case FileDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case FileUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case LiveHashAddTransaction:
		return i.SetTransactionID(transactionID), nil
	case LiveHashDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case ScheduleCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case ScheduleDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case ScheduleSignTransaction:
		return i.SetTransactionID(transactionID), nil
	case SystemDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case SystemUndeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenAssociateTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenBurnTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenDissociateTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenFreezeTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenGrantKycTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenMintTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenRevokeKycTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenUnfreezeTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case TokenWipeTransaction:
		return i.SetTransactionID(transactionID), nil
	case TopicCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case TopicDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case TopicMessageSubmitTransaction:
		return i.SetTransactionID(transactionID), nil
	case TopicUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case TransferTransaction:
		return i.SetTransactionID(transactionID), nil
	case *AccountCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *AccountDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *AccountUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *AccountAllowanceApproveTransaction:
		return i.SetTransactionID(transactionID), nil
	case *AccountAllowanceDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *ContractCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *ContractDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *ContractExecuteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *ContractUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *FileAppendTransaction:
		return i.SetTransactionID(transactionID), nil
	case *FileCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *FileDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *FileUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *LiveHashAddTransaction:
		return i.SetTransactionID(transactionID), nil
	case *LiveHashDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *ScheduleCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *ScheduleDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *ScheduleSignTransaction:
		return i.SetTransactionID(transactionID), nil
	case *SystemDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *SystemUndeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenAssociateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenBurnTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenDissociateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenFreezeTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenGrantKycTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenMintTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenRevokeKycTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenUnfreezeTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TokenWipeTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TopicCreateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TopicDeleteTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TopicMessageSubmitTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TopicUpdateTransaction:
		return i.SetTransactionID(transactionID), nil
	case *TransferTransaction:
		return i.SetTransactionID(transactionID), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetTransactionID(transaction interface{}) (TransactionID, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetTransactionID(), nil
	case AccountDeleteTransaction:
		return i.GetTransactionID(), nil
	case AccountUpdateTransaction:
		return i.GetTransactionID(), nil
	case AccountAllowanceApproveTransaction:
		return i.GetTransactionID(), nil
	case AccountAllowanceDeleteTransaction:
		return i.GetTransactionID(), nil
	case ContractCreateTransaction:
		return i.GetTransactionID(), nil
	case ContractDeleteTransaction:
		return i.GetTransactionID(), nil
	case ContractExecuteTransaction:
		return i.GetTransactionID(), nil
	case ContractUpdateTransaction:
		return i.GetTransactionID(), nil
	case FileAppendTransaction:
		return i.GetTransactionID(), nil
	case FileCreateTransaction:
		return i.GetTransactionID(), nil
	case FileDeleteTransaction:
		return i.GetTransactionID(), nil
	case FileUpdateTransaction:
		return i.GetTransactionID(), nil
	case LiveHashAddTransaction:
		return i.GetTransactionID(), nil
	case LiveHashDeleteTransaction:
		return i.GetTransactionID(), nil
	case ScheduleCreateTransaction:
		return i.GetTransactionID(), nil
	case ScheduleDeleteTransaction:
		return i.GetTransactionID(), nil
	case ScheduleSignTransaction:
		return i.GetTransactionID(), nil
	case SystemDeleteTransaction:
		return i.GetTransactionID(), nil
	case SystemUndeleteTransaction:
		return i.GetTransactionID(), nil
	case TokenAssociateTransaction:
		return i.GetTransactionID(), nil
	case TokenBurnTransaction:
		return i.GetTransactionID(), nil
	case TokenCreateTransaction:
		return i.GetTransactionID(), nil
	case TokenDeleteTransaction:
		return i.GetTransactionID(), nil
	case TokenDissociateTransaction:
		return i.GetTransactionID(), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionID(), nil
	case TokenFreezeTransaction:
		return i.GetTransactionID(), nil
	case TokenGrantKycTransaction:
		return i.GetTransactionID(), nil
	case TokenMintTransaction:
		return i.GetTransactionID(), nil
	case TokenRevokeKycTransaction:
		return i.GetTransactionID(), nil
	case TokenUnfreezeTransaction:
		return i.GetTransactionID(), nil
	case TokenUpdateTransaction:
		return i.GetTransactionID(), nil
	case TokenWipeTransaction:
		return i.GetTransactionID(), nil
	case TopicCreateTransaction:
		return i.GetTransactionID(), nil
	case TopicDeleteTransaction:
		return i.GetTransactionID(), nil
	case TopicMessageSubmitTransaction:
		return i.GetTransactionID(), nil
	case TopicUpdateTransaction:
		return i.GetTransactionID(), nil
	case TransferTransaction:
		return i.GetTransactionID(), nil
	case *AccountCreateTransaction:
		return i.GetTransactionID(), nil
	case *AccountDeleteTransaction:
		return i.GetTransactionID(), nil
	case *AccountUpdateTransaction:
		return i.GetTransactionID(), nil
	case *AccountAllowanceApproveTransaction:
		return i.GetTransactionID(), nil
	case *AccountAllowanceDeleteTransaction:
		return i.GetTransactionID(), nil
	case *ContractCreateTransaction:
		return i.GetTransactionID(), nil
	case *ContractDeleteTransaction:
		return i.GetTransactionID(), nil
	case *ContractExecuteTransaction:
		return i.GetTransactionID(), nil
	case *ContractUpdateTransaction:
		return i.GetTransactionID(), nil
	case *FileAppendTransaction:
		return i.GetTransactionID(), nil
	case *FileCreateTransaction:
		return i.GetTransactionID(), nil
	case *FileDeleteTransaction:
		return i.GetTransactionID(), nil
	case *FileUpdateTransaction:
		return i.GetTransactionID(), nil
	case *LiveHashAddTransaction:
		return i.GetTransactionID(), nil
	case *LiveHashDeleteTransaction:
		return i.GetTransactionID(), nil
	case *ScheduleCreateTransaction:
		return i.GetTransactionID(), nil
	case *ScheduleDeleteTransaction:
		return i.GetTransactionID(), nil
	case *ScheduleSignTransaction:
		return i.GetTransactionID(), nil
	case *SystemDeleteTransaction:
		return i.GetTransactionID(), nil
	case *SystemUndeleteTransaction:
		return i.GetTransactionID(), nil
	case *TokenAssociateTransaction:
		return i.GetTransactionID(), nil
	case *TokenBurnTransaction:
		return i.GetTransactionID(), nil
	case *TokenCreateTransaction:
		return i.GetTransactionID(), nil
	case *TokenDeleteTransaction:
		return i.GetTransactionID(), nil
	case *TokenDissociateTransaction:
		return i.GetTransactionID(), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionID(), nil
	case *TokenFreezeTransaction:
		return i.GetTransactionID(), nil
	case *TokenGrantKycTransaction:
		return i.GetTransactionID(), nil
	case *TokenMintTransaction:
		return i.GetTransactionID(), nil
	case *TokenRevokeKycTransaction:
		return i.GetTransactionID(), nil
	case *TokenUnfreezeTransaction:
		return i.GetTransactionID(), nil
	case *TokenUpdateTransaction:
		return i.GetTransactionID(), nil
	case *TokenWipeTransaction:
		return i.GetTransactionID(), nil
	case *TopicCreateTransaction:
		return i.GetTransactionID(), nil
	case *TopicDeleteTransaction:
		return i.GetTransactionID(), nil
	case *TopicMessageSubmitTransaction:
		return i.GetTransactionID(), nil
	case *TopicUpdateTransaction:
		return i.GetTransactionID(), nil
	case *TransferTransaction:
		return i.GetTransactionID(), nil
	default:
		return TransactionID{}, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSetTransactionMemo(transaction interface{}, transactionMemo string) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case AccountDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case AccountUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case AccountAllowanceApproveTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case AccountAllowanceDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case ContractCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case ContractDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case ContractExecuteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case ContractUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case FileAppendTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case FileCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case FileDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case FileUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case LiveHashAddTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case LiveHashDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case ScheduleCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case ScheduleDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case ScheduleSignTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case SystemDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case SystemUndeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenAssociateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenBurnTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenDissociateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenFreezeTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenGrantKycTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenMintTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenRevokeKycTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenUnfreezeTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TokenWipeTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TopicCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TopicDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TopicMessageSubmitTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TopicUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case TransferTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *AccountCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *AccountDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *AccountUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *AccountAllowanceApproveTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *AccountAllowanceDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *ContractCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *ContractDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *ContractExecuteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *ContractUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *FileAppendTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *FileCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *FileDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *FileUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *LiveHashAddTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *LiveHashDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *ScheduleCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *ScheduleDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *ScheduleSignTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *SystemDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *SystemUndeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenAssociateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenBurnTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenDissociateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenFreezeTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenGrantKycTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenMintTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenRevokeKycTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenUnfreezeTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TokenWipeTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TopicCreateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TopicDeleteTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TopicMessageSubmitTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TopicUpdateTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	case *TransferTransaction:
		return i.SetTransactionMemo(transactionMemo), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetTransactionMemo(transaction interface{}) (string, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetTransactionMemo(), nil
	case AccountDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case AccountUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case AccountAllowanceApproveTransaction:
		return i.GetTransactionMemo(), nil
	case AccountAllowanceDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case ContractCreateTransaction:
		return i.GetTransactionMemo(), nil
	case ContractDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case ContractExecuteTransaction:
		return i.GetTransactionMemo(), nil
	case ContractUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case FileAppendTransaction:
		return i.GetTransactionMemo(), nil
	case FileCreateTransaction:
		return i.GetTransactionMemo(), nil
	case FileDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case FileUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case LiveHashAddTransaction:
		return i.GetTransactionMemo(), nil
	case LiveHashDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case ScheduleCreateTransaction:
		return i.GetTransactionMemo(), nil
	case ScheduleDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case ScheduleSignTransaction:
		return i.GetTransactionMemo(), nil
	case SystemDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case SystemUndeleteTransaction:
		return i.GetTransactionMemo(), nil
	case TokenAssociateTransaction:
		return i.GetTransactionMemo(), nil
	case TokenBurnTransaction:
		return i.GetTransactionMemo(), nil
	case TokenCreateTransaction:
		return i.GetTransactionMemo(), nil
	case TokenDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case TokenDissociateTransaction:
		return i.GetTransactionMemo(), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case TokenFreezeTransaction:
		return i.GetTransactionMemo(), nil
	case TokenGrantKycTransaction:
		return i.GetTransactionMemo(), nil
	case TokenMintTransaction:
		return i.GetTransactionMemo(), nil
	case TokenRevokeKycTransaction:
		return i.GetTransactionMemo(), nil
	case TokenUnfreezeTransaction:
		return i.GetTransactionMemo(), nil
	case TokenUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case TokenWipeTransaction:
		return i.GetTransactionMemo(), nil
	case TopicCreateTransaction:
		return i.GetTransactionMemo(), nil
	case TopicDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case TopicMessageSubmitTransaction:
		return i.GetTransactionMemo(), nil
	case TopicUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case TransferTransaction:
		return i.GetTransactionMemo(), nil
	case *AccountCreateTransaction:
		return i.GetTransactionMemo(), nil
	case *AccountDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *AccountUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case *AccountAllowanceApproveTransaction:
		return i.GetTransactionMemo(), nil
	case *AccountAllowanceDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *ContractCreateTransaction:
		return i.GetTransactionMemo(), nil
	case *ContractDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *ContractExecuteTransaction:
		return i.GetTransactionMemo(), nil
	case *ContractUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case *FileAppendTransaction:
		return i.GetTransactionMemo(), nil
	case *FileCreateTransaction:
		return i.GetTransactionMemo(), nil
	case *FileDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *FileUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case *LiveHashAddTransaction:
		return i.GetTransactionMemo(), nil
	case *LiveHashDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *ScheduleCreateTransaction:
		return i.GetTransactionMemo(), nil
	case *ScheduleDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *ScheduleSignTransaction:
		return i.GetTransactionMemo(), nil
	case *SystemDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *SystemUndeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenAssociateTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenBurnTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenCreateTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenDissociateTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenFreezeTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenGrantKycTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenMintTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenRevokeKycTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenUnfreezeTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case *TokenWipeTransaction:
		return i.GetTransactionMemo(), nil
	case *TopicCreateTransaction:
		return i.GetTransactionMemo(), nil
	case *TopicDeleteTransaction:
		return i.GetTransactionMemo(), nil
	case *TopicMessageSubmitTransaction:
		return i.GetTransactionMemo(), nil
	case *TopicUpdateTransaction:
		return i.GetTransactionMemo(), nil
	case *TransferTransaction:
		return i.GetTransactionMemo(), nil
	default:
		return "", errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSetMaxTransactionFee(transaction interface{}, maxTransactionFee Hbar) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case AccountDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case AccountUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case AccountAllowanceApproveTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case AccountAllowanceDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case ContractCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case ContractDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case ContractExecuteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case ContractUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case FileAppendTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case FileCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case FileDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case FileUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case LiveHashAddTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case LiveHashDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case ScheduleCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case ScheduleDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case ScheduleSignTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case SystemDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case SystemUndeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenAssociateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenBurnTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenDissociateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenFreezeTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenGrantKycTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenMintTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenRevokeKycTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenUnfreezeTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TokenWipeTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TopicCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TopicDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TopicMessageSubmitTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TopicUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case TransferTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *AccountCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *AccountDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *AccountUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *AccountAllowanceApproveTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *AccountAllowanceDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *ContractCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *ContractDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *ContractExecuteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *ContractUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *FileAppendTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *FileCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *FileDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *FileUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *LiveHashAddTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *LiveHashDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *ScheduleCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *ScheduleDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *ScheduleSignTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *SystemDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *SystemUndeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenAssociateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenBurnTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenDissociateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenFreezeTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenGrantKycTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenMintTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenRevokeKycTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenUnfreezeTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TokenWipeTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TopicCreateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TopicDeleteTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TopicMessageSubmitTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TopicUpdateTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	case *TransferTransaction:
		return i.SetMaxTransactionFee(maxTransactionFee), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetMaxTransactionFee(transaction interface{}) (Hbar, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case AccountDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case AccountUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case AccountAllowanceApproveTransaction:
		return i.GetMaxTransactionFee(), nil
	case AccountAllowanceDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case ContractCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case ContractDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case ContractExecuteTransaction:
		return i.GetMaxTransactionFee(), nil
	case ContractUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case FileAppendTransaction:
		return i.GetMaxTransactionFee(), nil
	case FileCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case FileDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case FileUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case LiveHashAddTransaction:
		return i.GetMaxTransactionFee(), nil
	case LiveHashDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case ScheduleCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case ScheduleDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case ScheduleSignTransaction:
		return i.GetMaxTransactionFee(), nil
	case SystemDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case SystemUndeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenAssociateTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenBurnTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenDissociateTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenFreezeTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenGrantKycTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenMintTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenRevokeKycTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenUnfreezeTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case TokenWipeTransaction:
		return i.GetMaxTransactionFee(), nil
	case TopicCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case TopicDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case TopicMessageSubmitTransaction:
		return i.GetMaxTransactionFee(), nil
	case TopicUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case TransferTransaction:
		return i.GetMaxTransactionFee(), nil
	case *AccountCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *AccountDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *AccountUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *AccountAllowanceApproveTransaction:
		return i.GetMaxTransactionFee(), nil
	case *AccountAllowanceDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *ContractCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *ContractDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *ContractExecuteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *ContractUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *FileAppendTransaction:
		return i.GetMaxTransactionFee(), nil
	case *FileCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *FileDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *FileUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *LiveHashAddTransaction:
		return i.GetMaxTransactionFee(), nil
	case *LiveHashDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *ScheduleCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *ScheduleDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *ScheduleSignTransaction:
		return i.GetMaxTransactionFee(), nil
	case *SystemDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *SystemUndeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenAssociateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenBurnTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenDissociateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenFreezeTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenGrantKycTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenMintTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenRevokeKycTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenUnfreezeTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TokenWipeTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TopicCreateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TopicDeleteTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TopicMessageSubmitTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TopicUpdateTransaction:
		return i.GetMaxTransactionFee(), nil
	case *TransferTransaction:
		return i.GetMaxTransactionFee(), nil
	default:
		return Hbar{}, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSetTransactionValidDuration(transaction interface{}, transactionValidDuration time.Duration) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case AccountDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case AccountUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case AccountAllowanceApproveTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case AccountAllowanceDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case ContractCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case ContractDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case ContractExecuteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case ContractUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case FileAppendTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case FileCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case FileDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case FileUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case LiveHashAddTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case LiveHashDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case ScheduleCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case ScheduleDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case ScheduleSignTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case SystemDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case SystemUndeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenAssociateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenBurnTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenDissociateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenFreezeTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenGrantKycTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenMintTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenRevokeKycTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenUnfreezeTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TokenWipeTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TopicCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TopicDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TopicMessageSubmitTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TopicUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case TransferTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *AccountCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *AccountDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *AccountUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *AccountAllowanceApproveTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *AccountAllowanceDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *ContractCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *ContractDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *ContractExecuteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *ContractUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *FileAppendTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *FileCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *FileDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *FileUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *LiveHashAddTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *LiveHashDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *ScheduleCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *ScheduleDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *ScheduleSignTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *SystemDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *SystemUndeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenAssociateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenBurnTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenDissociateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenFreezeTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenGrantKycTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenMintTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenRevokeKycTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenUnfreezeTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TokenWipeTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TopicCreateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TopicDeleteTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TopicMessageSubmitTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TopicUpdateTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	case *TransferTransaction:
		return i.SetTransactionValidDuration(transactionValidDuration), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetTransactionValidDuration(transaction interface{}) (time.Duration, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case AccountDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case AccountUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case AccountAllowanceApproveTransaction:
		return i.GetTransactionValidDuration(), nil
	case AccountAllowanceDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case ContractCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case ContractDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case ContractExecuteTransaction:
		return i.GetTransactionValidDuration(), nil
	case ContractUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case FileAppendTransaction:
		return i.GetTransactionValidDuration(), nil
	case FileCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case FileDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case FileUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case LiveHashAddTransaction:
		return i.GetTransactionValidDuration(), nil
	case LiveHashDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case ScheduleCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case ScheduleDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case ScheduleSignTransaction:
		return i.GetTransactionValidDuration(), nil
	case SystemDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case SystemUndeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenAssociateTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenBurnTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenDissociateTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenFreezeTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenGrantKycTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenMintTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenRevokeKycTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenUnfreezeTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case TokenWipeTransaction:
		return i.GetTransactionValidDuration(), nil
	case TopicCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case TopicDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case TopicMessageSubmitTransaction:
		return i.GetTransactionValidDuration(), nil
	case TopicUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case TransferTransaction:
		return i.GetTransactionValidDuration(), nil
	case *AccountCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *AccountDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *AccountUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *AccountAllowanceApproveTransaction:
		return i.GetTransactionValidDuration(), nil
	case *AccountAllowanceDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *ContractCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *ContractDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *ContractExecuteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *ContractUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *FileAppendTransaction:
		return i.GetTransactionValidDuration(), nil
	case *FileCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *FileDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *FileUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *LiveHashAddTransaction:
		return i.GetTransactionValidDuration(), nil
	case *LiveHashDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *ScheduleCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *ScheduleDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *ScheduleSignTransaction:
		return i.GetTransactionValidDuration(), nil
	case *SystemDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *SystemUndeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenAssociateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenBurnTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenDissociateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenFreezeTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenGrantKycTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenMintTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenRevokeKycTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenUnfreezeTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TokenWipeTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TopicCreateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TopicDeleteTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TopicMessageSubmitTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TopicUpdateTransaction:
		return i.GetTransactionValidDuration(), nil
	case *TransferTransaction:
		return i.GetTransactionValidDuration(), nil
	default:
		return time.Duration(0), errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSetNodeAccountIDs(transaction interface{}, nodeAccountIDs []AccountID) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case AccountDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case AccountUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case AccountAllowanceApproveTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case AccountAllowanceDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case ContractCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case ContractDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case ContractExecuteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case ContractUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case FileAppendTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case FileCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case FileDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case FileUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case LiveHashAddTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case LiveHashDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case ScheduleCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case ScheduleDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case ScheduleSignTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case SystemDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case SystemUndeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenAssociateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenBurnTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenDissociateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenFreezeTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenGrantKycTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenMintTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenRevokeKycTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenUnfreezeTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TokenWipeTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TopicCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TopicDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TopicMessageSubmitTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TopicUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case TransferTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *AccountCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *AccountDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *AccountUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *AccountAllowanceApproveTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *AccountAllowanceDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *ContractCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *ContractDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *ContractExecuteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *ContractUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *FileAppendTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *FileCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *FileDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *FileUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *LiveHashAddTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *LiveHashDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *ScheduleCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *ScheduleDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *ScheduleSignTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *SystemDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *SystemUndeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenAssociateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenBurnTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenDissociateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenFreezeTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenGrantKycTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenMintTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenRevokeKycTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenUnfreezeTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TokenWipeTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TopicCreateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TopicDeleteTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TopicMessageSubmitTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TopicUpdateTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	case *TransferTransaction:
		return i.SetNodeAccountIDs(nodeAccountIDs), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetNodeAccountIDs(transaction interface{}) ([]AccountID, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case AccountDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case AccountUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case AccountAllowanceApproveTransaction:
		return i.GetNodeAccountIDs(), nil
	case AccountAllowanceDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case ContractCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case ContractDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case ContractExecuteTransaction:
		return i.GetNodeAccountIDs(), nil
	case ContractUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case FileAppendTransaction:
		return i.GetNodeAccountIDs(), nil
	case FileCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case FileDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case FileUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case LiveHashAddTransaction:
		return i.GetNodeAccountIDs(), nil
	case LiveHashDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case ScheduleCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case ScheduleDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case ScheduleSignTransaction:
		return i.GetNodeAccountIDs(), nil
	case SystemDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case SystemUndeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenAssociateTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenBurnTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenDissociateTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenFreezeTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenGrantKycTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenMintTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenRevokeKycTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenUnfreezeTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case TokenWipeTransaction:
		return i.GetNodeAccountIDs(), nil
	case TopicCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case TopicDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case TopicMessageSubmitTransaction:
		return i.GetNodeAccountIDs(), nil
	case TopicUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case TransferTransaction:
		return i.GetNodeAccountIDs(), nil
	case *AccountCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *AccountDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *AccountUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *AccountAllowanceApproveTransaction:
		return i.GetNodeAccountIDs(), nil
	case *AccountAllowanceDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *ContractCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *ContractDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *ContractExecuteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *ContractUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *FileAppendTransaction:
		return i.GetNodeAccountIDs(), nil
	case *FileCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *FileDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *FileUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *LiveHashAddTransaction:
		return i.GetNodeAccountIDs(), nil
	case *LiveHashDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *ScheduleCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *ScheduleDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *ScheduleSignTransaction:
		return i.GetNodeAccountIDs(), nil
	case *SystemDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *SystemUndeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenAssociateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenBurnTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenDissociateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenFreezeTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenGrantKycTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenMintTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenRevokeKycTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenUnfreezeTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TokenWipeTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TopicCreateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TopicDeleteTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TopicMessageSubmitTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TopicUpdateTransaction:
		return i.GetNodeAccountIDs(), nil
	case *TransferTransaction:
		return i.GetNodeAccountIDs(), nil
	default:
		return []AccountID{}, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetTransactionHash(transaction interface{}) ([]byte, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetTransactionHash()
	case AccountDeleteTransaction:
		return i.GetTransactionHash()
	case AccountUpdateTransaction:
		return i.GetTransactionHash()
	case AccountAllowanceApproveTransaction:
		return i.GetTransactionHash()
	case AccountAllowanceDeleteTransaction:
		return i.GetTransactionHash()
	case ContractCreateTransaction:
		return i.GetTransactionHash()
	case ContractDeleteTransaction:
		return i.GetTransactionHash()
	case ContractExecuteTransaction:
		return i.GetTransactionHash()
	case ContractUpdateTransaction:
		return i.GetTransactionHash()
	case FileAppendTransaction:
		return i.GetTransactionHash()
	case FileCreateTransaction:
		return i.GetTransactionHash()
	case FileDeleteTransaction:
		return i.GetTransactionHash()
	case FileUpdateTransaction:
		return i.GetTransactionHash()
	case LiveHashAddTransaction:
		return i.GetTransactionHash()
	case LiveHashDeleteTransaction:
		return i.GetTransactionHash()
	case ScheduleCreateTransaction:
		return i.GetTransactionHash()
	case ScheduleDeleteTransaction:
		return i.GetTransactionHash()
	case ScheduleSignTransaction:
		return i.GetTransactionHash()
	case SystemDeleteTransaction:
		return i.GetTransactionHash()
	case SystemUndeleteTransaction:
		return i.GetTransactionHash()
	case TokenAssociateTransaction:
		return i.GetTransactionHash()
	case TokenBurnTransaction:
		return i.GetTransactionHash()
	case TokenCreateTransaction:
		return i.GetTransactionHash()
	case TokenDeleteTransaction:
		return i.GetTransactionHash()
	case TokenDissociateTransaction:
		return i.GetTransactionHash()
	case TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionHash()
	case TokenFreezeTransaction:
		return i.GetTransactionHash()
	case TokenGrantKycTransaction:
		return i.GetTransactionHash()
	case TokenMintTransaction:
		return i.GetTransactionHash()
	case TokenRevokeKycTransaction:
		return i.GetTransactionHash()
	case TokenUnfreezeTransaction:
		return i.GetTransactionHash()
	case TokenUpdateTransaction:
		return i.GetTransactionHash()
	case TokenWipeTransaction:
		return i.GetTransactionHash()
	case TopicCreateTransaction:
		return i.GetTransactionHash()
	case TopicDeleteTransaction:
		return i.GetTransactionHash()
	case TopicMessageSubmitTransaction:
		return i.GetTransactionHash()
	case TopicUpdateTransaction:
		return i.GetTransactionHash()
	case TransferTransaction:
		return i.GetTransactionHash()
	case *AccountCreateTransaction:
		return i.GetTransactionHash()
	case *AccountDeleteTransaction:
		return i.GetTransactionHash()
	case *AccountUpdateTransaction:
		return i.GetTransactionHash()
	case *AccountAllowanceApproveTransaction:
		return i.GetTransactionHash()
	case *AccountAllowanceDeleteTransaction:
		return i.GetTransactionHash()
	case *ContractCreateTransaction:
		return i.GetTransactionHash()
	case *ContractDeleteTransaction:
		return i.GetTransactionHash()
	case *ContractExecuteTransaction:
		return i.GetTransactionHash()
	case *ContractUpdateTransaction:
		return i.GetTransactionHash()
	case *FileAppendTransaction:
		return i.GetTransactionHash()
	case *FileCreateTransaction:
		return i.GetTransactionHash()
	case *FileDeleteTransaction:
		return i.GetTransactionHash()
	case *FileUpdateTransaction:
		return i.GetTransactionHash()
	case *LiveHashAddTransaction:
		return i.GetTransactionHash()
	case *LiveHashDeleteTransaction:
		return i.GetTransactionHash()
	case *ScheduleCreateTransaction:
		return i.GetTransactionHash()
	case *ScheduleDeleteTransaction:
		return i.GetTransactionHash()
	case *ScheduleSignTransaction:
		return i.GetTransactionHash()
	case *SystemDeleteTransaction:
		return i.GetTransactionHash()
	case *SystemUndeleteTransaction:
		return i.GetTransactionHash()
	case *TokenAssociateTransaction:
		return i.GetTransactionHash()
	case *TokenBurnTransaction:
		return i.GetTransactionHash()
	case *TokenCreateTransaction:
		return i.GetTransactionHash()
	case *TokenDeleteTransaction:
		return i.GetTransactionHash()
	case *TokenDissociateTransaction:
		return i.GetTransactionHash()
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionHash()
	case *TokenFreezeTransaction:
		return i.GetTransactionHash()
	case *TokenGrantKycTransaction:
		return i.GetTransactionHash()
	case *TokenMintTransaction:
		return i.GetTransactionHash()
	case *TokenRevokeKycTransaction:
		return i.GetTransactionHash()
	case *TokenUnfreezeTransaction:
		return i.GetTransactionHash()
	case *TokenUpdateTransaction:
		return i.GetTransactionHash()
	case *TokenWipeTransaction:
		return i.GetTransactionHash()
	case *TopicCreateTransaction:
		return i.GetTransactionHash()
	case *TopicDeleteTransaction:
		return i.GetTransactionHash()
	case *TopicMessageSubmitTransaction:
		return i.GetTransactionHash()
	case *TopicUpdateTransaction:
		return i.GetTransactionHash()
	case *TransferTransaction:
		return i.GetTransactionHash()
	default:
		return nil, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetTransactionHashPerNode(transaction interface{}) (map[AccountID][]byte, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetTransactionHashPerNode()
	case AccountDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case AccountUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case AccountAllowanceApproveTransaction:
		return i.GetTransactionHashPerNode()
	case AccountAllowanceDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case ContractCreateTransaction:
		return i.GetTransactionHashPerNode()
	case ContractDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case ContractExecuteTransaction:
		return i.GetTransactionHashPerNode()
	case ContractUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case FileAppendTransaction:
		return i.GetTransactionHashPerNode()
	case FileCreateTransaction:
		return i.GetTransactionHashPerNode()
	case FileDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case FileUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case LiveHashAddTransaction:
		return i.GetTransactionHashPerNode()
	case LiveHashDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case ScheduleCreateTransaction:
		return i.GetTransactionHashPerNode()
	case ScheduleDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case ScheduleSignTransaction:
		return i.GetTransactionHashPerNode()
	case SystemDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case SystemUndeleteTransaction:
		return i.GetTransactionHashPerNode()
	case TokenAssociateTransaction:
		return i.GetTransactionHashPerNode()
	case TokenBurnTransaction:
		return i.GetTransactionHashPerNode()
	case TokenCreateTransaction:
		return i.GetTransactionHashPerNode()
	case TokenDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case TokenDissociateTransaction:
		return i.GetTransactionHashPerNode()
	case TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case TokenFreezeTransaction:
		return i.GetTransactionHashPerNode()
	case TokenGrantKycTransaction:
		return i.GetTransactionHashPerNode()
	case TokenMintTransaction:
		return i.GetTransactionHashPerNode()
	case TokenRevokeKycTransaction:
		return i.GetTransactionHashPerNode()
	case TokenUnfreezeTransaction:
		return i.GetTransactionHashPerNode()
	case TokenUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case TokenWipeTransaction:
		return i.GetTransactionHashPerNode()
	case TopicCreateTransaction:
		return i.GetTransactionHashPerNode()
	case TopicDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case TopicMessageSubmitTransaction:
		return i.GetTransactionHashPerNode()
	case TopicUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case TransferTransaction:
		return i.GetTransactionHashPerNode()
	case *AccountCreateTransaction:
		return i.GetTransactionHashPerNode()
	case *AccountDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *AccountUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case *AccountAllowanceApproveTransaction:
		return i.GetTransactionHashPerNode()
	case *AccountAllowanceDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *ContractCreateTransaction:
		return i.GetTransactionHashPerNode()
	case *ContractDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *ContractExecuteTransaction:
		return i.GetTransactionHashPerNode()
	case *ContractUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case *FileAppendTransaction:
		return i.GetTransactionHashPerNode()
	case *FileCreateTransaction:
		return i.GetTransactionHashPerNode()
	case *FileDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *FileUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case *LiveHashAddTransaction:
		return i.GetTransactionHashPerNode()
	case *LiveHashDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *ScheduleCreateTransaction:
		return i.GetTransactionHashPerNode()
	case *ScheduleDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *ScheduleSignTransaction:
		return i.GetTransactionHashPerNode()
	case *SystemDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *SystemUndeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenAssociateTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenBurnTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenCreateTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenDissociateTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenFreezeTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenGrantKycTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenMintTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenRevokeKycTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenUnfreezeTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case *TokenWipeTransaction:
		return i.GetTransactionHashPerNode()
	case *TopicCreateTransaction:
		return i.GetTransactionHashPerNode()
	case *TopicDeleteTransaction:
		return i.GetTransactionHashPerNode()
	case *TopicMessageSubmitTransaction:
		return i.GetTransactionHashPerNode()
	case *TopicUpdateTransaction:
		return i.GetTransactionHashPerNode()
	case *TransferTransaction:
		return i.GetTransactionHashPerNode()
	default:
		return nil, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSetMinBackoff(transaction interface{}, minBackoff time.Duration) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case AccountDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case AccountUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case AccountAllowanceApproveTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case AccountAllowanceDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case ContractCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case ContractDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case ContractExecuteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case ContractUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case FileAppendTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case FileCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case FileDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case FileUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case LiveHashAddTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case LiveHashDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case ScheduleCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case ScheduleDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case ScheduleSignTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case SystemDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case SystemUndeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenAssociateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenBurnTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenDissociateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenFreezeTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenGrantKycTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenMintTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenRevokeKycTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenUnfreezeTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TokenWipeTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TopicCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TopicDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TopicMessageSubmitTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TopicUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case TransferTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *AccountCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *AccountDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *AccountUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *AccountAllowanceApproveTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *AccountAllowanceDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *ContractCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *ContractDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *ContractExecuteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *ContractUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *FileAppendTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *FileCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *FileDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *FileUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *LiveHashAddTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *LiveHashDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *ScheduleCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *ScheduleDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *ScheduleSignTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *SystemDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *SystemUndeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenAssociateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenBurnTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenDissociateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenFreezeTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenGrantKycTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenMintTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenRevokeKycTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenUnfreezeTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TokenWipeTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TopicCreateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TopicDeleteTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TopicMessageSubmitTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TopicUpdateTransaction:
		return i.SetMinBackoff(minBackoff), nil
	case *TransferTransaction:
		return i.SetMinBackoff(minBackoff), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetMinBackoff(transaction interface{}) (time.Duration, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetMinBackoff(), nil
	case AccountDeleteTransaction:
		return i.GetMinBackoff(), nil
	case AccountUpdateTransaction:
		return i.GetMinBackoff(), nil
	case AccountAllowanceApproveTransaction:
		return i.GetMinBackoff(), nil
	case AccountAllowanceDeleteTransaction:
		return i.GetMinBackoff(), nil
	case ContractCreateTransaction:
		return i.GetMinBackoff(), nil
	case ContractDeleteTransaction:
		return i.GetMinBackoff(), nil
	case ContractExecuteTransaction:
		return i.GetMinBackoff(), nil
	case ContractUpdateTransaction:
		return i.GetMinBackoff(), nil
	case FileAppendTransaction:
		return i.GetMinBackoff(), nil
	case FileCreateTransaction:
		return i.GetMinBackoff(), nil
	case FileDeleteTransaction:
		return i.GetMinBackoff(), nil
	case FileUpdateTransaction:
		return i.GetMinBackoff(), nil
	case LiveHashAddTransaction:
		return i.GetMinBackoff(), nil
	case LiveHashDeleteTransaction:
		return i.GetMinBackoff(), nil
	case ScheduleCreateTransaction:
		return i.GetMinBackoff(), nil
	case ScheduleDeleteTransaction:
		return i.GetMinBackoff(), nil
	case ScheduleSignTransaction:
		return i.GetMinBackoff(), nil
	case SystemDeleteTransaction:
		return i.GetMinBackoff(), nil
	case SystemUndeleteTransaction:
		return i.GetMinBackoff(), nil
	case TokenAssociateTransaction:
		return i.GetMinBackoff(), nil
	case TokenBurnTransaction:
		return i.GetMinBackoff(), nil
	case TokenCreateTransaction:
		return i.GetMinBackoff(), nil
	case TokenDeleteTransaction:
		return i.GetMinBackoff(), nil
	case TokenDissociateTransaction:
		return i.GetMinBackoff(), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.GetMinBackoff(), nil
	case TokenFreezeTransaction:
		return i.GetMinBackoff(), nil
	case TokenGrantKycTransaction:
		return i.GetMinBackoff(), nil
	case TokenMintTransaction:
		return i.GetMinBackoff(), nil
	case TokenRevokeKycTransaction:
		return i.GetMinBackoff(), nil
	case TokenUnfreezeTransaction:
		return i.GetMinBackoff(), nil
	case TokenUpdateTransaction:
		return i.GetMinBackoff(), nil
	case TokenWipeTransaction:
		return i.GetMinBackoff(), nil
	case TopicCreateTransaction:
		return i.GetMinBackoff(), nil
	case TopicDeleteTransaction:
		return i.GetMinBackoff(), nil
	case TopicMessageSubmitTransaction:
		return i.GetMinBackoff(), nil
	case TopicUpdateTransaction:
		return i.GetMinBackoff(), nil
	case TransferTransaction:
		return i.GetMinBackoff(), nil
	case *AccountCreateTransaction:
		return i.GetMinBackoff(), nil
	case *AccountDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *AccountUpdateTransaction:
		return i.GetMinBackoff(), nil
	case *AccountAllowanceApproveTransaction:
		return i.GetMinBackoff(), nil
	case *AccountAllowanceDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *ContractCreateTransaction:
		return i.GetMinBackoff(), nil
	case *ContractDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *ContractExecuteTransaction:
		return i.GetMinBackoff(), nil
	case *ContractUpdateTransaction:
		return i.GetMinBackoff(), nil
	case *FileAppendTransaction:
		return i.GetMinBackoff(), nil
	case *FileCreateTransaction:
		return i.GetMinBackoff(), nil
	case *FileDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *FileUpdateTransaction:
		return i.GetMinBackoff(), nil
	case *LiveHashAddTransaction:
		return i.GetMinBackoff(), nil
	case *LiveHashDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *ScheduleCreateTransaction:
		return i.GetMinBackoff(), nil
	case *ScheduleDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *ScheduleSignTransaction:
		return i.GetMinBackoff(), nil
	case *SystemDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *SystemUndeleteTransaction:
		return i.GetMinBackoff(), nil
	case *TokenAssociateTransaction:
		return i.GetMinBackoff(), nil
	case *TokenBurnTransaction:
		return i.GetMinBackoff(), nil
	case *TokenCreateTransaction:
		return i.GetMinBackoff(), nil
	case *TokenDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *TokenDissociateTransaction:
		return i.GetMinBackoff(), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetMinBackoff(), nil
	case *TokenFreezeTransaction:
		return i.GetMinBackoff(), nil
	case *TokenGrantKycTransaction:
		return i.GetMinBackoff(), nil
	case *TokenMintTransaction:
		return i.GetMinBackoff(), nil
	case *TokenRevokeKycTransaction:
		return i.GetMinBackoff(), nil
	case *TokenUnfreezeTransaction:
		return i.GetMinBackoff(), nil
	case *TokenUpdateTransaction:
		return i.GetMinBackoff(), nil
	case *TokenWipeTransaction:
		return i.GetMinBackoff(), nil
	case *TopicCreateTransaction:
		return i.GetMinBackoff(), nil
	case *TopicDeleteTransaction:
		return i.GetMinBackoff(), nil
	case *TopicMessageSubmitTransaction:
		return i.GetMinBackoff(), nil
	case *TopicUpdateTransaction:
		return i.GetMinBackoff(), nil
	case *TransferTransaction:
		return i.GetMinBackoff(), nil
	default:
		return time.Duration(0), errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionSetMaxBackoff(transaction interface{}, maxBackoff time.Duration) (interface{}, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case AccountDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case AccountUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case AccountAllowanceApproveTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case AccountAllowanceDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case ContractCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case ContractDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case ContractExecuteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case ContractUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case FileAppendTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case FileCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case FileDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case FileUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case LiveHashAddTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case LiveHashDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case ScheduleCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case ScheduleDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case ScheduleSignTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case SystemDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case SystemUndeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenAssociateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenBurnTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenDissociateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenFreezeTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenGrantKycTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenMintTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenRevokeKycTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenUnfreezeTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TokenWipeTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TopicCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TopicDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TopicMessageSubmitTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TopicUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case TransferTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *AccountCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *AccountDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *AccountUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *AccountAllowanceApproveTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *AccountAllowanceDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *ContractCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *ContractDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *ContractExecuteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *ContractUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *FileAppendTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *FileCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *FileDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *FileUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *LiveHashAddTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *LiveHashDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *ScheduleCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *ScheduleDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *ScheduleSignTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *SystemDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *SystemUndeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenAssociateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenBurnTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenDissociateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenFreezeTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenGrantKycTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenMintTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenRevokeKycTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenUnfreezeTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TokenWipeTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TopicCreateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TopicDeleteTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TopicMessageSubmitTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TopicUpdateTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	case *TransferTransaction:
		return i.SetMaxBackoff(maxBackoff), nil
	default:
		return transaction, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionGetMaxBackoff(transaction interface{}) (time.Duration, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.GetMaxBackoff(), nil
	case AccountDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case AccountUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case AccountAllowanceApproveTransaction:
		return i.GetMaxBackoff(), nil
	case AccountAllowanceDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case ContractCreateTransaction:
		return i.GetMaxBackoff(), nil
	case ContractDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case ContractExecuteTransaction:
		return i.GetMaxBackoff(), nil
	case ContractUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case FileAppendTransaction:
		return i.GetMaxBackoff(), nil
	case FileCreateTransaction:
		return i.GetMaxBackoff(), nil
	case FileDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case FileUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case LiveHashAddTransaction:
		return i.GetMaxBackoff(), nil
	case LiveHashDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case ScheduleCreateTransaction:
		return i.GetMaxBackoff(), nil
	case ScheduleDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case ScheduleSignTransaction:
		return i.GetMaxBackoff(), nil
	case SystemDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case SystemUndeleteTransaction:
		return i.GetMaxBackoff(), nil
	case TokenAssociateTransaction:
		return i.GetMaxBackoff(), nil
	case TokenBurnTransaction:
		return i.GetMaxBackoff(), nil
	case TokenCreateTransaction:
		return i.GetMaxBackoff(), nil
	case TokenDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case TokenDissociateTransaction:
		return i.GetMaxBackoff(), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case TokenFreezeTransaction:
		return i.GetMaxBackoff(), nil
	case TokenGrantKycTransaction:
		return i.GetMaxBackoff(), nil
	case TokenMintTransaction:
		return i.GetMaxBackoff(), nil
	case TokenRevokeKycTransaction:
		return i.GetMaxBackoff(), nil
	case TokenUnfreezeTransaction:
		return i.GetMaxBackoff(), nil
	case TokenUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case TokenWipeTransaction:
		return i.GetMaxBackoff(), nil
	case TopicCreateTransaction:
		return i.GetMaxBackoff(), nil
	case TopicDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case TopicMessageSubmitTransaction:
		return i.GetMaxBackoff(), nil
	case TopicUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case TransferTransaction:
		return i.GetMaxBackoff(), nil
	case *AccountCreateTransaction:
		return i.GetMaxBackoff(), nil
	case *AccountDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *AccountUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case *AccountAllowanceApproveTransaction:
		return i.GetMaxBackoff(), nil
	case *AccountAllowanceDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *ContractCreateTransaction:
		return i.GetMaxBackoff(), nil
	case *ContractDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *ContractExecuteTransaction:
		return i.GetMaxBackoff(), nil
	case *ContractUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case *FileAppendTransaction:
		return i.GetMaxBackoff(), nil
	case *FileCreateTransaction:
		return i.GetMaxBackoff(), nil
	case *FileDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *FileUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case *LiveHashAddTransaction:
		return i.GetMaxBackoff(), nil
	case *LiveHashDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *ScheduleCreateTransaction:
		return i.GetMaxBackoff(), nil
	case *ScheduleDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *ScheduleSignTransaction:
		return i.GetMaxBackoff(), nil
	case *SystemDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *SystemUndeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenAssociateTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenBurnTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenCreateTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenDissociateTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenFreezeTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenGrantKycTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenMintTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenRevokeKycTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenUnfreezeTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case *TokenWipeTransaction:
		return i.GetMaxBackoff(), nil
	case *TopicCreateTransaction:
		return i.GetMaxBackoff(), nil
	case *TopicDeleteTransaction:
		return i.GetMaxBackoff(), nil
	case *TopicMessageSubmitTransaction:
		return i.GetMaxBackoff(), nil
	case *TopicUpdateTransaction:
		return i.GetMaxBackoff(), nil
	case *TransferTransaction:
		return i.GetMaxBackoff(), nil
	default:
		return time.Duration(0), errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionString(transaction interface{}) (string, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.String(), nil
	case AccountDeleteTransaction:
		return i.String(), nil
	case AccountUpdateTransaction:
		return i.String(), nil
	case AccountAllowanceApproveTransaction:
		return i.String(), nil
	case AccountAllowanceDeleteTransaction:
		return i.String(), nil
	case ContractCreateTransaction:
		return i.String(), nil
	case ContractDeleteTransaction:
		return i.String(), nil
	case ContractExecuteTransaction:
		return i.String(), nil
	case ContractUpdateTransaction:
		return i.String(), nil
	case FileAppendTransaction:
		return i.String(), nil
	case FileCreateTransaction:
		return i.String(), nil
	case FileDeleteTransaction:
		return i.String(), nil
	case FileUpdateTransaction:
		return i.String(), nil
	case LiveHashAddTransaction:
		return i.String(), nil
	case LiveHashDeleteTransaction:
		return i.String(), nil
	case ScheduleCreateTransaction:
		return i.String(), nil
	case ScheduleDeleteTransaction:
		return i.String(), nil
	case ScheduleSignTransaction:
		return i.String(), nil
	case SystemDeleteTransaction:
		return i.String(), nil
	case SystemUndeleteTransaction:
		return i.String(), nil
	case TokenAssociateTransaction:
		return i.String(), nil
	case TokenBurnTransaction:
		return i.String(), nil
	case TokenCreateTransaction:
		return i.String(), nil
	case TokenDeleteTransaction:
		return i.String(), nil
	case TokenDissociateTransaction:
		return i.String(), nil
	case TokenFeeScheduleUpdateTransaction:
		return i.String(), nil
	case TokenFreezeTransaction:
		return i.String(), nil
	case TokenGrantKycTransaction:
		return i.String(), nil
	case TokenMintTransaction:
		return i.String(), nil
	case TokenRevokeKycTransaction:
		return i.String(), nil
	case TokenUnfreezeTransaction:
		return i.String(), nil
	case TokenUpdateTransaction:
		return i.String(), nil
	case TokenWipeTransaction:
		return i.String(), nil
	case TopicCreateTransaction:
		return i.String(), nil
	case TopicDeleteTransaction:
		return i.String(), nil
	case TopicMessageSubmitTransaction:
		return i.String(), nil
	case TopicUpdateTransaction:
		return i.String(), nil
	case TransferTransaction:
		return i.String(), nil
	case *AccountCreateTransaction:
		return i.String(), nil
	case *AccountDeleteTransaction:
		return i.String(), nil
	case *AccountUpdateTransaction:
		return i.String(), nil
	case *AccountAllowanceApproveTransaction:
		return i.String(), nil
	case *AccountAllowanceDeleteTransaction:
		return i.String(), nil
	case *ContractCreateTransaction:
		return i.String(), nil
	case *ContractDeleteTransaction:
		return i.String(), nil
	case *ContractExecuteTransaction:
		return i.String(), nil
	case *ContractUpdateTransaction:
		return i.String(), nil
	case *FileAppendTransaction:
		return i.String(), nil
	case *FileCreateTransaction:
		return i.String(), nil
	case *FileDeleteTransaction:
		return i.String(), nil
	case *FileUpdateTransaction:
		return i.String(), nil
	case *LiveHashAddTransaction:
		return i.String(), nil
	case *LiveHashDeleteTransaction:
		return i.String(), nil
	case *ScheduleCreateTransaction:
		return i.String(), nil
	case *ScheduleDeleteTransaction:
		return i.String(), nil
	case *ScheduleSignTransaction:
		return i.String(), nil
	case *SystemDeleteTransaction:
		return i.String(), nil
	case *SystemUndeleteTransaction:
		return i.String(), nil
	case *TokenAssociateTransaction:
		return i.String(), nil
	case *TokenBurnTransaction:
		return i.String(), nil
	case *TokenCreateTransaction:
		return i.String(), nil
	case *TokenDeleteTransaction:
		return i.String(), nil
	case *TokenDissociateTransaction:
		return i.String(), nil
	case *TokenFeeScheduleUpdateTransaction:
		return i.String(), nil
	case *TokenFreezeTransaction:
		return i.String(), nil
	case *TokenGrantKycTransaction:
		return i.String(), nil
	case *TokenMintTransaction:
		return i.String(), nil
	case *TokenRevokeKycTransaction:
		return i.String(), nil
	case *TokenUnfreezeTransaction:
		return i.String(), nil
	case *TokenUpdateTransaction:
		return i.String(), nil
	case *TokenWipeTransaction:
		return i.String(), nil
	case *TopicCreateTransaction:
		return i.String(), nil
	case *TopicDeleteTransaction:
		return i.String(), nil
	case *TopicMessageSubmitTransaction:
		return i.String(), nil
	case *TopicUpdateTransaction:
		return i.String(), nil
	case *TransferTransaction:
		return i.String(), nil
	default:
		return "", errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionToBytes(transaction interface{}) ([]byte, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.ToBytes()
	case AccountDeleteTransaction:
		return i.ToBytes()
	case AccountUpdateTransaction:
		return i.ToBytes()
	case AccountAllowanceApproveTransaction:
		return i.ToBytes()
	case AccountAllowanceDeleteTransaction:
		return i.ToBytes()
	case ContractCreateTransaction:
		return i.ToBytes()
	case ContractDeleteTransaction:
		return i.ToBytes()
	case ContractExecuteTransaction:
		return i.ToBytes()
	case ContractUpdateTransaction:
		return i.ToBytes()
	case FileAppendTransaction:
		return i.ToBytes()
	case FileCreateTransaction:
		return i.ToBytes()
	case FileDeleteTransaction:
		return i.ToBytes()
	case FileUpdateTransaction:
		return i.ToBytes()
	case LiveHashAddTransaction:
		return i.ToBytes()
	case LiveHashDeleteTransaction:
		return i.ToBytes()
	case ScheduleCreateTransaction:
		return i.ToBytes()
	case ScheduleDeleteTransaction:
		return i.ToBytes()
	case ScheduleSignTransaction:
		return i.ToBytes()
	case SystemDeleteTransaction:
		return i.ToBytes()
	case SystemUndeleteTransaction:
		return i.ToBytes()
	case TokenAssociateTransaction:
		return i.ToBytes()
	case TokenBurnTransaction:
		return i.ToBytes()
	case TokenCreateTransaction:
		return i.ToBytes()
	case TokenDeleteTransaction:
		return i.ToBytes()
	case TokenDissociateTransaction:
		return i.ToBytes()
	case TokenFeeScheduleUpdateTransaction:
		return i.ToBytes()
	case TokenFreezeTransaction:
		return i.ToBytes()
	case TokenGrantKycTransaction:
		return i.ToBytes()
	case TokenMintTransaction:
		return i.ToBytes()
	case TokenRevokeKycTransaction:
		return i.ToBytes()
	case TokenUnfreezeTransaction:
		return i.ToBytes()
	case TokenUpdateTransaction:
		return i.ToBytes()
	case TokenWipeTransaction:
		return i.ToBytes()
	case TopicCreateTransaction:
		return i.ToBytes()
	case TopicDeleteTransaction:
		return i.ToBytes()
	case TopicMessageSubmitTransaction:
		return i.ToBytes()
	case TopicUpdateTransaction:
		return i.ToBytes()
	case TransferTransaction:
		return i.ToBytes()
	case *AccountCreateTransaction:
		return i.ToBytes()
	case *AccountDeleteTransaction:
		return i.ToBytes()
	case *AccountUpdateTransaction:
		return i.ToBytes()
	case *AccountAllowanceApproveTransaction:
		return i.ToBytes()
	case *AccountAllowanceDeleteTransaction:
		return i.ToBytes()
	case *ContractCreateTransaction:
		return i.ToBytes()
	case *ContractDeleteTransaction:
		return i.ToBytes()
	case *ContractExecuteTransaction:
		return i.ToBytes()
	case *ContractUpdateTransaction:
		return i.ToBytes()
	case *FileAppendTransaction:
		return i.ToBytes()
	case *FileCreateTransaction:
		return i.ToBytes()
	case *FileDeleteTransaction:
		return i.ToBytes()
	case *FileUpdateTransaction:
		return i.ToBytes()
	case *LiveHashAddTransaction:
		return i.ToBytes()
	case *LiveHashDeleteTransaction:
		return i.ToBytes()
	case *ScheduleCreateTransaction:
		return i.ToBytes()
	case *ScheduleDeleteTransaction:
		return i.ToBytes()
	case *ScheduleSignTransaction:
		return i.ToBytes()
	case *SystemDeleteTransaction:
		return i.ToBytes()
	case *SystemUndeleteTransaction:
		return i.ToBytes()
	case *TokenAssociateTransaction:
		return i.ToBytes()
	case *TokenBurnTransaction:
		return i.ToBytes()
	case *TokenCreateTransaction:
		return i.ToBytes()
	case *TokenDeleteTransaction:
		return i.ToBytes()
	case *TokenDissociateTransaction:
		return i.ToBytes()
	case *TokenFeeScheduleUpdateTransaction:
		return i.ToBytes()
	case *TokenFreezeTransaction:
		return i.ToBytes()
	case *TokenGrantKycTransaction:
		return i.ToBytes()
	case *TokenMintTransaction:
		return i.ToBytes()
	case *TokenRevokeKycTransaction:
		return i.ToBytes()
	case *TokenUnfreezeTransaction:
		return i.ToBytes()
	case *TokenUpdateTransaction:
		return i.ToBytes()
	case *TokenWipeTransaction:
		return i.ToBytes()
	case *TopicCreateTransaction:
		return i.ToBytes()
	case *TopicDeleteTransaction:
		return i.ToBytes()
	case *TopicMessageSubmitTransaction:
		return i.ToBytes()
	case *TopicUpdateTransaction:
		return i.ToBytes()
	case *TransferTransaction:
		return i.ToBytes()
	default:
		return nil, errors.New("(BUG) non-exhaustive switch statement")
	}
}

func TransactionExecute(transaction interface{}, client *Client) (TransactionResponse, error) { // nolint
	switch i := transaction.(type) {
	case AccountCreateTransaction:
		return i.Execute(client)
	case AccountDeleteTransaction:
		return i.Execute(client)
	case AccountUpdateTransaction:
		return i.Execute(client)
	case AccountAllowanceApproveTransaction:
		return i.Execute(client)
	case AccountAllowanceDeleteTransaction:
		return i.Execute(client)
	case ContractCreateTransaction:
		return i.Execute(client)
	case ContractDeleteTransaction:
		return i.Execute(client)
	case ContractExecuteTransaction:
		return i.Execute(client)
	case ContractUpdateTransaction:
		return i.Execute(client)
	case FileAppendTransaction:
		return i.Execute(client)
	case FileCreateTransaction:
		return i.Execute(client)
	case FileDeleteTransaction:
		return i.Execute(client)
	case FileUpdateTransaction:
		return i.Execute(client)
	case FreezeTransaction:
		return i.Execute(client)
	case LiveHashAddTransaction:
		return i.Execute(client)
	case LiveHashDeleteTransaction:
		return i.Execute(client)
	case ScheduleCreateTransaction:
		return i.Execute(client)
	case ScheduleDeleteTransaction:
		return i.Execute(client)
	case ScheduleSignTransaction:
		return i.Execute(client)
	case SystemDeleteTransaction:
		return i.Execute(client)
	case SystemUndeleteTransaction:
		return i.Execute(client)
	case TokenAssociateTransaction:
		return i.Execute(client)
	case TokenBurnTransaction:
		return i.Execute(client)
	case TokenCreateTransaction:
		return i.Execute(client)
	case TokenDeleteTransaction:
		return i.Execute(client)
	case TokenDissociateTransaction:
		return i.Execute(client)
	case TokenFeeScheduleUpdateTransaction:
		return i.Execute(client)
	case TokenFreezeTransaction:
		return i.Execute(client)
	case TokenGrantKycTransaction:
		return i.Execute(client)
	case TokenMintTransaction:
		return i.Execute(client)
	case TokenPauseTransaction:
		return i.Execute(client)
	case TokenRevokeKycTransaction:
		return i.Execute(client)
	case TokenUnfreezeTransaction:
		return i.Execute(client)
	case TokenUnpauseTransaction:
		return i.Execute(client)
	case TokenUpdateTransaction:
		return i.Execute(client)
	case TokenWipeTransaction:
		return i.Execute(client)
	case TopicCreateTransaction:
		return i.Execute(client)
	case TopicDeleteTransaction:
		return i.Execute(client)
	case TopicMessageSubmitTransaction:
		return i.Execute(client)
	case TopicUpdateTransaction:
		return i.Execute(client)
	case TransferTransaction:
		return i.Execute(client)
	case *AccountCreateTransaction:
		return i.Execute(client)
	case *AccountDeleteTransaction:
		return i.Execute(client)
	case *AccountUpdateTransaction:
		return i.Execute(client)
	case *AccountAllowanceApproveTransaction:
		return i.Execute(client)
	case *AccountAllowanceDeleteTransaction:
		return i.Execute(client)
	case *ContractCreateTransaction:
		return i.Execute(client)
	case *ContractDeleteTransaction:
		return i.Execute(client)
	case *ContractExecuteTransaction:
		return i.Execute(client)
	case *ContractUpdateTransaction:
		return i.Execute(client)
	case *FileAppendTransaction:
		return i.Execute(client)
	case *FileCreateTransaction:
		return i.Execute(client)
	case *FileDeleteTransaction:
		return i.Execute(client)
	case *FileUpdateTransaction:
		return i.Execute(client)
	case *FreezeTransaction:
		return i.Execute(client)
	case *LiveHashAddTransaction:
		return i.Execute(client)
	case *LiveHashDeleteTransaction:
		return i.Execute(client)
	case *ScheduleCreateTransaction:
		return i.Execute(client)
	case *ScheduleDeleteTransaction:
		return i.Execute(client)
	case *ScheduleSignTransaction:
		return i.Execute(client)
	case *SystemDeleteTransaction:
		return i.Execute(client)
	case *SystemUndeleteTransaction:
		return i.Execute(client)
	case *TokenAssociateTransaction:
		return i.Execute(client)
	case *TokenBurnTransaction:
		return i.Execute(client)
	case *TokenCreateTransaction:
		return i.Execute(client)
	case *TokenDeleteTransaction:
		return i.Execute(client)
	case *TokenDissociateTransaction:
		return i.Execute(client)
	case *TokenFeeScheduleUpdateTransaction:
		return i.Execute(client)
	case *TokenFreezeTransaction:
		return i.Execute(client)
	case *TokenGrantKycTransaction:
		return i.Execute(client)
	case *TokenMintTransaction:
		return i.Execute(client)
	case *TokenPauseTransaction:
		return i.Execute(client)
	case *TokenRevokeKycTransaction:
		return i.Execute(client)
	case *TokenUnfreezeTransaction:
		return i.Execute(client)
	case *TokenUnpauseTransaction:
		return i.Execute(client)
	case *TokenUpdateTransaction:
		return i.Execute(client)
	case *TokenWipeTransaction:
		return i.Execute(client)
	case *TopicCreateTransaction:
		return i.Execute(client)
	case *TopicDeleteTransaction:
		return i.Execute(client)
	case *TopicMessageSubmitTransaction:
		return i.Execute(client)
	case *TopicUpdateTransaction:
		return i.Execute(client)
	case *TransferTransaction:
		return i.Execute(client)
	default:
		return TransactionResponse{}, errors.New("(BUG) non-exhaustive switch statement")
	}
}

// ------------ Executable Functions ------------

func (tx *Transaction) shouldRetry(_ Executable, response interface{}) _ExecutionState {
	status := Status(response.(*services.TransactionResponse).NodeTransactionPrecheckCode)
	switch status {
	case StatusPlatformTransactionNotCreated, StatusPlatformNotActive, StatusBusy:
		return executionStateRetry
	case StatusTransactionExpired:
		return executionStateExpired
	case StatusOk:
		return executionStateFinished
	}

	return executionStateError
}

func (tx *Transaction) makeRequest() interface{} {
	index := tx.nodeAccountIDs._Length()*tx.transactionIDs.index + tx.nodeAccountIDs.index
	built, _ := tx._BuildTransaction(index)

	return built
}

func (tx *Transaction) advanceRequest() {
	tx.nodeAccountIDs._Advance()
	tx.signedTransactions._Advance()
}

func (tx *Transaction) getNodeAccountID() AccountID {
	return tx.nodeAccountIDs._GetCurrent().(AccountID)
}

func (tx *Transaction) mapStatusError(
	_ Executable,
	response interface{},
) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.TransactionResponse).NodeTransactionPrecheckCode),
		//NodeID: request.transaction.nodeAccountIDs,
		TxID: tx.GetTransactionID(),
	}
}

func (tx *Transaction) mapResponse(_ interface{}, nodeID AccountID, protoRequest interface{}) (interface{}, error) {
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

// Building empty object as "default" implementation. All inhertents must implement their own implementation.
func (tx *Transaction) getMethod(ch *_Channel) _Method {
	return _Method{}
}

// Building empty object as "default" implementation. All inhertents must implement their own implementation.
func (tx *Transaction) getName() string {
	return "transaction"
}

// Building empty object as "default" implementation. All inhertents must implement their own implementation.
func (tx *Transaction) validateNetworkOnIDs(client *Client) error {
	return errors.New("Function not implemented")
}

func (tx *Transaction) preFreezeWith(*Client) {
	// NO-OP
}

func (tx *Transaction) isTransaction() bool {
	return true
}

func (tx *Transaction) getTransactionIDAndMessage() (string, string) {
	return tx.GetTransactionID().String(), "transaction status received"
}

func (tx *Transaction) regenerateID(client *Client) bool {
	if !client.GetOperatorAccountID()._IsZero() && tx.regenerateTransactionID && !tx.transactionIDs.locked {
		tx.transactionIDs._Set(tx.transactionIDs.index, TransactionIDGenerate(client.GetOperatorAccountID()))
		return true
	}
	return false
}

func (tx *Transaction) execute(client *Client, e TransactionInterface) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if tx.freezeError != nil {
		return TransactionResponse{}, tx.freezeError
	}

	if !tx.IsFrozen() {
		_, err := tx.freezeWith(client, e)
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

	resp, err := _Execute(client, e)

	if err != nil {
		return TransactionResponse{
			TransactionID:  tx.GetTransactionID(),
			NodeID:         resp.(TransactionResponse).NodeID,
			ValidateStatus: true,
		}, err
	}

	return TransactionResponse{
		TransactionID:  tx.GetTransactionID(),
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
	}, nil
}

func (tx *Transaction) freezeWith(client *Client, e TransactionInterface) (TransactionInterface, error) {
	if tx.IsFrozen() {
		return tx, nil
	}

	e.preFreezeWith(client)

	tx._InitFee(client)
	if err := tx._InitTransactionID(client); err != nil {
		return tx, err
	}

	err := e.validateNetworkOnIDs(client)
	if err != nil {
		return &Transaction{}, err
	}
	body := e.build()

	return tx, _TransactionFreezeWith(tx, client, body)
}

func (tx *Transaction) schedule(e TransactionInterface) (*ScheduleCreateTransaction, error) {
	tx._RequireNotFrozen()

	scheduled, err := e.buildScheduled()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

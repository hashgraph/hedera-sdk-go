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
	// We introduce a boolean value to distinguish flow for signed tx vs unsigned transactions
	txIsSigned := true

	for i, transactionFromList := range list.TransactionList {
		var signedTransaction services.SignedTransaction
		var body services.TransactionBody

		// If the transaction is not signed/locked:
		if len(transactionFromList.SignedTransactionBytes) == 0 {
			txIsSigned = false
			if err := protobuf.Unmarshal(transactionFromList.BodyBytes, &body); err != nil { // nolint
				return Transaction{}, errors.Wrap(err, "error deserializing BodyBytes in TransactionFromBytes")
			}
		} else { // If the transaction is signed/locked
			if err := protobuf.Unmarshal(transactionFromList.SignedTransactionBytes, &signedTransaction); err != nil {
				return Transaction{}, errors.Wrap(err, "error deserializing SignedTransactionBytes in TransactionFromBytes")
			}
		}

		if txIsSigned {
			tx.signedTransactions = tx.signedTransactions._Push(&signedTransaction)

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

			if err := protobuf.Unmarshal(signedTransaction.GetBodyBytes(), &body); err != nil {
				return Transaction{}, errors.Wrap(err, "error deserializing BodyBytes in TransactionFromBytes")
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
		return *_ContractExecuteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ContractCreateInstance:
		return *_ContractCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ContractUpdateInstance:
		return *_ContractUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ContractDeleteInstance:
		return *_ContractDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoAddLiveHash:
		return *_LiveHashAddTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoCreateAccount:
		return *_AccountCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoDelete:
		return *_AccountDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoDeleteLiveHash:
		return *_LiveHashDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoTransfer:
		return *_TransferTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoUpdateAccount:
		return *_AccountUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoApproveAllowance:
		return *_AccountAllowanceApproveTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_CryptoDeleteAllowance:
		return *_AccountAllowanceDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_FileAppend:
		return *_FileAppendTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_FileCreate:
		return *_FileCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_FileDelete:
		return *_FileDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_FileUpdate:
		return *_FileUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_SystemDelete:
		return *_SystemDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_SystemUndelete:
		return *_SystemUndeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_Freeze:
		return *_FreezeTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ConsensusCreateTopic:
		return *_TopicCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ConsensusUpdateTopic:
		return *_TopicUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ConsensusDeleteTopic:
		return *_TopicDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ConsensusSubmitMessage:
		return *_TopicMessageSubmitTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenCreation:
		return *_TokenCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenFreeze:
		return *_TokenFreezeTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenUnfreeze:
		return *_TokenUnfreezeTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenGrantKyc:
		return *_TokenGrantKycTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenRevokeKyc:
		return *_TokenRevokeKycTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenDeletion:
		return *_TokenDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenUpdate:
		return *_TokenUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenMint:
		return *_TokenMintTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenBurn:
		return *_TokenBurnTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenWipe:
		return *_TokenWipeTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenAssociate:
		return *_TokenAssociateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenDissociate:
		return *_TokenDissociateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ScheduleCreate:
		return *_ScheduleCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ScheduleSign:
		return *_ScheduleSignTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_ScheduleDelete:
		return *_ScheduleDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenFeeScheduleUpdate:
		return *_TokenFeeScheduleUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenPause:
		return *_TokenPauseTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenUnpause:
		return *_TokenUnpauseTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_EthereumTransaction:
		return *_EthereumTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_UtilPrng:
		return *_PrngTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenUpdateNfts:
		return *_NewTokenUpdateNftsTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenReject:
		return *_TokenRejectTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_NodeCreate:
		return *_NodeCreateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_NodeUpdate:
		return *_NodeUpdateTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_NodeDelete:
		return *_NodeDeleteTransactionFromProtobuf(tx, first), nil
	case *services.TransactionBody_TokenAirdrop:
		return *_TokenAirdropTransactionFromProtobuf(tx, first), nil
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

func (tx *Transaction) _RequireFrozen() {
	if !tx.IsFrozen() {
		tx.freezeError = errTransactionIsNotFrozen
	}
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
	return tx.toBytes(tx)
}

func (tx *Transaction) toBytes(e TransactionInterface) ([]byte, error) {
	var pbTransactionList []byte
	var allTx []*services.Transaction
	var err error
	// If transaction is frozen, build all transactions and "signedTransactions"
	if tx.IsFrozen() {
		allTx, err = tx._BuildAllTransactions()
		tx.transactionIDs.locked = true
	} else { // Build only onlt "BodyBytes" for each transaction in the list
		allTx, err = tx.buildAllUnsignedTransactions(e)
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

func (tx *Transaction) buildAllUnsignedTransactions(e TransactionInterface) ([]*services.Transaction, error) {
	// All unsigned transactions would always be exactly 1
	allTx := make([]*services.Transaction, 0)
	if tx.nodeAccountIDs._IsEmpty() {
		t, err := tx.buildUnsignedTransaction(e, 0)
		if err != nil {
			return allTx, err
		}
		allTx = append(allTx, t)
	} else { // If we have set some node account ids, we have to make one transaction copy per node account
		for range tx.nodeAccountIDs.slice {
			t, err := tx.buildUnsignedTransaction(e, tx.nodeAccountIDs.index)
			tx.nodeAccountIDs._Advance()
			if err != nil {
				return allTx, err
			}
			allTx = append(allTx, t)
		}
	}
	return allTx, nil
}

func (tx *Transaction) buildUnsignedTransaction(e TransactionInterface, index int) (*services.Transaction, error) {
	body := e.build()
	if body.NodeAccountID == nil && !tx.nodeAccountIDs._IsEmpty() {
		body.NodeAccountID = tx.nodeAccountIDs._Get(index).(AccountID)._ToProtobuf()
	}

	bodyBytes, err := protobuf.Marshal(body)
	if err != nil {
		return &services.Transaction{}, errors.Wrap(err, "failed to update tx ID")
	}

	return &services.Transaction{BodyBytes: bodyBytes}, nil
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
func (tx *Transaction) signWithOperator(client *Client, e TransactionInterface) (TransactionInterface, error) { // nolint
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
	// We need to make sure the request is frozen
	tx._RequireFrozen()

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

// ------------ Executable Functions ------------
func (tx *Transaction) shouldRetry(_ Executable, response interface{}) _ExecutionState {
	status := Status(response.(*services.TransactionResponse).NodeTransactionPrecheckCode)

	retryableStatuses := map[Status]bool{
		StatusPlatformTransactionNotCreated: true,
		StatusPlatformNotActive:             true,
		StatusBusy:                          true,
		StatusThrottledAtConsensus:          true,
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

func (tx *Transaction) getLogID(transactionInterface Executable) string {
	timestamp := tx.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("%s:%d", transactionInterface.getName(), timestamp.UnixNano())
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

	if tx.grpcDeadline == nil {
		tx.grpcDeadline = client.requestTimeout
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

func (tx *Transaction) FreezeWith(client *Client, e TransactionInterface) (TransactionInterface, error) {
	return tx.freezeWith(client, e)
}

func (tx *Transaction) freezeWith(client *Client, e TransactionInterface) (TransactionInterface, error) { //nolint
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

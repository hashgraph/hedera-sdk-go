package hedera

import (
	"bytes"
	"crypto/sha512"
	"fmt"

	"github.com/pkg/errors"

	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

// Transaction contains the protobuf of a prepared transaction which can be signed and executed.

type ITransaction interface {
	_ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error)
}

type Transaction struct {
	nextNodeIndex        int
	nextTransactionIndex int
	maxRetry             int

	transactionFee           uint64
	memo                     string
	transactionValidDuration *time.Duration
	transactionID            TransactionID

	transactionIDs     []TransactionID
	transactions       []*proto.Transaction
	signedTransactions []*proto.SignedTransaction
	nodeIDs            []AccountID

	publicKeys         []PublicKey
	transactionSigners []TransactionSigner

	freezeError error

	maxBackoff *time.Duration
	minBackoff *time.Duration
}

func _NewTransaction() Transaction {
	duration := 120 * time.Second
	return Transaction{
		nextNodeIndex:            0,
		nextTransactionIndex:     0,
		maxRetry:                 10,
		transactionValidDuration: &duration,
		transactionIDs:           make([]TransactionID, 0),
		transactions:             make([]*proto.Transaction, 0),
		signedTransactions:       make([]*proto.SignedTransaction, 0),
		nodeIDs:                  make([]AccountID, 0),
		freezeError:              nil,
	}
}

func TransactionFromBytes(data []byte) (interface{}, error) { // nolint
	list := proto.TransactionList{}
	err := protobuf.Unmarshal(data, &list)
	if err != nil {
		return Transaction{}, errors.Wrap(err, "error deserializing from bytes to Transaction List")
	}

	tx := Transaction{
		nextNodeIndex:        0,
		nextTransactionIndex: 0,
		maxRetry:             10,
		transactionIDs:       make([]TransactionID, 0),
		transactions:         list.TransactionList,
		signedTransactions:   make([]*proto.SignedTransaction, 0),
		publicKeys:           make([]PublicKey, 0),
		transactionSigners:   make([]TransactionSigner, 0),
	}

	var first *proto.TransactionBody = nil

	for i, transaction := range list.TransactionList {
		var signedTransaction proto.SignedTransaction
		if err := protobuf.Unmarshal(transaction.SignedTransactionBytes, &signedTransaction); err != nil {
			return Transaction{}, errors.Wrap(err, "error deserializing SignedTransactionBytes in TransactionFromBytes")
		}

		tx.signedTransactions = append(tx.signedTransactions, &signedTransaction)

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

		var body proto.TransactionBody
		if err := protobuf.Unmarshal(signedTransaction.GetBodyBytes(), &body); err != nil {
			return Transaction{}, errors.Wrap(err, "error deserializing BodyBytes in TransactionFromBytes")
		}

		if first == nil {
			first = &body
		}
		var transactionID TransactionID
		var nodeAccountID AccountID
		if body.GetTransactionID() != nil {
			transactionID = _TransactionIDFromProtobuf(body.GetTransactionID())
		}

		if body.GetNodeAccountID() != nil {
			nodeAccountID = *_AccountIDFromProtobuf(body.GetNodeAccountID())
		}

		found := false

		for _, id := range tx.transactionIDs {
			if id.AccountID != nil && transactionID.AccountID != nil &&
				id.AccountID._Equals(*transactionID.AccountID) &&
				id.ValidStart != nil && transactionID.ValidStart != nil &&
				id.ValidStart.Equal(*transactionID.ValidStart) {
				found = true
				break
			}
		}

		if !found {
			tx.transactionIDs = append(tx.transactionIDs, transactionID)
		}

		for _, id := range tx.nodeIDs {
			if id._Equals(nodeAccountID) {
				found = true
				break
			}
		}

		if !found {
			tx.nodeIDs = append(tx.nodeIDs, nodeAccountID)
		}
	}

	if first == nil {
		return nil, errNoTransactionInBytes
	}

	switch first.Data.(type) {
	case *proto.TransactionBody_ContractCall:
		return _ContractExecuteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ContractCreateInstance:
		return _ContractCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ContractUpdateInstance:
		return _ContractUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ContractDeleteInstance:
		return _ContractDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoAddLiveHash:
		return _LiveHashAddTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoCreateAccount:
		return _AccountCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoDelete:
		return _AccountDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoDeleteLiveHash:
		return _LiveHashDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoTransfer:
		return _TransferTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoUpdateAccount:
		return _AccountUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_FileAppend:
		return _FileAppendTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_FileCreate:
		return _FileCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_FileDelete:
		return _FileDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_FileUpdate:
		return _FileUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_SystemDelete:
		return _SystemDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_SystemUndelete:
		return _SystemUndeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_Freeze:
		return _FreezeTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ConsensusCreateTopic:
		return _TopicCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ConsensusUpdateTopic:
		return _TopicUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ConsensusDeleteTopic:
		return _TopicDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ConsensusSubmitMessage:
		return _TopicMessageSubmitTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenCreation:
		return _TokenCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenFreeze:
		return _TokenFreezeTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenUnfreeze:
		return _TokenUnfreezeTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenGrantKyc:
		return _TokenGrantKycTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenRevokeKyc:
		return _TokenRevokeKycTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenDeletion:
		return _TokenDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenUpdate:
		return _TokenUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenMint:
		return _TokenMintTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenBurn:
		return _TokenBurnTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenWipe:
		return _TokenWipeTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenAssociate:
		return _TokenAssociateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenDissociate:
		return _TokenDissociateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ScheduleCreate:
		return _ScheduleCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ScheduleSign:
		return _ScheduleSignTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ScheduleDelete:
		return _ScheduleDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenFeeScheduleUpdate:
		return TokenFeeScheduleUpdateTransactionFromProtobuf(tx, first), nil
	default:
		return Transaction{}, errFailedToDeserializeBytes
	}
}

func (transaction *Transaction) GetSignatures() (map[AccountID]map[*PublicKey][]byte, error) {
	returnMap := make(map[AccountID]map[*PublicKey][]byte, len(transaction.nodeIDs))

	if len(transaction.signedTransactions) == 0 {
		return returnMap, nil
	}

	for i, nodeID := range transaction.nodeIDs {
		inner := make(map[*PublicKey][]byte, len(transaction.signedTransactions[i].SigMap.SigPair))

		for _, sigPair := range transaction.signedTransactions[i].SigMap.SigPair {
			key, err := PublicKeyFromBytes(sigPair.PubKeyPrefix)
			if err != nil {
				return make(map[AccountID]map[*PublicKey][]byte), err
			}
			switch sigPair.Signature.(type) {
			case *proto.SignaturePair_Contract:
				inner[&key] = sigPair.GetContract()
			case *proto.SignaturePair_Ed25519:
				inner[&key] = sigPair.GetEd25519()
			case *proto.SignaturePair_RSA_3072:
				inner[&key] = sigPair.GetRSA_3072()
			case *proto.SignaturePair_ECDSA_384:
				inner[&key] = sigPair.GetECDSA_384()
			}
		}

		returnMap[nodeID] = inner
	}

	return returnMap, nil
}

// func (transaction *Transaction) AddSignature(publicKey PublicKey, signature []byte) *Transaction {
//	transaction._RequireOneNodeAccountID()
//
//	if !transaction._IsFrozen() {
//		transaction._Freeze()
//	}
//
//	if transaction._KeyAlreadySigned(publicKey) {
//		return transaction
//	}
//
//	if len(transaction.signedTransactions) == 0 {
//		return transaction
//	}
//
//	transaction.transactions = make([]*proto.Transaction, 0)
//	transaction.publicKeys = append(transaction.publicKeys, publicKey)
//	transaction.transactionSigners = append(transaction.transactionSigners, nil)
//
//	for index := 0; index < len(transaction.signedTransactions); index++ {
//		transaction.signedTransactions[index].SigMap.SigPair = append(
//			transaction.signedTransactions[index].SigMap.SigPair,
//			publicKey._ToSignaturePairProtobuf(signature),
//		)
//	}
//
//	//transaction.signedTransactions[0].SigMap.SigPair = append(transaction.signedTransactions[0].SigMap.SigPair, publicKey._ToSignaturePairProtobuf(signature))
//	return transaction
//}

func (transaction *Transaction) GetTransactionHash() ([]byte, error) {
	hashes, err := transaction.GetTransactionHashPerNode()
	if err != nil {
		return []byte{}, err
	}

	return hashes[transaction.nodeIDs[0]], nil
}

func (transaction *Transaction) GetTransactionHashPerNode() (map[AccountID][]byte, error) {
	transactionHash := make(map[AccountID][]byte)

	if !transaction._IsFrozen() {
		return transactionHash, errTransactionIsNotFrozen
	}

	err := transaction._BuildAllTransactions()
	if err != nil {
		return transactionHash, err
	}

	for i, node := range transaction.nodeIDs {
		hash := sha512.New384()
		_, err := hash.Write(transaction.transactions[i].GetSignedTransactionBytes())
		if err != nil {
			return transactionHash, err
		}

		finalHash := hash.Sum(nil)

		transactionHash[node] = finalHash
	}

	return transactionHash, nil
}

func (transaction *Transaction) _InitFee(client *Client) {
	if client != nil && transaction.transactionFee == 0 {
		transaction.SetMaxTransactionFee(client.maxTransactionFee)
	}
}

func (transaction *Transaction) _InitTransactionID(client *Client) error {
	if len(transaction.transactionIDs) == 0 {
		if client != nil {
			if client.operator != nil {
				transaction.SetTransactionID(TransactionIDGenerate(client.operator.accountID))
			} else {
				return errNoClientOrTransactionID
			}
		} else {
			return errNoClientOrTransactionID
		}
	}

	transaction.transactionID = transaction.GetTransactionID()
	return nil
}

func (transaction *Transaction) _IsFrozen() bool {
	return len(transaction.signedTransactions) > 0
}

func (transaction *Transaction) _RequireNotFrozen() {
	if transaction._IsFrozen() {
		transaction.freezeError = errTransactionIsFrozen
	}
}

func (transaction *Transaction) _RequireOneNodeAccountID() {
	if len(transaction.nodeIDs) != 1 {
		panic("Transaction has more than one _Node ID set")
	}
}

func _TransactionFreezeWith(
	transaction *Transaction,
	client *Client,
	body *proto.TransactionBody,
) error {
	if len(transaction.nodeIDs) == 0 {
		if client != nil {
			transaction.nodeIDs = client.network._GetNodeAccountIDsForExecute()
		} else {
			return errNoClientOrTransactionIDOrNodeId
		}
	}

	for _, nodeAccountID := range transaction.nodeIDs {
		body.NodeAccountID = nodeAccountID._ToProtobuf()
		bodyBytes, err := protobuf.Marshal(body)
		if err != nil {
			// This should be unreachable
			// From the documentation this appears to only be possible if there are missing proto types
			panic(err)
		}

		transaction.signedTransactions = append(transaction.signedTransactions, &proto.SignedTransaction{
			BodyBytes: bodyBytes,
			SigMap: &proto.SignatureMap{
				SigPair: make([]*proto.SignaturePair, 0),
			},
		})
	}

	return nil
}

func (transaction *Transaction) _SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) {
	transaction.transactions = make([]*proto.Transaction, 0)
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, signer)
}

func (transaction *Transaction) _KeyAlreadySigned(
	pk PublicKey,
) bool {
	for _, key := range transaction.publicKeys {
		if key.String() == pk.String() {
			return true
		}
	}

	return false
}

func _TransactionShouldRetry(_ _Request, response _Response) _ExecutionState {
	switch Status(response.transaction.NodeTransactionPrecheckCode) {
	case StatusPlatformTransactionNotCreated, StatusBusy:
		return executionStateRetry
	case StatusOk:
		return executionStateFinished
	}

	return executionStateError
}

func _TransactionMakeRequest(request _Request) _ProtoRequest {
	index := len(request.transaction.nodeIDs)*request.transaction.nextTransactionIndex + request.transaction.nextNodeIndex
	_ = request.transaction._BuildTransaction(index)

	return _ProtoRequest{
		transaction: request.transaction.transactions[index],
	}
}

func _TransactionAdvanceRequest(request _Request) {
	length := len(request.transaction.nodeIDs)
	currentIndex := request.transaction.nextNodeIndex
	request.transaction.nextNodeIndex = (currentIndex + 1) % length
}

func _TransactionGetNodeAccountID(request _Request) AccountID {
	return request.transaction.nodeIDs[request.transaction.nextNodeIndex]
}

func _TransactionMapStatusError(
	request _Request,
	response _Response,
) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.transaction.NodeTransactionPrecheckCode),
		TxID:   request.transaction.GetTransactionID(),
	}
}

func _TransactionMapResponse(request _Request, _ _Response, nodeID AccountID, protoRequest _ProtoRequest) (_IntermediateResponse, error) {
	hash := sha512.New384()
	_, err := hash.Write(protoRequest.transaction.SignedTransactionBytes)
	if err != nil {
		return _IntermediateResponse{}, err
	}

	index := request.transaction.nextTransactionIndex
	request.transaction.nextTransactionIndex = (index + 1) % len(request.transaction.transactionIDs)

	return _IntermediateResponse{
		transaction: TransactionResponse{
			NodeID:        nodeID,
			TransactionID: request.transaction.transactionIDs[index],
			Hash:          hash.Sum(nil),
		},
	}, nil
}

func (transaction *Transaction) String() string {
	return fmt.Sprintf("%+v", transaction.signedTransactions[0])
}

func (transaction *Transaction) ToBytes() ([]byte, error) {
	if !transaction._IsFrozen() {
		return make([]byte, 0), errTransactionIsNotFrozen
	}

	err := transaction._BuildAllTransactions()
	if err != nil {
		return make([]byte, 0), err
	}

	pbTransactionList, lastError := protobuf.Marshal(&proto.TransactionList{
		TransactionList: transaction.transactions,
	})

	if lastError != nil {
		return make([]byte, 0), errors.Wrap(err, "error serializing transaction list")
	}

	return pbTransactionList, nil
}

func (transaction *Transaction) _SignTransaction(index int) {
	if len(transaction.signedTransactions[index].SigMap.SigPair) != 0 {
		for i, key := range transaction.publicKeys {
			if transaction.transactionSigners[i] != nil && bytes.Equal(transaction.signedTransactions[index].SigMap.SigPair[0].PubKeyPrefix, key.keyData) {
				return
			}
		}
	}

	bodyBytes := transaction.signedTransactions[index].GetBodyBytes()

	for i := 0; i < len(transaction.publicKeys); i++ {
		publicKey := transaction.publicKeys[i]
		signer := transaction.transactionSigners[i]

		if signer == nil {
			continue
		}

		transaction.signedTransactions[index].SigMap.SigPair = append(transaction.signedTransactions[index].SigMap.SigPair, publicKey._ToSignaturePairProtobuf(signer(bodyBytes)))
	}
}

func (transaction *Transaction) _BuildAllTransactions() error {
	for i := 0; i < len(transaction.signedTransactions); i++ {
		err := transaction._BuildTransaction(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (transaction *Transaction) _BuildTransaction(index int) error {
	if len(transaction.transactions) < index {
		for i := len(transaction.transactions); i < index; i++ {
			transaction.transactions = append(transaction.transactions, nil)
		}
	} else if len(transaction.transactions) > index &&
		transaction.transactions[index] != nil &&
		transaction.transactions[index].SignedTransactionBytes != nil {
		return nil
	}

	transaction._SignTransaction(index)

	data, err := protobuf.Marshal(transaction.signedTransactions[index])
	if err != nil {
		return errors.Wrap(err, "failed to serialize transactions for building")
	}

	transaction.transactions = append(transaction.transactions, &proto.Transaction{
		SignedTransactionBytes: data,
	})

	return nil
}

//
// Shared
//

func (transaction *Transaction) GetMaxTransactionFee() Hbar {
	return HbarFromTinybar(int64(transaction.transactionFee))
}

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction *Transaction) SetMaxTransactionFee(fee Hbar) *Transaction {
	transaction.transactionFee = uint64(fee.AsTinybar())
	return transaction
}

func (transaction *Transaction) GetTransactionMemo() string {
	return transaction.memo
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction *Transaction) SetTransactionMemo(memo string) *Transaction {
	transaction.memo = memo
	return transaction
}

func (transaction *Transaction) GetTransactionValidDuration() time.Duration {
	if transaction.transactionValidDuration != nil {
		return *transaction.transactionValidDuration
	}

	return 0
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction *Transaction) SetTransactionValidDuration(duration time.Duration) *Transaction {
	transaction.transactionValidDuration = &duration
	return transaction
}

func (transaction *Transaction) GetTransactionID() TransactionID {
	if len(transaction.transactionIDs) > 0 {
		return transaction.transactionIDs[transaction.nextTransactionIndex]
	}

	return TransactionID{}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction *Transaction) SetTransactionID(transactionID TransactionID) *Transaction {
	transaction.transactionIDs = []TransactionID{transactionID}
	return transaction
}

func (transaction *Transaction) GetNodeAccountIDs() []AccountID {
	if transaction.nodeIDs != nil {
		return transaction.nodeIDs
	}

	return make([]AccountID, 0)
}

// SetNodeAccountID sets the _Node AccountID for this Transaction.
func (transaction *Transaction) SetNodeAccountIDs(nodeID []AccountID) *Transaction {
	if transaction.nodeIDs == nil {
		transaction.nodeIDs = make([]AccountID, 0)
	}
	transaction.nodeIDs = append(transaction.nodeIDs, nodeID...)
	return transaction
}

func (transaction *Transaction) GetMaxRetry() int {
	return transaction.maxRetry
}

func (transaction *Transaction) SetMaxRetry(count int) *Transaction {
	transaction.maxRetry = count
	return transaction
}

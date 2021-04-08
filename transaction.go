package hedera

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"github.com/pkg/errors"

	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Transaction contains the protobuf of a prepared transaction which can be signed and executed.

type ITransaction interface {
	constructScheduleProtobuf() *proto.SchedulableTransactionBody
}

type Transaction struct {
	pbBody *proto.TransactionBody

	nextNodeIndex        int
	nextTransactionIndex int
	maxRetry             int

	transactionIDs     []TransactionID
	transactions       []*proto.Transaction
	signedTransactions []*proto.SignedTransaction
	nodeIDs            []AccountID

	freezeError error
}

func newTransaction() Transaction {
	return Transaction{
		pbBody: &proto.TransactionBody{
			TransactionValidDuration: durationToProtobuf(120 * time.Second),
		},
		nextNodeIndex:        0,
		nextTransactionIndex: 0,
		maxRetry:             10,
		transactionIDs:       make([]TransactionID, 0),
		transactions:         make([]*proto.Transaction, 0),
		signedTransactions:   make([]*proto.SignedTransaction, 0),
		nodeIDs:              make([]AccountID, 0),
		freezeError:          nil,
	}
}

func TransactionFromBytes(data []byte) (interface{}, error) {
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
		nodeIDs:              make([]AccountID, 0),
	}

	var first *proto.TransactionBody = nil

	for _, transaction := range list.TransactionList {
		var signedTransaction proto.SignedTransaction
		if err := protobuf.Unmarshal(transaction.SignedTransactionBytes, &signedTransaction); err != nil {
			return Transaction{}, errors.Wrap(err, "error deserializing SignedTransactionBytes in TransactionFromBytes")
		}

		tx.signedTransactions = append(tx.signedTransactions, &signedTransaction)

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
			transactionID = transactionIDFromProtobuf(body.GetTransactionID())
		}

		if body.GetNodeAccountID() != nil {
			nodeAccountID = accountIDFromProtobuf(body.GetNodeAccountID())
		}

		found := false

		for _, id := range tx.transactionIDs {
			if id.AccountID == transactionID.AccountID && id.ValidStart == transactionID.ValidStart {
				found = true
				break
			}
		}

		if !found {
			tx.transactionIDs = append(tx.transactionIDs, transactionID)
		}

		for _, id := range tx.nodeIDs {
			if id == nodeAccountID {
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
		return contractExecuteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ContractCreateInstance:
		return contractCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ContractUpdateInstance:
		return contractUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ContractDeleteInstance:
		return contractDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoAddLiveHash:
		return liveHashAddTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoCreateAccount:
		return accountCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoDelete:
		return accountDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoDeleteLiveHash:
		return liveHashDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoTransfer:
		return transferTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_CryptoUpdateAccount:
		return accountUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_FileAppend:
		return fileAppendTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_FileCreate:
		return fileCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_FileDelete:
		return fileDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_FileUpdate:
		return fileUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_SystemDelete:
		return systemDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_SystemUndelete:
		return systemUndeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_Freeze:
		return freezeTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ConsensusCreateTopic:
		return topicCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ConsensusUpdateTopic:
		return topicUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ConsensusDeleteTopic:
		return topicDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ConsensusSubmitMessage:
		return topicMessageSubmitTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenCreation:
		return tokenCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenFreeze:
		return tokenFreezeTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenUnfreeze:
		return tokenUnfreezeTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenGrantKyc:
		return tokenGrantKycTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenRevokeKyc:
		return tokenRevokeKycTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenDeletion:
		return tokenDeleteTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenUpdate:
		return tokenUpdateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenMint:
		return tokenMintTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenBurn:
		return tokenBurnTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenWipe:
		return tokenWipeTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenAssociate:
		return tokenAssociateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_TokenDissociate:
		return tokenDissociateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ScheduleCreate:
		return scheduleCreateTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ScheduleSign:
		return scheduleSignTransactionFromProtobuf(tx, first), nil
	case *proto.TransactionBody_ScheduleDelete:
		return scheduleDeleteTransactionFromProtobuf(tx, first), nil
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
				return make(map[AccountID]map[*PublicKey][]byte, 0), err
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

func (transaction *Transaction) AddSignature(publicKey PublicKey, signature []byte) *Transaction {
	transaction.requireOneNodeAccountID()

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	if len(transaction.signedTransactions) == 0 {
		return transaction
	}

	transaction.transactions = make([]*proto.Transaction, 0)

	for index := 0; index < len(transaction.signedTransactions); index++ {
		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	//transaction.signedTransactions[0].SigMap.SigPair = append(transaction.signedTransactions[0].SigMap.SigPair, publicKey.toSignaturePairProtobuf(signature))
	return transaction
}

func (transaction *Transaction) GetTransactionHash() ([]byte, error) {
	hashes, err := transaction.GetTransactionHashPerNode()
	if err != nil {
		return []byte{}, err
	}

	return hashes[transaction.nodeIDs[0]], nil
}

func (transaction *Transaction) GetTransactionHashPerNode() (map[AccountID][]byte, error) {
	transactionHash := make(map[AccountID][]byte)

	if !transaction.isFrozen() {
		return transactionHash, errTransactionIsNotFrozen
	}

	err := transaction.buildTransactions(len(transaction.signedTransactions))
	if err != nil {
		return transactionHash, err
	}

	for i, node := range transaction.nodeIDs {
		hash := sha512.New384()
		_, err := hash.Write(transaction.transactions[i].SignedTransactionBytes)
		if err != nil {
			return transactionHash, err
		}

		transactionHash[node] = []byte(hex.EncodeToString(hash.Sum(nil)))
	}

	return transactionHash, nil
}

func (transaction *Transaction) initFee(client *Client) {
	if client != nil && transaction.pbBody.TransactionFee == 0 {
		transaction.SetMaxTransactionFee(client.maxTransactionFee)
	}
}

func (transaction *Transaction) initTransactionID(client *Client) error {
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

	transaction.pbBody.TransactionID = transaction.GetTransactionID().toProtobuf()
	return nil
}

func (transaction *Transaction) isFrozen() bool {
	return len(transaction.signedTransactions) > 0
}

func (transaction *Transaction) requireNotFrozen() {
	if transaction.isFrozen() {
		transaction.freezeError = errTransactionIsFrozen
	}
}

func (transaction *Transaction) requireOneNodeAccountID() {
	if len(transaction.nodeIDs) != 1 {
		panic("Transaction has more than one node ID set")
	}
}

func transaction_freezeWith(
	transaction *Transaction,
	client *Client,
) error {
	if len(transaction.nodeIDs) == 0 {
		if client != nil {
			transaction.nodeIDs = client.network.getNodeAccountIDsForExecute()
		} else {
			return errNoClientOrTransactionIDOrNodeId
		}
	}

	for _, nodeAccountID := range transaction.nodeIDs {
		transaction.pbBody.NodeAccountID = nodeAccountID.toProtobuf()
		bodyBytes, err := protobuf.Marshal(transaction.pbBody)
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

func (transaction *Transaction) keyAlreadySigned(
	pk PublicKey,
) bool {
	if len(transaction.signedTransactions) > 0 {
		for _, pair := range transaction.signedTransactions[0].SigMap.SigPair {
			if bytes.HasPrefix(pk.keyData, pair.PubKeyPrefix) {
				return true
			}
		}
	}

	return false
}

func transaction_shouldRetry(_ request, response response) executionState {
	switch Status(response.transaction.NodeTransactionPrecheckCode) {
	case StatusPlatformTransactionNotCreated, StatusBusy:
		return executionStateRetry
	case StatusOk:
		return executionStateFinished
	}

	return executionStateError
}

func transaction_makeRequest(request request) protoRequest {
	index := len(request.transaction.nodeIDs)*request.transaction.nextTransactionIndex + request.transaction.nextNodeIndex
	_ = request.transaction.buildTransactions(index + 1)

	return protoRequest{
		transaction: request.transaction.transactions[index],
	}
}

func transaction_advanceRequest(request request) {
	length := len(request.transaction.nodeIDs)
	currentIndex := request.transaction.nextNodeIndex
	request.transaction.nextNodeIndex = (currentIndex + 1) % length
}

func transaction_getNodeAccountID(request request) AccountID {
	return request.transaction.nodeIDs[request.transaction.nextNodeIndex]
}

func transaction_mapStatusError(
	request request,
	response response,
) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.transaction.NodeTransactionPrecheckCode),
		TxID:   request.transaction.GetTransactionID(),
	}
}

func transaction_mapResponse(request request, _ response, nodeID AccountID, protoRequest protoRequest) (intermediateResponse, error) {
	hash := sha512.New384()
	_, err := hash.Write(protoRequest.transaction.SignedTransactionBytes)
	if err != nil {
		return intermediateResponse{}, err
	}

	index := request.transaction.nextTransactionIndex
	request.transaction.nextTransactionIndex = (index + 1) % len(request.transaction.transactionIDs)

	return intermediateResponse{
		transaction: TransactionResponse{
			NodeID:        nodeID,
			TransactionID: request.transaction.transactionIDs[index],
			Hash:          hash.Sum(nil),
		},
	}, nil
}

func (transaction *Transaction) String() string {
	return protobuf.MarshalTextString(transaction.signedTransactions[0])
}

func (transaction *Transaction) ToBytes() ([]byte, error) {
	if !transaction.isFrozen() {
		return make([]byte, 0), errTransactionIsNotFrozen
	}

	err := transaction.buildTransactions(len(transaction.signedTransactions))
	if err != nil {
		return make([]byte, 0), err
	}

	pbTransactionList, lastError := protobuf.Marshal(&proto.TransactionList{
		TransactionList: transaction.transactions,
	})

	if lastError != nil {
		return make([]byte, 0), errors.Wrap(err, "error serializing transaction list")
	} else {
		return pbTransactionList, nil
	}

}

func (transaction *Transaction) buildTransactions(untilIndex int) error {
	for i := len(transaction.transactions); i < untilIndex; i++ {
		data, err := protobuf.Marshal(transaction.signedTransactions[i])
		if err != nil {
			return errors.Wrap(err, "failed to serialize transactions for building")
		}

		transaction.transactions = append(transaction.transactions, &proto.Transaction{
			SignedTransactionBytes: data,
		})
	}

	return nil
}

//
// Shared
//

func (transaction *Transaction) GetMaxTransactionFee() Hbar {
	return HbarFromTinybar(int64(transaction.pbBody.TransactionFee))
}

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction *Transaction) SetMaxTransactionFee(fee Hbar) *Transaction {
	transaction.pbBody.TransactionFee = uint64(fee.AsTinybar())
	return transaction
}

func (transaction *Transaction) GetTransactionMemo() string {
	return transaction.pbBody.Memo
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction *Transaction) SetTransactionMemo(memo string) *Transaction {
	transaction.pbBody.Memo = memo
	return transaction
}

func (transaction *Transaction) GetTransactionValidDuration() time.Duration {
	if transaction.pbBody.TransactionValidDuration != nil {
		return durationFromProtobuf(transaction.pbBody.TransactionValidDuration)
	} else {
		return 0
	}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction *Transaction) SetTransactionValidDuration(duration time.Duration) *Transaction {
	transaction.pbBody.TransactionValidDuration = durationToProtobuf(duration)
	return transaction
}

func (transaction *Transaction) GetTransactionID() TransactionID {
	if len(transaction.transactionIDs) > 0 {
		return transaction.transactionIDs[transaction.nextTransactionIndex]
	} else {
		return TransactionID{}
	}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction *Transaction) SetTransactionID(transactionID TransactionID) *Transaction {
	transaction.transactionIDs = []TransactionID{transactionID}
	return transaction
}

func (transaction *Transaction) GetNodeAccountIDs() []AccountID {
	if transaction.nodeIDs != nil {
		return transaction.nodeIDs
	} else {
		return make([]AccountID, 0)
	}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
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

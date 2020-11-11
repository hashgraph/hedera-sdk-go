package hedera

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// Transaction contains the protobuf of a prepared transaction which can be signed and executed.
type Transaction struct {
	pbBody *proto.TransactionBody

	id TransactionID

	nextTransactionIndex int

	transactions []*proto.Transaction
	signatures   []*proto.SignatureMap
	nodeIDs      []AccountID
}

func newTransaction() Transaction {
	return Transaction{
		pbBody: &proto.TransactionBody{
			TransactionValidDuration: durationToProtobuf(120 * time.Second),
		},
		id:                   TransactionID{},
		nextTransactionIndex: 0,
		transactions:         make([]*proto.Transaction, 0),
		signatures:           make([]*proto.SignatureMap, 0),
		nodeIDs:              make([]AccountID, 0),
	}
}

func transactionFromProtobuf(transactions map[TransactionID]map[AccountID]*proto.Transaction, pb *proto.TransactionBody) Transaction {
	tx := Transaction{
		pbBody:               pb,
		id:                   transactionIDFromProtobuf(pb.TransactionID),
		nextTransactionIndex: 0,
		transactions:         make([]*proto.Transaction, 0),
		signatures:           make([]*proto.SignatureMap, 0),
		nodeIDs:              make([]AccountID, 0),
	}

	var protoTxs map[AccountID]*proto.Transaction
	for _, m := range transactions {
		protoTxs = m
		break
	}

	for nodeAccountID, protoTx := range protoTxs {
		tx.nodeIDs = append(tx.nodeIDs, nodeAccountID)
		tx.transactions = append(tx.transactions, protoTx)
		tx.signatures = append(tx.signatures, protoTx.GetSigMap())
	}

	return tx
}

func TransactionFromBytes(bytes []byte) (interface{}, error) {
	transactions := make(map[TransactionID]map[AccountID]*proto.Transaction)
	buf := protobuf.NewBuffer(bytes)
	var first *proto.TransactionBody = nil

	for {
		tx := proto.Transaction{}
		if err := buf.Unmarshal(&tx); err != nil {
			break
		}

		var txBody proto.TransactionBody
		if err := protobuf.Unmarshal(tx.GetBodyBytes(), &txBody); err != nil {
			return Transaction{}, err
		}

		if first == nil {
			first = &txBody
		}

		transactionID := transactionIDFromProtobuf(txBody.TransactionID)
		nodeAccountID := accountIDFromProtobuf(txBody.NodeAccountID)

		if _, ok := transactions[transactionID]; !ok {
			transactions[transactionID] = make(map[AccountID]*proto.Transaction)
		}

		transactions[transactionID][nodeAccountID] = &tx
	}

	if first == nil {
		return nil, errNoTransactionInBytes
	}

	switch first.Data.(type) {
	case *proto.TransactionBody_ContractCall:
		return contractExecuteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_ContractCreateInstance:
		return contractCreateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_ContractUpdateInstance:
		return contractUpdateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_ContractDeleteInstance:
		return contractDeleteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_CryptoAddLiveHash:
		return liveHashAddTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_CryptoCreateAccount:
		return accountCreateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_CryptoDelete:
		return accountDeleteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_CryptoDeleteLiveHash:
		return liveHashDeleteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_CryptoTransfer:
		return transferTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_CryptoUpdateAccount:
		return accountUpdateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_FileAppend:
		return fileAppendTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_FileCreate:
		return fileCreateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_FileDelete:
		return fileDeleteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_FileUpdate:
		return fileUpdateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_SystemDelete:
		return systemDeleteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_SystemUndelete:
		return systemUndeleteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_Freeze:
		return freezeTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_ConsensusCreateTopic:
		return topicCreateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_ConsensusUpdateTopic:
		return topicUpdateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_ConsensusDeleteTopic:
		return topicDeleteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_ConsensusSubmitMessage:
		return topicMessageSubmitTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenCreation:
		return tokenCreateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenFreeze:
		return tokenFreezeTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenUnfreeze:
		return tokenUnfreezeTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenGrantKyc:
		return tokenGrantKycTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenRevokeKyc:
		return tokenRevokeKycTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenDeletion:
		return tokenDeleteTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenUpdate:
		return tokenUpdateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenMint:
		return tokenMintTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenBurn:
		return tokenBurnTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenWipe:
		return tokenWipeTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenAssociate:
		return tokenAssociateTransactionFromProtobuf(transactions, first), nil
	case *proto.TransactionBody_TokenDissociate:
		return tokenDissociateTransactionFromProtobuf(transactions, first), nil
	default:
		return Transaction{}, errFailedToDeserializeBytes
	}
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

	if len(transaction.transactions) == 0 {
		return transactionHash, errTransactionIsNotFrozen
	}

	for i, node := range transaction.nodeIDs {
		data, err := protobuf.Marshal(transaction.transactions[i])
		if err != nil {
			// This should be unreachable
			// From the documentation this appears to only be possible if there are missing proto types
			return transactionHash, err
		}

		hash := sha512.New384()
		_, err = hash.Write(data)
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
	if transaction.pbBody.TransactionID == nil {
		if client.operator != nil {
			transaction.id = TransactionIDGenerate(client.operator.accountID)
			transaction.SetTransactionID(transaction.id)
		} else {
			return errNoClientOrTransactionID
		}
	}

	return nil
}

func (transaction *Transaction) isFrozen() bool {
	return len(transaction.transactions) > 0
}

func (transaction *Transaction) requireNotFrozen() {
	if transaction.isFrozen() {
		panic("Transaction is immutable; it has at least one signature or has been explicitly frozen\"")
	}
}

func transaction_freezeWith(
	transaction *Transaction,
	client *Client,
) error {
	if len(transaction.nodeIDs) == 0 {
		if client == nil {
			return errNoClientOrTransactionIDOrNodeId
		} else {
			transaction.nodeIDs = client.network.getNodeAccountIDsForExecute()
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

		sigmap := proto.SignatureMap{
			SigPair: make([]*proto.SignaturePair, 0),
		}
		transaction.signatures = append(transaction.signatures, &sigmap)
		transaction.transactions = append(transaction.transactions, &proto.Transaction{
			BodyBytes: bodyBytes,
			SigMap:    &sigmap,
		})
	}

	return nil
}

func (transaction *Transaction) keyAlreadySigned(
	pk PublicKey,
) bool {
	if len(transaction.signatures) > 0 {
		for _, pair := range transaction.signatures[0].SigPair {
			if bytes.HasPrefix(pk.keyData, pair.PubKeyPrefix) {
				return true
			}
		}
	}

	return false
}

func transaction_shouldRetry(status Status, _ response) bool {
	return status == StatusBusy
}

func transaction_makeRequest(request request) protoRequest {
	return protoRequest{
		transaction: request.transaction.transactions[request.transaction.nextTransactionIndex],
	}
}

func transaction_advanceRequest(request request) {
	length := len(request.transaction.transactions)
	currentIndex := request.transaction.nextTransactionIndex
	request.transaction.nextTransactionIndex = (currentIndex + 1) % length
}

func transaction_getNodeAccountID(request request) AccountID {
	return request.transaction.nodeIDs[request.transaction.nextTransactionIndex]
}

func transaction_mapResponseStatus(
	_ request,
	response response,
) Status {
	return Status(response.transaction.NodeTransactionPrecheckCode)
}

func transaction_mapResponse(request request, _ response, nodeID AccountID, protoRequest protoRequest) (intermediateResponse, error) {
	hash, err := protobuf.Marshal(protoRequest.transaction)
	if err != nil {
		return intermediateResponse{}, err
	}

	return intermediateResponse{
		transaction: TransactionResponse{
			NodeID:        nodeID,
			TransactionID: request.transaction.id,
			Hash:          hash,
		},
	}, nil
}

func (transaction *Transaction) String() string {
	return protobuf.MarshalTextString(transaction.transactions[0]) +
		protobuf.MarshalTextString(transaction.body())
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (transaction *Transaction) MarshalBinary() ([]byte, error) {
	return protobuf.Marshal(transaction.transactions[0])
}

func (transaction *Transaction) ToBytes() ([]byte, error) {
	return transaction.MarshalBinary()
}

// The protobuf stores the transaction body as raw bytes so we need to first
// decode what we have to inspect the Kind, TransactionID, and the NodeAccountID so we know how to
// properly execute it
func (transaction *Transaction) body() *proto.TransactionBody {
	transactionBody := new(proto.TransactionBody)
	err := protobuf.Unmarshal(transaction.transactions[0].GetBodyBytes(), transactionBody)
	if err != nil {
		// The bodyBytes inside of the transaction at this point have been verified and this should be impossible
		panic(err)
	}

	return transactionBody
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
	if transaction.pbBody.TransactionID != nil {
		return transactionIDFromProtobuf(transaction.pbBody.TransactionID)
	} else {
		return TransactionID{}
	}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction *Transaction) SetTransactionID(transactionID TransactionID) *Transaction {
	transaction.pbBody.TransactionID = transactionID.toProtobuf()
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

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

	// unfortunately; this is required to prevent setting the max TXFee if it is purposely set to 0
	// (for example, when .GetCost() is called)
	noTXFee bool

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
		noTXFee:              false,
		nextTransactionIndex: 0,
		transactions:         make([]*proto.Transaction, 0),
		signatures:           make([]*proto.SignatureMap, 0),
		nodeIDs:              make([]AccountID, 0),
	}
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (transaction *Transaction) UnmarshalBinary(txBytes []byte) error {
	transaction.transactions = make([]*proto.Transaction, 0)
	transaction.transactions = append(transaction.transactions, &proto.Transaction{})
	if err := protobuf.Unmarshal(txBytes, transaction.transactions[0]); err != nil {
		return err
	}

	var txBody proto.TransactionBody
	if err := protobuf.Unmarshal(transaction.transactions[0].GetBodyBytes(), &txBody); err != nil {
		return err
	}

	transaction.id = transactionIDFromProtobuf(txBody.TransactionID)

	return nil
}

func TransactionFromBytes(bytes []byte) Transaction {
	tx := Transaction{}
	(&tx).UnmarshalBinary(bytes)
	return tx
}

func (transaction *Transaction) ToBytes() ([]byte, error) {
	buf := protobuf.NewBuffer(make([]byte, 0))

	for _, tx := range transaction.transactions {
		err := buf.Marshal(tx)
		if err != nil {
			return buf.Bytes(), err
		}
	}

	return buf.Bytes(), nil
}

func (transaction *Transaction) GetTransactionHash() (map[AccountID][]byte, error) {
	transactionHash := make(map[AccountID][]byte)

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

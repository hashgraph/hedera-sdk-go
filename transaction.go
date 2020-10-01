package hedera

import (
	"bytes"
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

	transactions []proto.Transaction
	signatures   []proto.SignatureMap
	nodeIDs      []AccountID
}

func newTransaction() Transaction {
	return Transaction{
		pbBody: &proto.TransactionBody{
			TransactionValidDuration: durationToProto(120 * time.Second),
		},
		id:                   TransactionID{},
		noTXFee:              false,
		nextTransactionIndex: 0,
		transactions:         make([]proto.Transaction, 0),
		signatures:           make([]proto.SignatureMap, 0),
		nodeIDs:              make([]AccountID, 0),
	}
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (transaction *Transaction) UnmarshalBinary(txBytes []byte) error {
	transaction.transactions = make([]proto.Transaction, 0)
	transaction.transactions = append(transaction.transactions, proto.Transaction{})
	if err := protobuf.Unmarshal(txBytes, &transaction.transactions[0]); err != nil {
		return err
	}

	var txBody proto.TransactionBody
	if err := protobuf.Unmarshal(transaction.transactions[0].GetBodyBytes(), &txBody); err != nil {
		return err
	}

	transaction.id = transactionIDFromProto(txBody.TransactionID)

	return nil
}

func TransactionFromBytes(bytes []byte) Transaction {
	tx := Transaction{}
	(&tx).UnmarshalBinary(bytes)
	return tx
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *Transaction) Sign(
	privateKey PrivateKey,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
) *Transaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign, isFrozen, freezeWith)
}

func (transaction *Transaction) SignWithOperator(
	client *Client,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
) (*Transaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !isFrozen() {
		freezeWith(client)
	}

	return transaction.SignWith(client.operator.publicKey, client.operator.signer, isFrozen, freezeWith), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *Transaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
) *Transaction {
	if !isFrozen() {
		freezeWith(nil)
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index, tx := range transaction.transactions {
		signature := signer(tx.GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

func (transaction *Transaction) freezeWith(
	client *Client,
	isFrozen func() bool,
	onFreeze func(pbBody *proto.TransactionBody) bool,
) error {
	if client != nil {
		if transaction.pbBody.TransactionFee == 0 {
			transaction.SetMaxTransactionFee(client.maxTransactionFee)
		}

		if transaction.pbBody.TransactionID == nil {
			if client.operator != nil {
				transaction.SetTransactionID(NewTransactionID(client.operator.accountID))
			} else {
				return errNoClientOrTransactionID
			}
		}
	}

	if !onFreeze(transaction.pbBody) {
		return nil
	}

	if transaction.pbBody.TransactionID != nil && transaction.pbBody.NodeAccountID != nil {
		bodyBytes, err := protobuf.Marshal(transaction.pbBody)
		if err != nil {
			// This should be unreachable
			// From the documentation this appears to only be possible if there are missing proto types
			panic(err)
		}


		transaction.transactions = append(transaction.transactions, proto.Transaction{
			BodyBytes: bodyBytes,
		})

		return nil
	}

	if transaction.pbBody.TransactionID != nil && len(transaction.nodeIDs) > 0 {
		for _, id := range transaction.nodeIDs {
			transaction.pbBody.NodeAccountID = id.toProtobuf()
			bodyBytes, err := protobuf.Marshal(transaction.pbBody)
			if err != nil {
				// This should be unreachable
				// From the documentation this appears to only be possible if there are missing proto types
				panic(err)
			}

			transaction.signatures = append(transaction.signatures, proto.SignatureMap{})
			transaction.transactions = append(transaction.transactions, proto.Transaction{
				BodyBytes: bodyBytes,
			})
		}

		return nil
	}

    println("kkkkkkkkkkkkkkkk")
	if client != nil && transaction.pbBody.TransactionID != nil {
		size := client.getNumberOfNodesForTransaction()
        println("[freezeWith] size", size)

		for index := 0; index < size; index++ {
			node := client.getNextNode()
            println("[freezeWith] node", node.String())

			transaction.nodeIDs = append(transaction.nodeIDs, node)

			transaction.pbBody.NodeAccountID = node.toProtobuf()
			bodyBytes, err := protobuf.Marshal(transaction.pbBody)
			if err != nil {
				// This should be unreachable
				// From the documentation this appears to only be possible if there are missing proto types
				panic(err)
			}

			transaction.signatures = append(transaction.signatures, proto.SignatureMap{})
			transaction.transactions = append(transaction.transactions, proto.Transaction{
				BodyBytes: bodyBytes,
			})
            println("[freezeWith] transactions", len(transaction.transactions))
		}

        println("[freezeWith::end] transactions", len(transaction.transactions))
		return nil
	}

	return errNoClientOrTransactionIDOrNodeId
}

func defaultIsFrozen(transaction *Transaction) bool {
	return len(transaction.transactions) > 0
}

func (transaction *Transaction) requireNotFrozen(isFrozen func() bool) error {
	if isFrozen() {
		return errTransactionIsFrozen
	}

	return nil
}

func (transaction *Transaction) keyAlreadySigned(pk PublicKey) bool {
	if len(transaction.signatures) > 0 {
		for _, pair := range transaction.signatures[0].SigPair {
			if bytes.HasPrefix(pk.keyData, pair.PubKeyPrefix) {
				return true
			}
		}
	}

	return false
}

func (transaction *Transaction) shouldRetry(status Status) bool {
	return status == StatusBusy
}

func (transaction *Transaction) makeRequest() request {
	println("ffffffffffffffffF")
	println(len(transaction.transactions))
	println(transaction.nextTransactionIndex)
	return request{
		transaction: &transaction.transactions[transaction.nextTransactionIndex],
	}
}

func (transaction *Transaction) advanceRequest() {
	transaction.nextTransactionIndex++
}

func (transaction *Transaction) getNodeId(client *Client) AccountID {
	node := transaction.GetNodeID()

	if node.Shard == 0 && node.Realm == 0 && node.Account == 0 {
		return client.getNextNode()
	} else {
		return node
	}
}

func (transaction *Transaction) mapResponseStatus(response response) Status {
	return Status(response.transaction.NodeTransactionPrecheckCode)
}

func (transaction *Transaction) mapResponse(_ response, nodeID AccountID, request request) (intermediateResponse, error) {
	hash, err := protobuf.Marshal(request.transaction)
	if err != nil {
		return intermediateResponse{}, err
	}

	return intermediateResponse{
		transaction: TransactionResponse{
			NodeID:        nodeID,
			TransactionID: transaction.id,
			Hash:          hash,
		},
	}, nil
}

func (transaction *Transaction) isFrozen() bool {
	return len(transaction.transactions) > 0
}

// Execute executes the Transaction with the provided client
func (transaction *Transaction) execute(
	client *Client,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
	getMethod func(*channel) method,
) (TransactionResponse, error) {
	var isFrozen_ func() bool
	if isFrozen != nil {
		isFrozen_ = isFrozen
	} else {
		isFrozen_ = transaction.isFrozen
	}

	if !isFrozen_() {
		freezeWith(client)
        println("[Transaction.execute] transactions", len(transaction.transactions))
	}

	operatorID := client.GetOperatorID()
	transactionID := transaction.id

	if (operatorID.Shard != 0 && operatorID.Realm != 0 && operatorID.Account != 0) && (operatorID.Shard != transactionID.AccountID.Shard && operatorID.Realm != transactionID.AccountID.Realm && operatorID.Account != transactionID.AccountID.Account) {
		transaction.SignWith(client.GetOperatorKey(), client.operator.signer, isFrozen_, freezeWith)
	}

	response, err := execute(
		client,
		transaction.isFrozen,
		freezeWith,
		transaction.shouldRetry,
		transaction.makeRequest,
		transaction.advanceRequest,
		transaction.getNodeId,
		getMethod,
		transaction.mapResponseStatus,
		transaction.mapResponse,
	)

	return response.transaction, err
}

func (transaction *Transaction) String() string {
	return protobuf.MarshalTextString(&transaction.transactions[0]) +
		protobuf.MarshalTextString(transaction.body())
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (transaction *Transaction) MarshalBinary() ([]byte, error) {
	return protobuf.Marshal(&transaction.transactions[0])
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
		return durationFromProto(transaction.pbBody.TransactionValidDuration)
	} else {
		return 0
	}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction *Transaction) SetTransactionValidDuration(duration time.Duration) *Transaction {
	transaction.pbBody.TransactionValidDuration = durationToProto(duration)
	return transaction
}

func (transaction *Transaction) GetTransactionID() TransactionID {
	if transaction.pbBody.TransactionID != nil {
		return transactionIDFromProto(transaction.pbBody.TransactionID)
	} else {
		return TransactionID{}
	}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction *Transaction) SetTransactionID(transactionID TransactionID) *Transaction {
	transaction.pbBody.TransactionID = transactionID.toProtobuf()
	return transaction
}

func (transaction *Transaction) GetNodeID() AccountID {
	if transaction.pbBody.NodeAccountID != nil {
		return accountIDFromProto(transaction.pbBody.NodeAccountID)
	} else {
		return AccountID{}
	}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction *Transaction) SetNodeID(nodeAccountID AccountID) *Transaction {
	transaction.pbBody.NodeAccountID = nodeAccountID.toProtobuf()
	return transaction
}

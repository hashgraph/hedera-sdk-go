package hedera

import (
	"bytes"
	"google.golang.org/grpc/codes"
	"math"
	"math/rand"
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

	transactions []*proto.Transaction
	signatures   []*proto.SignatureMap
	nodeIDs      []AccountID
}

func newTransaction() Transaction {
	return Transaction{
		pbBody: &proto.TransactionBody{
			TransactionValidDuration: durationToProto(120 * time.Second),
		},
	}
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (transaction *Transaction) UnmarshalBinary(txBytes []byte) error {
	transaction.transactions = []*proto.Transaction{}
	transaction.transactions = append(transaction.transactions, &proto.Transaction{})
	if err := protobuf.Unmarshal(txBytes, transaction.transactions[0]); err != nil {
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
func (transaction Transaction) Sign(
	privateKey PrivateKey,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
) Transaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign, isFrozen, freezeWith)
}

func (transaction Transaction) SignWithOperator(
	client *Client,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
) (Transaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
		return Transaction{}, errClientOperatorSigning
	}

	if !isFrozen() {
		freezeWith(client)
	}

	return transaction.SignWith(client.operator.publicKey, client.operator.signer, isFrozen, freezeWith), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction Transaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
) Transaction {
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

func (transaction Transaction) freezeWith(
	client *Client,
	isFrozen func() bool,
	onFreeze func(pbBody *proto.TransactionBody) bool,
) (Transaction, error) {
	if client != nil {
		if !transaction.noTXFee && transaction.pbBody.TransactionFee == 0 {
			transaction.SetMaxTransactionFee(client.maxTransactionFee)
		}

		if transaction.pbBody.TransactionID == nil {
			if client.operator != nil {
				transaction.SetTransactionID(NewTransactionID(client.operator.accountID))
			} else {
				return Transaction{}, errNoClientOrTransactionID
			}
		}
	}

	if onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	if transaction.pbBody.TransactionID != nil && transaction.pbBody.NodeAccountID != nil {
		transaction.signatures = []*proto.SignatureMap{}
		transaction.transactions = []*proto.Transaction{}

		bodyBytes, err := protobuf.Marshal(transaction.pbBody)
		if err != nil {
			// This should be unreachable
			// From the documentation this appears to only be possible if there are missing proto types
			panic(err)
		}

		protoTransaction := proto.Transaction{
			BodyBytes: bodyBytes,
		}

		transaction.transactions = append(transaction.transactions, &protoTransaction)

		return transaction, nil
	}

	if transaction.pbBody.TransactionID != nil && len(transaction.nodeIDs) > 0 {
		transaction.signatures = []*proto.SignatureMap{}
		transaction.transactions = []*proto.Transaction{}

		for _, id := range transaction.nodeIDs {
			transaction.pbBody.NodeAccountID = id.toProtobuf()
			bodyBytes, err := protobuf.Marshal(transaction.pbBody)
			if err != nil {
				// This should be unreachable
				// From the documentation this appears to only be possible if there are missing proto types
				panic(err)
			}

			transaction.signatures = append(transaction.signatures, &proto.SignatureMap{})
			transaction.transactions = append(transaction.transactions, &proto.Transaction{
				BodyBytes: bodyBytes,
			})
		}

		return transaction, nil
	}

	if client != nil && transaction.pbBody.TransactionID != nil {
		size := client.getNumberOfNodesForTransaction()

		transaction.signatures = []*proto.SignatureMap{}
		transaction.transactions = []*proto.Transaction{}
		transaction.nodeIDs = []AccountID{}

		for index := 0; index < size; index++ {
			node := client.getNextNode()

			transaction.nodeIDs = append(transaction.nodeIDs, node.id)

			transaction.pbBody.NodeAccountID = node.id.toProtobuf()
			bodyBytes, err := protobuf.Marshal(transaction.pbBody)
			if err != nil {
				// This should be unreachable
				// From the documentation this appears to only be possible if there are missing proto types
				panic(err)
			}

			transaction.signatures = append(transaction.signatures, &proto.SignatureMap{})
			transaction.transactions = append(transaction.transactions, &proto.Transaction{
				BodyBytes: bodyBytes,
			})
		}

		return transaction, nil
	}

	return Transaction{}, errNoClientOrTransactionIDOrNodeId
}

func defaultIsFrozen(transaction Transaction) bool {
	return len(transaction.transactions) > 0
}

func (transaction Transaction) requireNotFrozen(isFrozen func() bool) error {
	if isFrozen() {
		return errTransactionIsFrozen
	}

	return nil
}

func (transaction Transaction) keyAlreadySigned(pk PublicKey) bool {
	if len(transaction.signatures) > 0 {
		for _, pair := range transaction.signatures[0].SigPair {
			if bytes.HasPrefix(pk.keyData, pair.PubKeyPrefix) {
				return true
			}
		}
	}

	return false
}

func (transaction Transaction) executeForResponse(
	client *Client,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
) (TransactionID, *proto.TransactionResponse, error) {
	if client.operator != nil {
		if _, err := transaction.SignWithOperator(client, isFrozen, freezeWith); err != nil {
			return TransactionID{}, nil, err
		}
	}

	transactionBody := transaction.body()
	id := transactionIDFromProto(transactionBody.TransactionID)

	nodeAccountID := accountIDFromProto(transactionBody.NodeAccountID)
	node := client.node(nodeAccountID)

	if node == nil {
		return id, nil, newErrLocalValidationf("NodeAccountID %v not found on Client", nodeAccountID)
	}

	methodName, err := getMethodName(transaction.body())

	if err != nil {
		return id, nil, err
	}

	validUntil := time.Now().Add(time.Duration(transactionBody.TransactionValidDuration.Seconds) * time.Second)
	resp := new(proto.TransactionResponse)
	length := len(transaction.transactions)

	for attempt := 0; true; attempt++ {
		tx := transaction.transactions[attempt%length]

		if attempt > 0 && time.Now().After(validUntil) {
			// Timed out
			break
		}

		if attempt > 0 {
			// After the first attempt, start an exponentially increasing delay
			delay := 500.0 * rand.Float64() * ((math.Pow(2, float64(attempt))) - 1)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		err := node.invoke(methodName, tx, resp)
		if err != nil {
			statusCode := err.(ErrHederaNetwork).StatusCode

			if statusCode != nil && (*statusCode == codes.Unavailable || *statusCode == codes.ResourceExhausted) {
				// try again on unavailable or ResourceExhausted
				continue
			}

			return id, resp, err
		}

		if resp.NodeTransactionPrecheckCode == proto.ResponseCodeEnum_BUSY {
			// Try again (in a flash) on BUSY
			continue
		}

		return id, resp, nil
	}

	// Timed out
	precheckCode := resp.NodeTransactionPrecheckCode
	if precheckCode == proto.ResponseCodeEnum_OK {
		precheckCode = proto.ResponseCodeEnum_TRANSACTION_EXPIRED
	}

	return id, nil, newErrHederaPreCheckStatus(transaction.id, Status(precheckCode))
}

// Execute executes the Transaction with the provided client
func (transaction Transaction) execute(
	client *Client,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
) (TransactionResponse, error) {
	id, resp, err := transaction.executeForResponse(client, isFrozen, freezeWith)

	if err != nil {
		return TransactionResponse{}, err
	}

	status := Status(resp.NodeTransactionPrecheckCode)

	if status.isExceptional(true) {
		// precheck failed
		return TransactionResponse{TransactionID: id}, newErrHederaPreCheckStatus(transaction.id, status)
	}

	// success
	return TransactionResponse{TransactionID: id}, nil
}

func (transaction Transaction) String() string {
	return protobuf.MarshalTextString(transaction.transactions[0]) +
		protobuf.MarshalTextString(transaction.body())
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (transaction Transaction) MarshalBinary() ([]byte, error) {
	return protobuf.Marshal(transaction.transactions[0])
}

func (transaction Transaction) ToBytes() ([]byte, error) {
	return transaction.MarshalBinary()
}

// The protobuf stores the transaction body as raw bytes so we need to first
// decode what we have to inspect the Kind, TransactionID, and the NodeAccountID so we know how to
// properly execute it
func (transaction Transaction) body() *proto.TransactionBody {
	transactionBody := new(proto.TransactionBody)
	err := protobuf.Unmarshal(transaction.transactions[0].GetBodyBytes(), transactionBody)
	if err != nil {
		// The bodyBytes inside of the transaction at this point have been verified and this should be impossible
		panic(err)
	}

	return transactionBody
}

// getMethodName returns the proto method name of the transaction body
func getMethodName(transactionBody *proto.TransactionBody) (string, error) {
	switch transactionBody.Data.(type) {
	case *proto.TransactionBody_CryptoCreateAccount:
		return "/proto.CryptoService/createAccount", nil

	case *proto.TransactionBody_CryptoTransfer:
		return "/proto.CryptoService/cryptoTransfer", nil

	case *proto.TransactionBody_CryptoUpdateAccount:
		return "/proto.CryptoService/updateAccount", nil

	case *proto.TransactionBody_CryptoDelete:
		return "/proto.CryptoService/cryptoDelete", nil

	// FileServices
	case *proto.TransactionBody_FileCreate:
		return "/proto.FileService/createFile", nil

	case *proto.TransactionBody_FileUpdate:
		return "/proto.FileService/updateFile", nil

	case *proto.TransactionBody_FileAppend:
		return "/proto.FileService/appendContent", nil

	case *proto.TransactionBody_FileDelete:
		return "/proto.FileService/deleteFile", nil

	// Contract
	case *proto.TransactionBody_ContractCreateInstance:
		return "/proto.SmartContractService/createContract", nil

	case *proto.TransactionBody_ContractDeleteInstance:
		return "/proto.SmartContractService/deleteContract", nil

	case *proto.TransactionBody_ContractUpdateInstance:
		return "/proto.SmartContractService/updateContract", nil

	case *proto.TransactionBody_ContractCall:
		return "/proto.SmartContractService/contractCallMethod", nil

	// System
	case *proto.TransactionBody_Freeze:
		return "/proto.FreezeService/freeze", nil

	case *proto.TransactionBody_SystemDelete:
		return "/proto.FileService/systemDelete", nil

	case *proto.TransactionBody_SystemUndelete:
		return "/proto.FileService/systemUndelete", nil

	// HCS
	case *proto.TransactionBody_ConsensusCreateTopic:
		return "/proto.ConsensusService/createTopic", nil
	case *proto.TransactionBody_ConsensusDeleteTopic:
		return "/proto.ConsensusService/deleteTopic", nil
	case *proto.TransactionBody_ConsensusUpdateTopic:
		return "/proto.ConsensusService/updateTopic", nil
	case *proto.TransactionBody_ConsensusSubmitMessage:
		return "/proto.ConsensusService/submitMessage", nil

	default:
		return "", newErrLocalValidationf("Could not find method name for: %T", transactionBody.Data)
	}

}

//
// Shared
//

func (transaction Transaction) GetMaxTransactionFee() Hbar {
	return HbarFromTinybar(int64(transaction.pbBody.TransactionFee))
}

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction Transaction) SetMaxTransactionFee(fee Hbar) Transaction {
	transaction.pbBody.TransactionFee = uint64(fee.AsTinybar())
	return transaction
}

func (transaction Transaction) GetTransactionMemo() string {
	return transaction.pbBody.Memo
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction Transaction) SetTransactionMemo(memo string) Transaction {
	transaction.pbBody.Memo = memo
	return transaction
}

func (transaction Transaction) GetTransactionValidDuration() time.Duration {
	return durationFromProto(transaction.pbBody.TransactionValidDuration)
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction Transaction) SetTransactionValidDuration(duration time.Duration) Transaction {
	transaction.pbBody.TransactionValidDuration = durationToProto(duration)
	return transaction
}

func (transaction Transaction) GetTransactionID() TransactionID {
	return transactionIDFromProto(transaction.pbBody.TransactionID)
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction Transaction) SetTransactionID(transactionID TransactionID) Transaction {
	transaction.pbBody.TransactionID = transactionID.toProtobuf()
	return transaction
}

func (transaction Transaction) GetNodeID() AccountID {
	return accountIDFromProto(transaction.pbBody.NodeAccountID)
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction Transaction) SetNodeID(nodeAccountID AccountID) Transaction {
	transaction.pbBody.NodeAccountID = nodeAccountID.toProtobuf()
	return transaction
}

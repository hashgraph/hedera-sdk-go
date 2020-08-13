package hedera

import (
	"bytes"
	"crypto/ed25519"
	"google.golang.org/grpc/codes"
	"math"
	"math/rand"
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// Transaction contains the protobuf of a prepared transaction which can be signed and executed.
type Transaction struct {
	pb *proto.Transaction
	id TransactionID
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (transaction *Transaction) UnmarshalBinary(txBytes []byte) error {
	transaction.pb = new(proto.Transaction)
	if err := protobuf.Unmarshal(txBytes, transaction.pb); err != nil {
		return err
	}

	var txBody proto.TransactionBody
	if err := protobuf.Unmarshal(transaction.pb.GetBodyBytes(), &txBody); err != nil {
		return err
	}

	transaction.id = transactionIDFromProto(txBody.TransactionID)

	return nil
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction Transaction) Sign(privateKey Ed25519PrivateKey) Transaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction Transaction) signWithOperator(operator operator) Transaction {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	var signedByOperator bool
	operatorPublicKey := operator.publicKey.keyData

	for _, sigPair := range transaction.pb.SigMap.SigPair {
		if bytes.Equal(sigPair.PubKeyPrefix, operatorPublicKey) {
			signedByOperator = true
			break
		}
	}

	if !signedByOperator {
		if operator.privateKey != nil {
			transaction.Sign(*operator.privateKey)
		} else {
			transaction.SignWith(operator.publicKey, operator.signer)
		}
	}

	return transaction
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction Transaction) SignWith(publicKey Ed25519PublicKey, signer TransactionSigner) Transaction {
	signature := signer(transaction.pb.GetBodyBytes())

	transaction.pb.SigMap.SigPair = append(transaction.pb.SigMap.SigPair, &proto.SignaturePair{
		PubKeyPrefix: publicKey.keyData,
		Signature:    &proto.SignaturePair_Ed25519{Ed25519: signature},
	})

	return transaction
}

// AppendSignature verifies provided signature and public key for corresponding transaction body and adds them
// to the Transaction's signature map
func (transaction Transaction) AppendSignature(publicKey Ed25519PublicKey, signature []byte) (*Transaction, error) {
	verifiedSignature := ed25519.Verify(publicKey.Bytes(), transaction.BodyBytes(), signature)

	if verifiedSignature != true {
		return nil, newErrSignatureVerification("invalid public key or signature provided")
	}

	transaction.pb.SigMap.SigPair = append(transaction.pb.SigMap.SigPair, &proto.SignaturePair{
		PubKeyPrefix: publicKey.keyData,
		Signature:    &proto.SignaturePair_Ed25519{Ed25519: signature},
	})

	return &transaction, nil
}

func (transaction Transaction) executeForResponse(client *Client) (TransactionID, *proto.TransactionResponse, error) {
	if client.operator != nil {
		transaction.signWithOperator(*client.operator)
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

	for attempt := 0; true; attempt++ {
		if attempt > 0 && time.Now().After(validUntil) {
			// Timed out
			break
		}

		if attempt > 0 {
			// After the first attempt, start an exponentially increasing delay
			delay := 500.0 * rand.Float64() * ((math.Pow(2, float64(attempt))) - 1)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		err := node.invoke(methodName, transaction.pb, resp)
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
func (transaction Transaction) Execute(client *Client) (TransactionID, error) {
	id, resp, err := transaction.executeForResponse(client)

	if err != nil {
		return id, err
	}

	status := Status(resp.NodeTransactionPrecheckCode)

	if status.isExceptional(true) {
		// precheck failed
		return id, newErrHederaPreCheckStatus(transaction.id, status)
	}

	// success
	return id, nil
}

func (transaction Transaction) String() string {
	return protobuf.MarshalTextString(transaction.pb) +
		protobuf.MarshalTextString(transaction.body())
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (transaction Transaction) MarshalBinary() ([]byte, error) {
	return protobuf.Marshal(transaction.pb)
}

// ID returns the transaction ID of the transaction
func (transaction Transaction) ID() TransactionID {
	// Provide an accessor function to prevent the user from mutating the
	// ID which would result in undefined behavior.
	return transaction.id
}

// BodyBytes returns the transaction body as raw bytes
func (transaction Transaction) BodyBytes() []byte {
	return transaction.pb.GetBodyBytes()
}

// The protobuf stores the transaction body as raw bytes so we need to first
// decode what we have to inspect the Kind, TransactionID, and the NodeAccountID so we know how to
// properly execute it
func (transaction Transaction) body() *proto.TransactionBody {
	transactionBody := new(proto.TransactionBody)
	err := protobuf.Unmarshal(transaction.pb.GetBodyBytes(), transactionBody)
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

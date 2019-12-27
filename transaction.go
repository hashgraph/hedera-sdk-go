package hedera

import (
	"bytes"
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
	"math"
	"math/rand"
	"time"
)

type Transaction struct {
	pb *proto.Transaction
}

func (transaction Transaction) Sign(privateKey Ed25519PrivateKey) Transaction {
	// TODO: Disallow duplicate [Sign] with the same private key

	signature := privateKey.Sign(transaction.pb.GetBodyBytes())

	transaction.pb.SigMap.SigPair = append(transaction.pb.SigMap.SigPair, &proto.SignaturePair{
		PubKeyPrefix: nil,
		Signature:    &proto.SignaturePair_Ed25519{Ed25519: signature},
	})

	return transaction
}

func (transaction Transaction) SignWith(publicKey Ed25519PublicKey, signer Signer) Transaction {
	signature := signer(transaction.pb.GetBodyBytes())

	transaction.pb.SigMap.SigPair = append(transaction.pb.SigMap.SigPair, &proto.SignaturePair{
		PubKeyPrefix: publicKey.keyData,
		Signature:    &proto.SignaturePair_Ed25519{Ed25519: signature},
	})

	return transaction
}

func (transaction Transaction) Execute(client *Client) (TransactionID, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	var signedByOperator bool
	operatorPublicKey := client.operator.privateKey.publicKey.keyData

	for _, sigPair := range transaction.pb.SigMap.SigPair {
		if bytes.Equal(sigPair.PubKeyPrefix, operatorPublicKey) {
			signedByOperator = true
			break
		}
	}

	if !signedByOperator {
		if client.operator.privateKey != nil {
			transaction.Sign(*client.operator.privateKey)
		} else {
			transaction.SignWith(client.operator.publicKey, client.operator.signer)
		}
	}

	transactionBody := transaction.body()
	id := transactionIDFromProto(transactionBody.TransactionID)

	nodeAccountID := accountIDFromProto(transactionBody.NodeAccountID)
	node := client.node(nodeAccountID)

	if node == nil {
		return id, fmt.Errorf("NodeAccountID %v not found on Client", nodeAccountID)
	}

	var methodName string

	// TODO: Add the rest of the types here
	switch transactionBody.Data.(type) {
	case *proto.TransactionBody_CryptoCreateAccount:
		methodName = "/proto.CryptoService/createAccount"

	case *proto.TransactionBody_CryptoTransfer:
		methodName = "/proto.CryptoService/cryptoTransfer"

	case *proto.TransactionBody_CryptoUpdateAccount:
		methodName = "/proto.CryptoService/updateAccount"

	// FileServices
	case *proto.TransactionBody_FileCreate:
		methodName = "/proto.FileService/createFile"

	case *proto.TransactionBody_FileUpdate:
		methodName = "/proto.FileService/updateFile"

	case *proto.TransactionBody_FileAppend:
		methodName = "/proto.FileService/appendFile"

	case *proto.TransactionBody_FileDelete:
		methodName = "/proto.FileService/deleteFile"

	// Contract
	case *proto.TransactionBody_ContractCreateInstance:
		methodName = "/proto.SmartContractService/createContract"

	case *proto.TransactionBody_ContractDeleteInstance:
		methodName = "/proto.SmartContractService/deleteContract"

	case *proto.TransactionBody_ContractUpdateInstance:
		methodName = "/proto.SmartContractService/updateContract"

	case *proto.TransactionBody_ContractCall:
		methodName = "/proto.SmartContractService/contractCallMethod"

	// System
	case *proto.TransactionBody_Freeze:
		methodName = "/proto.FreezeService/freeze"

	case *proto.TransactionBody_SystemDelete:
		methodName = "/proto.FileService/systemDelete"

	case *proto.TransactionBody_SystemUndelete:
		methodName = "/proto.FileService/systemUndelete"

	default:
		return id, fmt.Errorf("unimplemented: %T", transactionBody.Data)
	}

	validUntil := time.Now().Add(time.Duration(transactionBody.TransactionValidDuration.Seconds) * time.Second)
	resp := new(proto.TransactionResponse)

	for attempt := 0; true; attempt += 1 {
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
			return id, err
		}

		if resp.NodeTransactionPrecheckCode == proto.ResponseCodeEnum_BUSY {
			// Try again (in a flash) on BUSY
			continue
		}

		if isStatusExceptional(resp.NodeTransactionPrecheckCode, true) {
			return id, fmt.Errorf("%v", resp.NodeTransactionPrecheckCode)
		}

		return id, nil
	}

	// Timed out
	// TODO: Better error here?
	return id, fmt.Errorf("%v", resp.NodeTransactionPrecheckCode)
}

func (transaction Transaction) String() string {
	return protobuf.MarshalTextString(transaction.pb) +
		protobuf.MarshalTextString(transaction.body())
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

func isStatusExceptional(status proto.ResponseCodeEnum, unknownIsExceptional bool) bool {
	switch status {
	case proto.ResponseCodeEnum_SUCCESS, proto.ResponseCodeEnum_OK:
		return false

	case proto.ResponseCodeEnum_UNKNOWN, proto.ResponseCodeEnum_RECEIPT_NOT_FOUND, proto.ResponseCodeEnum_RECORD_NOT_FOUND:
		return unknownIsExceptional

	default:
		return true
	}
}

package hedera

import (
	"errors"
	"strings"
	"time"

	"github.com/hashgraph/hedera-sdk-go/hedera_proto"
)

const receiptRetryDelay = 500
const receiptInitialDelay = 1000

const prefixLen = 6

type ErrorTransactionValidation struct {
	Messages []string
	Err      error
}

func (e *ErrorTransactionValidation) Error() string {
	return "The following requirements were not met: \n" + strings.Join(e.Messages, "\n")
}

type TransactionBuilder interface {
	SetMaxTransactionFee(uint64) *TransactionBuilder
	SetMemo(string)
	validate() ErrorTransactionValidation
	Build() (Transaction, error)
}

type TransactionID struct {
	Account           AccountID
	ValidStartSeconds uint64
	ValidStartNanos   uint32
}

func generateTransactionID(accountID AccountID) TransactionID {
	now := time.Now()

	return TransactionID{
		accountID,
		uint64(now.Unix()),
		uint32(now.UnixNano() - (now.Unix() * 1e+9)),
	}
}

type Transaction struct {
	Kind   TransactionKind
	client *Client
	inner  hedera_proto.Transaction
}

func (transaction Transaction) AddSignature(signature []byte, publicKey Ed25519PublicKey) Transaction {
	signaturePair := hedera_proto.SignaturePair{
		PubKeyPrefix: publicKey.keyData,
		Signature: &hedera_proto.SignaturePair_Ed25519{
			Ed25519: signature,
		},
	}

	sigmap := transaction.inner.GetSigMap()

	if sigmap == nil {
		sigmap = &hedera_proto.SignatureMap{}
	}

	sigmap.SigPair = append(sigmap.SigPair, &signaturePair)

	transaction.inner.SigMap = sigmap

	return transaction
}

func (transaction Transaction) Sign(privateKey Ed25519PrivateKey) Transaction {
	signature := privateKey.Sign(transaction.inner.GetBodyBytes())

	return transaction.AddSignature(signature, privateKey.PublicKey())
}

func (transaction Transaction) getReceipt() (*TransactionReceipt, error) {
	return nil, nil
}

func (transaction Transaction) Execute() (*TransactionID, error) {
	// fixme: proper error handling
	if transaction.client == nil {
		return nil, errors.New("No client was provided on this transaction")
	}

	txID := generateTransactionID(transaction.client.operator.accountID)

	body := transaction.inner.GetBody()

	body.TransactionID = &hedera_proto.TransactionID{
		TransactionValidStart: &hedera_proto.Timestamp{
			Seconds: int64(txID.ValidStartSeconds),
			Nanos:   int32(txID.ValidStartNanos),
		},
		AccountID: &hedera_proto.AccountID{
			ShardNum:   int64(txID.Account.Shard),
			RealmNum:   int64(txID.Account.Realm),
			AccountNum: int64(txID.Account.Account),
		},
	}

	// todo: use response and handle precheck codes
	_, error := transaction.Kind.execute(*transaction.client, transaction.inner)

	// todo: handle result errors
	if error != nil {
		return nil, error
	}

	return &txID, nil
}

func (transaction Transaction) ExecuteForReceipt() (*TransactionReceipt, error) {
	_, err := transaction.Execute()

	if err != nil {
		return nil, err
	}

	// todo: add receipts
	return nil, nil
}

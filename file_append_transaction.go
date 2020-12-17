package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/pkg/errors"

	"time"
)

// FileAppendTransaction appends the given contents to the end of the file. If a file is too big to create with a single
// FileCreateTransaction, then it can be created with the first part of its contents, and then appended multiple times
// to create the entire file.
type FileAppendTransaction struct {
	Transaction
	pb        *proto.FileAppendTransactionBody
	maxChunks uint64
	contents  []byte
}

// NewFileAppendTransaction creates a FileAppendTransaction transaction which can be
// used to construct and execute a File Append Transaction.
func NewFileAppendTransaction() *FileAppendTransaction {
	pb := &proto.FileAppendTransactionBody{}

	transaction := FileAppendTransaction{
		pb:          pb,
		Transaction: newTransaction(),
		maxChunks:   10,
		contents:    make([]byte, 0),
	}

	return &transaction
}

func fileAppendTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FileAppendTransaction {
	return FileAppendTransaction{
		Transaction: transaction,
		pb:          pb.GetFileAppend(),
		maxChunks:   10,
		contents:    make([]byte, 0),
	}
}

// SetFileID sets the FileID of the file to which the bytes are appended to.
func (transaction *FileAppendTransaction) SetFileID(ID FileID) *FileAppendTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FileID = ID.toProtobuf()
	return transaction
}

func (transaction *FileAppendTransaction) GetFileID() FileID {
	return fileIDFromProtobuf(transaction.pb.GetFileID())
}

// SetContents sets the bytes to append to the contents of the file.
func (transaction *FileAppendTransaction) SetContents(contents []byte) *FileAppendTransaction {
	transaction.requireNotFrozen()
	transaction.contents = contents
	return transaction
}

func (transaction *FileAppendTransaction) GetContents() []byte {
	return transaction.contents
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func fileAppendTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getFile().AppendContent,
	}
}

func (transaction *FileAppendTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *FileAppendTransaction) Sign(
	privateKey PrivateKey,
) *FileAppendTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *FileAppendTransaction) SignWithOperator(
	client *Client,
) (*FileAppendTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *FileAppendTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileAppendTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	} else {
		transaction.transactions = make([]*proto.Transaction, 0)
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.signedTransactions); index++ {
		signature := signer(transaction.signedTransactions[index].GetBodyBytes())

		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *FileAppendTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	list, err := transaction.ExecuteAll(client)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.transactionIDs[transaction.nextTransactionIndex],
			NodeID:        list[0].NodeID,
			Hash:          make([]byte, 0),
		}, err
	}

	return list[0], nil
}

// ExecuteAll executes the all the Transactions with the provided client
func (transaction *FileAppendTransaction) ExecuteAll(
	client *Client,
) ([]TransactionResponse, error) {
	if client == nil || client.operator == nil {
		return []TransactionResponse{}, errNoClientProvided
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return []TransactionResponse{}, err
		}
	}

	transactionID := transaction.transactionIDs[0]

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	size := len(transaction.signedTransactions) / len(transaction.nodeIDs)
	list := make([]TransactionResponse, size)

	for i := 0; i < size; i++ {
		resp, err := execute(
			client,
			request{
				transaction: &transaction.Transaction,
			},
			transaction_shouldRetry,
			transaction_makeRequest,
			transaction_advanceRequest,
			transaction_getNodeAccountID,
			fileAppendTransaction_getMethod,
			transaction_mapResponseStatus,
			transaction_mapResponse,
		)

		if err != nil {
			return list, err
		}

		list[i] = resp.transaction

		_, err = NewTransactionReceiptQuery().
			SetNodeAccountIDs([]AccountID{resp.transaction.NodeID}).
			SetTransactionID(resp.transaction.TransactionID).
			Execute(client)
		if err != nil {
			return list, err
		}
	}

	return list, nil
}

func (transaction *FileAppendTransaction) Freeze() (*FileAppendTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FileAppendTransaction) FreezeWith(client *Client) (*FileAppendTransaction, error) {
	if len(transaction.nodeIDs) == 0 {
		if client == nil {
			return transaction, errNoClientOrTransactionIDOrNodeId
		} else {
			transaction.nodeIDs = client.network.getNodeAccountIDsForExecute()
		}
	}

	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	chunks := uint64((len(transaction.contents) + (chunkSize - 1)) / chunkSize)
	if chunks > transaction.maxChunks {
		return transaction, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: transaction.maxChunks,
		}
	}

	initialTransactionID := transaction.transactionIDs[0]
	nextTransactionID := initialTransactionID

	transaction.transactionIDs = []TransactionID{}
	transaction.transactions = []*proto.Transaction{}
	transaction.signedTransactions = []*proto.SignedTransaction{}

	for i := 0; uint64(i) < chunks; i += 1 {
		start := i * chunkSize
		end := start + chunkSize

		if end > len(transaction.contents) {
			end = len(transaction.contents)
		}

		transaction.transactionIDs = append(transaction.transactionIDs, nextTransactionID)

		transaction.pb.Contents = transaction.contents[start:end]

		transaction.pbBody.TransactionID = nextTransactionID.toProtobuf()
		transaction.pbBody.Data = &proto.TransactionBody_FileAppend{
			FileAppend: transaction.pb,
		}

		for _, nodeAccountID := range transaction.nodeIDs {
			transaction.pbBody.NodeAccountID = nodeAccountID.toProtobuf()

			bodyBytes, err := protobuf.Marshal(transaction.pbBody)
			if err != nil {
				return transaction, errors.Wrap(err, "error serializing body for file append")
			}

			transaction.signedTransactions = append(transaction.signedTransactions, &proto.SignedTransaction{
				BodyBytes: bodyBytes,
				SigMap:    &proto.SignatureMap{},
			})
		}

		nextTransactionID.ValidStart = nextTransactionID.ValidStart.Add(1 * time.Nanosecond)
	}

	return transaction, nil
}

func (transaction *FileAppendTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetMaxTransactionFee(fee Hbar) *FileAppendTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FileAppendTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetTransactionMemo(memo string) *FileAppendTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FileAppendTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetTransactionValidDuration(duration time.Duration) *FileAppendTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FileAppendTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetTransactionID(transactionID TransactionID) *FileAppendTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *FileAppendTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeAccountID sets the node AccountID for this FileAppendTransaction.
func (transaction *FileAppendTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileAppendTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *FileAppendTransaction) SetMaxRetry(count int) *FileAppendTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

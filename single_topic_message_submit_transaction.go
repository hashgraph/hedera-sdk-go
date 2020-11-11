package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type singleTopicMessageSubmitTransaction struct {
	Transaction
	pb *proto.ConsensusSubmitMessageTransactionBody
}

func newSingleTopicMessageSubmitTransaction(
	pbBody *proto.TransactionBody,
	pb *proto.ConsensusSubmitMessageTransactionBody,
	message []byte,
	chunkInfo *proto.ConsensusMessageChunkInfo,
) *singleTopicMessageSubmitTransaction {
	pb.ChunkInfo = chunkInfo
	pb.Message = message

	return &singleTopicMessageSubmitTransaction{
		Transaction: Transaction{
			pbBody: pbBody,
			id:     transactionIDFromProtobuf(pbBody.TransactionID),
		},
		pb: pb,
	}
}
func singleTopicMessageSubmitTransactionFromProtobuf(transactions map[TransactionID]map[AccountID]*proto.Transaction, pb *proto.TransactionBody) singleTopicMessageSubmitTransaction {
	return singleTopicMessageSubmitTransaction{
		Transaction: transactionFromProtobuf(transactions, pb),
		pb:          pb.GetConsensusSubmitMessage(),
	}

}
//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func singleTopicMessageSubmitTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getTopic().SubmitMessage,
	}
}

func (transaction *singleTopicMessageSubmitTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *singleTopicMessageSubmitTransaction) Sign(
	privateKey PrivateKey,
) *singleTopicMessageSubmitTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *singleTopicMessageSubmitTransaction) SignWithOperator(
	client *Client,
) (*singleTopicMessageSubmitTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *singleTopicMessageSubmitTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *singleTopicMessageSubmitTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.transactions); index++ {
		signature := signer(transaction.transactions[index].GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *singleTopicMessageSubmitTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	if len(transaction.Transaction.GetNodeAccountIDs()) == 0 {
		transaction.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	transactionID := transaction.id

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		singleTopicMessageSubmitTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.id,
		NodeID:        resp.transaction.NodeID,
	}, nil
}

func (transaction *singleTopicMessageSubmitTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ConsensusSubmitMessage{
		ConsensusSubmitMessage: transaction.pb,
	}

	return true
}

func (transaction *singleTopicMessageSubmitTransaction) Freeze() (*singleTopicMessageSubmitTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *singleTopicMessageSubmitTransaction) FreezeWith(client *Client) (*singleTopicMessageSubmitTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

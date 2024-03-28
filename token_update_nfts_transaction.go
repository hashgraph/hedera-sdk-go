package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TokenUpdateNfts struct {
	Transaction
	tokenID       *TokenID
	serialNumbers []int64
	metadata      []byte
}

func NewTokenUpdateNftsTransaction() *TokenUpdateNfts {
	return &TokenUpdateNfts{
		Transaction: _NewTransaction(),
	}
}

func _NewTokenUpdateNftsTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenUpdateNfts {
	return &TokenUpdateNfts{
		Transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenUpdateNfts().GetToken()),
		serialNumbers: func() []int64 {
			var serialNumbers []int64
			for _, serialNumber := range pb.GetTokenUpdateNfts().GetSerialNumbers() {
				serialNumbers = append(serialNumbers, int64(serialNumber))
			}
			return serialNumbers
		}(),
	}
}

// Getter for tokenID
func (t *TokenUpdateNfts) GetTokenID() *TokenID {
	return t.tokenID
}

// Setter for tokenID
func (t *TokenUpdateNfts) SetTokenID(tokenID TokenID) *TokenUpdateNfts {
	t._RequireNotFrozen()
	t.tokenID = &tokenID
	return t
}

// Getter for serialNumbers
func (t *TokenUpdateNfts) GetSerialNumbers() []int64 {
	return t.serialNumbers
}

// Setter for serialNumbers
func (t *TokenUpdateNfts) SetSerialNumbers(serialNumbers []int64) *TokenUpdateNfts {
	t._RequireNotFrozen()
	t.serialNumbers = serialNumbers
	return t
}

// Getter for metadata
func (t *TokenUpdateNfts) GetMetadata() []byte {
	return t.metadata
}

// Setter for metadata
func (t *TokenUpdateNfts) SetMetadata(metadata []byte) *TokenUpdateNfts {
	t._RequireNotFrozen()
	t.metadata = metadata
	return t
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenUpdateNfts) Sign(privateKey PrivateKey) *TokenUpdateNfts {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenUpdateNfts) SignWithOperator(client *Client) (*TokenUpdateNfts, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenUpdateNfts) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenUpdateNfts {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenUpdateNfts) AddSignature(publicKey PublicKey, signature []byte) *TokenUpdateNfts {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenUpdateNfts) SetGrpcDeadline(deadline *time.Duration) *TokenUpdateNfts {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenUpdateNfts) Freeze() (*TokenUpdateNfts, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenUpdateNfts) FreezeWith(client *Client) (*TokenUpdateNfts, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUpdateNfts.
func (tx *TokenUpdateNfts) SetMaxTransactionFee(fee Hbar) *TokenUpdateNfts {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenUpdateNfts) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenUpdateNfts {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenUpdateNfts.
func (tx *TokenUpdateNfts) SetTransactionMemo(memo string) *TokenUpdateNfts {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenUpdateNfts.
func (tx *TokenUpdateNfts) SetTransactionValidDuration(duration time.Duration) *TokenUpdateNfts {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenUpdateNfts) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenUpdateNfts.
func (tx *TokenUpdateNfts) SetTransactionID(transactionID TransactionID) *TokenUpdateNfts {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenUpdateNfts.
func (tx *TokenUpdateNfts) SetNodeAccountIDs(nodeID []AccountID) *TokenUpdateNfts {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenUpdateNfts) SetMaxRetry(count int) *TokenUpdateNfts {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenUpdateNfts) SetMaxBackoff(max time.Duration) *TokenUpdateNfts {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenUpdateNfts) SetMinBackoff(min time.Duration) *TokenUpdateNfts {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenUpdateNfts) SetLogLevel(level LogLevel) *TokenUpdateNfts {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenUpdateNfts) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenUpdateNfts) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

func (tx *TokenUpdateNfts) getName() string {
	return "TokenUpdateNfts"
}

func (tx *TokenUpdateNfts) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}
	return nil
}

func (tx *TokenUpdateNfts) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenUpdateNfts{
			TokenUpdateNfts: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenUpdateNfts) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenUpdateNfts{
			TokenUpdateNfts: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenUpdateNfts) buildProtoBody() *services.TokenUpdateNftsTransactionBody {
	body := &services.TokenUpdateNftsTransactionBody{}

	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}
	serialNumbers := make([]int64, 0)
	if len(tx.serialNumbers) != 0 {
		for _, serialNumber := range tx.serialNumbers {
			serialNumbers = append(serialNumbers, int64(serialNumber))
			body.SerialNumbers = serialNumbers
		}
	}
	body.Metadata = wrapperspb.Bytes(tx.metadata)
	return body
}

func (tx *TokenUpdateNfts) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UpdateNfts,
	}
}

func (tx *TokenUpdateNfts) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

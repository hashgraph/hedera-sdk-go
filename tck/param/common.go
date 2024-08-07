package param

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

type CommonTransactionParams struct {
	TransactionId            *string   `json:"transactionId"`
	MaxTransactionFee        *int64    `json:"maxTransactionFee"`
	ValidTransactionDuration *uint64   `json:"validTransactionDuration"`
	Memo                     *string   `json:"memo"`
	RegenerateTransactionId  *bool     `json:"regenerateTransactionId"`
	Signers                  *[]string `json:"signers"`
}

func (common *CommonTransactionParams) FillOutTransaction(transactionInterface hedera.TransactionInterface, transaction *hedera.Transaction, client *hedera.Client) {
	if common.TransactionId != nil {
		txId, _ := hedera.TransactionIdFromString(*common.TransactionId)
		transaction.SetTransactionID(txId)
	}

	if common.MaxTransactionFee != nil {
		transaction.SetMaxTransactionFee(hedera.HbarFromTinybar(*common.MaxTransactionFee))
	}

	if common.ValidTransactionDuration != nil {
		transaction.SetTransactionValidDuration(time.Duration(*common.ValidTransactionDuration) * time.Second)
	}

	if common.Memo != nil {
		transaction.SetTransactionMemo(*common.Memo)
	}

	if common.RegenerateTransactionId != nil {
		transaction.SetRegenerateTransactionID(*common.RegenerateTransactionId)
	}

	if common.Signers != nil {
		transaction.FreezeWith(client, transactionInterface)
		for _, signer := range *common.Signers {
			s, _ := hedera.PrivateKeyFromString(signer)
			transaction.Sign(s)
		}
	}
}

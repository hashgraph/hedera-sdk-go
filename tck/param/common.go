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

func (common *CommonTransactionParams) FillOutTransaction(transactionInterface hedera.TransactionInterface, client *hedera.Client) {
	if common.TransactionId != nil {
		txId, _ := hedera.TransactionIdFromString(*common.TransactionId)
		hedera.TransactionSetTransactionID(transactionInterface, txId)
	}

	if common.MaxTransactionFee != nil {
		hedera.TransactionSetMaxTransactionFee(transactionInterface, hedera.HbarFromTinybar(*common.MaxTransactionFee))
	}

	if common.ValidTransactionDuration != nil {
		hedera.TransactionSetTransactionValidDuration(transactionInterface, time.Duration(*common.ValidTransactionDuration)*time.Second)
	}

	if common.Memo != nil {
		hedera.TransactionSetTransactionMemo(transactionInterface, *common.Memo)
	}

	if common.RegenerateTransactionId != nil {
		hedera.TransactionSetTransactionID(transactionInterface, hedera.TransactionIDGenerate(client.GetOperatorAccountID()))
	}

	if common.Signers != nil {
		hedera.TransactionFreezeWith(transactionInterface, client)
		for _, signer := range *common.Signers {
			s, _ := hedera.PrivateKeyFromString(signer)
			hedera.TransactionSign(transactionInterface, s)
		}
	}
}

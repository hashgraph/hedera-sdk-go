package param

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

type CommonTransactionParams struct {
	TransactionId            *string   `json:"transactionId"`
	MaxTransactionFee        *int64    `json:"maxTransactionFee"`
	ValidTransactionDuration *uint64   `json:"validTransactionDuration"`
	Memo                     *string   `json:"memo"`
	RegenerateTransactionId  *bool     `json:"regenerateTransactionId"`
	Signers                  *[]string `json:"signers"`
}

func (common *CommonTransactionParams) FillOutTransaction(transactionInterface hiero.TransactionInterface, client *hiero.Client) {
	if common.TransactionId != nil {
		txId, _ := hiero.TransactionIdFromString(*common.TransactionId)
		hiero.TransactionSetTransactionID(transactionInterface, txId)
	}

	if common.MaxTransactionFee != nil {
		hiero.TransactionSetMaxTransactionFee(transactionInterface, hiero.HbarFromTinybar(*common.MaxTransactionFee))
	}

	if common.ValidTransactionDuration != nil {
		hiero.TransactionSetTransactionValidDuration(transactionInterface, time.Duration(*common.ValidTransactionDuration)*time.Second)
	}

	if common.Memo != nil {
		hiero.TransactionSetTransactionMemo(transactionInterface, *common.Memo)
	}

	if common.RegenerateTransactionId != nil {
		hiero.TransactionSetTransactionID(transactionInterface, hiero.TransactionIDGenerate(client.GetOperatorAccountID()))
	}

	if common.Signers != nil {
		hiero.TransactionFreezeWith(transactionInterface, client)
		for _, signer := range *common.Signers {
			s, _ := hiero.PrivateKeyFromString(signer)
			hiero.TransactionSign(transactionInterface, s)
		}
	}
}

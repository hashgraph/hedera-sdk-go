package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	// Our hypothetical primary service only knows the operator/sender's account ID and the recipient's accountID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	operatorPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	recipientAccountID := hedera.AccountID{Account: 3}

	// We create a client without a set operator
	client, err := hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	client.SetOperator(operatorAccountID, operatorPrivateKey)

	// We must manually construct a TransactionID with the accountID of the operator/sender
	// This is the account that will be charged the transaction fee
	txID := hedera.TransactionIDGenerate(operatorAccountID)

	// The following steps are required for manually signing
	transaction, err := hedera.NewTransferTransaction().
		// 1. Manually set the transaction ID
		SetTransactionID(txID).
		// 2. Add your sender and amount to be send
		AddHbarTransfer(operatorAccountID, hedera.NewHbar(-1)).
		// 3. add the recipient(s) and amount to be received
		AddHbarTransfer(recipientAccountID, hedera.NewHbar(1)).
		SetTransactionMemo("go sdk example multi_app_transfer/main.go").
		// 4. build the transaction using the client that does not have a set operator
		FreezeWith(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error freezing Transfer Transaction", err))
	}

	// Marshal your transaction to bytes
	txBytes, err := transaction.ToBytes()
	if err != nil {
		panic(fmt.Sprintf("%v : error converting transfer transaction to bytes", err))
	}

	fmt.Printf("Marshalled the unsigned transaction to bytes \n%v\n", txBytes)

	//
	// Send the bytes to the application or service that acts as a signer for your transactions
	//
	signedTxBytes, err := signingService(txBytes)

	if err != nil {
		panic(fmt.Sprintf("%v : error signing transfer transaction", err))
	}

	fmt.Printf("Received bytes for signed transaction: \n%v\n", signedTxBytes)

	// Unmarshal your bytes into the signed transaction
	var signedTx hedera.TransferTransaction
	tx, err := hedera.TransactionFromBytes(signedTxBytes)
	if err != nil {
		panic(fmt.Sprintf("%v : error converting bytes to transfer transaction", err))
	}

	// Converting from interface{} to TransferTransaction, if that's what we got
	switch t := tx.(type) {
	case hedera.TransferTransaction:
		signedTx = t
	default:
		panic("Did not receive `TransferTransaction` back from signed bytes")
	}

	// Execute the transaction
	response, err := signedTx.Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error executing the transfer transaction", err))
	}

	// Get the receipt of the transaction to check the status
	receipt, err := response.GetReceipt(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving transfer transaction receipt", err))
	}

	// If Status Success is returned then everything is good
	fmt.Printf("Crypto transfer status: %v\n", receipt.Status)
}

// signingService represents an offline service which knows the private keys needed for signing
// a transaction and returns the byte representation of the transaction
func signingService(txBytes []byte) ([]byte, error) {
	fmt.Println("\nSigning service has received the transaction")

	// Your signing service is aware of the operator's private key
	operatorPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		return txBytes, err
	}

	// Unmarshal the unsigned transaction's bytes
	var unsignedTx hedera.TransferTransaction
	tx, err := hedera.TransactionFromBytes(txBytes)
	if err != nil {
		return txBytes, err
	}

	// Converting from interface{} to TransferTransaction, if that's what we got
	switch t := tx.(type) {
	case hedera.TransferTransaction:
		unsignedTx = t
	default:
		panic("Did not receive `TransferTransaction` back from signed bytes")
	}

	fmt.Printf("The Signing service is signing the transaction with key: %v\n", operatorPrivateKey)

	// sign your unsigned transaction and marshal back to bytes
	return unsignedTx.
		Sign(operatorPrivateKey).
		ToBytes()
}

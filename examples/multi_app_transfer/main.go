package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {

	// Our hypothetical primary service only knows the operator/sender's account ID and the recipient's accountID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	recipientAccountID := hedera.AccountID{Account: 3}

	// We create a client without a set operator
	client := hedera.ClientForTestnet()

	// We must manually construct a TransactionID with the accountID of the operator/sender
	// This is the account that will be charged the transaction fee
	txID := hedera.NewTransactionID(operatorAccountID)

	// The following steps are required for manually signing
	transaction, err := hedera.NewCryptoTransferTransaction().
		// 1. Manually set the transaction ID
		SetTransactionID(txID).
		// 2. Add your sender and amount to be send
		AddSender(operatorAccountID, hedera.NewHbar(1)).
		// 3. add the recipient(s) and amount to be received
		AddRecipient(recipientAccountID, hedera.NewHbar(1)).
		SetTransactionMemo("go sdk example multi_app_transfer/main.go").
		// 4. build the transaction using the client that does not have a set operator
		Build(client)

	if err != nil {
		panic(err)
	}

	// marshal your transaction to bytes
	txBytes, err := transaction.MarshalBinary()

	if err != nil {
		panic(err)
	}

	fmt.Printf("marshalled the unsigned transaction to bytes \n%v\n", txBytes)

	//
	// Send the bytes to the application or service that acts as a signer for your transactions
	//
	signedTxBytes, err := signingService(txBytes)

	if err != nil {
		panic(err)
	}

	fmt.Printf("received bytes for signed transaction \n%v\n", signedTxBytes)

	// unmarshal your bytes into the signed transaction
	var signedTx hedera.Transaction
	err = signedTx.UnmarshalBinary(signedTxBytes)

	if err != nil {
		panic(err)
	}

	// execute the transaction
	txID, err = signedTx.Execute(client)

	if err != nil {
		panic(err)
	}

	// get the receipt of the transaction to check the status
	receipt, err := txID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	// if Status Success is returned then everything is good
	fmt.Printf("crypto transfer status: %v\n", receipt.Status)
}

// signingService represents an offline service which knows the private keys needed for signing
// a transaction and returns the byte representation of the transaction
func signingService(txBytes []byte) ([]byte, error) {
	fmt.Println("signing service has received the transaction")

	// Your signing service is aware of the operator's private key
	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		return txBytes, err
	}

	// unmarshal the unsigned transaction's bytes
	var unsignedTx hedera.Transaction
	err = unsignedTx.UnmarshalBinary(txBytes)

	if err != nil {
		return txBytes, err
	}

	fmt.Printf("The Signing service is signing the transaction with key %v\n", operatorPrivateKey)

	// sign your unsigned transaction and marshal back to bytes
	return unsignedTx.
		Sign(operatorPrivateKey).
		MarshalBinary()
}

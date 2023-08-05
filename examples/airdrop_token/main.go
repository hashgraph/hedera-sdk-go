package main

import (
	"encoding/csv"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"log"
	"os"
	"sync"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		log.Println(err.Error(), ": error creating client")
		return
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		log.Println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		log.Println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	// Open the file
	csvfile, err := os.Open("accounts.csv")
	if err != nil {
		log.Println("Couldn't open the csv file", err)
		return
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	var wg sync.WaitGroup

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err != nil {
			log.Println("Finished reading the CSV file.")
			break
		}

		// Increment WaitGroup counter
		wg.Add(1)

		// Start goroutine to send token
		go func(record []string) {
			// Decrement WaitGroup counter when the goroutine completes
			defer wg.Done()

			// Parse the account ID
			accountId, err := hedera.AccountIDFromString(record[0])
			if err != nil {
				log.Println("Invalid account ID:", record[0])
				return
			}

			// Transfer the token
			tokenId, _ := hedera.TokenIDFromString("0.0.1159074") // Replace with the Token ID you want to send
			_, err = hedera.NewTransferTransaction().
				AddTokenTransfer(tokenId, operatorAccountID, -100000000). // Send 0.001 of the token (token has 8 decimal places)
				AddTokenTransfer(tokenId, accountId, 100000000). // Receive 0.001 of the token
				Execute(client)
			if err != nil {
				log.Println("Failed to execute transaction for account ID:", record[0])
				return
			}

			log.Println("Successfully sent token to account ID:", record[0])
		}(record)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}

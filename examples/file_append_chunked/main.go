package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			println("not error", err.Error())
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")
	var operatorKey hedera.PrivateKey

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err = hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	_, err = readResources("file.txt")
	if err != nil {
		panic(err)
	}

	newFileResponse, err := hedera.NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte("Hello from Hedera.")).
		SetMaxTransactionFee(hedera.NewHbar(2)).
		Execute(client)
	if err != nil {
		panic(err)
	}

	newFileReceipt, err := newFileResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	var fileID hedera.FileID
	if newFileReceipt.FileID != nil {
		fileID = *newFileReceipt.FileID
		println("FileID: ", fileID.String())
	} else {
		panic("FileID is null")
	}

	var contents strings.Builder

	for i := 0; i <=4096; i++ {
		contents.WriteString("1")
	}

	fileAppend, err := hedera.NewFileAppendTransaction().
		SetNodeAccountIDs([]hedera.AccountID{newFileResponse.NodeID}).
		SetFileID(fileID).
		SetContents([]byte(contents.String())).
		SetMaxTransactionFee(hedera.NewHbar(1000)).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	fileAppendResponse, err := fileAppend.Execute(client)
	if err != nil {
		panic(err)
	}

	fileAppendReceipt, err := fileAppendResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	println(fileAppendReceipt.Status.String())
}

func readResources(filename string) (string, error) {
	var bigContents strings.Builder
	dat, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	buffer :=  bufio.NewReader(dat)
	if line, _, err :=buffer.ReadLine(); err != nil {
		return "", err
	} else {
		_, err := bigContents.Write(line)
		if err != nil {
			return "", err
		}
		for{
			line, _, err = buffer.ReadLine()
			if err != nil{
				if err == io.EOF{
					println("End of File")
					break
				} else {
					return "", err
				}
			}
			bigContents.Write(line)
			bigContents.WriteString("\n")
		}
	}

	return bigContents.String(), nil
}



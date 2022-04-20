package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	aliceKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}
	bobKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	charlieKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
	}

	transactionResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(aliceKey.PublicKey()).
		SetInitialBalance(hedera.NewHbar(5)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating account")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account creation receipt")
		return
	}

	aliceID := *transactionReceipt.AccountID

	transactionResponse, err = hedera.NewAccountCreateTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetKey(bobKey.PublicKey()).
		SetInitialBalance(hedera.NewHbar(5)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating second account")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving second account creation receipt")
		return
	}

	bobID := *transactionReceipt.AccountID

	transactionResponse, err = hedera.NewAccountCreateTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetKey(charlieKey.PublicKey()).
		SetInitialBalance(hedera.NewHbar(5)).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error creating second account")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving second account creation receipt")
		return
	}

	charlieID := *transactionReceipt.AccountID

	println("Alice's ID:", aliceID.String())
	println("Bob's ID:", bobID.String())
	println("Charlie's ID:", charlieID.String())
	println("Initial Balance:")
	err = printBalance(client, aliceID, bobID, charlieID, []hedera.AccountID{transactionResponse.NodeID})
	if err != nil {
		println(err.Error(), ": error retrieving balances")
		return
	}

	println("Approve an allowance of 2 Hbar with owner Alice and spender Bob")

	approvalFreeze, err := hedera.NewAccountAllowanceApproveTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		ApproveHbarAllowance(aliceID, bobID, hedera.NewHbar(2)).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account allowance approve transaction")
		return
	}

	approvalFreeze.Sign(aliceKey)

	transactionResponse, err = approvalFreeze.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account allowance approve transaction")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting account allowance receipt")
		return
	}

	err = printBalance(client, aliceID, bobID, charlieID, []hedera.AccountID{transactionResponse.NodeID})
	if err != nil {
		println(err.Error(), ": error retrieving balances")
		return
	}

	println("Transferring 1 Hbar from Alice to Charlie, but the transaction is signed _only_ by Bob (Bob is dipping into his allowance from Alice)")

	transferFreeze, err := hedera.NewTransferTransaction().
		AddApprovedHbarTransfer(aliceID, hedera.NewHbar(1).Negated(), true).
		AddHbarTransfer(charlieID, hedera.NewHbar(1)).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransactionID(hedera.TransactionIDGenerate(bobID)).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer transaction")
		return
	}

	transferFreeze.Sign(bobKey)

	transactionResponse, err = transferFreeze.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing transfer transaction")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving transfer transaction receipt")
		return
	}

	println("Transfer succeeded. Bob should now have 1 Hbar left in his allowance.")
	err = printBalance(client, aliceID, bobID, charlieID, []hedera.AccountID{transactionResponse.NodeID})
	if err != nil {
		println(err.Error(), ": error retrieving balances")
		return
	}

	println("Attempting to transfer 2 Hbar from Alice to Charlie using Bob's allowance.")
	println("This should fail, because there is only 1 Hbar left in Bob's allowance.")

	transferFreeze, err = hedera.NewTransferTransaction().
		AddApprovedHbarTransfer(aliceID, hedera.NewHbar(2).Negated(), true).
		AddHbarTransfer(charlieID, hedera.NewHbar(2)).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransactionID(hedera.TransactionIDGenerate(bobID)).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer transaction")
		return
	}

	transferFreeze.Sign(bobKey)

	transactionResponse, err = transferFreeze.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing transfer transaction")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ", Transfer failed as expected")
	}

	println("Adjusting Bob's allowance, increasing it by 2 Hbar. After this, Bob's allowance should be 3 Hbar.")

	allowanceAdjust, err := hedera.NewAccountAllowanceApproveTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		ApproveHbarAllowance(aliceID, bobID, hedera.NewHbar(2)).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account allowance adjust transaction")
		return
	}

	allowanceAdjust.Sign(aliceKey)

	transactionResponse, err = allowanceAdjust.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account allowance adjust transaction")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account allowance adjust receipt")
		return
	}

	err = printBalance(client, aliceID, bobID, charlieID, []hedera.AccountID{transactionResponse.NodeID})
	if err != nil {
		println(err.Error(), ": error retrieving balances")
		return
	}

	println("Attempting to transfer 2 Hbar from Alice to Charlie using Bob's allowance again.")
	println("This time it should succeed.")

	transferFreeze, err = hedera.NewTransferTransaction().
		AddApprovedHbarTransfer(aliceID, hedera.NewHbar(2).Negated(), true).
		AddHbarTransfer(charlieID, hedera.NewHbar(2)).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransactionID(hedera.TransactionIDGenerate(bobID)).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer transaction")
		return
	}

	transferFreeze.Sign(bobKey)

	transactionResponse, err = transferFreeze.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing transfer transaction")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ", error retrieving transfer transaction receipt")
		return
	}

	println("Transfer succeeded.")
	err = printBalance(client, aliceID, bobID, charlieID, []hedera.AccountID{transactionResponse.NodeID})
	if err != nil {
		println(err.Error(), ": error retrieving balances")
		return
	}

	println("Deleting Bob's allowance")

	approvalFreeze, err = hedera.NewAccountAllowanceApproveTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		ApproveHbarAllowance(aliceID, bobID, hedera.ZeroHbar).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing account allowance approve transaction")
		return
	}

	approvalFreeze.Sign(aliceKey)

	transactionResponse, err = approvalFreeze.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account allowance approve transaction")
		return
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting account allowance receipt")
		return
	}

	println("If Bob tries to use his allowance it should fail.")

	transferFreeze, err = hedera.NewTransferTransaction().
		AddApprovedHbarTransfer(aliceID, hedera.NewHbar(1).Negated(), true).
		AddHbarTransfer(charlieID, hedera.NewHbar(1)).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransactionID(hedera.TransactionIDGenerate(bobID)).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing transfer transaction")
		return
	}

	transferFreeze.Sign(bobKey)

	transactionResponse, err = transferFreeze.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing transfer transaction")
		return
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": Error just like expected")
	}

	println("\nCleaning up")

	accountDelete, err := hedera.NewAccountDeleteTransaction().
		SetAccountID(aliceID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing alice's account deletion")
		return
	}

	accountDelete.Sign(aliceKey)

	transactionResponse, err = accountDelete.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing alice's account deletion")
		return
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving alice's account deletion receipt")
		return
	}

	accountDelete, err = hedera.NewAccountDeleteTransaction().
		SetAccountID(bobID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing bob's account deletion")
		return
	}

	accountDelete.Sign(bobKey)

	transactionResponse, err = accountDelete.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing bob's account deletion")
		return
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving bob's account deletion receipt")
		return
	}

	accountDelete, err = hedera.NewAccountDeleteTransaction().
		SetAccountID(charlieID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		println(err.Error(), ": error freezing charlie's account deletion")
		return
	}

	accountDelete.Sign(charlieKey)

	transactionResponse, err = accountDelete.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing charlie's account deletion")
		return
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving charlie's account deletion receipt")
		return
	}

	err = client.Close()
	if err != nil {
		println(err.Error(), ": error closing client")
		return
	}
}

func printBalance(client *hedera.Client, alice hedera.AccountID, bob hedera.AccountID, charlie hedera.AccountID, nodeID []hedera.AccountID) error {
	println()

	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(alice).
		SetNodeAccountIDs(nodeID).
		Execute(client)
	if err != nil {
		return err
	}
	println("Alice's balance:", balance.Hbars.String())

	balance, err = hedera.NewAccountBalanceQuery().
		SetAccountID(bob).
		SetNodeAccountIDs(nodeID).
		Execute(client)
	if err != nil {
		return err
	}
	println("Bob's balance:", balance.Hbars.String())

	balance, err = hedera.NewAccountBalanceQuery().
		SetAccountID(charlie).
		SetNodeAccountIDs(nodeID).
		Execute(client)
	if err != nil {
		return err
	}
	println("Charlie's balance:", balance.Hbars.String())

	println()
	return nil
}

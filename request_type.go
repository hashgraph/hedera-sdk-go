package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"fmt"
)

type RequestType uint32

const (
	// UNSPECIFIED - Need to keep first value as unspecified because first element is ignored and not parsed (0 is ignored by parser)
	RequestTypeNone RequestType = 0
	// crypto transfe
	RequestTypeCryptoTransfer RequestType = 1
	// crypto update account
	RequestTypeCryptoUpdate RequestType = 2
	// crypto delete account
	RequestTypeCryptoDelete RequestType = 3
	// Add a livehash to a crypto account
	RequestTypeCryptoAddLiveHash RequestType = 4
	// Delete a livehash from a crypto account
	RequestTypeCryptoDeleteLiveHash RequestType = 5
	// Smart Contract Call
	RequestTypeContractCall RequestType = 6
	// Smart Contract Create Contract
	RequestTypeContractCreate RequestType = 7
	// Smart Contract update contract
	RequestTypeContractUpdate RequestType = 8
	// File Operation create file
	RequestTypeFileCreate RequestType = 9
	// File Operation append file
	RequestTypeFileAppend RequestType = 10
	// File Operation update file
	RequestTypeFileUpdate RequestType = 11
	// File Operation delete file
	RequestTypeFileDelete RequestType = 12
	// crypto get account balance
	RequestTypeCryptoGetAccountBalance RequestType = 13
	// crypto get account record
	RequestTypeCryptoGetAccountRecords RequestType = 14
	// Crypto get info
	RequestTypeCryptoGetInfo RequestType = 15
	// Smart Contract Call
	RequestTypeContractCallLocal RequestType = 16
	// Smart Contract get info
	RequestTypeContractGetInfo RequestType = 17
	// Smart Contract, get the byte code
	RequestTypeContractGetBytecode RequestType = 18
	// Smart Contract, get by _Solidity ID
	RequestTypeGetBySolidityID RequestType = 19
	// Smart Contract, get by key
	RequestTypeGetByKey RequestType = 20
	// Get a live hash from a crypto account
	RequestTypeCryptoGetLiveHash RequestType = 21
	// Crypto, get the stakers for the _Node
	RequestTypeCryptoGetStakers RequestType = 22
	// File Operations get file contents
	RequestTypeFileGetContents RequestType = 23
	// File Operations get the info of the file
	RequestTypeFileGetInfo RequestType = 24
	// Crypto get the transaction records
	RequestTypeTransactionGetRecord RequestType = 25
	// Contract get the transaction records
	RequestTypeContractGetRecords RequestType = 26
	// crypto create account
	RequestTypeCryptoCreate RequestType = 27
	// system delete file
	RequestTypeSystemDelete RequestType = 28
	// system undelete file
	RequestTypeSystemUndelete RequestType = 29
	// delete contract
	RequestTypeContractDelete RequestType = 30
	// freeze
	RequestTypeFreeze RequestType = 31
	// Create Tx Record
	RequestTypeCreateTransactionRecord RequestType = 32
	// Crypto Auto Renew
	RequestTypeCryptoAccountAutoRenew RequestType = 33
	// Contract Auto Renew
	RequestTypeContractAutoRenew RequestType = 34
	// Get Version
	RequestTypeGetVersionInfo RequestType = 35
	// Transaction Get Receipt
	RequestTypeTransactionGetReceipt RequestType = 36
	// Create Topic
	RequestTypeConsensusCreateTopic RequestType = 50
	// Update Topic
	RequestTypeConsensusUpdateTopic RequestType = 51
	// Delete Topic
	RequestTypeConsensusDeleteTopic RequestType = 52
	// Get Topic information
	RequestTypeConsensusGetTopicInfo RequestType = 53
	// Submit message to topic
	RequestTypeConsensusSubmitMessage RequestType = 54
	RequestTypeUncheckedSubmit        RequestType = 55
	// Create Token
	RequestTypeTokenCreate RequestType = 56
	// Get Token information
	RequestTypeTokenGetInfo RequestType = 58
	// Freeze Account
	RequestTypeTokenFreezeAccount RequestType = 59
	// Unfreeze Account
	RequestTypeTokenUnfreezeAccount RequestType = 60
	// Grant KYC to Account
	RequestTypeTokenGrantKycToAccount RequestType = 61
	// Revoke KYC from Account
	RequestTypeTokenRevokeKycFromAccount RequestType = 62
	// Delete Token
	RequestTypeTokenDelete RequestType = 63
	// Update Token
	RequestTypeTokenUpdate RequestType = 64
	// Mint tokens to treasury
	RequestTypeTokenMint RequestType = 65
	// Burn tokens from treasury
	RequestTypeTokenBurn RequestType = 66
	// Wipe token amount from Account holder
	RequestTypeTokenAccountWipe RequestType = 67
	// Associate tokens to an account
	RequestTypeTokenAssociateToAccount RequestType = 68
	// Dissociate tokens from an account
	RequestTypeTokenDissociateFromAccount RequestType = 69
	// Create Scheduled Transaction
	RequestTypeScheduleCreate RequestType = 70
	// Delete Scheduled Transaction
	RequestTypeScheduleDelete RequestType = 71
	// Sign Scheduled Transaction
	RequestTypeScheduleSign RequestType = 72
	// Get Scheduled Transaction Information
	RequestTypeScheduleGetInfo RequestType = 73
)

// String() returns a string representation of the status
func (requestType RequestType) String() string { // nolint
	switch requestType {
	case RequestTypeNone:
		return "NONE"
	case RequestTypeCryptoTransfer:
		return "CRYPTO_TRANSFER"
	case RequestTypeCryptoUpdate:
		return "CRYPTO_UPDATE"
	case RequestTypeCryptoDelete:
		return "CRYPTO_DELETE"
	case RequestTypeCryptoAddLiveHash:
		return "CRYPTO_ADD_LIVE_HASH"
	case RequestTypeCryptoDeleteLiveHash:
		return "CRYPTO_DELETE_LIVE_HASH"
	case RequestTypeContractCall:
		return "CONTRACT_CALL"
	case RequestTypeContractCreate:
		return "CONTRACT_CREATE"
	case RequestTypeContractUpdate:
		return "CONTRACT_UPDATE"
	case RequestTypeFileCreate:
		return "FILE_CREATE"
	case RequestTypeFileAppend:
		return "FILE_APPEND"
	case RequestTypeFileUpdate:
		return "FILE_UPDATE"
	case RequestTypeFileDelete:
		return "FILE_DELETE"
	case RequestTypeCryptoGetAccountBalance:
		return "CRYPTO_GET_ACCOUNT_BALANCE"
	case RequestTypeCryptoGetAccountRecords:
		return "CRYPTO_GET_ACCOUNT_RECORDS"
	case RequestTypeCryptoGetInfo:
		return "CRYPTO_GET_INFO"
	case RequestTypeContractCallLocal:
		return "CONTRACT_CALL_LOCAL"
	case RequestTypeContractGetInfo:
		return "CONTRACT_GET_INFO"
	case RequestTypeContractGetBytecode:
		return "CONTRACT_GET_BYTECODE"
	case RequestTypeGetBySolidityID:
		return "GET_BY_SOLIDITY_ID"
	case RequestTypeGetByKey:
		return "GET_BY_KEY"
	case RequestTypeCryptoGetLiveHash:
		return "CRYPTO_GET_LIVE_HASH"
	case RequestTypeCryptoGetStakers:
		return "CRYPTO_GET_STAKERS"
	case RequestTypeFileGetContents:
		return "FILE_GET_CONTENTS"
	case RequestTypeFileGetInfo:
		return "FILE_GET_INFO"
	case RequestTypeTransactionGetRecord:
		return "TRANSACTION_GET_RECORD"
	case RequestTypeContractGetRecords:
		return "CONTRACT_GET_RECORDS"
	case RequestTypeCryptoCreate:
		return "CRYPTO_CREATE"
	case RequestTypeSystemDelete:
		return "SYSTEM_DELETE"
	case RequestTypeSystemUndelete:
		return "SYSTEM_UNDELETE"
	case RequestTypeContractDelete:
		return "CONTRACT_DELETE"
	case RequestTypeFreeze:
		return "FREEZE"
	case RequestTypeCreateTransactionRecord:
		return "CREATE_TRANSACTION_RECORD"
	case RequestTypeCryptoAccountAutoRenew:
		return "CRYPTO_ACCOUNT_AUTO_RENEW"
	case RequestTypeContractAutoRenew:
		return "CONTRACT_AUTO_RENEW"
	case RequestTypeGetVersionInfo:
		return "GET_VERSION_INFO"
	case RequestTypeTransactionGetReceipt:
		return "TRANSACTION_GET_RECEIPT"
	case RequestTypeConsensusCreateTopic:
		return "CONSENSUS_CREATE_TOPIC"
	case RequestTypeConsensusUpdateTopic:
		return "CONSENSUS_UPDATE_TOPIC"
	case RequestTypeConsensusDeleteTopic:
		return "CONSENSUS_DELETE_TOPIC"
	case RequestTypeConsensusGetTopicInfo:
		return "CONSENSUS_GET_TOPIC_INFO"
	case RequestTypeConsensusSubmitMessage:
		return "CONSENSUS_SUBMIT_MESSAGE"
	case RequestTypeUncheckedSubmit:
		return "UNCHECKED_SUBMIT"
	case RequestTypeTokenCreate:
		return "TOKEN_CREATE"
	case RequestTypeTokenGetInfo:
		return "TOKEN_GET_INFO"
	case RequestTypeTokenFreezeAccount:
		return "TOKEN_FREEZE_ACCOUNT"
	case RequestTypeTokenUnfreezeAccount:
		return "TOKEN_UNFREEZE_ACCOUNT"
	case RequestTypeTokenGrantKycToAccount:
		return "TOKEN_GRANT_KYC_TO_ACCOUNT"
	case RequestTypeTokenRevokeKycFromAccount:
		return "TOKEN_REVOKE_KYC_TO_ACCOUNT"
	case RequestTypeTokenDelete:
		return "TOKEN_DELETE"
	case RequestTypeTokenUpdate:
		return "TOKEN_UPDATE"
	case RequestTypeTokenMint:
		return "TOKEN_MINT"
	case RequestTypeTokenBurn:
		return "TOKEN_BURN"
	case RequestTypeTokenAccountWipe:
		return "TOKEN_ACCOUNT_WIPE"
	case RequestTypeTokenAssociateToAccount:
		return "TOKEN_ASSOCIATE_TO_ACCOUNT"
	case RequestTypeTokenDissociateFromAccount:
		return "TOKEN_DISSOCIATE_FROM_ACCOUNT"
	case RequestTypeScheduleCreate:
		return "SCHEDULE_CREATE"
	case RequestTypeScheduleDelete:
		return "SCHEDULE_DELETE"
	case RequestTypeScheduleSign:
		return "SCHEDULE_SIGN"
	case RequestTypeScheduleGetInfo:
		return "SCHEDULE_GET_INFO"
	}

	panic(fmt.Sprintf("unreachable: RequestType.String() switch statement is non-exhaustive. RequestType: %v", uint32(requestType)))
}

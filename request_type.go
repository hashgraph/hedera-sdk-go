package hiero

// SPDX-License-Identifier: Apache-2.0

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
	// Get Token Account Nft Information
	RequestTypeTokenGetAccountNftInfos RequestType = 74
	// Get Token Nft Information
	RequestTypeTokenGetNftInfo RequestType = 75
	// Get Token Nft List Information
	RequestTypeTokenGetNftInfos RequestType = 76
	// Update a token's custom fee schedule, if permissible
	RequestTypeTokenFeeScheduleUpdate RequestType = 77
	// Get execution time(s) by TransactionID, if available
	RequestTypeNetworkGetExecutionTime RequestType = 78
	// Pause the Token
	RequestTypeTokenPause RequestType = 79
	// Unpause the Token
	RequestTypeTokenUnpause RequestType = 80
	// Approve allowance for a spender relative to the owner account
	RequestTypeCryptoApproveAllowance RequestType = 81
	// Deletes granted allowances on owner account
	RequestTypeCryptoDeleteAllowance RequestType = 82
	// Gets all the information about an account, including balance and allowances
	RequestTypeGetAccountDetails RequestType = 83
	// Ethereum Transaction
	RequestTypeEthereumTransaction RequestType = 84
	// Updates the staking info at the end of staking period
	RequestTypeNodeStakeUpdate RequestType = 85
	// Generates a pseudorandom number
	RequestTypePrng RequestType = 86
	// Get a record for a transaction (lasts 180 seconds)
	RequestTypeTransactionGetFastRecord RequestType = 87
	// Update the metadata of one or more NFT's of a specific token type
	RequestTypeTokenUpdateNfts RequestType = 88
	// Create a new node
	RequestTypeNodeCreate RequestType = 89
	// Update an existing node
	RequestTypeNodeUpdate RequestType = 90
	// Delete a node
	RequestTypeNodeDelete RequestType = 91
	// Transfer token balances to treasury
	RequestTypeTokenReject RequestType = 92
	// Airdrop tokens to accounts
	RequestTypeTokenAirdrop RequestType = 93
	// Remove pending airdrops from state
	RequestTypeTokenCancelAirdrop RequestType = 94
	// Claim pending airdrops
	RequestTypeTokenClaimAirdrop RequestType = 95
	// TSS Messages for a candidate roster
	RequestTypeTssMessage RequestType = 96
	// Vote on TSS validity
	RequestTypeTssVote RequestType = 97
	// Node's signature of block hash using TSS private share
	RequestTypeTssShareSignature RequestType = 98
	// Submit node public TSS encryption key
	RequestTypeTssEncryptionKey RequestType = 99
	// Submit signature of state root hash
	RequestTypeStateSignatureTransaction RequestType = 100
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
	case RequestTypeTokenGetAccountNftInfos:
		return "TOKEN_GET_ACCOUNT_NFT_INFOS"
	case RequestTypeTokenGetNftInfo:
		return "TOKEN_GET_NFT_INFO"
	case RequestTypeTokenGetNftInfos:
		return "TOKEN_GET_NFT_INFOS"
	case RequestTypeTokenFeeScheduleUpdate:
		return "TOKEN_FEE_SCHEDULE_UPDATE"
	case RequestTypeNetworkGetExecutionTime:
		return "NETWORK_GET_EXECUTION_TIME"
	case RequestTypeTokenPause:
		return "TOKEN_PAUSE"
	case RequestTypeTokenUnpause:
		return "TOKEN_UNPAUSE"
	case RequestTypeCryptoApproveAllowance:
		return "CRYPTO_APPROVE_ALLOWANCE"
	case RequestTypeCryptoDeleteAllowance:
		return "CRYPTO_DELETE_ALLOWANCE"
	case RequestTypeGetAccountDetails:
		return "GET_ACCOUNT_DETAILS"
	case RequestTypeEthereumTransaction:
		return "ETHEREUM_TRANSACTION"
	case RequestTypeNodeStakeUpdate:
		return "NODE_STAKE_UPDATE"
	case RequestTypePrng:
		return "PRNG"
	case RequestTypeTransactionGetFastRecord:
		return "TRANSACTION_GET_FAST_RECORD"
	case RequestTypeTokenUpdateNfts:
		return "TOKEN_UPDATE_NFTS"
	case RequestTypeNodeCreate:
		return "NODE_CREATE"
	case RequestTypeNodeUpdate:
		return "NODE_UPDATE"
	case RequestTypeNodeDelete:
		return "NODE_DELETE"
	case RequestTypeTokenReject:
		return "TOKEN_REJECT"
	case RequestTypeTokenAirdrop:
		return "TOKEN_AIRDROP"
	case RequestTypeTokenCancelAirdrop:
		return "TOKEN_CANCEL_AIRDROP"
	case RequestTypeTokenClaimAirdrop:
		return "TOKEN_CLAIM_AIRDROP"
	case RequestTypeTssMessage:
		return "TSS_MESSAGE"
	case RequestTypeTssVote:
		return "TSS_VOTE"
	case RequestTypeTssShareSignature:
		return "TSS_SHARE_SIGNATURE"
	case RequestTypeTssEncryptionKey:
		return "TSS_ENCRYPTION_KEY"
	case RequestTypeStateSignatureTransaction:
		return "STATE_SIGNATURE_TRANSACTION"
	}

	panic(fmt.Sprintf("unreachable: RequestType.String() switch statement is non-exhaustive. RequestType: %v", uint32(requestType)))
}

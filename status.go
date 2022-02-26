package hedera

import "fmt"

type Status uint32

const (
	StatusOk                                                       Status = 0
	StatusInvalidTransaction                                       Status = 1
	StatusPayerAccountNotFound                                     Status = 2
	StatusInvalidNodeAccount                                       Status = 3
	StatusTransactionExpired                                       Status = 4
	StatusInvalidTransactionStart                                  Status = 5
	StatusInvalidTransactionDuration                               Status = 6
	StatusInvalidSignature                                         Status = 7
	StatusMemoTooLong                                              Status = 8
	StatusInsufficientTxFee                                        Status = 9
	StatusInsufficientPayerBalance                                 Status = 10
	StatusDuplicateTransaction                                     Status = 11
	StatusBusy                                                     Status = 12
	StatusNotSupported                                             Status = 13
	StatusInvalidFileID                                            Status = 14
	StatusInvalidAccountID                                         Status = 15
	StatusInvalidContractID                                        Status = 16
	StatusInvalidTransactionID                                     Status = 17
	StatusReceiptNotFound                                          Status = 18
	StatusRecordNotFound                                           Status = 19
	StatusInvalidSolidityID                                        Status = 20
	StatusUnknown                                                  Status = 21
	StatusSuccess                                                  Status = 22
	StatusFailInvalid                                              Status = 23
	StatusFailFee                                                  Status = 24
	StatusFailBalance                                              Status = 25
	StatusKeyRequired                                              Status = 26
	StatusBadEncoding                                              Status = 27
	StatusInsufficientAccountBalance                               Status = 28
	StatusInvalidSolidityAddress                                   Status = 29
	StatusInsufficientGas                                          Status = 30
	StatusContractSizeLimitExceeded                                Status = 31
	StatusLocalCallModificationException                           Status = 32
	StatusContractRevertExecuted                                   Status = 33
	StatusContractExecutionException                               Status = 34
	StatusInvalidReceivingNodeAccount                              Status = 35
	StatusMissingQueryHeader                                       Status = 36
	StatusAccountUpdateFailed                                      Status = 37
	StatusInvalidKeyEncoding                                       Status = 38
	StatusNullSolidityAddress                                      Status = 39
	StatusContractUpdateFailed                                     Status = 40
	StatusInvalidQueryHeader                                       Status = 41
	StatusInvalidFeeSubmitted                                      Status = 42
	StatusInvalidPayerSignature                                    Status = 43
	StatusKeyNotProvided                                           Status = 44
	StatusInvalidExpirationTime                                    Status = 45
	StatusNoWaclKey                                                Status = 46
	StatusFileContentEmpty                                         Status = 47
	StatusInvalidAccountAmounts                                    Status = 48
	StatusEmptyTransactionBody                                     Status = 49
	StatusInvalidTransactionBody                                   Status = 50
	StatusInvalidSignatureTypeMismatchingKey                       Status = 51
	StatusInvalidSignatureCountMismatchingKey                      Status = 52
	StatusEmptyLiveHashBody                                        Status = 53
	StatusEmptyLiveHash                                            Status = 54
	StatusEmptyLiveHashKeys                                        Status = 55
	StatusInvalidLiveHashSize                                      Status = 56
	StatusEmptyQueryBody                                           Status = 57
	StatusEmptyLiveHashQuery                                       Status = 58
	StatusLiveHashNotFound                                         Status = 59
	StatusAccountIDDoesNotExist                                    Status = 60
	StatusLiveHashAlreadyExists                                    Status = 61
	StatusInvalidFileWacl                                          Status = 62
	StatusSerializationFailed                                      Status = 63
	StatusTransactionOversize                                      Status = 64
	StatusTransactionTooManyLayers                                 Status = 65
	StatusContractDeleted                                          Status = 66
	StatusPlatformNotActive                                        Status = 67
	StatusKeyPrefixMismatch                                        Status = 68
	StatusPlatformTransactionNotCreated                            Status = 69
	StatusInvalidRenewalPeriod                                     Status = 70
	StatusInvalidPayerAccountID                                    Status = 71
	StatusAccountDeleted                                           Status = 72
	StatusFileDeleted                                              Status = 73
	StatusAccountRepeatedInAccountAmounts                          Status = 74
	StatusSettingNegativeAccountBalance                            Status = 75
	StatusObtainerRequired                                         Status = 76
	StatusObtainerSameContractID                                   Status = 77
	StatusObtainerDoesNotExist                                     Status = 78
	StatusModifyingImmutableContract                               Status = 79
	StatusFileSystemException                                      Status = 80
	StatusAutorenewDurationNotInRange                              Status = 81
	StatusErrorDecodingBytestring                                  Status = 82
	StatusContractFileEmpty                                        Status = 83
	StatusContractBytecodeEmpty                                    Status = 84
	StatusInvalidInitialBalance                                    Status = 85
	StatusInvalidReceiveRecordThreshold                            Status = 86
	StatusInvalidSendRecordThreshold                               Status = 87
	StatusAccountIsNotGenesisAccount                               Status = 88
	StatusPayerAccountUnauthorized                                 Status = 89
	StatusInvalidFreezeTransactionBody                             Status = 90
	StatusFreezeTransactionBodyNotFound                            Status = 91
	StatusTransferListSizeLimitExceeded                            Status = 92
	StatusResultSizeLimitExceeded                                  Status = 93
	StatusNotSpecialAccount                                        Status = 94
	StatusContractNegativeGas                                      Status = 95
	StatusContractNegativeValue                                    Status = 96
	StatusInvalidFeeFile                                           Status = 97
	StatusInvalidExchangeRateFile                                  Status = 98
	StatusInsufficientLocalCallGas                                 Status = 99
	StatusEntityNotAllowedToDelete                                 Status = 100
	StatusAuthorizationFailed                                      Status = 101
	StatusFileUploadedProtoInvalid                                 Status = 102
	StatusFileUploadedProtoNotSavedToDisk                          Status = 103
	StatusFeeScheduleFilePartUploaded                              Status = 104
	StatusExchangeRateChangeLimitExceeded                          Status = 105
	StatusMaxContractStorageExceeded                               Status = 106
	StatusTransferAccountSameAsDeleteAccount                       Status = 107
	StatusTotalLedgerBalanceInvalid                                Status = 108
	StatusExpirationReductionNotAllowed                            Status = 110
	StatusMaxGasLimitExceeded                                      Status = 111
	StatusMaxFileSizeExceeded                                      Status = 112
	StatusReceiverSigRequired                                      Status = 113
	StatusInvalidTopicID                                           Status = 150
	StatusInvalidAdminKey                                          Status = 155
	StatusInvalidSubmitKey                                         Status = 156
	StatusUnauthorized                                             Status = 157
	StatusInvalidTopicMessage                                      Status = 158
	StatusInvalidAutorenewAccount                                  Status = 159
	StatusAutorenewAccountNotAllowed                               Status = 160
	StatusTopicExpired                                             Status = 162
	StatusInvalidChunkNumber                                       Status = 163
	StatusInvalidChunkTransactionID                                Status = 164
	StatusAccountFrozenForToken                                    Status = 165
	StatusTokensPerAccountLimitExceeded                            Status = 166
	StatusInvalidTokenID                                           Status = 167
	StatusInvalidTokenDecimals                                     Status = 168
	StatusInvalidTokenInitialSupply                                Status = 169
	StatusInvalidTreasuryAccountForToken                           Status = 170
	StatusInvalidTokenSymbol                                       Status = 171
	StatusTokenHasNoFreezeKey                                      Status = 172
	StatusTransfersNotZeroSumForToken                              Status = 173
	StatusMissingTokenSymbol                                       Status = 174
	StatusTokenSymbolTooLong                                       Status = 175
	StatusAccountKycNotGrantedForToken                             Status = 176
	StatusTokenHasNoKycKey                                         Status = 177
	StatusInsufficientTokenBalance                                 Status = 178
	StatusTokenWasDeleted                                          Status = 179
	StatusTokenHasNoSupplyKey                                      Status = 180
	StatusTokenHasNoWipeKey                                        Status = 181
	StatusInvalidTokenMintAmount                                   Status = 182
	StatusInvalidTokenBurnAmount                                   Status = 183
	StatusTokenNotAssociatedToAccount                              Status = 184
	StatusCannotWipeTokenTreasuryAccount                           Status = 185
	StatusInvalidKycKey                                            Status = 186
	StatusInvalidWipeKey                                           Status = 187
	StatusInvalidFreezeKey                                         Status = 188
	StatusInvalidSupplyKey                                         Status = 189
	StatusMissingTokenName                                         Status = 190
	StatusTokenNameTooLong                                         Status = 191
	StatusInvalidWipingAmount                                      Status = 192
	StatusTokenIsImmutable                                         Status = 193
	StatusTokenAlreadyAssociatedToAccount                          Status = 194
	StatusTransactionRequiresZeroTokenBalances                     Status = 195
	StatusAccountIsTreasury                                        Status = 196
	StatusTokenIDRepeatedInTokenList                               Status = 197
	StatusTokenTransferListSizeLimitExceeded                       Status = 198
	StatusEmptyTokenTransferBody                                   Status = 199
	StatusEmptyTokenTransferAccountAmounts                         Status = 200
	StatusInvalidScheduleID                                        Status = 201
	StatusScheduleIsImmutable                                      Status = 202
	StatusInvalidSchedulePayerID                                   Status = 203
	StatusInvalidScheduleAccountID                                 Status = 204
	StatusNoNewValidSignatures                                     Status = 205
	StatusUnresolvableRequiredSigners                              Status = 206
	StatusScheduledTransactionNotInWhitelist                       Status = 207
	StatusSomeSignaturesWereInvalid                                Status = 208
	StatusTransactionIDFieldNotAllowed                             Status = 209
	StatusIdenticalScheduleAlreadyCreated                          Status = 210
	StatusInvalidZeroByteInString                                  Status = 211
	StatusScheduleAlreadyDeleted                                   Status = 212
	StatusScheduleAlreadyExecuted                                  Status = 213
	StatusMessageSizeTooLarge                                      Status = 214
	StatusOperationRepeatedInBucketGroups                          Status = 215
	StatusBucketCapacityOverflow                                   Status = 216
	StatusNodeCapacityNotSufficientForOperation                    Status = 217
	StatusBucketHasNoThrottleGroups                                Status = 218
	StatusThrottleGroupHasZeroOpsPerSec                            Status = 219
	StatusSuccessButMissingExpectedOperation                       Status = 220
	StatusUnparseableThrottleDefinitions                           Status = 221
	StatusInvalidThrottleDefinitions                               Status = 222
	StatusAccountExpiredAndPendingRemoval                          Status = 223
	StatusInvalidTokenMaxSupply                                    Status = 224
	StatusInvalidTokenNftSerialNumber                              Status = 225
	StatusInvalidNftID                                             Status = 226
	StatusMetadataTooLong                                          Status = 227
	StatusBatchSizeLimitExceeded                                   Status = 228
	StatusInvalidQueryRange                                        Status = 229
	StatusFractionDividesByZero                                    Status = 230
	StatusInsufficientPayerBalanceForCustomFee                     Status = 231
	StatusCustomFeesListTooLong                                    Status = 232
	StatusInvalidCustomFeeCollector                                Status = 233
	StatusInvalidTokenIDInCustomFees                               Status = 234
	StatusTokenNotAssociatedToFeeCollector                         Status = 235
	StatusTokenMaxSupplyReached                                    Status = 236
	StatusSenderDoesNotOwnNftSerialNo                              Status = 237
	StatusCustomFeeNotFullySpecified                               Status = 238
	StatusCustomFeeMustBePositive                                  Status = 239
	StatusTokenHasNoFeeScheduleKey                                 Status = 240
	StatusCustomFeeOutsideNumericRange                             Status = 241
	StatusRoyaltyFractionCannotExceedOne                           Status = 242
	StatusFractionalFeeMaxAmountLessThanMinAmount                  Status = 243
	StatusCustomScheduleAlreadyHasNoFees                           Status = 244
	StatusCustomFeeDenominationMustBeFungibleCommon                Status = 245
	StatusCustomFractionalFeeOnlyAllowedForFungibleCommon          Status = 246
	StatusInvalidCustomFeeScheduleKey                              Status = 247
	StatusInvalidTokenMintMetadata                                 Status = 248
	StatusInvalidTokenBurnMetadata                                 Status = 249
	StatusCurrentTreasuryStillOwnsNfts                             Status = 250
	StatusAccountStillOwnsNfts                                     Status = 251
	StatusTreasuryMustOwnBurnedNft                                 Status = 252
	StatusAccountDoesNotOwnWipedNft                                Status = 253
	StatusAccountAmountTransfersOnlyAllowedForFungibleCommon       Status = 254
	StatusMaxNftsInPriceRegimeHaveBeenMinted                       Status = 255
	StatusPayerAccountDeleted                                      Status = 256
	StatusCustomFeeChargingExceededMaxRecursionDepth               Status = 257
	StatusCustomFeeChargingExceededMaxAccountAmounts               Status = 258
	StatusInsufficientSenderAccountBalanceForCustomFee             Status = 259
	StatusSerialNumberLimitReached                                 Status = 260
	StatusCustomRoyaltyFeeOnlyAllowedForNonFungibleUnique          Status = 261
	StatusNoRemainingAutomaticAssociations                         Status = 262
	StatusExistingAutomaticAssociationsExceedGivenLimit            Status = 263
	StatusRequestedNumAutomaticAssociationsExceedsAssociationLimit Status = 264
	StatusTokenIsPaused                                            Status = 265
	StatusTokenHasNoPauseKey                                       Status = 266
	StatusInvalidPauseKey                                          Status = 267
	StatusFreezeUpdateFileDoesNotExist                             Status = 268
	StatusFreezeUpdateFileHashDoesNotMatch                         Status = 269
	StatusNoUpgradeHasBeenPrepared                                 Status = 270
	StatusNoFreezeIsScheduled                                      Status = 271
	StatusUpdateFileHashChangedSincePrepareUpgrade                 Status = 272
	StatusFreezeStartTimeMustBeFuture                              Status = 273
	StatusPreparedUpdateFileIsImmutable                            Status = 274
	StatusFreezeAlreadyScheduled                                   Status = 275
	StatusFreezeUpgradeInProgress                                  Status = 276
	StatusUpdateFileIDDoesNotMatchPrepared                         Status = 277
	StatusUpdateFileHashDoesNotMatchPrepared                       Status = 278
	StatusConsensusGasExhausted                                    Status = 279
	StatusRevertedSuccess                                          Status = 280
	StatusMaxStorageInPriceRegimeHasBeenUsed                       Status = 281
	StatusInvalidAliasKey                                          Status = 282
	StatusUnexpectedTokenDecimals                                  Status = 283
	StatusInvalidProxyAccountID                                    Status = 284
	StatusInvalidTransferAccountID                                 Status = 285
	StatusInvalidFeeCollectorAccountID                             Status = 286
	StatusAliasIsImmutable                                         Status = 287
	StatusSpenderAccountSameAsOwner                                Status = 288
	StatusAmountExceedsTokenMaxSupply                              Status = 289
	StatusNegativeAllowanceAmount                                  Status = 290
	StatusCannotApproveForAllFungibleCommon                        Status = 291
	StatusSpenderDoesNotHaveAllowance                              Status = 292
	StatusAmountExceedsAllowance                                   Status = 293
	StatusMaxAllowancesExceeded                                    Status = 294
	StatusEmptyAllowances                                          Status = 295
	StatusSpenderAccountRepeatedInAllowance                        Status = 296
	StatusRepeatedSerialNumsInNftAllowances                        Status = 297
	StatusFungibleTokenInNftAllowances                             Status = 298
	StatusNftInFungibleTokenAllowances                             Status = 299
	StatusInvalidAllowanceOwnerID                                  Status = 300
	StatusInvalidAllowanceSpenderID                                Status = 301
)

// String() returns a string representation of the status
func (status Status) String() string { // nolint
	switch status {
	case StatusOk:
		return "OK"
	case StatusInvalidTransaction:
		return "INVALID_TRANSACTION"
	case StatusPayerAccountNotFound:
		return "PAYER_ACCOUNT_NOT_FOUND"
	case StatusInvalidNodeAccount:
		return "INVALID_NODE_ACCOUNT"
	case StatusTransactionExpired:
		return "TRANSACTION_EXPIRED"
	case StatusInvalidTransactionStart:
		return "INVALID_TRANSACTION_START"
	case StatusInvalidTransactionDuration:
		return "INVALID_TRANSACTION_DURATION"
	case StatusInvalidSignature:
		return "INVALID_SIGNATURE"
	case StatusMemoTooLong:
		return "MEMO_TOO_LONG"
	case StatusInsufficientTxFee:
		return "INSUFFICIENT_TX_FEE"
	case StatusInsufficientPayerBalance:
		return "INSUFFICIENT_PAYER_BALANCE"
	case StatusDuplicateTransaction:
		return "DUPLICATE_TRANSACTION"
	case StatusBusy:
		return "BUSY"
	case StatusNotSupported:
		return "NOT_SUPPORTED"
	case StatusInvalidFileID:
		return "INVALID_FILE_ID"
	case StatusInvalidAccountID:
		return "INVALID_ACCOUNT_ID"
	case StatusInvalidContractID:
		return "INVALID_CONTRACT_ID"
	case StatusInvalidTransactionID:
		return "INVALID_TRANSACTION_ID"
	case StatusReceiptNotFound:
		return "RECEIPT_NOT_FOUND"
	case StatusRecordNotFound:
		return "RECORD_NOT_FOUND"
	case StatusInvalidSolidityID:
		return "INVALID_SOLIDITY_ID"
	case StatusUnknown:
		return "UNKNOWN"
	case StatusSuccess:
		return "SUCCESS"
	case StatusFailInvalid:
		return "FAIL_INVALID"
	case StatusFailFee:
		return "FAIL_FEE"
	case StatusFailBalance:
		return "FAIL_BALANCE"
	case StatusKeyRequired:
		return "KEY_REQUIRED"
	case StatusBadEncoding:
		return "BAD_ENCODING"
	case StatusInsufficientAccountBalance:
		return "INSUFFICIENT_ACCOUNT_BALANCE"
	case StatusInvalidSolidityAddress:
		return "INVALID_SOLIDITY_ADDRESS"
	case StatusInsufficientGas:
		return "INSUFFICIENT_GAS"
	case StatusContractSizeLimitExceeded:
		return "CONTRACT_SIZE_LIMIT_EXCEEDED"
	case StatusLocalCallModificationException:
		return "LOCAL_CALL_MODIFICATION_EXCEPTION"
	case StatusContractRevertExecuted:
		return "CONTRACT_REVERT_EXECUTED"
	case StatusContractExecutionException:
		return "CONTRACT_EXECUTION_EXCEPTION"
	case StatusInvalidReceivingNodeAccount:
		return "INVALID_RECEIVING_NODE_ACCOUNT"
	case StatusMissingQueryHeader:
		return "MISSING_QUERY_HEADER"
	case StatusAccountUpdateFailed:
		return "ACCOUNT_UPDATE_FAILED"
	case StatusInvalidKeyEncoding:
		return "INVALID_KEY_ENCODING"
	case StatusNullSolidityAddress:
		return "NULL_SOLIDITY_ADDRESS"
	case StatusContractUpdateFailed:
		return "CONTRACT_UPDATE_FAILED"
	case StatusInvalidQueryHeader:
		return "INVALID_QUERY_HEADER"
	case StatusInvalidFeeSubmitted:
		return "INVALID_FEE_SUBMITTED"
	case StatusInvalidPayerSignature:
		return "INVALID_PAYER_SIGNATURE"
	case StatusKeyNotProvided:
		return "KEY_NOT_PROVIDED"
	case StatusInvalidExpirationTime:
		return "INVALID_EXPIRATION_TIME"
	case StatusNoWaclKey:
		return "NO_WACL_KEY"
	case StatusFileContentEmpty:
		return "FILE_CONTENT_EMPTY"
	case StatusInvalidAccountAmounts:
		return "INVALID_ACCOUNT_AMOUNTS"
	case StatusEmptyTransactionBody:
		return "EMPTY_TRANSACTION_BODY"
	case StatusInvalidTransactionBody:
		return "INVALID_TRANSACTION_BODY"
	case StatusInvalidSignatureTypeMismatchingKey:
		return "INVALID_SIGNATURE_TYPE_MISMATCHING_KEY"
	case StatusInvalidSignatureCountMismatchingKey:
		return "INVALID_SIGNATURE_COUNT_MISMATCHING_KEY"
	case StatusEmptyLiveHashBody:
		return "EMPTY_LIVE_HASH_BODY"
	case StatusEmptyLiveHash:
		return "EMPTY_LIVE_HASH"
	case StatusEmptyLiveHashKeys:
		return "EMPTY_LIVE_HASH_KEYS"
	case StatusInvalidLiveHashSize:
		return "INVALID_LIVE_HASH_SIZE"
	case StatusEmptyQueryBody:
		return "EMPTY_QUERY_BODY"
	case StatusEmptyLiveHashQuery:
		return "EMPTY_LIVE_HASH_QUERY"
	case StatusLiveHashNotFound:
		return "LIVE_HASH_NOT_FOUND"
	case StatusAccountIDDoesNotExist:
		return "ACCOUNT_ID_DOES_NOT_EXIST"
	case StatusLiveHashAlreadyExists:
		return "LIVE_HASH_ALREADY_EXISTS"
	case StatusInvalidFileWacl:
		return "INVALID_FILE_WACL"
	case StatusSerializationFailed:
		return "SERIALIZATION_FAILED"
	case StatusTransactionOversize:
		return "TRANSACTION_OVERSIZE"
	case StatusTransactionTooManyLayers:
		return "TRANSACTION_TOO_MANY_LAYERS"
	case StatusContractDeleted:
		return "CONTRACT_DELETED"
	case StatusPlatformNotActive:
		return "PLATFORM_NOT_ACTIVE"
	case StatusKeyPrefixMismatch:
		return "KEY_PREFIX_MISMATCH"
	case StatusPlatformTransactionNotCreated:
		return "PLATFORM_TRANSACTION_NOT_CREATED"
	case StatusInvalidRenewalPeriod:
		return "INVALID_RENEWAL_PERIOD"
	case StatusInvalidPayerAccountID:
		return "INVALID_PAYER_ACCOUNT_ID"
	case StatusAccountDeleted:
		return "ACCOUNT_DELETED"
	case StatusFileDeleted:
		return "FILE_DELETED"
	case StatusAccountRepeatedInAccountAmounts:
		return "ACCOUNT_REPEATED_IN_ACCOUNT_AMOUNTS"
	case StatusSettingNegativeAccountBalance:
		return "SETTING_NEGATIVE_ACCOUNT_BALANCE"
	case StatusObtainerRequired:
		return "OBTAINER_REQUIRED"
	case StatusObtainerSameContractID:
		return "OBTAINER_SAME_CONTRACT_ID"
	case StatusObtainerDoesNotExist:
		return "OBTAINER_DOES_NOT_EXIST"
	case StatusModifyingImmutableContract:
		return "MODIFYING_IMMUTABLE_CONTRACT"
	case StatusFileSystemException:
		return "FILE_SYSTEM_EXCEPTION"
	case StatusAutorenewDurationNotInRange:
		return "AUTORENEW_DURATION_NOT_IN_RANGE"
	case StatusErrorDecodingBytestring:
		return "ERROR_DECODING_BYTESTRING"
	case StatusContractFileEmpty:
		return "CONTRACT_FILE_EMPTY"
	case StatusContractBytecodeEmpty:
		return "CONTRACT_BYTECODE_EMPTY"
	case StatusInvalidInitialBalance:
		return "INVALID_INITIAL_BALANCE"
	case StatusInvalidReceiveRecordThreshold:
		return "INVALID_RECEIVE_RECORD_THRESHOLD"
	case StatusInvalidSendRecordThreshold:
		return "INVALID_SEND_RECORD_THRESHOLD"
	case StatusAccountIsNotGenesisAccount:
		return "ACCOUNT_IS_NOT_GENESIS_ACCOUNT"
	case StatusPayerAccountUnauthorized:
		return "PAYER_ACCOUNT_UNAUTHORIZED"
	case StatusInvalidFreezeTransactionBody:
		return "INVALID_FREEZE_TRANSACTION_BODY"
	case StatusFreezeTransactionBodyNotFound:
		return "FREEZE_TRANSACTION_BODY_NOT_FOUND"
	case StatusTransferListSizeLimitExceeded:
		return "TRANSFER_LIST_SIZE_LIMIT_EXCEEDED"
	case StatusResultSizeLimitExceeded:
		return "RESULT_SIZE_LIMIT_EXCEEDED"
	case StatusNotSpecialAccount:
		return "NOT_SPECIAL_ACCOUNT"
	case StatusContractNegativeGas:
		return "CONTRACT_NEGATIVE_GAS"
	case StatusContractNegativeValue:
		return "CONTRACT_NEGATIVE_VALUE"
	case StatusInvalidFeeFile:
		return "INVALID_FEE_FILE"
	case StatusInvalidExchangeRateFile:
		return "INVALID_EXCHANGE_RATE_FILE"
	case StatusInsufficientLocalCallGas:
		return "INSUFFICIENT_LOCAL_CALL_GAS"
	case StatusEntityNotAllowedToDelete:
		return "ENTITY_NOT_ALLOWED_TO_DELETE"
	case StatusAuthorizationFailed:
		return "AUTHORIZATION_FAILED"
	case StatusFileUploadedProtoInvalid:
		return "FILE_UPLOADED_PROTO_INVALID"
	case StatusFileUploadedProtoNotSavedToDisk:
		return "FILE_UPLOADED_PROTO_NOT_SAVED_TO_DISK"
	case StatusFeeScheduleFilePartUploaded:
		return "FEE_SCHEDULE_FILE_PART_UPLOADED"
	case StatusExchangeRateChangeLimitExceeded:
		return "EXCHANGE_RATE_CHANGE_LIMIT_EXCEEDED"
	case StatusMaxContractStorageExceeded:
		return "MAX_CONTRACT_STORAGE_EXCEEDED"
	case StatusTransferAccountSameAsDeleteAccount:
		return "TRANSFER_ACCOUNT_SAME_AS_DELETE_ACCOUNT"
	case StatusTotalLedgerBalanceInvalid:
		return "TOTAL_LEDGER_BALANCE_INVALID"
	case StatusExpirationReductionNotAllowed:
		return "EXPIRATION_REDUCTION_NOT_ALLOWED"
	case StatusMaxGasLimitExceeded:
		return "MAX_GAS_LIMIT_EXCEEDED"
	case StatusMaxFileSizeExceeded:
		return "MAX_FILE_SIZE_EXCEEDED"
	case StatusReceiverSigRequired:
		return "RECEIVER_SIG_REQUIRED"
	case StatusInvalidTopicID:
		return "INVALID_TOPIC_ID"
	case StatusInvalidAdminKey:
		return "INVALID_ADMIN_KEY"
	case StatusInvalidSubmitKey:
		return "INVALID_SUBMIT_KEY"
	case StatusUnauthorized:
		return "UNAUTHORIZED"
	case StatusInvalidTopicMessage:
		return "INVALID_TOPIC_MESSAGE"
	case StatusInvalidAutorenewAccount:
		return "INVALID_AUTORENEW_ACCOUNT"
	case StatusAutorenewAccountNotAllowed:
		return "AUTORENEW_ACCOUNT_NOT_ALLOWED"
	case StatusTopicExpired:
		return "TOPIC_EXPIRED"
	case StatusInvalidChunkNumber:
		return "INVALID_CHUNK_NUMBER"
	case StatusInvalidChunkTransactionID:
		return "INVALID_CHUNK_TRANSACTION_ID"
	case StatusAccountFrozenForToken:
		return "ACCOUNT_FROZEN_FOR_TOKEN"
	case StatusTokensPerAccountLimitExceeded:
		return "TOKENS_PER_ACCOUNT_LIMIT_EXCEEDED"
	case StatusInvalidTokenID:
		return "INVALID_TOKEN_ID"
	case StatusInvalidTokenDecimals:
		return "INVALID_TOKEN_DECIMALS"
	case StatusInvalidTokenInitialSupply:
		return "INVALID_TOKEN_INITIAL_SUPPLY"
	case StatusInvalidTreasuryAccountForToken:
		return "INVALID_TREASURY_ACCOUNT_FOR_TOKEN"
	case StatusInvalidTokenSymbol:
		return "INVALID_TOKEN_SYMBOL"
	case StatusTokenHasNoFreezeKey:
		return "TOKEN_HAS_NO_FREEZE_KEY"
	case StatusTransfersNotZeroSumForToken:
		return "TRANSFERS_NOT_ZERO_SUM_FOR_TOKEN"
	case StatusMissingTokenSymbol:
		return "MISSING_TOKEN_SYMBOL"
	case StatusTokenSymbolTooLong:
		return "TOKEN_SYMBOL_TOO_LONG"
	case StatusAccountKycNotGrantedForToken:
		return "ACCOUNT_KYC_NOT_GRANTED_FOR_TOKEN"
	case StatusTokenHasNoKycKey:
		return "TOKEN_HAS_NO_KYC_KEY"
	case StatusInsufficientTokenBalance:
		return "INSUFFICIENT_TOKEN_BALANCE"
	case StatusTokenWasDeleted:
		return "TOKEN_WAS_DELETED"
	case StatusTokenHasNoSupplyKey:
		return "TOKEN_HAS_NO_SUPPLY_KEY"
	case StatusTokenHasNoWipeKey:
		return "TOKEN_HAS_NO_WIPE_KEY"
	case StatusInvalidTokenMintAmount:
		return "INVALID_TOKEN_MINT_AMOUNT"
	case StatusInvalidTokenBurnAmount:
		return "INVALID_TOKEN_BURN_AMOUNT"
	case StatusTokenNotAssociatedToAccount:
		return "TOKEN_NOT_ASSOCIATED_TO_ACCOUNT"
	case StatusCannotWipeTokenTreasuryAccount:
		return "CANNOT_WIPE_TOKEN_TREASURY_ACCOUNT"
	case StatusInvalidKycKey:
		return "INVALID_KYC_KEY"
	case StatusInvalidWipeKey:
		return "INVALID_WIPE_KEY"
	case StatusInvalidFreezeKey:
		return "INVALID_FREEZE_KEY"
	case StatusInvalidSupplyKey:
		return "INVALID_SUPPLY_KEY"
	case StatusMissingTokenName:
		return "MISSING_TOKEN_NAME"
	case StatusTokenNameTooLong:
		return "TOKEN_NAME_TOO_LONG"
	case StatusInvalidWipingAmount:
		return "INVALID_WIPING_AMOUNT"
	case StatusTokenIsImmutable:
		return "TOKEN_IS_IMMUTABLE"
	case StatusTokenAlreadyAssociatedToAccount:
		return "TOKEN_ALREADY_ASSOCIATED_TO_ACCOUNT"
	case StatusTransactionRequiresZeroTokenBalances:
		return "TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES"
	case StatusAccountIsTreasury:
		return "ACCOUNT_IS_TREASURY"
	case StatusTokenIDRepeatedInTokenList:
		return "TOKEN_ID_REPEATED_IN_TOKEN_LIST"
	case StatusTokenTransferListSizeLimitExceeded:
		return "TOKEN_TRANSFER_LIST_SIZE_LIMIT_EXCEEDED"
	case StatusEmptyTokenTransferBody:
		return "EMPTY_TOKEN_TRANSFER_BODY"
	case StatusEmptyTokenTransferAccountAmounts:
		return "EMPTY_TOKEN_TRANSFER_ACCOUNT_AMOUNTS"
	case StatusInvalidScheduleID:
		return "INVALID_SCHEDULE_ID"
	case StatusScheduleIsImmutable:
		return "SCHEDULE_IS_IMMUTABLE"
	case StatusInvalidSchedulePayerID:
		return "INVALID_SCHEDULE_PAYER_ID"
	case StatusInvalidScheduleAccountID:
		return "INVALID_SCHEDULE_ACCOUNT_ID"
	case StatusNoNewValidSignatures:
		return "NO_NEW_VALID_SIGNATURES"
	case StatusUnresolvableRequiredSigners:
		return "UNRESOLVABLE_REQUIRED_SIGNERS"
	case StatusScheduledTransactionNotInWhitelist:
		return "SCHEDULED_TRANSACTION_NOT_IN_WHITELIST"
	case StatusSomeSignaturesWereInvalid:
		return "SOME_SIGNATURES_WERE_INVALID"
	case StatusTransactionIDFieldNotAllowed:
		return "TRANSACTION_ID_FIELD_NOT_ALLOWED"
	case StatusIdenticalScheduleAlreadyCreated:
		return "IDENTICAL_SCHEDULE_ALREADY_CREATED"
	case StatusInvalidZeroByteInString:
		return "INVALID_ZERO_BYTE_IN_STRING"
	case StatusScheduleAlreadyDeleted:
		return "SCHEDULE_ALREADY_DELETED"
	case StatusScheduleAlreadyExecuted:
		return "SCHEDULE_ALREADY_EXECUTED"
	case StatusMessageSizeTooLarge:
		return "MESSAGE_SIZE_TOO_LARGE"
	case StatusOperationRepeatedInBucketGroups:
		return "OPERATION_REPEATED_IN_BUCKET_GROUPS"
	case StatusBucketCapacityOverflow:
		return "BUCKET_CAPACITY_OVERFLOW"
	case StatusNodeCapacityNotSufficientForOperation:
		return "NODE_CAPACITY_NOT_SUFFICIENT_FOR_OPERATION"
	case StatusBucketHasNoThrottleGroups:
		return "BUCKET_HAS_NO_THROTTLE_GROUPS"
	case StatusThrottleGroupHasZeroOpsPerSec:
		return "THROTTLE_GROUP_HAS_ZERO_OPS_PER_SEC"
	case StatusSuccessButMissingExpectedOperation:
		return "SUCCESS_BUT_MISSING_EXPECTED_OPERATION"
	case StatusUnparseableThrottleDefinitions:
		return "UNPARSEABLE_THROTTLE_DEFINITIONS"
	case StatusInvalidThrottleDefinitions:
		return "INVALID_THROTTLE_DEFINITIONS"
	case StatusAccountExpiredAndPendingRemoval:
		return "ACCOUNT_EXPIRED_AND_PENDING_REMOVAL"
	case StatusInvalidTokenMaxSupply:
		return "INVALID_TOKEN_MAX_SUPPLY"
	case StatusInvalidTokenNftSerialNumber:
		return "INVALID_TOKEN_NFT_SERIAL_NUMBER"
	case StatusInvalidNftID:
		return "INVALID_NFT_ID"
	case StatusMetadataTooLong:
		return "METADATA_TOO_LONG"
	case StatusBatchSizeLimitExceeded:
		return "BATCH_SIZE_LIMIT_EXCEEDED"
	case StatusInvalidQueryRange:
		return "INVALID_QUERY_RANGE"
	case StatusFractionDividesByZero:
		return "FRACTION_DIVIDES_BY_ZERO"
	case StatusInsufficientPayerBalanceForCustomFee:
		return "INSUFFICIENT_PAYER_BALANCE_FOR_CUSTOM_FEE"
	case StatusCustomFeesListTooLong:
		return "CUSTOM_FEES_LIST_TOO_LONG"
	case StatusInvalidCustomFeeCollector:
		return "INVALID_CUSTOM_FEE_COLLECTOR"
	case StatusInvalidTokenIDInCustomFees:
		return "INVALID_TOKEN_ID_IN_CUSTOM_FEES"
	case StatusTokenNotAssociatedToFeeCollector:
		return "TOKEN_NOT_ASSOCIATED_TO_FEE_COLLECTOR"
	case StatusTokenMaxSupplyReached:
		return "TOKEN_MAX_SUPPLY_REACHED"
	case StatusSenderDoesNotOwnNftSerialNo:
		return "SENDER_DOES_NOT_OWN_NFT_SERIAL_NO"
	case StatusCustomFeeNotFullySpecified:
		return "CUSTOM_FEE_NOT_FULLY_SPECIFIED"
	case StatusCustomFeeMustBePositive:
		return "CUSTOM_FEE_MUST_BE_POSITIVE"
	case StatusTokenHasNoFeeScheduleKey:
		return "TOKEN_HAS_NO_FEE_SCHEDULE_KEY"
	case StatusCustomFeeOutsideNumericRange:
		return "CUSTOM_FEE_OUTSIDE_NUMERIC_RANGE"
	case StatusRoyaltyFractionCannotExceedOne:
		return "ROYALTY_FRACTION_CANNOT_EXCEED_ONE"
	case StatusFractionalFeeMaxAmountLessThanMinAmount:
		return "FRACTIONAL_FEE_MAX_AMOUNT_LESS_THAN_MIN_AMOUNT"
	case StatusCustomScheduleAlreadyHasNoFees:
		return "CUSTOM_SCHEDULE_ALREADY_HAS_NO_FEES"
	case StatusCustomFeeDenominationMustBeFungibleCommon:
		return "CUSTOM_FEE_DENOMINATION_MUST_BE_FUNGIBLE_COMMON"
	case StatusCustomFractionalFeeOnlyAllowedForFungibleCommon:
		return "CUSTOM_FRACTIONAL_FEE_ONLY_ALLOWED_FOR_FUNGIBLE_COMMON"
	case StatusInvalidCustomFeeScheduleKey:
		return "INVALID_CUSTOM_FEE_SCHEDULE_KEY"
	case StatusInvalidTokenMintMetadata:
		return "INVALID_TOKEN_MINT_METADATA"
	case StatusInvalidTokenBurnMetadata:
		return "INVALID_TOKEN_BURN_METADATA"
	case StatusCurrentTreasuryStillOwnsNfts:
		return "CURRENT_TREASURY_STILL_OWNS_NFTS"
	case StatusAccountStillOwnsNfts:
		return "ACCOUNT_STILL_OWNS_NFTS"
	case StatusTreasuryMustOwnBurnedNft:
		return "TREASURY_MUST_OWN_BURNED_NFT"
	case StatusAccountDoesNotOwnWipedNft:
		return "ACCOUNT_DOES_NOT_OWN_WIPED_NFT"
	case StatusAccountAmountTransfersOnlyAllowedForFungibleCommon:
		return "ACCOUNT_AMOUNT_TRANSFERS_ONLY_ALLOWED_FOR_FUNGIBLE_COMMON"
	case StatusMaxNftsInPriceRegimeHaveBeenMinted:
		return "MAX_NFTS_IN_PRICE_REGIME_HAVE_BEEN_MINTED"
	case StatusPayerAccountDeleted:
		return "PAYER_ACCOUNT_DELETED"
	case StatusCustomFeeChargingExceededMaxRecursionDepth:
		return "CUSTOM_FEE_CHARGING_EXCEEDED_MAX_RECURSION_DEPTH"
	case StatusCustomFeeChargingExceededMaxAccountAmounts:
		return "CUSTOM_FEE_CHARGING_EXCEEDED_MAX_ACCOUNT_AMOUNTS"
	case StatusInsufficientSenderAccountBalanceForCustomFee:
		return "INSUFFICIENT_SENDER_ACCOUNT_BALANCE_FOR_CUSTOM_FEE"
	case StatusSerialNumberLimitReached:
		return "SERIAL_NUMBER_LIMIT_REACHED"
	case StatusCustomRoyaltyFeeOnlyAllowedForNonFungibleUnique:
		return "CUSTOM_ROYALTY_FEE_ONLY_ALLOWED_FOR_NON_FUNGIBLE_UNIQUE"
	case StatusNoRemainingAutomaticAssociations:
		return "NO_REMAINING_AUTOMATIC_ASSOCIATIONS"
	case StatusExistingAutomaticAssociationsExceedGivenLimit:
		return "EXISTING_AUTOMATIC_ASSOCIATIONS_EXCEED_GIVEN_LIMIT"
	case StatusRequestedNumAutomaticAssociationsExceedsAssociationLimit:
		return "REQUESTED_NUM_AUTOMATIC_ASSOCIATIONS_EXCEEDS_ASSOCIATION_LIMIT"
	case StatusTokenIsPaused:
		return "TOKEN_IS_PAUSED"
	case StatusTokenHasNoPauseKey:
		return "TOKEN_HAS_NO_PAUSE_KEY"
	case StatusInvalidPauseKey:
		return "INVALID_PAUSE_KEY"
	case StatusFreezeUpdateFileDoesNotExist:
		return "FREEZE_UPDATE_FILE_DOES_NOT_EXIST"
	case StatusFreezeUpdateFileHashDoesNotMatch:
		return "FREEZE_UPDATE_FILE_HASH_DOES_NOT_MATCH"
	case StatusNoUpgradeHasBeenPrepared:
		return "NO_UPGRADE_HAS_BEEN_PREPARED"
	case StatusNoFreezeIsScheduled:
		return "NO_FREEZE_IS_SCHEDULED"
	case StatusUpdateFileHashChangedSincePrepareUpgrade:
		return "UPDATE_FILE_HASH_CHANGED_SINCE_PREPARE_UPGRADE"
	case StatusFreezeStartTimeMustBeFuture:
		return "FREEZE_START_TIME_MUST_BE_FUTURE"
	case StatusPreparedUpdateFileIsImmutable:
		return "PREPARED_UPDATE_FILE_IS_IMMUTABLE"
	case StatusFreezeAlreadyScheduled:
		return "FREEZE_ALREADY_SCHEDULED"
	case StatusFreezeUpgradeInProgress:
		return "FREEZE_UPGRADE_IN_PROGRESS"
	case StatusUpdateFileIDDoesNotMatchPrepared:
		return "UPDATE_FILE_ID_DOES_NOT_MATCH_PREPARED"
	case StatusUpdateFileHashDoesNotMatchPrepared:
		return "UPDATE_FILE_HASH_DOES_NOT_MATCH_PREPARED"
	case StatusConsensusGasExhausted:
		return "CONSENSUS_GAS_EXHAUSTED"
	case StatusRevertedSuccess:
		return "REVERTED_SUCCESS"
	case StatusMaxStorageInPriceRegimeHasBeenUsed:
		return "MAX_STORAGE_IN_PRICE_REGIME_HAS_BEEN_USED"
	case StatusInvalidAliasKey:
		return "INVALID_ALIAS_KEY"
	case StatusUnexpectedTokenDecimals:
		return "UNEXPECTED_TOKEN_DECIMALS"
	case StatusInvalidProxyAccountID:
		return "INVALID_PROXY_ACCOUNT_ID"
	case StatusInvalidTransferAccountID:
		return "INVALID_TRANSFER_ACCOUNT_ID"
	case StatusInvalidFeeCollectorAccountID:
		return "INVALID_FEE_COLLECTOR_ACCOUNT_ID"
	case StatusAliasIsImmutable:
		return "ALIAS_IS_IMMUTABLE"
	case StatusSpenderAccountSameAsOwner:
		return "SPENDER_ACCOUNT_SAME_AS_OWNER"
	case StatusAmountExceedsTokenMaxSupply:
		return "AMOUNT_EXCEEDS_TOKEN_MAX_SUPPLY"
	case StatusNegativeAllowanceAmount:
		return "NEGATIVE_ALLOWANCE_AMOUNT"
	case StatusCannotApproveForAllFungibleCommon:
		return "CANNOT_APPROVE_FOR_ALL_FUNGIBLE_COMMON"
	case StatusSpenderDoesNotHaveAllowance:
		return "SPENDER_DOES_NOT_HAVE_ALLOWANCE"
	case StatusAmountExceedsAllowance:
		return "AMOUNT_EXCEEDS_ALLOWANCE"
	case StatusMaxAllowancesExceeded:
		return "MAX_ALLOWANCES_EXCEEDED"
	case StatusEmptyAllowances:
		return "EMPTY_ALLOWANCES"
	case StatusSpenderAccountRepeatedInAllowance:
		return "SPENDER_ACCOUNT_REPEATED_IN_ALLOWANCES"
	case StatusRepeatedSerialNumsInNftAllowances:
		return "REPEATED_SERIAL_NUMS_IN_NFT_ALLOWANCES"
	case StatusFungibleTokenInNftAllowances:
		return "FUNGIBLE_TOKEN_IN_NFT_ALLOWANCES"
	case StatusNftInFungibleTokenAllowances:
		return "NFT_IN_FUNGIBLE_TOKEN_ALLOWANCES"
	case StatusInvalidAllowanceOwnerID:
		return "INVALID_ALLOWANCE_OWNER_ID"
	case StatusInvalidAllowanceSpenderID:
		return "INVALID_ALLOWANCE_SPENDER_ID"
	}

	panic(fmt.Sprintf("unreachable: Status.String() switch statement is non-exhaustive. Status: %v", uint32(status)))
}

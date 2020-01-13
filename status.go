package hedera

type Status uint32

const (
	Ok                                  Status = 0
	InvalidTransaction                  Status = 1
	PayerAccountNotFound                Status = 2
	InvalidNodeAccount                  Status = 3
	TransactionExpired         Status = 4
	InvalidTransactionStart    Status = 5
	InvalidTransactionDuration Status = 6
	InvalidSignature         Status = 7
	MemoTooLong              Status = 8
	InsufficientTxFee        Status = 9
	InsufficientPayerBalance   Status = 10
	DuplicateTransaction       Status = 11
	Busy                       Status = 12
	NotSupported               Status = 13
	InvalidFileID              Status = 14
	InvalidAccountID           Status = 15
	InvalidContractID          Status = 16
	InvalidTransactionID       Status = 17
	ReceiptNotFound            Status = 18
	RecordNotFound             Status = 19
	InvalidSolidityID          Status = 20
	Unknown                    Status = 21
	Success                    Status = 22
	FailInvalid                Status = 23
	FailFee                    Status = 24
	FailBalance                Status = 25
	KeyRequired                Status = 26
	BadEncoding                Status = 27
	InsufficientAccountBalance Status = 28
	InvalidSolidityAddress     Status = 29
	InsufficientGas            Status = 30
	ContractSizeLimitExceeded           Status = 31
	LocalCallModificationException      Status = 32
	ContractRevertExecuted              Status = 33
	ContractExecutionException          Status = 34
	InvalidReceivingNodeAccount         Status = 35
	MissingQueryHeader                  Status = 36
	AccountUpdateFailed                 Status = 37
	InvalidKeyEncoding                  Status = 38
	NullSolidityAddress                 Status = 39
	ContractUpdateFailed                Status = 40
	InvalidQueryHeader                  Status = 41
	InvalidFeeSubmitted                 Status = 42
	InvalidPayerSignature               Status = 43
	KeyNotProvided                      Status = 44
	InvalidExpirationTime               Status = 45
	NoWaclKey                           Status = 46
	FileContentEmpty                    Status = 47
	InvalidAccountAmounts               Status = 48
	EmptyTransactionBody                Status = 49
	InvalidTransactionBody              Status = 50
	InvalidSignatureTypeMismatchingKey  Status = 51
	InvalidSignatureCountMismatchingKey Status = 52
	EmptyClaimBody                      Status = 53
	EmptyClaimHash                      Status = 54
	EmptyClaimKeys                      Status = 55
	InvalidClaimHashSize                Status = 56
	EmptyQueryBody                      Status = 57
	EmptyClaimQuery                     Status = 58
	ClaimNotFound                       Status = 59
	AccountIdDoesNotExist               Status = 60
	ClaimAlreadyExists              Status = 61
	InvalidFileWacl                 Status = 62
	SerializationFailed             Status = 63
	TransactionOversize             Status = 64
	TransactionTooManyLayers        Status = 65
	ContractDeleted                 Status = 66
	PlatformNotActive               Status = 67
	KeyPrefixMismatch               Status = 68
	PlatformTransactionNotCreated   Status = 69
	InvalidRenewalPeriod            Status = 70
	InvalidPayerAccountID           Status = 71
	AccountDeleted                  Status = 72
	FileDeleted                     Status = 73
	AccountRepeatedInAccountAmounts Status = 74
	SettingNegativeAccountBalance   Status = 75
	ObtainerRequired                Status = 76
	ObtainerSameContractID          Status = 77
	ObtainerDoesNotExist            Status = 78
	ModifyingImmutableContract      Status = 79
	FileSystemException             Status = 80
	AutorenewDurationNotInRange     Status = 81
	ErrorDecodingBytestring         Status = 82
	ContractFileEmpty               Status = 83
	ContractBytecodeEmpty           Status = 84
	InvalidInitialBalance           Status = 85
	InvalidReceiveRecordThreshold   Status = 86
	InvalidSendRecordThreshold      Status = 87
	AccountIsNotGenesisAccount          Status = 88
	PayerAccountUnauthorized            Status = 89
	InvalidFreezeTransactionBody        Status = 90
	FreezeTransactionBodyNotFound       Status = 91
	TransferListSizeLimitExceeded       Status = 92
	ResultSizeLimitExceeded             Status = 93
	NotSpecialAccount                   Status = 94
	ContractNegativeGas                 Status = 95
	ContractNegativeValue               Status = 96
	InvalidFeeFile                      Status = 97
	InvalidExchangeRateFile             Status = 98
	InsufficientLocalCallGas            Status = 99
	EntityNotAllowedToDelete            Status = 100
	AuthorizationFailed                 Status = 101
	FileUploadedProtoInvalid            Status = 102
	FileUploadedProtoNotSavedToDisk     Status = 103
	FeeScheduleFilePartUploaded         Status = 104
	ExchangeRateChangeLimitExceeded     Status = 105
	MaxContractStorageExceeded          Status = 106
	TransaferAccountSameAsDeleteAccount Status = 107
	TotalLedgerBalanceInvalid           Status = 108
	ExpirationReductionNotAllowed       Status = 110
	MaxGasLimitExceeded                 Status = 111
	MaxFileSizeExceeded                 Status = 112
)

func (status Status) String() string {
	switch status {
	case Ok:
		return "OK"
	case InvalidTransaction:
		return "INVALID_TRANSACTION"
	case PayerAccountNotFound:
		return "PAYER_ACCOUNT_NOT_FOUND"
	case InvalidNodeAccount:
		return "INVALID_NODE_ACCOUNT"
	case TransactionExpired:
		return "TRANSACTION_EXPIRED"
	case InvalidTransactionStart:
		return "INVALID_TRANSACTION_START"
	case InvalidTransactionDuration:
		return "INVALID_TRANSACTION_DURATION"
	case InvalidSignature:
		return "INVALID_SIGNATURE"
	case MemoTooLong:
		return "MEMO_TOO_LONG"
	case InsufficientTxFee:
		return "INSUFFICIENT_TX_FEE"
	case InsufficientPayerBalance:
		return "INSUFFICIENT_PAYER_BALANCE"
	case DuplicateTransaction:
		return "DUPLICATE_TRANSACTION"
	case Busy:
		return "BUSY"
	case NotSupported:
		return "NOT_SUPPORTED"
	case InvalidFileID:
		return "INVALID_FILE_ID"
	case InvalidAccountID:
		return "INVALID_ACCOUNT_ID"
	case InvalidContractID:
		return "INVALID_CONTRACT_ID"
	case InvalidTransactionID:
		return "INVALID_TRANSACTION_ID"
	case ReceiptNotFound:
		return "RECEIPT_NOT_FOUND"
	case RecordNotFound:
		return "RECORD_NOT_FOUND"
	case InvalidSolidityID:
		return "INVALID_SOLIDITY_ID"
	case Unknown:
		return "UNKNOWN"
	case Success:
		return "SUCCESS"
	case FailInvalid:
		return "FAIL_INVALID"
	case FailFee:
		return "FAIL_FEE"
	case FailBalance:
		return "FAIL_BALANCE"
	case KeyRequired:
		return "KEY_REQUIRED"
	case BadEncoding:
		return "BAD_ENCODING"
	case InsufficientAccountBalance:
		return "INSUFFICIENT_ACCOUNT_BALANCE"
	case InvalidSolidityAddress:
		return "INVALID_SOLIDITY_ADDRESS"
	case InsufficientGas:
		return "INSUFFICIENT_GAS"
	case ContractSizeLimitExceeded:
		return "CONTRACT_SIZE_LIMIT_EXCEEDED"
	case LocalCallModificationException:
		return "LOCAL_CALL_MODIFICATION_EXCEPTION"
	case ContractRevertExecuted:
		return "CONTRACT_REVERT_EXECUTED"
	case ContractExecutionException:
		return "CONTRACT_EXECUTION_EXCEPTION"
	case InvalidReceivingNodeAccount:
		return "INVALID_RECEIVING_NODE_ACCOUNT"
	case MissingQueryHeader:
		return "MISSING_QUERY_HEADER"
	case AccountUpdateFailed:
		return "ACCOUNT_UPDATE_FAILED"
	case InvalidKeyEncoding:
		return "INVALID_KEY_ENCODING"
	case NullSolidityAddress:
		return "NULL_SOLIDITY_ADDRESS"
	case ContractUpdateFailed:
		return "CONTRACT_UPDATE_FAILED"
	case InvalidQueryHeader:
		return "INVALID_QUERY_HEADER"
	case InvalidFeeSubmitted:
		return "INVALID_FEE_SUBMITTED"
	case InvalidPayerSignature:
		return "INVALID_PAYER_SIGNATURE"
	case KeyNotProvided:
		return "KEY_NOT_PROVIDED"
	case InvalidExpirationTime:
		return "INVALID_EXPIRATION_TIME"
	case NoWaclKey:
		return "NO_WACL_KEY"
	case FileContentEmpty:
		return "FILE_CONTENT_EMPTY"
	case InvalidAccountAmounts:
		return "INVALID_ACCOUNT_AMOUNTS"
	case EmptyTransactionBody:
		return "EMPTY_TRANSACTION_BODY"
	case InvalidTransactionBody:
		return "INVALID_TRANSACTION_BODY"
	case InvalidSignatureTypeMismatchingKey:
		return "INVALID_SIGNATURE_TYPE_MISMATCHING_KEY"
	case InvalidSignatureCountMismatchingKey:
		return "INVALID_SIGNATURE_COUNT_MISMATCHING_KEY"
	case EmptyClaimBody:
		return "EMPTY_CLAIM_BODY"
	case EmptyClaimHash:
		return "EMPTY_CLAIM_HASH"
	case EmptyClaimKeys:
		return "EMPTY_CLAIM_KEYS"
	case InvalidClaimHashSize:
		return "INVALID_CLAIM_HASH_SIZE"
	case EmptyQueryBody:
		return "EMPTY_QUERY_BODY"
	case EmptyClaimQuery:
		return "EMPTY_CLAIM_QUERY"
	case ClaimNotFound:
		return "CLAIM_NOT_FOUND"
	case AccountIdDoesNotExist:
		return "ACCOUNT_ID_DOES_NOT_EXIST"
	case ClaimAlreadyExists:
		return "CLAIM_ALREADY_EXISTS"
	case InvalidFileWacl:
		return "INVALID_FILE_WACL"
	case SerializationFailed:
		return "SERIALIZATION_FAILED"
	case TransactionOversize:
		return "TRANSACTION_OVERSIZE"
	case TransactionTooManyLayers:
		return "TRANSACTION_TOO_MANY_LAYERS"
	case ContractDeleted:
		return "CONTRACT_DELETED"
	case PlatformNotActive:
		return "PLATFORM_NOT_ACTIVE"
	case KeyPrefixMismatch:
		return "KEY_PREFIX_MISMATCH"
	case PlatformTransactionNotCreated:
		return "PLATFORM_TRANSACTION_NOT_CREATED"
	case InvalidRenewalPeriod:
		return "INVALID_RENEWAL_PERIOD"
	case InvalidPayerAccountID:
		return "INVALID_PAYER_ACCOUNT_ID"
	case AccountDeleted:
		return "ACCOUNT_DELETED"
	case FileDeleted:
		return "FILE_DELETED"
	case AccountRepeatedInAccountAmounts:
		return "ACCOUNT_REPEATED_IN_ACCOUNT_AMOUNTS"
	case SettingNegativeAccountBalance:
		return "SETTING_NEGATIVE_ACCOUNT_BALANCE"
	case ObtainerRequired:
		return "OBTAINER_REQUIRED"
	case ObtainerSameContractID:
		return "OBTAINER_SAME_CONTRACT_ID"
	case ObtainerDoesNotExist:
		return "OBTAINER_DOES_NOT_EXIST"
	case ModifyingImmutableContract:
		return "MODIFYING_IMMUTABLE_CONTRACT"
	case FileSystemException:
		return "FILE_SYSTEM_EXCEPTION"
	case AutorenewDurationNotInRange:
		return "AUTORENEW_DURATION_NOT_IN_RANGE"
	case ErrorDecodingBytestring:
		return "ERROR_DECODING_BYTESTRING"
	case ContractFileEmpty:
		return "CONTRACT_FILE_EMPTY"
	case ContractBytecodeEmpty:
		return "CONTRACT_BYTECODE_EMPTY"
	case InvalidInitialBalance:
		return "INVALID_INITIAL_BALANCE"
	case InvalidReceiveRecordThreshold:
		return "INVALID_RECEIVE_RECORD_THRESHOLD"
	case InvalidSendRecordThreshold:
		return "INVALID_SEND_RECORD_THRESHOLD"
	case AccountIsNotGenesisAccount:
		return "ACCOUNT_IS_NOT_GENESIS_ACCOUNT"
	case PayerAccountUnauthorized:
		return "PAYER_ACCOUNT_UNAUTHORIZED"
	case InvalidFreezeTransactionBody:
		return "INVALID_FREEZE_TRANSACTION_BODY"
	case FreezeTransactionBodyNotFound:
		return "FREEZE_TRANSACTION_BODY_NOT_FOUND"
	case TransferListSizeLimitExceeded:
		return "TRANSFER_LIST_SIZE_LIMIT_EXCEEDED"
	case ResultSizeLimitExceeded:
		return "RESULT_SIZE_LIMIT_EXCEEDED"
	case NotSpecialAccount:
		return "NOT_SPECIAL_ACCOUNT"
	case ContractNegativeGas:
		return "CONTRACT_NEGATIVE_GAS"
	case ContractNegativeValue:
		return "CONTRACT_NEGATIVE_VALUE"
	case InvalidFeeFile:
		return "INVALID_FEE_FILE"
	case InvalidExchangeRateFile:
		return "INVALID_EXCHANGE_RATE_FILE"
	case InsufficientLocalCallGas:
		return "INSUFFICIENT_LOCAL_CALL_GAS"
	case EntityNotAllowedToDelete:
		return "ENTITY_NOT_ALLOWED_TO_DELETE"
	case AuthorizationFailed:
		return "AUTHORIZATION_FAILED"
	case FileUploadedProtoInvalid:
		return "FILE_UPLOADED_PROTO_INVALID"
	case FileUploadedProtoNotSavedToDisk:
		return "FILE_UPLOADED_PROTO_NOT_SAVED_TO_DISK"
	case FeeScheduleFilePartUploaded:
		return "FEE_SCHEDULE_FILE_PART_UPLOADED"
	case ExchangeRateChangeLimitExceeded:
		return "EXCHANGE_RATE_CHANGE_LIMIT_EXCEEDED"
	case MaxContractStorageExceeded:
		return "MAX_CONTRACT_STORAGE_EXCEEDED"
	case TransaferAccountSameAsDeleteAccount:
		return "TRANSAFER_ACCOUNT_SAME_AS_DELETE_ACCOUNT"
	case TotalLedgerBalanceInvalid:
		return "TOTAL_LEDGER_BALANCE_INVALID"
	case ExpirationReductionNotAllowed:
		return "EXPIRATION_REDUCTION_NOT_ALLOWED"
	case MaxGasLimitExceeded:
		return "MAX_GAS_LIMIT_EXCEEDED"
	case MaxFileSizeExceeded:
		return "MAX_FILE_SIZE_EXCEEDED"
	}

	panic("unreacahble: Status.String() switch statement is non-exhaustive")
}

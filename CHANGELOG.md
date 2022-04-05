## v2.13.0-beta.1

### Added

* `AccountAllowanceDeleteTransaction`
* `ContractFunctionResult.[gas|hbarAmount|contractFunctionParametersBytes]`
* `AccountAllowanceExample`
* `ScheduleTransferExample`

### Deprecated

* `AccountAllowanceAdjustTransaction.revokeTokenNftAllowance()` with no replacement.
* `AccountAllowanceApproveTransaction.AddHbarApproval()`, use `ApproveHbarAllowance()` instead.
* `AccountAllowanceApproveTransaction.ApproveTokenApproval()`, use `GrantTokenNftAllowance()` instead.
* `AccountAllowanceApproveTransaction.ApproveTokenNftApproval()`, use `ApproveTokenNftAllowance()` instead.

## v2.12.0

### Added

* `AccountInfoFlowVerify[Signature|Transaction]()`
* `Client.[Set|Get]NodeMinReadmitPeriod()`
* Support for using any node from the entire network upon execution
  if node account IDs have no been locked for the request.
* Support for all integer widths for `ContractFunction[Result|Selector|Params]`

### Fixed

* Ledger ID checksums
* `TransactionFromBytes()` should validate all the transaction bodies are the same

### Changed

* Network behavior to follow a more standard approach (remove the sorting we
  used to do).

## v2.12.0-beta.1

### Added

* `AccountInfoFlowVerify[Signature|Transaction]()`
* `Client.[Set|Get]NodeMinReadmitPeriod()`
* Support for using any node from the entire network upon execution
  if node account IDs have no been locked for the request.
* Support for all integer widths for `ContractFunction[Result|Selector|Params]`

### Fixed

* Ledger ID checksums
* `TransactionFromBytes()` should validate all the transaction bodies are the same

### Changed

* Network behavior to follow a more standard approach (remove the sorting we
  used to do).

## v2.11.0

### Added

* `ContractCreateFlow`
* `Query.[Set|Get]PaymentTransactionID`
* Verbose logging using zerolog
* `*[Transaction|Query].[Set|Get]GrpcDeadline()`
* `TransactionRecord.[hbar|Token|TokenNft]AllowanceAdjustments`
* `TransferTransaction.AddApproved[Hbar|Token|Nft]Transfer()`
* `AccountAllowanceApproveTransaction.Approve[Hbar|Token|TokenNft]Allowance()`
* `AccountAllowanceAdjustTransaction.[Grant|Revoke][Hbar|Token|TokenNft]Allowance()`
* `AccountAllowanceAdjustTransaction.[Grant|Revoke]TokenNftAllowanceAllSerials()`

### Fixed

* `HbarAllowance.OwnerAccountID`, wasn't being set.
* Min/max backoff for nodes should start at 8s to 60s
* The current backoff for nodes should be used when sorting inside of network
  meaning nodes with a smaller current backoff will be prioritized
* `TopicMessageQuery` start time should have a default

### Deprecated

* `AccountUpdateTransaction.[Set|Get]AliasKey`

### Removed

* `Account[Approve|Adjust]AllowanceTransaction.Add[Hbar|Token|TokenNft]AllowanceWithOwner()`

## v2.11.0-beta.1

### Added

* `ContractCreateFlow`
* `Account[Approve|Adjust]AllowanceTransaction.add[Hbar|Token|TokenNft]AllowanceWithOwner()`
* `Query.[Set|Get]PaymentTransactionID`
* Verbose logging using zerolog
* `*[Transaction|Query].[Set|Get]GrpcDeadline()`

### Fixed

* `HbarAllowance.OwnerAccountID`, wasn't being set.
* Min/max backoff for nodes should start at 8s to 60s
* The current backoff for nodes should be used when sorting inside of network
  meaning nodes with a smaller current backoff will be prioritized

### Deprecated

* `AccountUpdateTransaction.[Set|Get]AliasKey`

## v2.10.0

### Added

* `owner` field to `*Allowance`.
* Added free `AddressBookQuery`.

### Fixed

* Changed mirror node port to correct one, 443.
* Occasional ECDSA invalid length error.
* ContractIDFromString() now sets EvmAddress correctly to nil, when evm address is not detected

## v2.10.0-beta.1

### Added

* `owner` field to `*Allowance`.
* Added free `AddressBookQuery`.

### Fixed

* Changed mirror node port to correct one, 443.

## v2.9.0

### Added

* CREATE2 Solidity addresses can now be represented by a `ContractId` with `EvmAddress` set.
* `ContractId.FromEvmAddress()`
* `ContractFunctionResult.StateChanges`
* `ContractFunctionResult.EvmAddress`
* `ContractStateChange`
* `StorageChange`
* New response codes.
* `ChunkedTransaction.[Set|Get]ChunkSize()`, and changed default chunk size for `FileAppendTransaction` to 2048.
* `AccountAllowance[Adjust|Approve]Transaction`
* `AccountInfo.[hbar|token|tokenNft]Allowances`
* `[Hbar|Token|TokenNft]Allowance`
* `[Hbar|Token|TokenNft]Allowance`
* `TransferTransaction.set[Hbar|Token|TokenNft]TransferApproval()`

### Fixed

* Requests not cycling though nodes.
* Free queries not attempting to retry on different nodes.

### Deprecated

* `ContractId.FromSolidityAddress()`, use `ContractId.FromEvmAddress()` instead.
* `ContractFunctionResult.CreatedContractIDs`.

## v2.9.0-beta.2

### Added

* CREATE2 Solidity addresses can now be represented by a `ContractId` with `EvmAddress` set.
* `ContractId.FromEvmAddress()`
* `ContractFunctionResult.StateChanges`
* `ContractFunctionResult.EvmAddress`
* `ContractStateChange`
* `StorageChange`
* New response codes.
* `ChunkedTransaction.[Set|Get]ChunkSize()`, and changed default chunk size for `FileAppendTransaction` to 2048.
* `AccountAllowance[Adjust|Approve]Transaction`
* `AccountInfo.[hbar|token|tokenNft]Allowances`
* `[Hbar|Token|TokenNft]Allowance`
* `[Hbar|Token|TokenNft]Allowance`
* `TransferTransaction.set[Hbar|Token|TokenNft]TransferApproval()`

### Fixed

* Requests not cycling though nodes.
* Free queries not attempting to retry on different nodes.

### Deprecated

* `ContractId.FromSolidityAddress()`, use `ContractId.FromEvmAddress()` instead.
* `ContractFunctionResult.CreatedContractIDs`.

## v2.9.0-beta.1

### Added

* CREATE2 Solidity addresses can now be represented by a `ContractId` with `EvmAddress` set.
* `ContractId.FromEvmAddress()`
* `ContractFunctionResult.StateChanges`
* `ContractFunctionResult.EvmAddress`
* `ContractStateChange`
* `StorageChange`
* New response codes.
* `ChunkedTransaction.[Set|Get]ChunkSize()`, and changed default chunk size for `FileAppendTransaction` to 2048.
* `AccountAllowance[Adjust|Approve]Transaction`
* `AccountInfo.[hbar|token|tokenNft]Allowances`
* `[Hbar|Token|TokenNft]Allowance`
* `[Hbar|Token|TokenNft]Allowance`
* `TransferTransaction.set[Hbar|Token|TokenNft]TransferApproval()`

### Fixed

* Requests not cycling though nodes.
* Free queries not attempting to retry on different nodes.

### Deprecated

* `ContractId.FromSolidityAddress()`, use `ContractId.FromEvmAddress()` instead.
* `ContractFunctionResult.CreatedContractIDs`.

## v2.8.0

### Added

* Support for regenerating transaction IDs on demand if a request
  responses with `TRANSACITON_EXPIRED`

## v2.8.0-beta.1

### Added

* Support for regenerating transaction IDs on demand if a request
  responses with `TRANSACITON_EXPIRED`

## v2.7.0

### Added

* `AccountId.AliasKey`, including `AccountId.[From]String()` support.
* `[PublicKey|PrivateKey].ToAccountId()`.
* `AliasKey` fields in `TransactionRecord` and `AccountInfo`.
* `Nonce` field in `TransactionId`, including `TransactionId.[set|get]Nonce()`
* `Children` fields in `TransactionRecord` and `TransactionReceipt`
* `Duplicates` field in `TransactionReceipt`
* `[TransactionReceiptQuery|TransactionRecordQuery].[Set|Get]IncludeChildren()`
* `TransactionReceiptQuery.[Set|Get]IncludeDuplicates()`
* New response codes.
* Support for ECDSA SecP256K1 keys.
* `PrivateKeyGenerate[ED25519|ECDSA]()`
* `[Private|Public]KeyFrom[Bytes|String][DER|ED25519|ECDSA]()`
* `[Private|Public]Key.[Bytes|String][Raw|DER]()`
* `DelegateContractId`
* `*Id.[from|to]SolidityAddress()`

### Deprecated

* `PrivateKeyGenerate()`, use `PrivateKeyGenerate[ED25519|ECDSA]()` instead.

## v2.7.0-beta.1

### Added

* `AccountId.AliasKey`, including `AccountId.[From]String()` support.
* `[PublicKey|PrivateKey].ToAccountId()`.
* `AliasKey` fields in `TransactionRecord` and `AccountInfo`.
* `Nonce` field in `TransactionId`, including `TransactionId.[set|get]Nonce()`
* `Children` fields in `TransactionRecord` and `TransactionReceipt`
* `Duplicates` field in `TransactionReceipt`
* `[TransactionReceiptQuery|TransactionRecordQuery].[Set|Get]IncludeChildren()`
* `TransactionReceiptQuery.[Set|Get]IncludeDuplicates()`
* New response codes.
* Support for ECDSA SecP256K1 keys.
* `PrivateKeyGenerate[ED25519|ECDSA]()`
* `[Private|Public]KeyFrom[Bytes|String][DER|ED25519|ECDSA]()`
* `[Private|Public]Key.[Bytes|String][Raw|DER]()`

### Deprecated

* `PrivateKeyGenerate()`, use `PrivateKeyGenerate[ED25519|ECDSA]()` instead.

## v2.6.0

### Added

* New smart contract response codes

### Deprecated

* `ContractCallQuery.[Set|Get]MaxResultSize()`
* `ContractUpdateTransaction.[Set|Get]ByteCodeFileID()`

## v2.6.0-beta.1

### Added

* New smart contract response codes

### Deprecated

* `ContractCallQuery.[Set|Get]MaxResultSize()`
* `ContractUpdateTransaction.[Set|Get]ByteCodeFileID()`

## v2.5.1

### Fixed

* `TransferTransaction.GetTokenTransfers()`
* `TransferTransaction.AddTokenTransfer()`
* Persistent error not being handled correctly
* `TransactionReceiptQuery` should return even on a bad status codes.
  Only *.GetReceipt()` should error on non `SUCCESS` status codes

## v2.5.1-beta.1

### Changed

* Refactored and updated node account ID handling to err whenever a node account ID of 0.0.0 is being set

## v2.5.0-beta.1

### Deprecated

* `ContractCallQuery.[Set|Get]MaxResultSize()`
* `ContractUpdateTransaction.[Set|Get]ByteCodeFileID()`

## v2.5.0

### Fixed

* `TransactionReceiptQuery` should fill out `TransactionReceipt` even when a bad `Status` is returned

## v2.4.1

### Fixed

* `TransferTransaction` should serialize the transfers list deterministically

## v2.4.0

### Added

* Support for toggling TLS for both mirror network and services network

## v2.3.0

### Added

* `FreezeType`
* `FreezeTransaction.[get|set]FreezeType()`

## v2.3.0-beta 1

### Added

* Support for HIP-24 (token pausing)
    * `TokenInfo.PauseKey`
    * `TokenInfo.PauseStatus`
    * `TokenCreateTransaction.PauseKey`
    * `TokenUpdateTransaction.PauseKey`
    * `TokenPauseTransaction`
    * `TokenUnpauseTransaction`

## v2.2.0

### Added

* Support for automatic token associations
    * `TransactionRecord.AutomaticTokenAssociations`
    * `AccountInfo.MaxAutomaticTokenAssociations`
    * `AccountCreateTransaction.MaxAutomaticTokenAssociations`
    * `AccountUpdateTransaction.MaxAutomaticTokenAssociations`
    * `TokenRelationship.AutomaticAssociation`
    * `TokenAssociation`
* `Transaction*` helper methods - should make it easier to use the result of `TransactionFromBytes()`

### Fixed

* TLS now properly confirms certificate hashes
* `TokenUpdateTransaction.GetExpirationTime()` returns the correct time
* Several `*.Get*()` methods required a parameter similiar to `*.Set*()`
  This has been changed completely instead of deprecated because we treated this as hard bug
* Several `nil` dereference issues related to to/from protobuf conversions

## v2.2.0-beta.1

### Added

* Support for automatic token associations
    * `TransactionRecord.AutomaticTokenAssociations`
    * `AccountInfo.MaxAutomaticTokenAssociations`
    * `AccountCreateTransaction.MaxAutomaticTokenAssociations`
    * `AccountUpdateTransaction.MaxAutomaticTokenAssociations`
    * `TokenRelationship.AutomaticAssociation`
    * `TokenAssociation`

## v2.1.16

### Added

* Support for TLS
* Setters which follow the builder pattern to `Custom*Fee`
* `Client.[min|max]Backoff()` support

### Deprecated

* `TokenNftInfoQuery.ByNftID()` - use `TokenNftInfoQuery.SetNftID()` instead
* `TokenNftInfoQuery.[By|Set|Get]AccountId()` with no replacement
* `TokenNftInfoQuery.[By|Set|Get]TokenId()` with no replacement
* `TokenNftInfoQuery.[Set|Get]Start()` with no replacement
* `TokenNftInfoQuery.[Set|Get]End()` with no replacement

## v2.1.15

### Fixed

* `AssessedCustomFee.PayerAccountIDs` was misspelled

## v2.1.14

### Added

* Support for `CustomRoyaltyFee`
* Support for `AssessedCustomFee.payerAccountIds`

### Fixed

* `nil` dereference issues within `*.validateNetworkIDs()`

## v2.1.13

### Added

* Implement `Client.pingAll()`
* Implement `Client.SetAutoChecksumValidation()` which validates all entity ID checksums on requests before executing

### Fixed

* nil dereference errors when decoding invalid PEM files

## v2.1.12

### Added

* Updated `Status` with new response codes
* Support for `Hbar.[from|to]String()` to be reversible

## v2.1.11

### Removed

* `*.AddCustomFee()` use `*.SetCustomFees()` instead

### Changes

* Update `Status` with new codes

### Fixes

* `PrivateKey.LegacyDerive()` should correctly handle indicies

## v2.1.11-beta.1

### Added

* Support for NFTS
    * Creating NFT tokens
    * Minting NFTs
    * Burning NFTs
    * Transfering NFTs
    * Wiping NFTs
    * Query NFT information
* Support for Custom Fees on tokens:
    * Setting custom fees on a token
    * Updating custom fees on an existing token

## v2.1.10

### Added

* All requests should retry on gRPC error `INTERNAL` if the message contains `RST_STREAM`
* `AccountBalance.Tokens` as a replacement for `AccountBalance.Token`
* `AccountBalance.TokenDecimals`
* All transactions will now `sign-on-demand` which should result in improved performance

### Fixed

* `TopicMessageQuery` not calling `Unsubscribe` when a stream is cancelled
* `TopicMessageQuery` should add 1 nanosecond to the `StartTime` of the last received message
* `TopicMessageQuery` allocate space for entire chunked message ahead of time
  for retries
* `TokenDeleteTransaction.SetTokenID()` incorrectly setting `tokenID` resulting in `GetTokenID()` always returning an empty `TokenID`
* `TransferTransaction.GetTokenTransfers()` incorrectly setting an empty value

### Deprecated

* `AccountBalance.Token` use `AccountBalance.Tokens` instead

## v2.1.9

### Fixed

* `Client.SetMirroNetwork()` producing a nil pointer exception on next use of a mirror network
* Mirror node TLS no longer producing nil pointer exception

## v2.1.8

### Added

* Support TLS for mirror node connections.
* Support for entity ID checksums which are validated whenever a request begins execution.
  This includes the IDs within the request, the account ID within the transaction ID, and
  query responses will contain entity IDs with a checksum for the network the query was executed on.

### Fixed

* `TransactionTransaction.AddHbarTransfer()` incorrectly determine total transfer per account ID

## v2.1.7

### Fixed

* `TopicMessageQuery.MaxBackoff` was not being used at all
* `TopicMessageQuery.Limit` was being incorrectly update with full `TopicMessages` rather than per chunk
* `TopicMessageQuery.StartTime` was not being updated each time a message was received
* `TopicMessageQuery.CompletionHandler` was be called at incorrect times
* Removed the use of locks and `sync.Map` within `TopicMessageQuery` as it is unncessary
* Added default logging to `ErrorHandler` and `CompletionHandler`

## v2.1.6

* Support for `MaxBackoff`, `MaxAttempts`, `RetryHandler`, and `CompletionHandler` in `TopicMessageQuery`
* Default logging behavior to `TopicMessageQuery` if an error handler or completion handler was not set

### Fixed

* Renamed `ScheduleInfo.Signers` -> `ScheduleInfo.Signatories`
* `TopicMessageQuery` retry handling; this should retry on more gRPC errors
* `TopicMessageQuery` max retry timeout; before this would could wait up to 4m with no feedback
* `durationFromProtobuf()` incorrectly calculation duration
* `*Token.GetAutoRenewPeriod()` and `*Token.GetExpirationTime()` nil dereference
* `Hbar.As()` using multiplication instead of division, and should return a `float64`

### Added

* Exposed `Hbar.Negated()`

## v2.1.5

###

* Scheduled transaction support: `ScheduleCreateTransaction`, `ScheduleDeleteTransaction`, and `ScheduleSignTransaction`
* Non-Constant Time Comparison of HMACs [NCC-E001154-006]
* Decreased `CHUNK_SIZE` 4096->1024 and increased default max chunks 10->20

## v2.1.5-beta.5

### Fixed

* Non-Constant Time Comparison of HMACs [NCC-E001154-006]
* Decreased `CHUNK_SIZE` 4096->1024 and increased default max chunks 10->20
* Renamed `ScheduleInfo.GetTransaction()` -> `ScheduleInfo.getScheduledTransaction()`

## v2.1.5-beta.4

### Fixed

* `Transaction.Schedule()` should error when scheduling un-scheduable tranasctions

### Removed

* `nonce` from `TransactionID`
* `ScheduleTransactionBody` - should not be part of the public API

## v2.1.5-beta.3

### Fixed

* `Transaction[Receipt|Record]Query` should not error for status `IDENTICAL_SCHEDULE_ALREADY_CREATED`
  because the other fields on the receipt are present with that status.
* `ErrHederaReceiptStatus` should print `exception receipt status ...` instead of
  `exception precheck status ...`

## v2.1.5-beta.2

### Fixed

* Executiong should retry on status `PLATFORM_TRANSACTION_NOT_CREATED`
* Error handling throughout the SDK
    * A precheck error shoudl be returned when the exceptional status is in the header
    * A receipt error should be returned when the exceptional status is in the receipt
* `TransactionRecordQuery` should retry on node precheck code `OK` only if we're not
  getting cost of query.
* `Transaction[Receipt|Record]Query` should retry on both `RECEIPT_NOT_FOUND` and
  `RECORD_NOT_FOUND` status codes when node precheck code is `OK`

## v2.1.5-beta.1

### Fixed

* Updated scheduled transaction to use new HAPI porotubfs

### Removed
* `ScheduleCreateTransaction.AddScheduledSignature()`
* `ScheduleCreateTransaction.GetScheduledSignatures()`
* `ScheduleSignTransaction.addScheduledSignature()`
* `ScheduleSignTransaction.GetScheduledSignatures()`

## v2.x

### Added

* Support for scheduled transactions.
    * `ScheduleCreateTransaction` - Create a new scheduled transaction
    * `ScheduleSignTransaction` - Sign an existing scheduled transaction on the network
    * `ScheduleDeleteTransaction` - Delete a scheduled transaction
    * `ScheduleInfoQuery` - Query the info including `bodyBytes` of a scheduled transaction
    * `ScheduleId`
* Support for scheduled and nonce in `TransactionId`
    * `TransactionIdWithNonce()` - Supports creating transaction ID with random bytes.
    * `TransactionId.[Set|Get]Scheduled()` - Supports scheduled transaction IDs.
* `TransactionIdWithValidStart()`

### Fixed

* Updated protobufs [#120](https://github.com/hashgraph/hedera-sdk-go/issues/120)

### Deprecate

* `NewTransactionId()` - Use `TransactionIdWithValidStart()` instead.

## v2.0.0

### Changes

* All requests support getter methods as well as setters.
* All requests support multiple node account IDs being set.
* `TransactionFromBytes()` supports multiple node account IDs and existing
  signatures.
* All requests support a max retry count using `SetMaxRetry()`


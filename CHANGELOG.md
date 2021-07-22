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
 

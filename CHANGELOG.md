## v2.1.5

###

 * Scheduled transaction support: `ScheduleCreateTransaction`, `ScheduleSignTransaction`, and `ScheduleSignTransaction`
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
 

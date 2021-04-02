## Unreleased

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
 

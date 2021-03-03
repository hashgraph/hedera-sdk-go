## Unreleased

### General Changes
 * Updated `TransactionId` to support new being `scheduled` and being constructed
   from `nonce`.

### Added
 * `ScheduleCreateTransaction` - Create a new scheduled transaction
 * `ScheduleSignTransaction` - Sign an existing scheduled transaction on the network
 * `ScheduleDeleteTransaction` - Delete a scheduled transaction
 * `ScheduleInfoQuery` - Query the info including `bodyBytes` of a scheduled transaction
 * `ScheduleId`

## v2.0.0

### General Changes

 * All requests support getter methods as well as setters.
 * All requests support multiple node account IDs being set.
 * `TransactionFromBytes()` supports multiple node account IDs and existing
    signatures.
 * All requests support a max retry count using `SetMaxRetry()`
 

# ABI package for using uint256

> Regarding this [Issue]: hedera-sdk do not support uint256 conversion from big.Int.
> Solution is usage of go-ethereum abi library.

[Issue]: https://github.com/hashgraph/hedera-sdk-go/issues/534


## Example

```
AddUint256(abi.U256(big.NewInt(some number)))
```

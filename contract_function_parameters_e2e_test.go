//go:build all || e2e
// +build all e2e

package hedera

import (
	"math/big"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	contractID   ContractID
	deployedOnce sync.Once
	intTypeMap   = map[string]intTypeData{"int256": {(*ContractFunctionParameters).AddInt256,
		func(result *ContractFunctionResult) []byte { return result.GetInt256(0) },
		"returnInt256",
	},
		"int248": {(*ContractFunctionParameters).AddInt248,
			func(result *ContractFunctionResult) []byte { return result.GetInt248(0) },
			"returnInt248",
		},
		"int240": {(*ContractFunctionParameters).AddInt240,
			func(result *ContractFunctionResult) []byte { return result.GetInt240(0) },
			"returnInt240",
		},
		"int232": {(*ContractFunctionParameters).AddInt232,
			func(result *ContractFunctionResult) []byte { return result.GetInt232(0) },
			"returnInt232",
		},
		"int224": {(*ContractFunctionParameters).AddInt224,
			func(result *ContractFunctionResult) []byte { return result.GetInt224(0) },
			"returnInt224",
		},
		"int216": {(*ContractFunctionParameters).AddInt216,
			func(result *ContractFunctionResult) []byte { return result.GetInt216(0) },
			"returnInt216",
		},
		"int208": {(*ContractFunctionParameters).AddInt208,
			func(result *ContractFunctionResult) []byte { return result.GetInt208(0) },
			"returnInt208",
		},
		"int200": {(*ContractFunctionParameters).AddInt200,
			func(result *ContractFunctionResult) []byte { return result.GetInt200(0) },
			"returnInt200",
		},
		"int192": {(*ContractFunctionParameters).AddInt192,
			func(result *ContractFunctionResult) []byte { return result.GetInt192(0) },
			"returnInt192",
		},
		"int184": {(*ContractFunctionParameters).AddInt184,
			func(result *ContractFunctionResult) []byte { return result.GetInt184(0) },
			"returnInt184",
		},
		"int176": {(*ContractFunctionParameters).AddInt176,
			func(result *ContractFunctionResult) []byte { return result.GetInt176(0) },
			"returnInt176",
		},
		"int168": {(*ContractFunctionParameters).AddInt168,
			func(result *ContractFunctionResult) []byte { return result.GetInt168(0) },
			"returnInt168",
		},
		"int160": {(*ContractFunctionParameters).AddInt160,
			func(result *ContractFunctionResult) []byte { return result.GetInt160(0) },
			"returnInt160",
		},
		"int152": {(*ContractFunctionParameters).AddInt152,
			func(result *ContractFunctionResult) []byte { return result.GetInt152(0) },
			"returnInt152",
		},
		"int144": {(*ContractFunctionParameters).AddInt144,
			func(result *ContractFunctionResult) []byte { return result.GetInt144(0) },
			"returnInt144",
		},
		"int136": {(*ContractFunctionParameters).AddInt136,
			func(result *ContractFunctionResult) []byte { return result.GetInt136(0) },
			"returnInt136",
		},
		"int128": {(*ContractFunctionParameters).AddInt128,
			func(result *ContractFunctionResult) []byte { return result.GetInt128(0) },
			"returnInt128",
		},
		"int120": {(*ContractFunctionParameters).AddInt120,
			func(result *ContractFunctionResult) []byte { return result.GetInt120(0) },
			"returnInt120",
		},
		"int112": {(*ContractFunctionParameters).AddInt112,
			func(result *ContractFunctionResult) []byte { return result.GetInt112(0) },
			"returnInt112",
		},
		"int104": {(*ContractFunctionParameters).AddInt104,
			func(result *ContractFunctionResult) []byte { return result.GetInt104(0) },
			"returnInt104",
		},
		"int96": {(*ContractFunctionParameters).AddInt96,
			func(result *ContractFunctionResult) []byte { return result.GetInt96(0) },
			"returnInt96",
		},
		"int88": {(*ContractFunctionParameters).AddInt88,
			func(result *ContractFunctionResult) []byte { return result.GetInt88(0) },
			"returnInt88",
		},
		"int80": {(*ContractFunctionParameters).AddInt80,
			func(result *ContractFunctionResult) []byte { return result.GetInt80(0) },
			"returnInt80",
		},
		"int72": {(*ContractFunctionParameters).AddInt72,
			func(result *ContractFunctionResult) []byte { return result.GetInt72(0) },
			"returnInt72",
		},
		"uint256": {(*ContractFunctionParameters).AddUint256,
			func(result *ContractFunctionResult) []byte { return result.GetUint256(0) },
			"returnUint256",
		},
		"uint248": {(*ContractFunctionParameters).AddUint248,
			func(result *ContractFunctionResult) []byte { return result.GetUint248(0) },
			"returnUint248",
		},
		"uint240": {(*ContractFunctionParameters).AddUint240,
			func(result *ContractFunctionResult) []byte { return result.GetUint240(0) },
			"returnUint240",
		},
		"uint232": {(*ContractFunctionParameters).AddUint232,
			func(result *ContractFunctionResult) []byte { return result.GetUint232(0) },
			"returnUint232",
		},
		"uint224": {(*ContractFunctionParameters).AddUint224,
			func(result *ContractFunctionResult) []byte { return result.GetUint224(0) },
			"returnUint224",
		},
		"uint216": {(*ContractFunctionParameters).AddUint216,
			func(result *ContractFunctionResult) []byte { return result.GetUint216(0) },
			"returnUint216",
		},
		"uint208": {(*ContractFunctionParameters).AddUint208,
			func(result *ContractFunctionResult) []byte { return result.GetUint208(0) },
			"returnUint208",
		},
		"uint200": {(*ContractFunctionParameters).AddUint200,
			func(result *ContractFunctionResult) []byte { return result.GetUint200(0) },
			"returnUint200",
		},
		"uint192": {(*ContractFunctionParameters).AddUint192,
			func(result *ContractFunctionResult) []byte { return result.GetInt192(0) },
			"returnUint192",
		},
		"uint184": {(*ContractFunctionParameters).AddUint184,
			func(result *ContractFunctionResult) []byte { return result.GetUint184(0) },
			"returnUint184",
		},
		"uint176": {(*ContractFunctionParameters).AddUint176,
			func(result *ContractFunctionResult) []byte { return result.GetUint176(0) },
			"returnUint176",
		},
		"uint168": {(*ContractFunctionParameters).AddUint168,
			func(result *ContractFunctionResult) []byte { return result.GetUint168(0) },
			"returnUint168",
		},
		"uint160": {(*ContractFunctionParameters).AddUint160,
			func(result *ContractFunctionResult) []byte { return result.GetUint160(0) },
			"returnUint160",
		},
		"uint152": {(*ContractFunctionParameters).AddUint152,
			func(result *ContractFunctionResult) []byte { return result.GetUint152(0) },
			"returnUint152",
		},
		"uint144": {(*ContractFunctionParameters).AddUint144,
			func(result *ContractFunctionResult) []byte { return result.GetUint144(0) },
			"returnUint144",
		},
		"uint136": {(*ContractFunctionParameters).AddUint136,
			func(result *ContractFunctionResult) []byte { return result.GetUint136(0) },
			"returnUint136",
		},
		"uint128": {(*ContractFunctionParameters).AddUint128,
			func(result *ContractFunctionResult) []byte { return result.GetUint128(0) },
			"returnUint128",
		},
		"uint120": {(*ContractFunctionParameters).AddUint120,
			func(result *ContractFunctionResult) []byte { return result.GetUint120(0) },
			"returnUint120",
		},
		"uint112": {(*ContractFunctionParameters).AddUint112,
			func(result *ContractFunctionResult) []byte { return result.GetUint112(0) },
			"returnUint112",
		},
		"uint104": {(*ContractFunctionParameters).AddUint104,
			func(result *ContractFunctionResult) []byte { return result.GetUint104(0) },
			"returnUint104",
		},
		"uint96": {(*ContractFunctionParameters).AddUint96,
			func(result *ContractFunctionResult) []byte { return result.GetUint96(0) },
			"returnUint96",
		},
		"uint88": {(*ContractFunctionParameters).AddUint88,
			func(result *ContractFunctionResult) []byte { return result.GetUint88(0) },
			"returnUint88",
		},
		"uint80": {(*ContractFunctionParameters).AddUint80,
			func(result *ContractFunctionResult) []byte { return result.GetUint80(0) },
			"returnUint80",
		},
		"uint72": {(*ContractFunctionParameters).AddUint72,
			func(result *ContractFunctionResult) []byte { return result.GetUint72(0) },
			"returnUint72",
		},
	}
)

type intTypeData struct {
	fnAdd     func(parameters *ContractFunctionParameters, value []byte) *ContractFunctionParameters
	fnExtract func(result *ContractFunctionResult) []byte
	fnName    string
}

func intType(t *testing.T, env IntegrationTestEnv, intType string, value string) {
	data, ok := intTypeMap[intType]
	require.True(t, ok)

	valueBigInt, ok := new(big.Int).SetString(value, 10)
	require.True(t, ok)

	contractCall, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).
		SetFunction(data.fnName, data.fnAdd(NewContractFunctionParameters(), toTwosComplementFromBigInt(valueBigInt))).
		Execute(env.Client)
	require.NoError(t, err)
	resultBigInt := new(big.Int)
	if strings.Contains(intType, "uint") {
		resultBigInt = new(big.Int).SetBytes(data.fnExtract(&contractCall))
	} else {
		resultBigInt = toBigIntFromTwosComplement(data.fnExtract(&contractCall))
	}

	require.Equal(t, valueBigInt.String(), resultBigInt.String())
}

func toBigIntFromTwosComplement(data []byte) *big.Int {
	isNegative := data[0]&0x80 == 0x80

	// If the number is positive, just use SetBytes.
	if !isNegative {
		c := big.NewInt(0)
		c.SetBytes(data)
		return c
	}

	// If the number is negative, calculate the two's complement.
	c := big.NewInt(0)
	c.SetBytes(data)

	c.Sub(c, big.NewInt(1))

	mask := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(len(data)*8)), nil)
	mask.Sub(mask, big.NewInt(1))
	c.Xor(c, mask)
	c.Neg(c)

	return c
}
func toTwosComplementFromBigInt(value *big.Int) []byte { 
	// First, get the bytes of the absolute value of the number.
	absBytes := value.Bytes()

	// If the number is positive or zero, pad the bytes with zeros.
	if value.Sign() >= 0 {
		result := make([]byte, 32-len(absBytes))
		return append(result, absBytes...)
	}

	// If the number is negative, we need to calculate the two's complement.
	result := make([]byte, 32)
	for i := range result {
		result[i] = 0xff
	}

	for i, b := range absBytes {
		result[len(result)-len(absBytes)+i] -= b
	}

	for i := len(result) - 1; i >= 0; i-- {
		result[i]++
		if result[i] != 0 {
			break
		}
	}

	return result
}

func deployContract(env IntegrationTestEnv) (*ContractID, error) {
	var result *ContractID
	var err error

	deployedOnce.Do(func() {
		result, err = performDeploy(env)
		if err == nil {
			contractID = *result
		}
	})
	return result, err
}

func performDeploy(env IntegrationTestEnv) (*ContractID, error) {
	bytecode := []byte(`0x608060405234801561001057600080fd5b50611851806100206000396000f3fe608060405234801561001057600080fd5b50600436106103fc5760003560e01c806388b7e6f511610215578063c503772d11610125578063de9fb484116100b8578063e713cda811610087578063e713cda814610cb2578063f4e490f514610cd5578063f6e877f414610cf6578063f8293f6e14610d1c578063ffb8050114610d3e57600080fd5b8063de9fb48414610c12578063e05e91e014610c3f578063e066de5014610c66578063e0f53e2414610c8c57600080fd5b8063cdb9e4e8116100f4578063cdb9e4e814610b7b578063d79d4d4014610ba1578063dade0c0b14610bc7578063dbb04ed914610be957600080fd5b8063c503772d14610ae6578063c6c18a1c14610b06578063c7d8b87e14610b30578063cbd2e6a514610b5657600080fd5b8063a75761f1116101a8578063b8da8d1611610177578063b8da8d1614610a2b578063b989c7ee14610a51578063ba945bdb14610a72578063bb6b524314610a98578063bd90536a14610abe57600080fd5b8063a75761f1146109a2578063b2db404a146109c3578063b4e3e7b1146109e4578063b834bfe914610a0a57600080fd5b80639b1794ae116101e45780639b1794ae14610914578063a08b9f671461093a578063a1bda1221461095b578063a401d60d1461097c57600080fd5b806388b7e6f514610890578063923f5edf146108b157806394cd7c80146108d257806398508ba3146108f357600080fd5b80633b45e6e01161031057806364e008c1116102a357806372a06b4d1161027257806372a06b4d146107e1578063796a27ea146108025780637d0dc262146108285780637ec32d8414610849578063881c8fb71461086a57600080fd5b806364e008c11461075257806368ef4466146107735780636a54715c1461079457806370a5cb81146107b557600080fd5b806344e7b037116102df57806344e7b037146106de578063545e21131461070457806359adb2df146105b3578063628bc3ef1461073157600080fd5b80633b45e6e0146106535780633e1a2771146106745780633f396e6714610695578063407b899b146106bd57600080fd5b806311ec6c901161039357806322937ea91161036257806322937ea9146105b35780632ef16e8e146105cf5780632f47a40d146105f05780632f6c1bb41461061157806333520ec31461063257600080fd5b806311ec6c901461052a578063129ed5da1461054b57806312cd95a114610571578063189cea8e1461059257600080fd5b806308123e09116103cf57806308123e09146104a15780630a958dc8146104c757806310d54553146104e8578063118b84151461050957600080fd5b8063017fa10b14610401578063021d88ab1461042f578063037454301461045a57806306ac6fe11461047b575b600080fd5b61041261040f36600461126e565b90565b6040516001600160801b0390911681526020015b60405180910390f35b61043d61040f366004611680565b6040516bffffffffffffffffffffffff9091168152602001610426565b61046861040f366004610dd6565b604051600c9190910b8152602001610426565b61048961040f366004611247565b6040516001600160781b039091168152602001610426565b6104af61040f366004611592565b60405166ffffffffffffff9091168152602001610426565b6104d561040f3660046110cb565b60405160049190910b8152602001610426565b6104f661040f366004610e82565b60405160119190910b8152602001610426565b61051761040f366004611071565b604051601e9190910b8152602001610426565b61053861040f366004610ee5565b60405160139190910b8152602001610426565b61055961040f3660046112e3565b6040516001600160981b039091168152602001610426565b61057f61040f366004610ea3565b60405160129190910b8152602001610426565b6105a061040f366004610f48565b60405160169190910b8152602001610426565b6105c161040f366004611092565b604051908152602001610426565b6105dd61040f366004611050565b604051601d9190910b8152602001610426565b6105fe61040f3660046111b1565b604051600a9190910b8152602001610426565b61061f61040f366004610fab565b60405160199190910b8152602001610426565b61064061040f366004610fcc565b604051601a9190910b8152602001610426565b61066161040f366004610e61565b60405160109190910b8152602001610426565b61068261040f36600461100e565b604051601c9190910b8152602001610426565b6106a361040f3660046115e1565b60405168ffffffffffffffffff9091168152602001610426565b6106cb61040f36600461112e565b60405160079190910b8152602001610426565b6106ec61040f36600461132c565b6040516001600160a01b039091168152602001610426565b610717610712366004611170565b610d62565b60408051600093840b81529190920b602082015201610426565b61073f61040f366004611190565b60405160099190910b8152602001610426565b61076061040f366004610f69565b60405160179190910b8152602001610426565b61078161040f366004610f06565b60405160149190910b8152602001610426565b6107a261040f36600461110d565b60405160069190910b8152602001610426565b6107c361040f3660046111f3565b6040516cffffffffffffffffffffffffff9091168152602001610426565b6107ef61040f366004611170565b60405160009190910b8152602001610426565b61081061040f36600461143d565b6040516001600160d81b039091168152602001610426565b61083661040f366004610fed565b604051601b9190910b8152602001610426565b61085761040f36600461102f565b60405160029190910b8152602001610426565b61087861040f36600461137a565b6040516001600160b01b039091168152602001610426565b61089e61040f366004610e40565b604051600f9190910b8152602001610426565b6108bf61040f366004610f27565b60405160159190910b8152602001610426565b6108e061040f366004610f8a565b60405160189190910b8152602001610426565b61090161040f366004610dfe565b604051600d9190910b8152602001610426565b61092261040f3660046113a1565b6040516001600160b81b039091168152602001610426565b61094861040f3660046111d2565b604051600b9190910b8152602001610426565b61096961040f366004610ec4565b60405160019190910b8152602001610426565b61098a61040f3660046114d5565b6040516001600160f01b039091168152602001610426565b6109b061040f36600461130a565b60405161ffff9091168152602001610426565b6109d161040f3660046110ec565b60405160059190910b8152602001610426565b6109f261040f366004611220565b6040516001600160701b039091168152602001610426565b610a1861040f366004610e1f565b604051600e9190910b8152602001610426565b610a3961040f366004611353565b6040516001600160a81b039091168152602001610426565b610a5f61040f3660046110aa565b60405160039190910b8152602001610426565b610a8061040f366004611295565b6040516001600160881b039091168152602001610426565b610aa661040f3660046113ef565b6040516001600160c81b039091168152602001610426565b610ad1610acc366004611092565b610d7a565b60408051928352602083019190915201610426565b610af461040f36600461160a565b60405160ff9091168152602001610426565b610b1461040f366004611655565b6040516affffffffffffffffffffff9091168152602001610426565b610b3e61040f3660046114fc565b6040516001600160f81b039091168152602001610426565b610b6461040f36600461156c565b60405165ffffffffffff9091168152602001610426565b610b8961040f366004611464565b6040516001600160e01b039091168152602001610426565b610baf61040f366004611416565b6040516001600160d01b039091168152602001610426565b610bda610bd5366004611523565b610d89565b604051610426939291906116ac565b610bf761040f36600461162b565b60405169ffffffffffffffffffff9091168152602001610426565b610c25610c203660046110cb565b610dc7565b60408051600493840b81529190920b602082015201610426565b610c4d61040f3660046115b9565b60405167ffffffffffffffff9091168152602001610426565b610c7461040f36600461148b565b6040516001600160e81b039091168152602001610426565b610c9a61040f3660046113c8565b6040516001600160c01b039091168152602001610426565b610cc061040f366004611523565b60405163ffffffff9091168152602001610426565b610ce361040f36600461114f565b60405160089190910b8152602001610426565b610d0461040f3660046112bc565b6040516001600160901b039091168152602001610426565b610d2a61040f3660046114b2565b60405162ffffff9091168152602001610426565b610d4c61040f366004611547565b60405164ffffffffff9091168152602001610426565b60008082610d718160146117a9565b91509150915091565b60008082610d7181600161171e565b600080606083610d9a6001826117e0565b6040805180820190915260028152614f4b60f01b602082015291945063ffffffff16925090509193909250565b60008082610d7181600161175f565b600060208284031215610de7578081fd5b813580600c0b8114610df7578182fd5b9392505050565b600060208284031215610e0f578081fd5b813580600d0b8114610df7578182fd5b600060208284031215610e30578081fd5b813580600e0b8114610df7578182fd5b600060208284031215610e51578081fd5b813580600f0b8114610df7578182fd5b600060208284031215610e72578081fd5b81358060100b8114610df7578182fd5b600060208284031215610e93578081fd5b81358060110b8114610df7578182fd5b600060208284031215610eb4578081fd5b81358060120b8114610df7578182fd5b600060208284031215610ed5578081fd5b81358060010b8114610df7578182fd5b600060208284031215610ef6578081fd5b81358060130b8114610df7578182fd5b600060208284031215610f17578081fd5b81358060140b8114610df7578182fd5b600060208284031215610f38578081fd5b81358060150b8114610df7578182fd5b600060208284031215610f59578081fd5b81358060160b8114610df7578182fd5b600060208284031215610f7a578081fd5b81358060170b8114610df7578182fd5b600060208284031215610f9b578081fd5b81358060180b8114610df7578182fd5b600060208284031215610fbc578081fd5b81358060190b8114610df7578182fd5b600060208284031215610fdd578081fd5b813580601a0b8114610df7578182fd5b600060208284031215610ffe578081fd5b813580601b0b8114610df7578182fd5b60006020828403121561101f578081fd5b813580601c0b8114610df7578182fd5b600060208284031215611040578081fd5b81358060020b8114610df7578182fd5b600060208284031215611061578081fd5b813580601d0b8114610df7578182fd5b600060208284031215611082578081fd5b813580601e0b8114610df7578182fd5b6000602082840312156110a3578081fd5b5035919050565b6000602082840312156110bb578081fd5b81358060030b8114610df7578182fd5b6000602082840312156110dc578081fd5b81358060040b8114610df7578182fd5b6000602082840312156110fd578081fd5b81358060050b8114610df7578182fd5b60006020828403121561111e578081fd5b81358060060b8114610df7578182fd5b60006020828403121561113f578081fd5b81358060070b8114610df7578182fd5b600060208284031215611160578081fd5b81358060080b8114610df7578182fd5b600060208284031215611181578081fd5b813580820b8114610df7578182fd5b6000602082840312156111a1578081fd5b81358060090b8114610df7578182fd5b6000602082840312156111c2578081fd5b813580600a0b8114610df7578182fd5b6000602082840312156111e3578081fd5b813580600b0b8114610df7578182fd5b600060208284031215611204578081fd5b81356cffffffffffffffffffffffffff81168114610df7578182fd5b600060208284031215611231578081fd5b81356001600160701b0381168114610df7578182fd5b600060208284031215611258578081fd5b81356001600160781b0381168114610df7578182fd5b60006020828403121561127f578081fd5b81356001600160801b0381168114610df7578182fd5b6000602082840312156112a6578081fd5b81356001600160881b0381168114610df7578182fd5b6000602082840312156112cd578081fd5b81356001600160901b0381168114610df7578182fd5b6000602082840312156112f4578081fd5b81356001600160981b0381168114610df7578182fd5b60006020828403121561131b578081fd5b813561ffff81168114610df7578182fd5b60006020828403121561133d578081fd5b81356001600160a01b0381168114610df7578182fd5b600060208284031215611364578081fd5b81356001600160a81b0381168114610df7578182fd5b60006020828403121561138b578081fd5b81356001600160b01b0381168114610df7578182fd5b6000602082840312156113b2578081fd5b81356001600160b81b0381168114610df7578182fd5b6000602082840312156113d9578081fd5b81356001600160c01b0381168114610df7578182fd5b600060208284031215611400578081fd5b81356001600160c81b0381168114610df7578182fd5b600060208284031215611427578081fd5b81356001600160d01b0381168114610df7578182fd5b60006020828403121561144e578081fd5b81356001600160d81b0381168114610df7578182fd5b600060208284031215611475578081fd5b81356001600160e01b0381168114610df7578182fd5b60006020828403121561149c578081fd5b81356001600160e81b0381168114610df7578182fd5b6000602082840312156114c3578081fd5b813562ffffff81168114610df7578182fd5b6000602082840312156114e6578081fd5b81356001600160f01b0381168114610df7578182fd5b60006020828403121561150d578081fd5b81356001600160f81b0381168114610df7578182fd5b600060208284031215611534578081fd5b813563ffffffff81168114610df7578182fd5b600060208284031215611558578081fd5b813564ffffffffff81168114610df7578182fd5b60006020828403121561157d578081fd5b813565ffffffffffff81168114610df7578182fd5b6000602082840312156115a3578081fd5b813566ffffffffffffff81168114610df7578182fd5b6000602082840312156115ca578081fd5b813567ffffffffffffffff81168114610df7578182fd5b6000602082840312156115f2578081fd5b813568ffffffffffffffffff81168114610df7578182fd5b60006020828403121561161b578081fd5b813560ff81168114610df7578182fd5b60006020828403121561163c578081fd5b813569ffffffffffffffffffff81168114610df7578182fd5b600060208284031215611666578081fd5b81356affffffffffffffffffffff81168114610df7578182fd5b600060208284031215611691578081fd5b81356bffffffffffffffffffffffff81168114610df7578182fd5b63ffffffff841681526000602067ffffffffffffffff851681840152606060408401528351806060850152825b818110156116f5578581018301518582016080015282016116d9565b818111156117065783608083870101525b50601f01601f19169290920160800195945050505050565b600080821280156001600160ff1b038490038513161561174057611740611805565b600160ff1b839003841281161561175957611759611805565b50500190565b60008160040b8360040b82821282647fffffffff0382138115161561178657611786611805565b82647fffffffff190382128116156117a0576117a0611805565b50019392505050565b600081810b83820b82821282607f038213811516156117ca576117ca611805565b82607f190382128116156117a0576117a0611805565b600063ffffffff838116908316818110156117fd576117fd611805565b039392505050565b634e487b7160e01b600052601160045260246000fdfea26469706673582212205364c98c43bf562527dfb30742be8777928382a6eda15d701082fc80093ac7f364736f6c63430008040033`)

	fileCreate, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey.PublicKey()).
		Execute(env.Client)
	if err != nil {
		return nil, err
	}
	fileCreate.SetValidateStatus(true)
	receipt, err := fileCreate.GetReceipt(env.Client)
	if err != nil {
		return nil, err
	}
	fileAppend, err := NewFileAppendTransaction().SetFileID(*receipt.FileID).SetContents(bytecode).Execute(env.Client)
	if err != nil {
		return nil, err
	}
	fileAppend.SetValidateStatus(true)
	_, err = fileAppend.GetReceipt(env.Client)
	if err != nil {
		return nil, err
	}
	contractCreate, err := NewContractCreateTransaction().
		SetBytecodeFileID(*receipt.FileID).
		SetGas(10000000).Execute(env.Client)
	if err != nil {
		return nil, err
	}
	contractCreate.SetValidateStatus(true)
	contractReceipt, err := contractCreate.GetReceipt(env.Client)
	if err != nil {
		return nil, err
	}
	return contractReceipt.ContractID, nil
	//return &ContractID{0, 0, 1034, nil, nil}, nil
}

func TestUint8Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint8(0)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint8", NewContractFunctionParameters().AddUint8(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint8(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint8Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint8(255)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint8", NewContractFunctionParameters().AddUint8(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint8(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint16Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint16(0)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint16", NewContractFunctionParameters().AddUint16(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint16(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint16Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint16(65535)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint16", NewContractFunctionParameters().AddUint16(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint16(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint24Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint32(0)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint24", NewContractFunctionParameters().AddUint24(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint24(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint24Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint32(16777215)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint24", NewContractFunctionParameters().AddUint24(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint24(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint32Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint32(0)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint32", NewContractFunctionParameters().AddUint32(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint32(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint32Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint32(4294967295)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint32", NewContractFunctionParameters().AddUint32(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint32(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint40Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint64(0)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint40", NewContractFunctionParameters().AddUint40(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint40(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint40Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint64(109951162777)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint40", NewContractFunctionParameters().AddUint40(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint40(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint48Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint64(0)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint48", NewContractFunctionParameters().AddUint48(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint48(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint48Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint64(281474976710655)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint48", NewContractFunctionParameters().AddUint48(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint48(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint56Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint64(0)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint56", NewContractFunctionParameters().AddUint56(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint56(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint56Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint64(72057594037927935)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint56", NewContractFunctionParameters().AddUint56(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint56(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint64Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint64(0)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint64", NewContractFunctionParameters().AddUint64(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint64(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint64Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint64(9223372036854775807)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnUint64", NewContractFunctionParameters().AddUint64(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint64(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint72Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint72", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint72Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint72", "4722366482869645213695")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint80Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint80", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint80Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint80", "1208925819614629174706175")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint88Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint88", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint88Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint88", "309485009821345068724781055")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint96Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint96", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint96Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint96", "79228162514264337593543950335")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint104Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint104", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint104Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint104", "20282409603651670423947251286015")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestUint112Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint112", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint112Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint112", "5192296858534827628530496329220095")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint120Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint120", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint120Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint120", "1329227995784915872903807060280344575")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint128Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint128", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint128Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint128", "340282366920938463463374607431768211455")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint136Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint136", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint136Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint136", "87112285931760246646623899502532662132735")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint144Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint144", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint144Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint144", "22300745198530623141535718272648361505980415")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint152Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint152", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint152Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint152", "5708990770823839524233143877797980545530986495")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint160Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint160", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint160Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint168", "1461501637330902918203684832716283019655932542975")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint168Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint168", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint168Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint168", "374144419156711147060143317175368453031918731001855")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint176Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint176", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint176Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint176", "95780971304118053647396689196894323976171195136475135")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint184Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint184", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint184Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint192", "24519928653854221733733552434404946937899825954937634815")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint192Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint192", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint192Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint192", "6277101735386680763835789423207666416102355444464034512895")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint200Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint200", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint200Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint200", "1606938044258990275541962092341162602522202993782792835301375")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint208Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint208", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint208Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint208", "411376139330301510538742295639337626245683966408394965837152255")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint216Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint224", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint216Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint224", "105312291668557186697918027683670432318895095400549111254310977535")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint224Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint224", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint224Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint224", "26959946667150639794667015087019630673637144422540572481103610249215")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint232Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint232", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint232Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint232", "6901746346790563787434755862277025452451108972170386555162524223799295")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint240Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint240", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint240Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint240", "1766847064778384329583297500742918515827483896875618958121606201292619775")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint248Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint248", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint248Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint248", "452312848583266388373324160190187140051835877600158453279131187530910662655")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint256Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint256", "0")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestUint256Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "uint256", "115792089237316195423570985008687907853269984665640564039457584007913129639935")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt8Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int8(-128)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt8", NewContractFunctionParameters().AddInt8(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt8(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt8Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int8(127)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt8", NewContractFunctionParameters().AddInt8(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt8(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt16Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int16(-32768)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt16", NewContractFunctionParameters().AddInt16(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt16(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt16Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int16(32767)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt16", NewContractFunctionParameters().AddInt16(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt16(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt24Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int32(-8388608)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt24", NewContractFunctionParameters().AddInt24(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt24(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt24Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int32(8388607)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt24", NewContractFunctionParameters().AddInt24(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt24(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt32Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int32(-2147483648)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt32", NewContractFunctionParameters().AddInt32(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt32(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt32Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int32(2147483647)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt32", NewContractFunctionParameters().AddInt32(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt32(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt40Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(-549755813888)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt40", NewContractFunctionParameters().AddInt40(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(value), contractCal.GetInt40(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt40Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(549755813887)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt40", NewContractFunctionParameters().AddInt40(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(value), contractCal.GetInt40(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt48Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(-140737488355328)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt48", NewContractFunctionParameters().AddInt48(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(value), contractCal.GetInt48(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt48Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(140737488355327)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt48", NewContractFunctionParameters().AddInt48(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(value), contractCal.GetInt48(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt56Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(-36028797018963968)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt56", NewContractFunctionParameters().AddInt56(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt56(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt56Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(36028797018963967)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt56", NewContractFunctionParameters().AddInt56(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt56(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt64Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(-9223372036854775808)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt64", NewContractFunctionParameters().AddInt64(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt64(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestInt64Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(9223372036854775807)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt64", NewContractFunctionParameters().AddInt64(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt64(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt72Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int72", "-2361183241434822606848")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt72Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int72", "2361183241434822606847")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt80Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int80", "-604462909807314587353088")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt80Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int80", "604462909807314587353087")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt88Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int88", "-154742504910672534362390528")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt88Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int88", "154742504910672534362390527")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt96Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int96", "-39614081257132168796771975168")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt96Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int96", "39614081257132168796771975167")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt104Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int104", "-10141204801825835211973625643008")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt104Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int104", "10141204801825835211973625643007")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt112Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int112", "-2596148429267413814265248164610048")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt112Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int112", "2596148429267413814265248164610047")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt120Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int120", "-664613997892457936451903530140172288")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt120Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int120", "664613997892457936451903530140172287")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt128Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int128", "-170141183460469231731687303715884105728")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt128Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int128", "170141183460469231731687303715884105727")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt136Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int136", "-43556142965880123323311949751266331066368")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt136Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int136", "43556142965880123323311949751266331066367")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt144Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int144", "-11150372599265311570767859136324180752990208")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt144Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int144", "11150372599265311570767859136324180752990207")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt152Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int152", "-2854495385411919762116571938898990272765493248")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt152Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int152", "2854495385411919762116571938898990272765493247")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt160Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int160", "-730750818665451459101842416358141509827966271488")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt160Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int160", "730750818665451459101842416358141509827966271487")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt168Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int168", "-187072209578355573530071658587684226515959365500928")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt168Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int168", "187072209578355573530071658587684226515959365500927")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt176Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int176", "-47890485652059026823698344598447161988085597568237568")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt176Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int176", "47890485652059026823698344598447161988085597568237567")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt184Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int184", "-12259964326927110866866776217202473468949912977468817408")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt184Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int184", "12259964326927110866866776217202473468949912977468817407")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt192Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int192", "-3138550867693340381917894711603833208051177722232017256448")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt192Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int192", "3138550867693340381917894711603833208051177722232017256447")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt200Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int200", "-803469022129495137770981046170581301261101496891396417650688")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt200Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int200", "803469022129495137770981046170581301261101496891396417650687")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt208Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int208", "-205688069665150755269371147819668813122841983204197482918576128")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt208Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int208", "205688069665150755269371147819668813122841983204197482918576127")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt216Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int216", "-52656145834278593348959013841835216159447547700274555627155488768")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt216Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int216", "52656145834278593348959013841835216159447547700274555627155488767")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt224Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int224", "-13479973333575319897333507543509815336818572211270286240551805124608")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt224Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int224", "13479973333575319897333507543509815336818572211270286240551805124607")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt232Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int232", "-3450873173395281893717377931138512726225554486085193277581262111899648")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt232Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int232", "3450873173395281893717377931138512726225554486085193277581262111899647")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt240Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int240", "-883423532389192164791648750371459257913741948437809479060803100646309888")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt240Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int240", "883423532389192164791648750371459257913741948437809479060803100646309887")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt248Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int248", "-226156424291633194186662080095093570025917938800079226639565593765455331328")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt248Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int248", "226156424291633194186662080095093570025917938800079226639565593765455331327")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt256Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int256", "-57896044618658097711785492504343953926634992332820282019728792003956564819968")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestInt256Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	intType(t, env, "int256", "57896044618658097711785492504343953926634992332820282019728792003956564819967")
	err := CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestMultipleInt8(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int8(-128)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnInt8Multiple", NewContractFunctionParameters().AddInt8(value)).SetQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt8(0))
	require.Equal(t, int8(-108), contractCal.GetInt8(1))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestMultipleInt40(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := int64(549755813885)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnMultipleInt40", NewContractFunctionParameters().AddInt40(value)).SetQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(549755813885), contractCal.GetInt40(0))
	require.Equal(t, int64(549755813886), contractCal.GetInt40(1))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestMultipleInt256(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value, ok := new(big.Int).SetString("-57896044618658097711785492504343953926634992332820282019728792003956564819968", 10)
	require.True(t, ok)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnMultipleInt256", NewContractFunctionParameters().AddInt256(toTwosComplementFromBigInt(value))).SetQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, toBigIntFromTwosComplement(contractCal.GetInt256(0)))
	require.Equal(t, value.Add(value, big.NewInt(1)), toBigIntFromTwosComplement(contractCal.GetInt256(1)))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestMultipleTypes(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := uint32(4294967295)
	contractCal, err := NewContractCallQuery().SetGas(15000000).
		SetContractID(contractID).SetFunction("returnMultipleTypeParams", NewContractFunctionParameters().AddUint32(value)).SetQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint32(0))
	require.Equal(t, uint64(4294967294), contractCal.GetUint64(1))
	require.Equal(t, "OK", contractCal.GetString(2))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

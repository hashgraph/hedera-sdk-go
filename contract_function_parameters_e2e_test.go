//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const gas = 50_000

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
	valueBigIntCopy := new(big.Int).Set(valueBigInt)

	contractCall, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).
		SetQueryPayment(NewHbar(20)).
		SetFunction(data.fnName, data.fnAdd(NewContractFunctionParameters(), To256BitBytes(valueBigInt))).
		Execute(env.Client)

	require.NoError(t, err)
	resultBigInt := new(big.Int)
	if strings.Contains(intType, "uint") {
		resultBigInt = new(big.Int).SetBytes(data.fnExtract(&contractCall))
	} else {
		value := new(big.Int).SetBytes(data.fnExtract(&contractCall))
		resultBigInt = ToSigned256(value)
	}

	require.Equal(t, valueBigIntCopy.String(), resultBigInt.String())
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
	bytecode := []byte(`0x608060405234801561001057600080fd5b50611f83806100206000396000f3fe608060405234801561001057600080fd5b50600436106104805760003560e01c806381dbe13e11610257578063bb6b524311610146578063dbb04ed9116100c3578063e713cda811610087578063e713cda814610df6578063f4e490f514610e19578063f6e877f414610e3a578063f8293f6e14610e60578063ffb8050114610e8257600080fd5b8063dbb04ed914610d2d578063de9fb48414610d56578063e05e91e014610d83578063e066de5014610daa578063e0f53e2414610dd057600080fd5b8063cbd2e6a51161010a578063cbd2e6a514610c9a578063cdb9e4e814610cbf578063d1b10ad71461066d578063d79d4d4014610ce5578063dade0c0b14610d0b57600080fd5b8063bb6b524314610bdc578063bd90536a14610c02578063c503772d14610c2a578063c6c18a1c14610c4a578063c7d8b87e14610c7457600080fd5b8063a1bda122116101d4578063b4e3e7b111610198578063b4e3e7b114610b28578063b834bfe914610b4e578063b8da8d1614610b6f578063b989c7ee14610b95578063ba945bdb14610bb657600080fd5b8063a1bda12214610a9f578063a401d60d14610ac0578063a75761f114610ae6578063aa80ca2e14610733578063b2db404a14610b0757600080fd5b8063923f5edf1161021b578063923f5edf146109f557806394cd7c8014610a1657806398508ba314610a375780639b1794ae14610a58578063a08b9f6714610a7e57600080fd5b806381dbe13e14610978578063827147ce146105cf578063881c8fb71461099357806388b7e6f5146109b9578063909c5b24146109da57600080fd5b806338fa665811610373578063628bc3ef116102f057806372a06b4d116102b457806372a06b4d146108ef578063796a27ea146109105780637d0dc262146109365780637ec32d84146109575780637f8082f71461073357600080fd5b8063628bc3ef1461083f57806364e008c11461086057806368ef4466146108815780636a54715c146108a257806370a5cb81146108c357600080fd5b806344e7b0371161033757806344e7b0371461068957806348d848d0146107d95780634bbc9a67146107f7578063545e21131461081257806359adb2df1461066d57600080fd5b806338fa6658146107335780633b45e6e01461074e5780633e1a27711461076f5780633f396e6714610790578063407b899b146107b857600080fd5b8063129ed5da116104015780632421101f116103c55780632421101f146106895780632ef16e8e146106af5780632f47a40d146106d05780632f6c1bb4146106f157806333520ec31461071257600080fd5b8063129ed5da146105ea57806312cd95a114610610578063189cea8e146106315780631d1145621461065257806322937ea91461066d57600080fd5b80630a958dc8116104485780630a958dc81461054b57806310d545531461056c578063118b84151461058d57806311ec6c90146105ae578063126bc815146105cf57600080fd5b8063017fa10b14610485578063021d88ab146104b357806303745430146104de57806306ac6fe1146104ff57806308123e0914610525575b600080fd5b6104966104933660046117ac565b90565b6040516001600160801b0390911681526020015b60405180910390f35b6104c1610493366004611b97565b6040516bffffffffffffffffffffffff90911681526020016104aa565b6104ec610493366004611333565b604051600c9190910b81526020016104aa565b61050d610493366004611785565b6040516001600160781b0390911681526020016104aa565b610533610493366004611aa9565b60405166ffffffffffffff90911681526020016104aa565b610559610493366004611609565b60405160049190910b81526020016104aa565b61057a6104933660046113d8565b60405160119190910b81526020016104aa565b61059b6104933660046115c7565b604051601e9190910b81526020016104aa565b6105bc61049336600461143b565b60405160139190910b81526020016104aa565b6105dd6104933660046112f8565b6040516104aa9190611d82565b6105f8610493366004611821565b6040516001600160981b0390911681526020016104aa565b61061e6104933660046113f9565b60405160129190910b81526020016104aa565b61063f61049336600461149e565b60405160169190910b81526020016104aa565b610660610493366004610fbd565b6040516104aa9190611c0e565b61067b6104933660046112e0565b6040519081526020016104aa565b610697610493366004610f9a565b6040516001600160a01b0390911681526020016104aa565b6106bd6104933660046115a6565b604051601d9190910b81526020016104aa565b6106de6104933660046116ef565b604051600a9190910b81526020016104aa565b6106ff610493366004611501565b60405160199190910b81526020016104aa565b610720610493366004611522565b604051601a9190910b81526020016104aa565b6107416104933660046110ec565b6040516104aa9190611c95565b61075c6104933660046113b7565b60405160109190910b81526020016104aa565b61077d610493366004611564565b604051601c9190910b81526020016104aa565b61079e610493366004611af8565b60405168ffffffffffffffffff90911681526020016104aa565b6107c661049336600461166c565b60405160079190910b81526020016104aa565b6107e76104933660046112c6565b60405190151581526020016104aa565b61080561049336600461105e565b6040516104aa9190611c5b565b6108256108203660046116ae565b610ea6565b60408051600093840b81529190920b6020820152016104aa565b61084d6104933660046116ce565b60405160099190910b81526020016104aa565b61086e6104933660046114bf565b60405160179190910b81526020016104aa565b61088f61049336600461145c565b60405160149190910b81526020016104aa565b6108b061049336600461164b565b60405160069190910b81526020016104aa565b6108d1610493366004611731565b6040516cffffffffffffffffffffffffff90911681526020016104aa565b6108fd6104933660046116ae565b60405160009190910b81526020016104aa565b61091e610493366004611954565b6040516001600160d81b0390911681526020016104aa565b610944610493366004611543565b604051601b9190910b81526020016104aa565b610965610493366004611585565b60405160029190910b81526020016104aa565b610986610493366004611224565b6040516104aa9190611d2e565b6109a1610493366004611891565b6040516001600160b01b0390911681526020016104aa565b6109c7610493366004611396565b604051600f9190910b81526020016104aa565b6109e8610493366004611173565b6040516104aa9190611ccd565b610a0361049336600461147d565b60405160159190910b81526020016104aa565b610a246104933660046114e0565b60405160189190910b81526020016104aa565b610a45610493366004611354565b604051600d9190910b81526020016104aa565b610a666104933660046118b8565b6040516001600160b81b0390911681526020016104aa565b610a8c610493366004611710565b604051600b9190910b81526020016104aa565b610aad61049336600461141a565b60405160019190910b81526020016104aa565b610ace6104933660046119ec565b6040516001600160f01b0390911681526020016104aa565b610af4610493366004611848565b60405161ffff90911681526020016104aa565b610b1561049336600461162a565b60405160059190910b81526020016104aa565b610b3661049336600461175e565b6040516001600160701b0390911681526020016104aa565b610b5c610493366004611375565b604051600e9190910b81526020016104aa565b610b7d61049336600461186a565b6040516001600160a81b0390911681526020016104aa565b610ba36104933660046115e8565b60405160039190910b81526020016104aa565b610bc46104933660046117d3565b6040516001600160881b0390911681526020016104aa565b610bea610493366004611906565b6040516001600160c81b0390911681526020016104aa565b610c15610c103660046112e0565b610ebe565b604080519283526020830191909152016104aa565b610c38610493366004611b21565b60405160ff90911681526020016104aa565b610c58610493366004611b6c565b6040516affffffffffffffffffffff90911681526020016104aa565b610c82610493366004611a13565b6040516001600160f81b0390911681526020016104aa565b610ca8610493366004611a83565b60405165ffffffffffff90911681526020016104aa565b610ccd61049336600461197b565b6040516001600160e01b0390911681526020016104aa565b610cf361049336600461192d565b6040516001600160d01b0390911681526020016104aa565b610d1e610d19366004611a3a565b610ecd565b6040516104aa93929190611d95565b610d3b610493366004611b42565b60405169ffffffffffffffffffff90911681526020016104aa565b610d69610d64366004611609565b610f0b565b60408051600493840b81529190920b6020820152016104aa565b610d91610493366004611ad0565b60405167ffffffffffffffff90911681526020016104aa565b610db86104933660046119a2565b6040516001600160e81b0390911681526020016104aa565b610dde6104933660046118df565b6040516001600160c01b0390911681526020016104aa565b610e04610493366004611a3a565b60405163ffffffff90911681526020016104aa565b610e2761049336600461168d565b60405160089190910b81526020016104aa565b610e486104933660046117fa565b6040516001600160901b0390911681526020016104aa565b610e6e6104933660046119c9565b60405162ffffff90911681526020016104aa565b610e90610493366004611a5e565b60405164ffffffffff90911681526020016104aa565b60008082610eb5816014611ead565b91509150915091565b60008082610eb5816001611e22565b600080606083610ede600182611ee4565b6040805180820190915260028152614f4b60f01b602082015291945063ffffffff16925090509193909250565b60008082610eb5816001611e63565b80358015158114610f2a57600080fd5b919050565b600082601f830112610f3f578081fd5b813567ffffffffffffffff811115610f5957610f59611f1f565b610f6c601f8201601f1916602001611dcd565b818152846020838601011115610f80578283fd5b816020850160208301379081016020019190915292915050565b600060208284031215610fab578081fd5b8135610fb681611f35565b9392505050565b60006020808385031215610fcf578182fd5b823567ffffffffffffffff811115610fe5578283fd5b8301601f81018513610ff5578283fd5b803561100861100382611dfe565b611dcd565b80828252848201915084840188868560051b8701011115611027578687fd5b8694505b8385101561105257803561103e81611f35565b83526001949094019391850191850161102b565b50979650505050505050565b60006020808385031215611070578182fd5b823567ffffffffffffffff811115611086578283fd5b8301601f81018513611096578283fd5b80356110a461100382611dfe565b80828252848201915084840188868560051b87010111156110c3578687fd5b8694505b83851015611052576110d881610f1a565b8352600194909401939185019185016110c7565b600060208083850312156110fe578182fd5b823567ffffffffffffffff811115611114578283fd5b8301601f81018513611124578283fd5b803561113261100382611dfe565b80828252848201915084840188868560051b8701011115611151578687fd5b8694505b83851015611052578035835260019490940193918501918501611155565b60006020808385031215611185578182fd5b823567ffffffffffffffff8082111561119c578384fd5b818501915085601f8301126111af578384fd5b81356111bd61100382611dfe565b80828252858201915085850189878560051b88010111156111dc578788fd5b875b84811015611215578135868111156111f457898afd5b6112028c8a838b0101610f2f565b85525092870192908701906001016111de565b50909998505050505050505050565b60006020808385031215611236578182fd5b823567ffffffffffffffff8082111561124d578384fd5b818501915085601f830112611260578384fd5b813561126e61100382611dfe565b80828252858201915085850189878560051b880101111561128d578788fd5b875b84811015611215578135868111156112a557898afd5b6112b38c8a838b0101610f2f565b855250928701929087019060010161128f565b6000602082840312156112d7578081fd5b610fb682610f1a565b6000602082840312156112f1578081fd5b5035919050565b600060208284031215611309578081fd5b813567ffffffffffffffff81111561131f578182fd5b61132b84828501610f2f565b949350505050565b600060208284031215611344578081fd5b813580600c0b8114610fb6578182fd5b600060208284031215611365578081fd5b813580600d0b8114610fb6578182fd5b600060208284031215611386578081fd5b813580600e0b8114610fb6578182fd5b6000602082840312156113a7578081fd5b813580600f0b8114610fb6578182fd5b6000602082840312156113c8578081fd5b81358060100b8114610fb6578182fd5b6000602082840312156113e9578081fd5b81358060110b8114610fb6578182fd5b60006020828403121561140a578081fd5b81358060120b8114610fb6578182fd5b60006020828403121561142b578081fd5b81358060010b8114610fb6578182fd5b60006020828403121561144c578081fd5b81358060130b8114610fb6578182fd5b60006020828403121561146d578081fd5b81358060140b8114610fb6578182fd5b60006020828403121561148e578081fd5b81358060150b8114610fb6578182fd5b6000602082840312156114af578081fd5b81358060160b8114610fb6578182fd5b6000602082840312156114d0578081fd5b81358060170b8114610fb6578182fd5b6000602082840312156114f1578081fd5b81358060180b8114610fb6578182fd5b600060208284031215611512578081fd5b81358060190b8114610fb6578182fd5b600060208284031215611533578081fd5b813580601a0b8114610fb6578182fd5b600060208284031215611554578081fd5b813580601b0b8114610fb6578182fd5b600060208284031215611575578081fd5b813580601c0b8114610fb6578182fd5b600060208284031215611596578081fd5b81358060020b8114610fb6578182fd5b6000602082840312156115b7578081fd5b813580601d0b8114610fb6578182fd5b6000602082840312156115d8578081fd5b813580601e0b8114610fb6578182fd5b6000602082840312156115f9578081fd5b81358060030b8114610fb6578182fd5b60006020828403121561161a578081fd5b81358060040b8114610fb6578182fd5b60006020828403121561163b578081fd5b81358060050b8114610fb6578182fd5b60006020828403121561165c578081fd5b81358060060b8114610fb6578182fd5b60006020828403121561167d578081fd5b81358060070b8114610fb6578182fd5b60006020828403121561169e578081fd5b81358060080b8114610fb6578182fd5b6000602082840312156116bf578081fd5b813580820b8114610fb6578182fd5b6000602082840312156116df578081fd5b81358060090b8114610fb6578182fd5b600060208284031215611700578081fd5b813580600a0b8114610fb6578182fd5b600060208284031215611721578081fd5b813580600b0b8114610fb6578182fd5b600060208284031215611742578081fd5b81356cffffffffffffffffffffffffff81168114610fb6578182fd5b60006020828403121561176f578081fd5b81356001600160701b0381168114610fb6578182fd5b600060208284031215611796578081fd5b81356001600160781b0381168114610fb6578182fd5b6000602082840312156117bd578081fd5b81356001600160801b0381168114610fb6578182fd5b6000602082840312156117e4578081fd5b81356001600160881b0381168114610fb6578182fd5b60006020828403121561180b578081fd5b81356001600160901b0381168114610fb6578182fd5b600060208284031215611832578081fd5b81356001600160981b0381168114610fb6578182fd5b600060208284031215611859578081fd5b813561ffff81168114610fb6578182fd5b60006020828403121561187b578081fd5b81356001600160a81b0381168114610fb6578182fd5b6000602082840312156118a2578081fd5b81356001600160b01b0381168114610fb6578182fd5b6000602082840312156118c9578081fd5b81356001600160b81b0381168114610fb6578182fd5b6000602082840312156118f0578081fd5b81356001600160c01b0381168114610fb6578182fd5b600060208284031215611917578081fd5b81356001600160c81b0381168114610fb6578182fd5b60006020828403121561193e578081fd5b81356001600160d01b0381168114610fb6578182fd5b600060208284031215611965578081fd5b81356001600160d81b0381168114610fb6578182fd5b60006020828403121561198c578081fd5b81356001600160e01b0381168114610fb6578182fd5b6000602082840312156119b3578081fd5b81356001600160e81b0381168114610fb6578182fd5b6000602082840312156119da578081fd5b813562ffffff81168114610fb6578182fd5b6000602082840312156119fd578081fd5b81356001600160f01b0381168114610fb6578182fd5b600060208284031215611a24578081fd5b81356001600160f81b0381168114610fb6578182fd5b600060208284031215611a4b578081fd5b813563ffffffff81168114610fb6578182fd5b600060208284031215611a6f578081fd5b813564ffffffffff81168114610fb6578182fd5b600060208284031215611a94578081fd5b813565ffffffffffff81168114610fb6578182fd5b600060208284031215611aba578081fd5b813566ffffffffffffff81168114610fb6578182fd5b600060208284031215611ae1578081fd5b813567ffffffffffffffff81168114610fb6578182fd5b600060208284031215611b09578081fd5b813568ffffffffffffffffff81168114610fb6578182fd5b600060208284031215611b32578081fd5b813560ff81168114610fb6578182fd5b600060208284031215611b53578081fd5b813569ffffffffffffffffffff81168114610fb6578182fd5b600060208284031215611b7d578081fd5b81356affffffffffffffffffffff81168114610fb6578182fd5b600060208284031215611ba8578081fd5b81356bffffffffffffffffffffffff81168114610fb6578182fd5b60008151808452815b81811015611be857602081850181015186830182015201611bcc565b81811115611bf95782602083870101525b50601f01601f19169290920160200192915050565b6020808252825182820181905260009190848201906040850190845b81811015611c4f5783516001600160a01b031683529284019291840191600101611c2a565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015611c4f578351151583529284019291840191600101611c77565b6020808252825182820181905260009190848201906040850190845b81811015611c4f57835183529284019291840191600101611cb1565b6000602080830181845280855180835260408601915060408160051b8701019250838701855b82811015611d2157603f19888603018452611d0f858351611bc3565b94509285019290850190600101611cf3565b5092979650505050505050565b6000602080830181845280855180835260408601915060408160051b8701019250838701855b82811015611d2157603f19888603018452611d70858351611bc3565b94509285019290850190600101611d54565b602081526000610fb66020830184611bc3565b63ffffffff8416815267ffffffffffffffff83166020820152606060408201526000611dc46060830184611bc3565b95945050505050565b604051601f8201601f1916810167ffffffffffffffff81118282101715611df657611df6611f1f565b604052919050565b600067ffffffffffffffff821115611e1857611e18611f1f565b5060051b60200190565b600080821280156001600160ff1b0384900385131615611e4457611e44611f09565b600160ff1b8390038412811615611e5d57611e5d611f09565b50500190565b60008160040b8360040b82821282647fffffffff03821381151615611e8a57611e8a611f09565b82647fffffffff19038212811615611ea457611ea4611f09565b50019392505050565b600081810b83820b82821282607f03821381151615611ece57611ece611f09565b82607f19038212811615611ea457611ea4611f09565b600063ffffffff83811690831681811015611f0157611f01611f09565b039392505050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114611f4a57600080fd5b5056fea264697066735822122027163c9c7a018e3f491b10f71ff4861efc506503e9f39bd3fc08dc44e99cd34c64736f6c63430008040033`)

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
		SetGas(500000).Execute(env.Client)
	if err != nil {
		return nil, err
	}
	contractCreate.SetValidateStatus(true)
	contractReceipt, err := contractCreate.GetReceipt(env.Client)
	if err != nil {
		return nil, err
	}
	return contractReceipt.ContractID, nil
}

func TestUint8Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint8(0)
	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(contractID).SetFunction("returnUint8", NewContractFunctionParameters().AddUint8(value)).Execute(env.Client)
	require.NoError(t, err)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint8", NewContractFunctionParameters().AddUint8(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint8(0))

}
func TestUint8Max(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint8(255)
	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(contractID).SetFunction("returnUint8", NewContractFunctionParameters().AddUint8(value)).Execute(env.Client)
	require.NoError(t, err)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint8", NewContractFunctionParameters().AddUint8(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint8(0))
}

func TestUint16Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint16(0)
	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(contractID).SetFunction("returnUint16", NewContractFunctionParameters().AddUint16(value)).Execute(env.Client)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint16", NewContractFunctionParameters().AddUint16(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint16(0))

}
func TestUint16Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint16(65535)
	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(contractID).SetFunction("returnUint16", NewContractFunctionParameters().AddUint16(value)).Execute(env.Client)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint16", NewContractFunctionParameters().AddUint16(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint16(0))

}

func TestUint24Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint32(0)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint24", NewContractFunctionParameters().AddUint24(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint24(0))

}
func TestUint24Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint32(16777215)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint24", NewContractFunctionParameters().AddUint24(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint24(0))

}
func TestUint32Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint32(0)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint32", NewContractFunctionParameters().AddUint32(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint32(0))
}
func TestUint32Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint32(4294967295)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint32", NewContractFunctionParameters().AddUint32(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint32(0))
}

func TestUint40Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint64(0)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint40", NewContractFunctionParameters().AddUint40(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint40(0))

}
func TestUint40Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint64(109951162777)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint40", NewContractFunctionParameters().AddUint40(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint40(0))

}

func TestUint48Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint64(0)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint48", NewContractFunctionParameters().AddUint48(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint48(0))

}
func TestUint48Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint64(281474976710655)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint48", NewContractFunctionParameters().AddUint48(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint48(0))

}

func TestUint56Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint64(0)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint56", NewContractFunctionParameters().AddUint56(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint56(0))

}
func TestUint56Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint64(72057594037927935)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint56", NewContractFunctionParameters().AddUint56(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint56(0))

}

func TestUint64Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint64(0)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint64", NewContractFunctionParameters().AddUint64(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint64(0))

}
func TestUint64Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint64(9223372036854775807)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint64", NewContractFunctionParameters().AddUint64(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint64(0))

}

func TestUint72Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint72", "0")

}
func TestUint72Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint72", "4722366482869645213695")

}

func TestUint80Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint80", "0")

}
func TestUint80Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint80", "1208925819614629174706175")

}
func TestUint88Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint88", "0")

}
func TestUint88Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint88", "309485009821345068724781055")

}
func TestUint96Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint96", "0")

}
func TestUint96Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint96", "79228162514264337593543950335")

}
func TestUint104Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint104", "0")

}
func TestUint104Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint104", "20282409603651670423947251286015")

}

func TestUint112Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint112", "0")

}
func TestUint112Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint112", "5192296858534827628530496329220095")

}
func TestUint120Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint120", "0")

}
func TestUint120Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint120", "1329227995784915872903807060280344575")

}
func TestUint128Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint128", "0")

}
func TestUint128Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint128", "340282366920938463463374607431768211455")

}
func TestUint136Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint136", "0")

}
func TestUint136Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint136", "87112285931760246646623899502532662132735")

}
func TestUint144Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint144", "0")

}
func TestUint144Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint144", "22300745198530623141535718272648361505980415")

}
func TestUint152Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint152", "0")

}
func TestUint152Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint152", "5708990770823839524233143877797980545530986495")

}
func TestUint160Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint160", "0")

}
func TestUint160Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint160", "1461501637330902918203684832716283019655932542975")

}
func TestUint168Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint168", "0")

}
func TestUint168Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint168", "374144419156711147060143317175368453031918731001855")

}
func TestUint176Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint176", "0")

}
func TestUint176Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint176", "95780971304118053647396689196894323976171195136475135")

}
func TestUint184Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint184", "0")

}
func TestUint184Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint184", "24519928653854221733733552434404946937899825954937634815")

}
func TestUint192Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint192", "0")

}
func TestUint192Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint192", "6277101735386680763835789423207666416102355444464034512895")

}
func TestUint200Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint200", "0")

}
func TestUint200Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint200", "1606938044258990275541962092341162602522202993782792835301375")

}
func TestUint208Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint208", "0")

}
func TestUint208Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint208", "411376139330301510538742295639337626245683966408394965837152255")

}
func TestUint216Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint216", "0")

}
func TestUint216Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint216", "105312291668557186697918027683670432318895095400549111254310977535")

}
func TestUint224Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint224", "0")

}
func TestUint224Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint224", "26959946667150639794667015087019630673637144422540572481103610249215")

}
func TestUint232Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint232", "0")

}
func TestUint232Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint232", "6901746346790563787434755862277025452451108972170386555162524223799295")

}
func TestUint240Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint240", "0")

}
func TestUint240Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint240", "1766847064778384329583297500742918515827483896875618958121606201292619775")

}
func TestUint248Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint248", "0")

}
func TestUint248Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint248", "452312848583266388373324160190187140051835877600158453279131187530910662655")

}
func TestUint256Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint256", "0")

}
func TestUint256Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "uint256", "115792089237316195423570985008687907853269984665640564039457584007913129639935")

}

func TestInt8Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int8(-128)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt8", NewContractFunctionParameters().AddInt8(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt8(0))

}
func TestInt8Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int8(127)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt8", NewContractFunctionParameters().AddInt8(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt8(0))

}

func TestInt16Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int16(-32768)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt16", NewContractFunctionParameters().AddInt16(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt16(0))

}
func TestInt16Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int16(32767)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt16", NewContractFunctionParameters().AddInt16(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt16(0))

}

func TestInt24Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int32(-8388608)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt24", NewContractFunctionParameters().AddInt24(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt24(0))

}
func TestInt24Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int32(8388607)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt24", NewContractFunctionParameters().AddInt24(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt24(0))

}
func TestInt32Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int32(-2147483648)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt32", NewContractFunctionParameters().AddInt32(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt32(0))

}
func TestInt32Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int32(2147483647)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt32", NewContractFunctionParameters().AddInt32(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt32(0))

}

func TestInt40Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(-549755813888)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt40", NewContractFunctionParameters().AddInt40(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(value), contractCal.GetInt40(0))

}
func TestInt40Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(549755813887)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt40", NewContractFunctionParameters().AddInt40(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(value), contractCal.GetInt40(0))

}

func TestInt48Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(-140737488355328)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt48", NewContractFunctionParameters().AddInt48(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(value), contractCal.GetInt48(0))

}
func TestInt48Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(140737488355327)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt48", NewContractFunctionParameters().AddInt48(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(value), contractCal.GetInt48(0))

}

func TestInt56Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(-36028797018963968)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt56", NewContractFunctionParameters().AddInt56(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt56(0))

}
func TestInt56Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(36028797018963967)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt56", NewContractFunctionParameters().AddInt56(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt56(0))

}

func TestInt64Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(-9223372036854775808)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt64", NewContractFunctionParameters().AddInt64(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt64(0))

}
func TestInt64Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(9223372036854775807)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt64", NewContractFunctionParameters().AddInt64(value)).SetMaxQueryPayment(NewHbar(20)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt64(0))

}

func TestInt72Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int72", "-2361183241434822606848")

}

func TestInt72Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int72", "2361183241434822606847")

}

func TestInt80Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int80", "-604462909807314587353088")

}

func TestInt80Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int80", "604462909807314587353087")

}

func TestInt88Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int88", "-154742504910672534362390528")

}

func TestInt88Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int88", "154742504910672534362390527")

}

func TestInt96Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int96", "-39614081257132168796771975168")

}

func TestInt96Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int96", "39614081257132168796771975167")

}

func TestInt104Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int104", "-10141204801825835211973625643008")

}

func TestInt104Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int104", "10141204801825835211973625643007")

}

func TestInt112Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int112", "-2596148429267413814265248164610048")

}

func TestInt112Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int112", "2596148429267413814265248164610047")

}

func TestInt120Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int120", "-664613997892457936451903530140172288")

}

func TestInt120Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int120", "664613997892457936451903530140172287")

}

func TestInt128Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int128", "-170141183460469231731687303715884105728")

}

func TestInt128Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int128", "170141183460469231731687303715884105727")

}

func TestInt136Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int136", "-43556142965880123323311949751266331066368")

}

func TestInt136Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int136", "43556142965880123323311949751266331066367")

}

func TestInt144Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int144", "-11150372599265311570767859136324180752990208")

}

func TestInt144Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int144", "11150372599265311570767859136324180752990207")

}

func TestInt152Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int152", "-2854495385411919762116571938898990272765493248")

}

func TestInt152Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int152", "2854495385411919762116571938898990272765493247")

}

func TestInt160Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int160", "-730750818665451459101842416358141509827966271488")

}

func TestInt160Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int160", "730750818665451459101842416358141509827966271487")

}

func TestInt168Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int168", "-187072209578355573530071658587684226515959365500928")

}

func TestInt168Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int168", "187072209578355573530071658587684226515959365500927")

}

func TestInt176Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int176", "-47890485652059026823698344598447161988085597568237568")

}

func TestInt176Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int176", "47890485652059026823698344598447161988085597568237567")

}

func TestInt184Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int184", "-12259964326927110866866776217202473468949912977468817408")

}

func TestInt184Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int184", "12259964326927110866866776217202473468949912977468817407")

}

func TestInt192Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int192", "-3138550867693340381917894711603833208051177722232017256448")

}

func TestInt192Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int192", "3138550867693340381917894711603833208051177722232017256447")

}

func TestInt200Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int200", "-803469022129495137770981046170581301261101496891396417650688")

}

func TestInt200Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int200", "803469022129495137770981046170581301261101496891396417650687")

}

func TestInt208Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int208", "-205688069665150755269371147819668813122841983204197482918576128")

}

func TestInt208Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int208", "205688069665150755269371147819668813122841983204197482918576127")

}

func TestInt216Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int216", "-52656145834278593348959013841835216159447547700274555627155488768")

}

func TestInt216Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int216", "52656145834278593348959013841835216159447547700274555627155488767")

}

func TestInt224Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int224", "-13479973333575319897333507543509815336818572211270286240551805124608")

}

func TestInt224Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int224", "13479973333575319897333507543509815336818572211270286240551805124607")

}

func TestInt232Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int232", "-3450873173395281893717377931138512726225554486085193277581262111899648")

}

func TestInt232Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int232", "3450873173395281893717377931138512726225554486085193277581262111899647")

}

func TestInt240Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int240", "-883423532389192164791648750371459257913741948437809479060803100646309888")

}

func TestInt240Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int240", "883423532389192164791648750371459257913741948437809479060803100646309887")

}

func TestInt248Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int248", "-226156424291633194186662080095093570025917938800079226639565593765455331328")

}

func TestInt248Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int248", "226156424291633194186662080095093570025917938800079226639565593765455331327")

}

func TestInt256Min(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int256", "-57896044618658097711785492504343953926634992332820282019728792003956564819968")

}

func TestInt256Max(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	intType(t, env, "int256", "57896044618658097711785492504343953926634992332820282019728792003956564819967")

}

func TestMultipleInt8(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int8(-128)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt8Multiple", NewContractFunctionParameters().AddInt8(value)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetInt8(0))
	require.Equal(t, int8(-108), contractCal.GetInt8(1))

}
func TestMultipleInt40(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := int64(549755813885)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnMultipleInt40", NewContractFunctionParameters().AddInt40(value)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(549755813885), contractCal.GetInt40(0))
	require.Equal(t, int64(549755813886), contractCal.GetInt40(1))

}
func TestMultipleInt256(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value, ok := new(big.Int).SetString("-123", 10)
	require.True(t, ok)
	valueTwos := To256BitBytes(value)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnMultipleInt256", NewContractFunctionParameters().AddInt256(valueTwos)).Execute(env.Client)
	require.NoError(t, err)
	value1, ok := new(big.Int).SetString("-123", 10)
	require.True(t, ok)
	value2, ok := new(big.Int).SetString("-122", 10)
	require.True(t, ok)
	require.Equal(t, To256BitBytes(value1), contractCal.GetInt256(0))
	require.Equal(t, To256BitBytes(value2), contractCal.GetInt256(1))
}

func TestMultipleTypes(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := uint32(4294967295)
	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnMultipleTypeParams", NewContractFunctionParameters().AddUint32(value)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetUint32(0))
	require.Equal(t, uint64(4294967294), contractCal.GetUint64(1))
	require.Equal(t, "OK", contractCal.GetString(2))

}

func TestBigInt256(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value, ok := new(big.Int).SetString("-123", 10)
	require.True(t, ok)

	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnInt256", NewContractFunctionParameters().AddInt256BigInt(value)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetBigInt(0))

}

func TestBigUint256(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value, ok := new(big.Int).SetString("123", 10)
	require.True(t, ok)

	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnUint256", NewContractFunctionParameters().AddUint256BigInt(value)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetBigInt(0))

}

func TestMultiplBigInt256(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value, ok := new(big.Int).SetString("-123", 10)
	require.True(t, ok)

	contractCal, err := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnMultipleInt256", NewContractFunctionParameters().AddInt256BigInt(value)).Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, contractCal.GetBigInt(0))
	require.Equal(t, new(big.Int).Add(value, big.NewInt(1)), contractCal.GetBigInt(1))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
func TestString(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := "Test"

	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnString", NewContractFunctionParameters().AddString(value))
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, result.GetString(0))
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestStringArray(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	deployContract(env)
	value := []string{"Test1", "Test2"}

	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnStringArr", NewContractFunctionParameters().AddStringArray(value))
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	parsedResult, _ := result.GetResult("string[]")
	strArr := parsedResult.([]string)
	require.Equal(t, value[0], strArr[0])
	require.Equal(t, value[1], strArr[1])

}
func TestAddress(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := "1234567890123456789012345678901234567890"
	params, err := NewContractFunctionParameters().AddAddress(value)
	require.NoError(t, err)
	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(contractID).SetFunction("returnAddress", params).Execute(env.Client)
	require.NoError(t, err)

	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnAddress", params)
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, hex.EncodeToString(result.GetAddress(0)))

}

func TestAddressArray(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := []string{"1234567890123456789012345678901234567890", "1234567890123456789012345678901234567891"}
	params, err := NewContractFunctionParameters().AddAddressArray(value)
	require.NoError(t, err)
	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnAddressArr", params)
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	addArr, err := result.GetResult("address[]")
	require.NoError(t, err)
	addresses := addArr.([]Address)
	require.Equal(t, value[0], strings.TrimPrefix(addresses[0].String(), "0x"))
	require.Equal(t, value[1], strings.TrimPrefix(addresses[1].String(), "0x"))

}

func TestBoolean(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := true

	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnBoolean", NewContractFunctionParameters().AddBool(value))
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, result.GetBool(0))

}

func TestBytes(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := []byte("Test")

	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnBytes", NewContractFunctionParameters().AddBytes(value))
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, value, result.GetBytes(0))

}

func TestBytesArray(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := [][]byte{[]byte("Test1"), []byte("Test2")}

	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnBytesArr", NewContractFunctionParameters().AddBytesArray(value))
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	bytesArrInterface, err := result.GetResult("bytes[]")
	require.NoError(t, err)
	require.Equal(t, value, bytesArrInterface.([][]uint8))

}

func TestBytes32(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)
	value := [32]byte{}
	copy(value[:], []byte("Test"))

	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).
		SetContractID(contractID).SetFunction("returnBytes32", NewContractFunctionParameters().AddBytes32(value))
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(value[:], result.GetBytes32(0)))

}

func TestBytes32Array(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	deployContract(env)

	value := [][]byte{
		[]byte("Test1"),
		[]byte("Test2"),
	}
	var expected1 [32]byte
	var expected2 [32]byte
	copy(expected1[len(expected1)-len(value[0]):], value[0])
	copy(expected2[len(expected2)-len(value[1]):], value[1])

	contractCal := NewContractCallQuery().SetGas(gas).SetQueryPayment(NewHbar(1)).SetQueryPayment(NewHbar(11)).
		SetContractID(contractID).SetFunction("returnBytes32Arr", NewContractFunctionParameters().AddBytes32Array(value))
	result, err := contractCal.Execute(env.Client)
	require.NoError(t, err)
	bytes32ArrInterface, err := result.GetResult("bytes32[]")
	require.NoError(t, err)
	require.Equal(t, expected1, bytes32ArrInterface.([][32]byte)[0])
	require.Equal(t, expected2, bytes32ArrInterface.([][32]byte)[1])

}

func TestContractNonces(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	bytecode := []byte(`6080604052348015600f57600080fd5b50604051601a90603b565b604051809103906000f0801580156035573d6000803e3d6000fd5b50506047565b605c8061009483390190565b603f806100556000396000f3fe6080604052600080fdfea2646970667358221220a20122cbad3457fedcc0600363d6e895f17048f5caa4afdab9e655123737567d64736f6c634300081200336080604052348015600f57600080fd5b50603f80601d6000396000f3fe6080604052600080fdfea264697066735822122053dfd8835e3dc6fedfb8b4806460b9b7163f8a7248bac510c6d6808d9da9d6d364736f6c63430008120033`)
	fileCreate, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey.PublicKey()).
		SetContents(bytecode).
		Execute(env.Client)
	require.NoError(t, err)
	fileCreate.SetValidateStatus(true)
	receipt, err := fileCreate.GetReceipt(env.Client)
	require.NoError(t, err)
	require.Equal(t, StatusSuccess, receipt.Status)
	contractCreate, err := NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(100000).
		SetBytecodeFileID(*receipt.FileID).
		SetContractMemo("[e2e::ContractADeploysContractBInConstructor]").
		Execute(env.Client)
	require.NoError(t, err)
	contractCreate.SetValidateStatus(true)
	record, err := contractCreate.GetRecord(env.Client)
	require.NoError(t, err)
	require.Equal(t, StatusSuccess, record.Receipt.Status)
	require.Equal(t, int64(2), record.CallResult.ContractNonces[0].Nonce)
	require.Equal(t, int64(1), record.CallResult.ContractNonces[1].Nonce)
}

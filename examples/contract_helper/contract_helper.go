package contract_helper

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/pkg/errors"
)

type ContractHelper struct {
	ContractID             hedera.ContractID
	stepResultValidators   map[int32]func(hedera.ContractFunctionResult) bool
	stepParameterSuppliers map[int32]func() *hedera.ContractFunctionParameters
	stepPayableAmounts     map[int32]*hedera.Hbar
	stepSigners            map[int32][]hedera.PrivateKey
	stepFeePayers          map[int32]*hedera.AccountID
	stepLogic              map[int32]func(address string)
}

func NewContractHelper(bytecode []byte, constructorParameters hedera.ContractFunctionParameters, client *hedera.Client) *ContractHelper {
	response, err := hedera.NewContractCreateFlow().
		SetBytecode(bytecode).
		SetGas(8000000).
		SetMaxChunks(30).
		SetConstructorParameters(&constructorParameters).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err := response.GetReceipt(client)
	if err != nil {
		panic(err)
	}
	if receipt.ContractID != nil {
		return &ContractHelper{
			ContractID:             *receipt.ContractID,
			stepResultValidators:   make(map[int32]func(hedera.ContractFunctionResult) bool),
			stepParameterSuppliers: make(map[int32]func() *hedera.ContractFunctionParameters),
			stepPayableAmounts:     make(map[int32]*hedera.Hbar),
			stepSigners:            make(map[int32][]hedera.PrivateKey),
			stepFeePayers:          make(map[int32]*hedera.AccountID),
			stepLogic:              make(map[int32]func(address string)),
		}
	}

	return &ContractHelper{}
}

func (this *ContractHelper) SetResultValidatorForStep(stepIndex int32, validator func(hedera.ContractFunctionResult) bool) *ContractHelper {
	this.stepResultValidators[stepIndex] = validator
	return this
}

func (this *ContractHelper) SetParameterSupplierForStep(stepIndex int32, supplier func() *hedera.ContractFunctionParameters) *ContractHelper {
	this.stepParameterSuppliers[stepIndex] = supplier
	return this
}

func (this *ContractHelper) SetPayableAmountForStep(stepIndex int32, amount hedera.Hbar) *ContractHelper {
	this.stepPayableAmounts[stepIndex] = &amount
	return this
}

func (this *ContractHelper) AddSignerForStep(stepIndex int32, signer hedera.PrivateKey) *ContractHelper {
	if _, ok := this.stepSigners[stepIndex]; ok {
		this.stepSigners[stepIndex] = append(this.stepSigners[stepIndex], signer)
	} else {
		this.stepSigners[stepIndex] = make([]hedera.PrivateKey, 0)
		this.stepSigners[stepIndex] = append(this.stepSigners[stepIndex], signer)
	}

	return this
}

func (this *ContractHelper) SetFeePayerForStep(stepIndex int32, account hedera.AccountID, accountKey hedera.PrivateKey) *ContractHelper {
	this.stepFeePayers[stepIndex] = &account
	return this.AddSignerForStep(stepIndex, accountKey)
}

func (this *ContractHelper) SetStepLogic(stepIndex int32, specialFunction func(address string)) *ContractHelper {
	this.stepLogic[stepIndex] = specialFunction
	return this
}

func (this *ContractHelper) GetResultValidator(stepIndex int32) func(hedera.ContractFunctionResult) bool {
	if _, ok := this.stepResultValidators[stepIndex]; ok {
		return this.stepResultValidators[stepIndex]
	}

	return func(result hedera.ContractFunctionResult) bool {
		responseStatus := hedera.Status(result.GetInt32(0))
		isValid := responseStatus == hedera.StatusSuccess
		if !isValid {
			println("Encountered invalid response status", responseStatus.String())
		}
		return isValid
	}
}

func (this *ContractHelper) GetParameterSupplier(stepIndex int32) func() *hedera.ContractFunctionParameters {
	if _, ok := this.stepParameterSuppliers[stepIndex]; ok {
		return this.stepParameterSuppliers[stepIndex]
	}

	return func() *hedera.ContractFunctionParameters {
		return nil
	}
}

func (this *ContractHelper) GetPayableAmount(stepIndex int32) *hedera.Hbar {
	return this.stepPayableAmounts[stepIndex]
}

func (this *ContractHelper) GetSigners(stepIndex int32) []hedera.PrivateKey {
	if _, ok := this.stepSigners[stepIndex]; ok {
		return this.stepSigners[stepIndex]
	}

	return []hedera.PrivateKey{}
}

func (this *ContractHelper) ExecuteSteps(firstStep int32, lastStep int32, client *hedera.Client) (*ContractHelper, error) {
	for stepIndex := firstStep; stepIndex <= lastStep; stepIndex++ {
		println("Attempting to execuite step", stepIndex)

		transaction := hedera.NewContractExecuteTransaction().
			SetContractID(this.ContractID).
			SetGas(10000000)

		payableAmount := this.GetPayableAmount(stepIndex)
		if payableAmount != nil {
			transaction.SetPayableAmount(*payableAmount)
		}

		functionName := "step" + strconv.Itoa(int(stepIndex))
		parameters := this.GetParameterSupplier(stepIndex)()
		if parameters != nil {
			transaction.SetFunction(functionName, parameters)
		} else {
			transaction.SetFunction(functionName, nil)
		}

		if feePayerAccountID, ok := this.stepFeePayers[stepIndex]; ok {
			transaction.SetTransactionID(hedera.TransactionIDGenerate(*feePayerAccountID))
		}

		frozen, err := transaction.FreezeWith(client)
		if err != nil {
			return &ContractHelper{}, err
		}
		for _, signer := range this.GetSigners(stepIndex) {
			frozen.Sign(signer)
		}

		response, err := frozen.Execute(client)
		if err != nil {
			return &ContractHelper{}, err
		}

		record, err := response.GetRecord(client)
		if err != nil {
			return &ContractHelper{}, err
		}

		functionResult, err := record.GetContractExecuteResult()
		if err != nil {
			return &ContractHelper{}, err
		}

		if this.stepLogic[stepIndex] != nil {
			address := functionResult.GetAddress(1)
			if function, exists := this.stepLogic[stepIndex]; exists && function != nil {
				function(hex.EncodeToString(address))
			}
		}

		if this.GetResultValidator(stepIndex)(functionResult) {
			fmt.Printf("Step %d completed, and returned valid result. (TransactionId %s)", stepIndex, record.TransactionID.String())
		} else {
			return &ContractHelper{}, errors.New(fmt.Sprintf("Step %d returned invalid result", stepIndex))
		}
	}

	return this, nil
}

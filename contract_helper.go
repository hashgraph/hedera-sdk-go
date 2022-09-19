package hedera

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type ContractHelper struct {
	contractID             ContractID
	stepResultValidators   map[int32]func(ContractFunctionResult) bool
	stepParameterSuppliers map[int32]func() *ContractFunctionParameters
	stepPayableAmounts     map[int32]*Hbar
	stepSigners            map[int32][]PrivateKey
	stepFeePayers          map[int32]*AccountID
	nodeAccountIDs         []AccountID
}

type JsonObject struct {
	Object string `json:"object"`
}

func NewContractHelper(bytecode string, constructorParameters ContractFunctionParameters, client *Client) *ContractHelper {
	response, err := NewContractCreateFlow().
		SetBytecodeWithString(bytecode).
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
			contractID:             *receipt.ContractID,
			stepResultValidators:   make(map[int32]func(ContractFunctionResult) bool),
			stepParameterSuppliers: make(map[int32]func() *ContractFunctionParameters),
			stepPayableAmounts:     make(map[int32]*Hbar),
			stepSigners:            make(map[int32][]PrivateKey),
			stepFeePayers:          make(map[int32]*AccountID),
			nodeAccountIDs:         make([]AccountID, 0),
		}
	}

	return &ContractHelper{}
}

func GetJsonResource(jsonBytes []byte) (*JsonObject, error) {
	var jsonObject JsonObject
	err := json.Unmarshal(jsonBytes, &jsonObject)
	if err != nil {
		return nil, err
	}

	return &jsonObject, nil
}

func (this *ContractHelper) SetResultValidatorForStep(stepIndex int32, validator func(ContractFunctionResult) bool) *ContractHelper {
	this.stepResultValidators[stepIndex] = validator
	return this
}

func (this *ContractHelper) SetParameterSupplierForStep(stepIndex int32, supplier func() *ContractFunctionParameters) *ContractHelper {
	this.stepParameterSuppliers[stepIndex] = supplier
	return this
}

func (this *ContractHelper) SetPayableAmountForStep(stepIndex int32, amount Hbar) *ContractHelper {
	this.stepPayableAmounts[stepIndex] = &amount
	return this
}

func (this *ContractHelper) AddSignerForStep(stepIndex int32, signer PrivateKey) *ContractHelper {
	if _, ok := this.stepSigners[stepIndex]; ok {
		this.stepSigners[stepIndex] = append(this.stepSigners[stepIndex], signer)
	} else {
		this.stepSigners[stepIndex] = make([]PrivateKey, 0)
		this.stepSigners[stepIndex] = append(this.stepSigners[stepIndex], signer)
	}

	return this
}

func (this *ContractHelper) SetFeePayerForStep(stepIndex int32, account AccountID, accountKey PrivateKey) *ContractHelper {
	this.stepFeePayers[stepIndex] = &account
	return this.AddSignerForStep(stepIndex, accountKey)
}

func (this *ContractHelper) GetResultValidator(stepIndex int32) func(ContractFunctionResult) bool {
	if _, ok := this.stepResultValidators[stepIndex]; ok {
		return this.stepResultValidators[stepIndex]
	}

	return func(result ContractFunctionResult) bool {
		responseStatus := Status(result.GetInt32(0))
		isValid := responseStatus == StatusSuccess
		if !isValid {
			println("Encountered invalid response status", responseStatus.String())
		}
		return isValid
	}
}

func (this *ContractHelper) GetParameterSupplier(stepIndex int32) func() *ContractFunctionParameters {
	if _, ok := this.stepParameterSuppliers[stepIndex]; ok {
		return this.stepParameterSuppliers[stepIndex]
	}

	return func() *ContractFunctionParameters {
		return nil
	}
}

func (this *ContractHelper) GetPayableAmount(stepIndex int32) *Hbar {
	return this.stepPayableAmounts[stepIndex]
}

func (this *ContractHelper) GetSigners(stepIndex int32) []PrivateKey {
	if _, ok := this.stepSigners[stepIndex]; ok {
		return this.stepSigners[stepIndex]
	}

	return []PrivateKey{}
}

func (this *ContractHelper) SetNodeAccountIDs(accountIDs []AccountID) *ContractHelper {
	this.nodeAccountIDs = accountIDs
	return this
}

func (this *ContractHelper) GetNodeAccountIDs() []AccountID {
	return this.nodeAccountIDs
}

func (this *ContractHelper) ExecuteSteps(firstStep int32, lastStep int32, client *Client) (*ContractHelper, error) {
	for stepIndex := firstStep; stepIndex <= lastStep; stepIndex++ {
		println("Attempting to execuite step", stepIndex)

		transaction := NewContractExecuteTransaction().
			SetContractID(this.contractID).
			SetGas(10000000)
		if len(this.nodeAccountIDs) > 0 {
			transaction.SetNodeAccountIDs(this.nodeAccountIDs)
		}

		payableAmount := this.GetPayableAmount(stepIndex)
		if payableAmount != nil {
			println("pay")
			transaction.SetPayableAmount(*payableAmount)
		}

		functionName := "step" + string(stepIndex)
		println(functionName)
		parameters := this.GetParameterSupplier(stepIndex)()
		if parameters != nil {
			transaction.SetFunction(functionName, parameters)
		} else {
			transaction.SetFunction(functionName, nil)
		}

		feePayerAccountID := this.stepFeePayers[stepIndex]
		if feePayerAccountID != nil {
			println("payer")
			transaction.SetTransactionID(TransactionIDGenerate(*feePayerAccountID))
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
			println("execute")
			return &ContractHelper{}, err
		}

		record, err := response.GetRecord(client)
		if err != nil {
			println("record")
			return &ContractHelper{}, err
		}

		functionResult, err := record.GetContractExecuteResult()
		if err != nil {
			return &ContractHelper{}, err
		}

		if this.GetResultValidator(stepIndex)(functionResult) {
			fmt.Printf("Step %d completed, and returned valid result. (TransactionId %s)", stepIndex, record.TransactionID.String())
		} else {
			return &ContractHelper{}, errors.New(fmt.Sprintf("Step %d returned invalid result", stepIndex))
		}
	}

	return this, nil
}

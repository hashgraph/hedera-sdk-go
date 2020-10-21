package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAccountStakersQuery_Execute(t *testing.T) {
	client, err := ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	_, err = NewAccountStakersQuery().
		SetAccountID(client.GetOperatorID()).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.Error(t, err)
}

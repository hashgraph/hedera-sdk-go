package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
)

func main() {
	var client *hedera.Client
	var err error

	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	client.SetOperator(operatorAccountID, operatorKey)

	resp, err := hedera.NewTokenCreateTransaction().
	SetTokenName("ffff").
	SetTokenSymbol("F").
	SetTokenMemo("fnord").
	SetTokenType(hedera.TokenTypeNonFungibleUnique).
	SetTreasuryAccountID(client.GetOperatorAccountID()).
	SetAdminKey(client.GetOperatorPublicKey()).
	SetFreezeKey(client.GetOperatorPublicKey()).
	SetWipeKey(client.GetOperatorPublicKey()).
	SetKycKey(client.GetOperatorPublicKey()).
	SetSupplyKey(client.GetOperatorPublicKey()).
	SetCustomFees([]hedera.Fee{
		hedera.CustomRoyaltyFee{
			CustomFee: hedera.CustomFee{
				FeeCollectorAccountID: &operatorAccountID,
			},
			Numerator:   1,
			Denominator: 20,
			FallbackFee: &hedera.CustomFixedFee{
				CustomFee: hedera.CustomFee{
					FeeCollectorAccountID: &operatorAccountID,
				},
				Amount: 10,
			},
		},
	}).
	SetFreezeDefault(false).
	Execute(client)

	if err != nil {
		panic(err)
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		panic(err);
	}
}
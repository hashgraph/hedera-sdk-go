package hedera

//func TestTokenFeeScheduleUpdateTransaction_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	resp, err := NewTokenCreateTransaction().
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetTokenName("ffff").
//		SetTokenSymbol("F").
//		SetDecimals(3).
//		SetInitialSupply(1000000).
//		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		SetFreezeKey(env.Client.GetOperatorPublicKey()).
//		SetWipeKey(env.Client.GetOperatorPublicKey()).
//		SetKycKey(env.Client.GetOperatorPublicKey()).
//		SetSupplyKey(env.Client.GetOperatorPublicKey()).
//		SetFreezeDefault(false).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	receipt, err := resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	tokenID := *receipt.TokenID
//
//	resp, err = NewTokenUpdateTransaction().
//		SetTokenID(tokenID).
//		SetTokenSymbol("A").
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	customFee := CustomFee{
//		Fee:                   FixedFee{
//			Amount:              0,
//			DenominationTokenID: &tokenID,
//		},
//		FeeCollectorAccountID: &env.OperatorID,
//	}
//
//	resp, err = NewTokenFeeScheduleUpdateTransaction().
//		SetTokenID(tokenID).
//		AddCustomFee(customFee).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	info, err := NewTokenInfoQuery().
//		SetTokenID(tokenID).
//		SetNodeAccountIDs([]AccountID{resp.NodeID}).
//		Execute(env.Client)
//	assert.NoError(t, err)
//	assert.True(t, len(info.CustomFees) > 0)
//}

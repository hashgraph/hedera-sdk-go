package hedera

type AccountBalance struct {
	Hbar Hbar
	Token *[]TokenBalance
}

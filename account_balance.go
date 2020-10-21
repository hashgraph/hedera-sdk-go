package hedera

type AccountBalance struct {
	Hbar Hbar
	Token map[TokenID]int64
}

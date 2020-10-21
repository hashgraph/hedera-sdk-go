package hedera

type AccountBalance struct {
	Hbars Hbar
	Token map[TokenID]uint64
}

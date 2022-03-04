package hedera

import "github.com/pkg/errors"

type VerifySignatureFlow struct {
	Transaction
	key            *PublicKey
	accountID      *AccountID
	nodeAccountIDs []AccountID
}

func NewVerifySignatureFlow() *VerifySignatureFlow {
	transaction := VerifySignatureFlow{
		Transaction: _NewTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(20))

	return &transaction
}

func (query *VerifySignatureFlow) SetAccountID(id AccountID) *VerifySignatureFlow {
	query._RequireNotFrozen()
	query.accountID = &id
	return query
}

func (query *VerifySignatureFlow) GetBytecode() AccountID {
	if query.accountID != nil {
		return *query.accountID
	}

	return AccountID{}
}

func (query *VerifySignatureFlow) SetKey(key PublicKey) *VerifySignatureFlow {
	query._RequireNotFrozen()
	query.key = &key
	return query
}

func (query *VerifySignatureFlow) GetKey() PublicKey {
	if query.key != nil {
		return *query.key
	}

	return PublicKey{}
}

func (query *VerifySignatureFlow) _CreateAccountInfoQuery(client *Client) *AccountInfoQuery {
	if client == nil {
		return &AccountInfoQuery{}
	}
	if query.accountID != nil {
		accountInfoQuery := NewAccountInfoQuery().
			SetAccountID(*query.accountID).
			SetMaxQueryPayment(NewHbar(1)).
			SetQueryPayment(HbarFromTinybar(25))

		if len(query.nodeAccountIDs) > 0 {
			accountInfoQuery.SetNodeAccountIDs(query.nodeAccountIDs)
		}

		return accountInfoQuery
	}

	return &AccountInfoQuery{}
}

func (query *VerifySignatureFlow) Execute(client *Client) (bool, error) {
	accountInfo, err := query._CreateAccountInfoQuery(client).
		Execute(client)
	if err != nil {
		return false, err
	}

	key := accountInfo.Key

	switch keyType := key.(type) {
	case *PublicKey:
		if keyType.String() == query.key.String() {
			return true, nil
		}

		return false, nil
	case PublicKey:
		if keyType.String() == query.key.String() {
			return true, nil
		}

		return false, nil
	case *KeyList:
		found := false
		for _, k := range keyType.keys {
			switch keyType := k.(type) {
			case *PublicKey:
				if keyType.String() == query.key.String() {
					found = true
				}
			case PublicKey:
				if keyType.String() == query.key.String() {
					found = true
				}
			}
		}

		return found, nil
	default:
		return false, errors.New("unsupported key type")
	}
}

func (query *VerifySignatureFlow) SetNodeAccountIDs(nodeID []AccountID) *VerifySignatureFlow {
	query._RequireNotFrozen()
	query.nodeAccountIDs = nodeID
	return query
}

func (query *VerifySignatureFlow) GetNodeAccountIDs() []AccountID {
	return query.nodeAccountIDs
}

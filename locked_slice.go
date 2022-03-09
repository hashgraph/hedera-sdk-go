package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type _LockedSlice struct {
	slice  []interface{}
	locked bool
	index  int
}

func _NewLockedSlice() *_LockedSlice {
	return &_LockedSlice{
		slice:  make([]interface{}, 0),
		locked: false,
		index:  0,
	}
}

func (s *_LockedSlice) _SetSlice(slic []interface{}) (*_LockedSlice, error) { //nolint
	if s.locked {
		return &_LockedSlice{}, errLockedSlice
	}
	s.slice = slic
	s.index = 0

	return s, nil
}

func (s *_LockedSlice) _PushNodeAccountIDs(items ...AccountID) (*_LockedSlice, error) { //nolint
	if s.locked {
		return &_LockedSlice{}, errLockedSlice
	}
	for _, k := range items {
		if len(s.slice) > 0 {
			switch s.slice[0].(type) { //nolint
			case AccountID:
				s.slice = append(s.slice, k)
			}
		} else {
			s.slice = append(s.slice, k)
		}
	}

	return s, nil
}

func (s *_LockedSlice) _PushTransactionIDs(items ...TransactionID) (*_LockedSlice, error) { //nolint
	if s.locked {
		return &_LockedSlice{}, errLockedSlice
	}
	for _, k := range items {
		if len(s.slice) > 0 {
			switch s.slice[0].(type) { //nolint
			case TransactionID:
				s.slice = append(s.slice, k)
			}
		} else {
			s.slice = append(s.slice, k)
		}
	}

	return s, nil
}

func (s *_LockedSlice) _PushTransactions(items ...*services.Transaction) (*_LockedSlice, error) { //nolint
	if s.locked {
		return &_LockedSlice{}, errLockedSlice
	}
	for _, k := range items {
		if len(s.slice) > 0 {
			switch s.slice[0].(type) { //nolint
			case *services.Transaction:
				s.slice = append(s.slice, k)
			}
		} else {
			s.slice = append(s.slice, k)
		}
	}

	return s, nil
}

func (s *_LockedSlice) _PushSignedTransactions(items ...*services.SignedTransaction) (*_LockedSlice, error) { //nolint
	if s.locked {
		return &_LockedSlice{}, errLockedSlice
	}
	for _, k := range items {
		if len(s.slice) > 0 {
			switch s.slice[0].(type) { //nolint
			case *services.SignedTransaction:
				s.slice = append(s.slice, k)
			}
		} else {
			s.slice = append(s.slice, k)
		}
	}

	return s, nil
}

func (s *_LockedSlice) _Clear() (*_LockedSlice, error) { //nolint
	if s.locked {
		return &_LockedSlice{}, errLockedSlice
	}
	s.slice = make([]interface{}, 0)
	return s, nil
}

func (s *_LockedSlice) _Get(index int) interface{} { //nolint
	return s.slice[index]
}

func (s *_LockedSlice) _Set(index int, item interface{}) (*_LockedSlice, error) { //nolint
	if s.locked {
		return &_LockedSlice{}, errLockedSlice
	}
	if len(s.slice) > 0 {
		switch i := item.(type) {
		case TransactionID:
			switch s.slice[0].(type) { //nolint
			case TransactionID:
				s.slice[index] = i
			}
		case *services.Transaction:
			switch s.slice[0].(type) { //nolint
			case *services.Transaction:
				s.slice[index] = i
			}
		case *services.SignedTransaction:
			switch s.slice[0].(type) { //nolint
			case *services.SignedTransaction:
				s.slice[index] = i
			}
		case AccountID:
			switch s.slice[0].(type) { //nolint
			case AccountID:
				s.slice[index] = i
			}
		}
	} else {
		s.slice = append(s.slice, item)
	}

	return s, nil
}

func (s *_LockedSlice) _SetIfAbsent(index int32, item interface{}) (*_LockedSlice, error) { //nolint
	if s.locked {
		return &_LockedSlice{}, errLockedSlice
	}
	if int32(s._Length()) > index {
		if s.slice[index] == nil {
			s.slice[index] = item
		}
	}

	return s, nil
}

func (s *_LockedSlice) _GetNext() interface{} { //nolint
	return s._Get(s._Advance())
}

func (s *_LockedSlice) _GetCurrent() interface{} { //nolint
	index := s.index - 1
	if index < 0 {
		index = s._Length() - 1
	}

	return s._Get(index)
}

func (s *_LockedSlice) _Advance() int { //nolint
	index := s.index
	s.index = (s.index + 1) % len(s.slice)
	return index
}

func (s *_LockedSlice) _IsEmpty() bool { //nolint
	return len(s.slice) == 0
}

func (s *_LockedSlice) _Length() int { //nolint
	return len(s.slice)
}

func (s *_LockedSlice) _GetTransactionIDs() []TransactionID { //nolint
	temp := make([]TransactionID, 0)
	if s._Length() > 0 {
		for _, k := range s.slice {
			switch i := k.(type) { //nolint
			case TransactionID:
				temp = append(temp, i)
			}
		}
	}

	return temp
}

func (s *_LockedSlice) _GetTransactions() []*services.Transaction { //nolint
	temp := make([]*services.Transaction, 0)
	if s._Length() > 0 {
		for _, k := range s.slice {
			switch i := k.(type) { //nolint
			case *services.Transaction:
				temp = append(temp, i)
			}
		}
	}

	return temp
}

func (s *_LockedSlice) _GetSignedTransactions() []*services.SignedTransaction { //nolint
	temp := make([]*services.SignedTransaction, 0)
	if s._Length() > 0 {
		for _, k := range s.slice {
			switch i := k.(type) { //nolint
			case *services.SignedTransaction:
				temp = append(temp, i)
			}
		}
	}

	return temp
}

func (s *_LockedSlice) _GetNodeAccountIDs() []AccountID { //nolint
	temp := make([]AccountID, 0)
	if s._Length() > 0 {
		for _, k := range s.slice {
			switch i := k.(type) { //nolint
			case AccountID:
				temp = append(temp, i)
			}
		}
	}

	return temp
}

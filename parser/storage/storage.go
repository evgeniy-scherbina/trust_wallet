package storage

import lib_eth "github.com/evgeniy-scherbina/trust_wallet/lib/eth"

type InMemoryStorage struct {
	subscribedAddresses map[string]struct{}

	addrToTxs map[string][]*lib_eth.Transaction
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		subscribedAddresses: make(map[string]struct{}),
		addrToTxs:           make(map[string][]*lib_eth.Transaction),
	}
}

func (s *InMemoryStorage) Subscribe(address string) {
	s.subscribedAddresses[address] = struct{}{}
}

func (s *InMemoryStorage) GetTransactions(address string) []*lib_eth.Transaction {
	return s.addrToTxs[address]
}

func (s *InMemoryStorage) subscribed(address string) bool {
	_, ok := s.subscribedAddresses[address]
	return ok
}

func (s *InMemoryStorage) ProcessTx(tx *lib_eth.Transaction) {
	if s.subscribed(tx.From) {
		s.addrToTxs[tx.From] = append(s.addrToTxs[tx.From], tx)
	}

	if s.subscribed(tx.To) {
		s.addrToTxs[tx.To] = append(s.addrToTxs[tx.To], tx)
	}
}

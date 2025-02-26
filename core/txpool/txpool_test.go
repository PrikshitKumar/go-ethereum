package txpool

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/event"
	"github.com/stretchr/testify/assert"
)

// Mock subpool that filters all transactions as valid.
type mockSubpool struct{}

func (m *mockSubpool) Init(blockNumber uint64, header *types.Header, reserver AddressReserver) error {
	return nil
}
func (m *mockSubpool) Filter(tx *types.Transaction) bool {
	return true
}
func (m *mockSubpool) Add(txs []*types.Transaction, sync bool) []error {
	errs := make([]error, len(txs))
	for i := range txs {
		errs[i] = nil // Assuming all transactions are accepted
	}
	return errs
}
func (m *mockSubpool) Clear()                    {}
func (m *mockSubpool) Remove(txHash common.Hash) {}
func (m *mockSubpool) Contains(txHash common.Hash) bool {
	return false
}
func (m *mockSubpool) Pending(filter PendingFilter) map[common.Address][]*LazyTransaction {
	return make(map[common.Address][]*LazyTransaction)
}
func (m *mockSubpool) Close() error {
	return nil
}
func (m *mockSubpool) Promote()                                           {}
func (m *mockSubpool) Demote()                                            {}
func (m *mockSubpool) Reset(oldHead *types.Header, newHead *types.Header) {}
func (m *mockSubpool) SetGasTip(gasTip *big.Int)                          {}

func (m *mockSubpool) Content() (map[common.Address][]*types.Transaction, map[common.Address][]*types.Transaction) {
	return make(map[common.Address][]*types.Transaction), make(map[common.Address][]*types.Transaction)
}
func (m *mockSubpool) ContentFrom(addr common.Address) ([]*types.Transaction, []*types.Transaction) {
	return nil, nil
}
func (m *mockSubpool) Get(txHash common.Hash) *types.Transaction {
	return nil
}
func (m *mockSubpool) GetBlobs(hashes []common.Hash) ([]*kzg4844.Blob, []*kzg4844.Proof) {
	return nil, nil
}
func (m *mockSubpool) Has(txHash common.Hash) bool {
	return false
}
func (m *mockSubpool) Nonce(addr common.Address) uint64 {
	return 0
}
func (m *mockSubpool) Stats() (int, int) {
	return 0, 0
}
func (m *mockSubpool) Status(hash common.Hash) TxStatus {
	return 0
}
func (m *mockSubpool) SubscribeTransactions(ch chan<- core.NewTxsEvent, flag bool) event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		<-quit
		return nil
	})
}

func TestTxPool_Add(t *testing.T) {
	// mock subpool
	subpool := &mockSubpool{}

	// TxPool with a single subpool
	pool := &TxPool{
		subpools: []SubPool{subpool},
	}

	// Test accounts initialization
	var (
		// Blocked Account
		blockedKey, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")

		// Valid Account
		validKey, _ = crypto.GenerateKey()

		// Recipient Account
		toKey, _ = crypto.GenerateKey()
		toAddr   = crypto.PubkeyToAddress(toKey.PublicKey)
	)

	chainID := big.NewInt(1) // Dummmy chain ID

	blockedTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     0,
		GasTipCap: big.NewInt(2e9),
		GasFeeCap: big.NewInt(3e9),
		Gas:       21000,
		To:        &toAddr,
		Value:     big.NewInt(1),
		Data:      nil,
	})

	validTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     0,
		GasTipCap: big.NewInt(2e9),
		GasFeeCap: big.NewInt(3e9),
		Gas:       21000,
		To:        &toAddr,
		Value:     big.NewInt(1),
		Data:      nil,
	})

	// Sign transactions
	signer := types.LatestSignerForChainID(chainID)

	signedBlockedTx, err := types.SignTx(blockedTx, signer, blockedKey)
	if err != nil {
		panic(fmt.Sprintf("Signing blocked transaction failed: %v", err))
	}
	signedValidTx, err := types.SignTx(validTx, signer, validKey)
	if err != nil {
		panic(fmt.Sprintf("Signing valid transaction failed: %v", err))
	}

	// Run test cases
	tests := []struct {
		name        string
		txs         []*types.Transaction
		expectError bool
	}{
		{"Blocked Address", []*types.Transaction{signedBlockedTx}, true},
		{"Valid Address", []*types.Transaction{signedValidTx}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := pool.Add(tt.txs, true)
			if tt.expectError {
				expectedErr := fmt.Sprintf("transaction from blacklisted address %s rejected", blockedAddress.Hex())
				assert.EqualError(t, errs[0], expectedErr)
			} else {
				assert.NoError(t, errs[0])
			}
		})
	}
}

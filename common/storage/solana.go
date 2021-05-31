package storage

import "fmt"

func (storage *Storage) SetSolanaRecentBlock(round int, blockHash []byte) error {
	return storage.setValue([]byte(fmt.Sprintf("solana_recent_block_%d", round)), blockHash)
}

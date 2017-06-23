package core

import (
	"sync"

	"github.com/exced/blockchain/crypto"
)

// Blockchain is a list of blocks.
type Blockchain struct {
	Mutex  *sync.Mutex
	Blocks []*Block
}

// NewBlockchain build a new blockchain from the genesis block.
func NewBlockchain() *Blockchain {
	return &Blockchain{Mutex: &sync.Mutex{}, Blocks: []*Block{genesisBlock}}
}

// Save Blockchain
func (bc *Blockchain) Save(file string) error {
	return Save(file, bc)
}

func (bc Blockchain) Len() int           { return len(bc.Blocks) }
func (bc Blockchain) Swap(i, j int)      { bc.Blocks[i], bc.Blocks[j] = bc.Blocks[j], bc.Blocks[i] }
func (bc Blockchain) Less(i, j int) bool { return bc.Blocks[i].Index < bc.Blocks[j].Index }

// GetLastBlock retrieves the last block of the blockchain
func (bc *Blockchain) GetLastBlock() *Block {
	return bc.Blocks[bc.Len()-1]
}

func (bc *Blockchain) getGenesis() *Block {
	return bc.Blocks[0]
}

// Append given tail blockchain to received blockchain
func (bc *Blockchain) Append(tail *Blockchain) {
	for _, block := range tail.Blocks {
		bc.Blocks = append(bc.Blocks, block)
	}
}

// AppendBlock append a new block at the end of the blockchain
func (bc *Blockchain) AppendBlock(b *Block) {
	bc.Mutex.Lock()
	bc.Blocks = append(bc.Blocks, b)
	bc.Mutex.Unlock()
}

// IsValid tests if all blocks of the blockchain are valid
func (bc *Blockchain) IsValid() bool {
	if bc.Blocks[0] != genesisBlock {
		return false
	}
	pb := bc.Blocks[0]
	for i := 1; i < bc.Len(); i++ {
		if bc.Blocks[i].IsValid(pb) {
			pb = bc.Blocks[i]
		} else {
			return false
		}
	}
	return true
}

// IsBlockValid tests if given block is valid with the last block of the blockchain
func (bc *Blockchain) IsBlockValid(b *Block) bool {
	return b.IsValid(bc.GetLastBlock())
}

// Link given block to last block of received blockchain : set the previousHash field to previous Hash
func (bc *Blockchain) Link(b *Block) {
	b.Link(bc.GetLastBlock())
}

// Mine looks for a nonce for the last block of received blockchain
func (bc *Blockchain) Mine() {
	bc.Mutex.Lock()
	bc.GetLastBlock().Mine()
	bc.Mutex.Unlock()
}

// GenNext returns the next block to work on for miners.
func (bc *Blockchain) GenNext(t *Transactions) *Block {
	bc.Mutex.Lock()
	res := bc.GetLastBlock().GenNext(t)
	bc.Mutex.Unlock()
	return res
}

// Fetch returns a blockchain fetching our received blockchain and other given blockchain. Does nothing
// if other blockchain is not valid.
func (bc *Blockchain) Fetch(other *Blockchain) *Blockchain {
	if other.Len() < bc.Len() {
		return bc
	}
	if !other.IsValid() {
		return bc
	}
	for i, block := range bc.Blocks {
		if block != other.Blocks[i] {
			return bc
		}
	}
	return other
}

// PoW returns true if Proof of Work is done for received blockchain and given difficulty.
func (bc *Blockchain) PoW(difficulty int) bool {
	bc.Mutex.Lock()
	res := crypto.MatchHash(bc.GetLastBlock().ToHash(), difficulty)
	bc.Mutex.Unlock()
	return res
}

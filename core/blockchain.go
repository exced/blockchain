package core

// Blockchain is a list of blocks.
type Blockchain []*Block

// NewBlockchain build a new blockchain from the genesis block.
func NewBlockchain() *Blockchain {
	return &Blockchain{genesisBlock}
}

// Save Blockchain
func (bc *Blockchain) Save(file string) error {
	return Save(file, bc)
}

func (bc Blockchain) Len() int           { return len(bc) }
func (bc Blockchain) Swap(i, j int)      { bc[i], bc[j] = bc[j], bc[i] }
func (bc Blockchain) Less(i, j int) bool { return bc[i].Index < bc[j].Index }

// GetLastBlock retrieves the last block of the blockchain
func (bc *Blockchain) GetLastBlock() *Block {
	return (*bc)[bc.Len()-1]
}

func (bc *Blockchain) getGenesis() *Block {
	return (*bc)[0]
}

// Append given tail to received blockchain
func (bc *Blockchain) Append(tail *Blockchain) {
	for _, block := range *tail {
		*bc = append(*bc, block)
	}
}

// AppendBlock append a new block at the end of the blockchain
func (bc *Blockchain) AppendBlock(b *Block) {
	*bc = append(*bc, b)
}

// IsValid tests if all blocks of the blockchain are valid
func (bc *Blockchain) IsValid() bool {
	if (*bc)[0] != genesisBlock {
		return false
	}
	pb := (*bc)[0]
	for i := 1; i < bc.Len(); i++ {
		if (*bc)[i].IsValid(pb) {
			pb = (*bc)[i]
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

// Mine looks for a nonce for the last block of received blockchain to satisfy given difficulty
func (bc *Blockchain) Mine(difficulty int) *Block {
	return bc.GetLastBlock().Mine(difficulty)
}

func (bc *Blockchain) GenNext(transactions []*Transaction) *Block {
	return bc.GetLastBlock().GenNext(transactions)
}

func (bc *Blockchain) Fetch(other *Blockchain) {
	if other.Len() < bc.Len() {
		return
	}
	if !other.IsValid() {
		return
	}
	for i, block := range *other {
		if block != (*bc)[i] {
			return
		}
	}
	copy(*bc, *other)
}

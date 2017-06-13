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

// OpenBlockchainFile open blockchain file
func OpenBlockchainFile(file string) error {
	var blockchain = &Blockchain{}
	return Load(file, blockchain)
}

func (bc Blockchain) Len() int           { return len(bc) }
func (bc Blockchain) Swap(i, j int)      { bc[i], bc[j] = bc[j], bc[i] }
func (bc Blockchain) Less(i, j int) bool { return bc[i].Index < bc[j].Index }

func (bc *Blockchain) getLastBlock() *Block {
	return (*bc)[bc.Len()-1]
}

func (bc *Blockchain) getGenesis() *Block {
	return (*bc)[0]
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
		if (*bc)[i].isValid(pb) {
			pb = (*bc)[i]
		} else {
			return false
		}
	}
	return true
}

func (bc *Blockchain) IsTransactionValid(t *Transaction) bool {
	return bc.getLastBlock().IsTransactionValid(t)
}

// Mine looks for a nonce for the last block of received blockchain to satisfy given difficulty
func (bc *Blockchain) Mine(difficulty int) *Block {
	return bc.getLastBlock().Mine(difficulty)
}

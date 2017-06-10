package core

// Blockchain is a list of blocks.
type Blockchain []*Block

// NewBlockchain build a new blockchain from the genesis block.
func NewBlockchain() *Blockchain {
	return &Blockchain{genesisBlock}
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

// AddBlock add a new block at the end of the blockchain
func (bc *Blockchain) AddBlock(data string) {
	nb := bc.getLastBlock().genNext(data)
	*bc = append(*bc, nb)
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

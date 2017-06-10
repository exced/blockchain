package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

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

// GenBlockchainFile generate Blockchain copy and store it in a file
func GenBlockchainFile(path string) error {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
		})
	return ioutil.WriteFile(path, pemdata, 0644)
}

func OpenBlockchainFile(path string) (*Blockchain, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(f)
}

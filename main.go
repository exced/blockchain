package main

import (
	"fmt"

	"github.com/exced/simple-blockchain/core"
)

func main() {
	blockchain := core.NewBlockchain()
	fmt.Println((*blockchain)[0])
	blockchain.AddBlock("hey")
	fmt.Println((*blockchain)[1])
	fmt.Println(*blockchain)
	fmt.Println("is valid: %#v ", blockchain.IsValid())
}

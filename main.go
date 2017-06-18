package main

import (
	"fmt"

	"github.com/exced/blockchain/core"
)

func main() {
	bc := core.NewBlockchain()
	fmt.Println(bc.GetLastBlock().ToHash())
}
